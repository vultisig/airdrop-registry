package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/pkg/address"
	clientAsynq "github.com/vultisig/airdrop-registry/pkg/asynq"
	"github.com/vultisig/airdrop-registry/pkg/balance"
	"github.com/vultisig/airdrop-registry/pkg/price"
	"github.com/vultisig/airdrop-registry/pkg/utils"
	"gorm.io/gorm"
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

	cycleID32, err := strconv.ParseUint(p.Cycle, 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseUint failed: %v", err)
	}
	cycleID := uint(cycleID32)

	cycle, err := services.GetCycleByID(cycleID)
	if err != nil {
		return fmt.Errorf("services.GetCycleByID failed: %v", err)
	}

	endTime := cycle.CreatedAt
	startTime := endTime.Add(-2 * time.Hour)

	log.Printf("Calculating points for Vault: ecdsa=%s, eddsa=%s, cycleID=%s, timeRange=%v to %v",
		p.ECDSA, p.EDDSA, p.Cycle, startTime, endTime)

	vaults, err := services.GetAllVaults()
	if err != nil {
		return fmt.Errorf("services.GetAllVaults failed: %v", err)
	}

	totalUSD := 0.0
	for _, vault := range vaults {
		averageBalance, err := services.GetAverageBalanceForTimeRange(vault.ECDSA, vault.EDDSA, startTime, endTime)
		if err != nil {
			return fmt.Errorf("services.GetAverageBalanceForTimeRange failed for vault %s: %v", vault.ECDSA, err)
		}
		totalUSD += averageBalance
	}

	averageBalance, err := services.GetAverageBalanceForTimeRange(p.ECDSA, p.EDDSA, startTime, endTime)
	if err != nil {
		return fmt.Errorf("services.GetAverageBalanceForTimeRange failed: %v", err)
	}

	share := utils.CalculateShare(averageBalance, totalUSD)

	point := &models.Point{
		ECDSA:   p.ECDSA,
		EDDSA:   p.EDDSA,
		Balance: averageBalance,
		Share:   share,
		CycleID: cycleID,
	}

	err = services.SavePoint(point)
	if err != nil {
		return fmt.Errorf("services.SavePoint failed: %v", err)
	}

	log.Printf("Point share for Vault: ecdsa=%s, eddsa=%s is %f for cycle %s", p.ECDSA, p.EDDSA, share, p.Cycle)

	result := map[string]interface{}{
		"share":   share,
		"balance": averageBalance,
		"cycle":   cycleID,
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

// Parent task processors

func ProcessBalanceFetchParentTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	client := clientAsynq.AsynqClient

	vaults, err := services.GetAllVaults()
	if err != nil {
		return fmt.Errorf("failed to get vaults: %v", err)
	}

	for _, vault := range vaults {
		addresses, err := address.GenerateSupportedChainAddresses(vault.ECDSA, vault.HexChainCode, vault.EDDSA)
		if err != nil {
			return fmt.Errorf("failed to generate addresses: %v", err)
		}

		for chain, addr := range addresses {
			if err := EnqueueBalanceFetchTask(client, vault.ECDSA, vault.EDDSA, chain, addr); err != nil {
				return fmt.Errorf("failed to enqueue balance fetch task: %v", err)
			}
		}
	}

	return nil
}

func ProcessPointsCalculationParentTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	client := clientAsynq.AsynqClient

	currentCycle, err := services.GetCurrentCycle()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get current cycle: %w", err)
	}

	now := time.Now()
	if currentCycle != nil {
		timeSinceLastCycle := now.Sub(currentCycle.CreatedAt)

		if timeSinceLastCycle < 24*time.Hour {
			return fmt.Errorf("too soon for a new cycle: the most recent cycle was created %v ago: %w", timeSinceLastCycle, asynq.SkipRetry)
		}

		r := rand.New(rand.NewSource(now.UnixNano()))
		randomValue := r.Float64()

		if timeSinceLastCycle < 36*time.Hour {
			if randomValue > 0.2 { // 20% chance
				return fmt.Errorf("randomly decided not to create a new cycle (20%% chance): the most recent cycle was created %v ago: %w", timeSinceLastCycle, asynq.SkipRetry)
			}
		} else {
			if randomValue > 0.5 { // 50% chance
				return fmt.Errorf("randomly decided not to create a new cycle (50%% chance): the most recent cycle was created %v ago: %w", timeSinceLastCycle, asynq.SkipRetry)
			}
		}
	}

	newCycle, err := services.CreateCycle(&models.Cycle{CreatedAt: now})
	if err != nil {
		return fmt.Errorf("failed to create new cycle: %w", err)
	}

	vaults, err := services.GetAllVaults()
	if err != nil {
		return fmt.Errorf("failed to get vaults: %w", err)
	}

	var enqueueErrors []error
	for _, vault := range vaults {
		if err := EnqueuePointsCalculationTask(client, vault.ECDSA, vault.EDDSA, fmt.Sprintf("%d", newCycle.ID)); err != nil {
			enqueueErrors = append(enqueueErrors, err)
		}
	}

	result := map[string]interface{}{
		"cycle_id":         newCycle.ID,
		"cycle_created_at": newCycle.CreatedAt,
		"vaults_processed": len(vaults),
		"enqueue_errors":   len(enqueueErrors),
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %w", err)
	}

	if _, err := t.ResultWriter().Write(resultBytes); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %w", err)
	}

	if len(enqueueErrors) > 0 {
		return fmt.Errorf("failed to enqueue some points calculation tasks: %v", enqueueErrors)
	}

	return nil
}

func ProcessPriceFetchParentTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	client := clientAsynq.AsynqClient

	pairs, err := services.GetUniqueActiveChainTokenPairs()
	if err != nil {
		return fmt.Errorf("failed to get active chain-token pairs: %v", err)
	}

	for _, pair := range pairs {
		if err := EnqueuePriceFetchTask(client, pair.Chain, pair.Token); err != nil {
			return fmt.Errorf("failed to enqueue price fetch task: %v", err)
		}
	}

	return nil
}
