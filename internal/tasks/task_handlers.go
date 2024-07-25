package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/pkg/address"
	clientAsynq "github.com/vultisig/airdrop-registry/pkg/asynq"
	"github.com/vultisig/airdrop-registry/pkg/balance"
	"github.com/vultisig/airdrop-registry/pkg/price"
)

func ProcessBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p BalanceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching balance for vault: ecdsa=%s, eddsa=%s, chain=%s, address=%s", p.ECDSA, p.EDDSA, p.Chain, p.Address)

	balanceAmount, err := balance.FetchBalanceOfAddress(p.Chain, p.Address)
	if err != nil {
		return fmt.Errorf("services.GetBalanceOfAddress failed: %v", err)
	}

	token, err := balance.GetBaseTokenByChain(p.Chain)
	if err != nil {
		return fmt.Errorf("balances.GetBaseTokenByChain failed: %v", err)
	}

	b := &models.Balance{
		ECDSA:   p.ECDSA,
		EDDSA:   p.EDDSA,
		Chain:   p.Chain,
		Address: p.Address,
		Balance: balanceAmount,
		Token:   token,
		Date:    time.Now().Unix(),
	}

	if balanceAmount > 0 {
		_, err = services.SaveBalanceWithLatestPrice(b)
		if err != nil {
			return fmt.Errorf("services.SaveBalance failed: %v", err)
		}
	}

	log.Printf("Balance for vault: ecdsa=%s, eddsa=%s, chain=%s, address=%s is %f", p.ECDSA, p.EDDSA, p.Chain, p.Address, balanceAmount)

	result := map[string]interface{}{
		"native_balance": balanceAmount,
	}

	if p.Chain == "ethereum" || p.Chain == "bsc" || p.Chain == "optimism" {
		tokenBalances, err := balance.FetchTokensWithBalance(p.Address, p.Chain)
		if err != nil {
			return fmt.Errorf("balance.FetchTokensWithBalance failed: %v", err)
		}

		tokenAddresses := make([]string, 0, len(tokenBalances))
		for tokenAddress := range tokenBalances {
			tokenAddresses = append(tokenAddresses, tokenAddress)
		}

		tokenInfoMap, err := balance.GetTokenInfo(tokenAddresses, p.Chain)
		if err != nil {
			return fmt.Errorf("GetTokenInfo failed: %v", err)
		}

		for tokenAddress, tokenBalance := range tokenBalances {
			tokenBalanceFloat, err := strconv.ParseFloat(tokenBalance, 64)
			if err != nil {
				return fmt.Errorf("strconv.ParseFloat failed: %v", err)
			}

			tokenInfo, ok := tokenInfoMap[tokenAddress]
			if !ok {
				return fmt.Errorf("token info not found for address: %s", tokenAddress)
			}

			adjustedBalance := tokenBalanceFloat / math.Pow(10, float64(tokenInfo.Decimals))

			tb := &models.Balance{
				ECDSA:   p.ECDSA,
				EDDSA:   p.EDDSA,
				Chain:   p.Chain,
				Address: p.Address,
				Balance: adjustedBalance,
				Token:   tokenAddress,
				Date:    time.Now().Unix(),
			}

			_, err = services.SaveBalanceWithLatestPrice(tb)
			if err != nil {
				return fmt.Errorf("services.SaveBalance failed: %v", err)
			}
		}
		result["token_balances"] = tokenBalances
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}

	if _, err := t.ResultWriter().Write(resultBytes); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %v", err)
	}

	return nil
}

func ProcessPointsCalculationTask(ctx context.Context, t *asynq.Task) error {
	var p PointsCalculationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v, %w", err, asynq.SkipRetry)
	}
	log.Printf("Calculating points for Vault: ecdsa=%s, eddsa=%s", p.ECDSA, p.EDDSA)

	vaults, err := services.GetAllVaults()
	if err != nil {
		return fmt.Errorf("services.GetAllVaults failed: %v", err)
	}

	// @TODO: Check when the latest scan was done instead of last 2 hours
	timeSince := time.Now().Add(-2 * time.Hour)

	totalUSD := 0.0
	for _, vault := range vaults {
		averageBalance, err := services.GetAverageBalanceSince(vault.ECDSA, vault.EDDSA, timeSince)
		if err != nil {
			return fmt.Errorf("services.GetLatestBalancesByVaultKeys failed: %v", err)
		}
		totalUSD += averageBalance
	}

	averageBalance, err := services.GetAverageBalanceSince(p.ECDSA, p.EDDSA, timeSince)
	if err != nil {
		return fmt.Errorf("services.GetLatestBalancesByVaultKeys failed: %v", err)
	}

	share := float64((averageBalance / totalUSD) * 100)

	point := &models.Point{
		ECDSA:   p.ECDSA,
		EDDSA:   p.EDDSA,
		Balance: averageBalance,
		Share:   share,
	}

	err = services.SavePoint(point)
	if err != nil {
		return fmt.Errorf("services.SavePoint failed: %v", err)
	}

	log.Printf("Point share for Vault: ecdsa=%s, eddsa=%s is %f", p.ECDSA, p.EDDSA, share)

	result := map[string]interface{}{
		"share":   share,
		"balance": averageBalance,
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}

	if _, err := t.ResultWriter().Write(resultBytes); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %v", err)
	}

	return nil
}

func ProcessPriceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p PriceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching coin price for token: chain=%s, token=%s", p.Chain, p.Token)

	usdPrice, err := price.FetchPrice(p.Chain, p.Token)
	if err != nil {
		return fmt.Errorf("price.FetchPrice failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Price for token: chain=%s, token=%s is %f\n", p.Chain, p.Token, usdPrice)

	price := &models.Price{
		Chain:  p.Chain,
		Token:  p.Token,
		Price:  usdPrice,
		Date:   time.Now().Unix(),
		Source: "coingecko",
	}

	err = services.SavePrice(price)
	if err != nil {
		return fmt.Errorf("services.SavePrice failed: %v", err)
	}

	result := map[string]interface{}{
		"price": usdPrice,
	}
	resultBytes, err := json.Marshal(result)
	if _, err := t.ResultWriter().Write([]byte(resultBytes)); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %v", err)
	}

	return nil
}

func ProcessPriceFetchAllActivePairsTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	asynqClient := clientAsynq.AsynqClient

	pairs, err := services.GetUniqueActiveChainTokenPairs()
	if err != nil {
		return fmt.Errorf("services.GetUniqueActiveChainTokenPairs failed: %v", err)
	}

	for _, pair := range pairs {
		fmt.Printf("Enqueuing price fetch task for pair: chain=%s, token=%s\n", pair.Chain, pair.Token)
		err := EnqueuePriceFetchTask(asynqClient, pair.Chain, pair.Token)
		if err != nil {
			return fmt.Errorf("failed to enqueue price fetch tasks: %v", err)
		}
		fmt.Printf("Enqueued price fetch task for pair: chain=%s, token=%s\n", pair.Chain, pair.Token)
	}

	log.Printf("Enqueued price fetch tasks for all active pairs")

	return nil
}

func ProcessBalanceFetchAllTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	asynqClient := clientAsynq.AsynqClient

	vaults, err := services.GetAllVaults()
	if err != nil {
		return fmt.Errorf("services.GetVaults failed: %v", err)
	}

	for _, vault := range vaults {
		addresses, err := address.GenerateSupportedChainAddresses(vault.ECDSA, vault.HexChainCode, vault.EDDSA)
		if err != nil {
			return fmt.Errorf("address.GenerateSupportedChainAddresses failed: %v", err)
		}

		for chain, address := range addresses {
			fmt.Printf("Enqueuing balance fetch task for vault: ecdsa=%s, eddsa=%s, chain=%s, address=%s\n", vault.ECDSA, vault.EDDSA, chain, address)
			err := EnqueueBalanceFetchTask(asynqClient, vault.ECDSA, vault.EDDSA, chain, address)
			if err != nil {
				return fmt.Errorf("failed to enqueue balance fetch tasks: %v", err)
			}
		}
		fmt.Printf("Enqueued balance fetch tasks for vault: ecdsa=%s, eddsa=%s\n", vault.ECDSA, vault.EDDSA)
	}

	log.Printf("Enqueued balance fetch tasks for all vaults")

	return nil
}
