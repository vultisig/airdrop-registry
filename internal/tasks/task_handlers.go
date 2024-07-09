package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	clientAsynq "github.com/vultisig/airdrop-registry/pkg/asynq"
	"github.com/vultisig/airdrop-registry/pkg/balance"
	"github.com/vultisig/airdrop-registry/pkg/price"
)

func ProcessBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p BalanceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching balance for vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s", p.ECCDSA, p.EDDSA, p.Chain, p.Address)

	balanceAmount, err := balance.FetchBalanceOfAddress(p.Chain, p.Address)
	if err != nil {
		return fmt.Errorf("services.GetBalanceOfAddress failed: %v: %w", err, asynq.SkipRetry)
	}

	token, err := balance.GetBaseTokenByChain(p.Chain)
	if err != nil {
		return fmt.Errorf("balances.GetBaseTokenByChain failed: %v: %w", err, asynq.SkipRetry)
	}

	b := &models.Balance{
		ECDSA:   p.ECCDSA,
		EDDSA:   p.EDDSA,
		Chain:   p.Chain,
		Address: p.Address,
		Balance: balanceAmount,
		Token:   token,
		Date:    time.Now().Unix(),
	}

	err = services.SaveBalance(b)
	if err != nil {
		return fmt.Errorf("services.SaveBalance failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Balance for vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s is %f", p.ECCDSA, p.EDDSA, p.Chain, p.Address, balanceAmount)

	result := map[string]interface{}{
		"balance": balanceAmount,
	}
	resultBytes, err := json.Marshal(result)
	if _, err := t.ResultWriter().Write([]byte(resultBytes)); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

func ProcessPointsCalculationTask(ctx context.Context, t *asynq.Task) error {
	var p PointsCalculationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Calculating points for Vault: eccdsa=%s, eddsa=%s", p.ECCDSA, p.EDDSA)
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
		return fmt.Errorf("services.SavePrice failed: %v: %w", err, asynq.SkipRetry)
	}

	result := map[string]interface{}{
		"price": usdPrice,
	}
	resultBytes, err := json.Marshal(result)
	if _, err := t.ResultWriter().Write([]byte(resultBytes)); err != nil {
		return fmt.Errorf("t.ResultWriter.Write failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

func ProcessPriceFetchAllActivePairsTask(ctx context.Context, t *asynq.Task) error {
	clientAsynq.Initialize()
	asynqClient := clientAsynq.AsynqClient
	fmt.Println(&clientAsynq.AsynqClient)

	pairs, err := services.GetUniqueActiveChainTokenPairs()
	if err != nil {
		return fmt.Errorf("!!4 services.GetUniqueActiveChainTokenPairs failed: %v", err)
	}

	for _, pair := range pairs {
		fmt.Printf("amogus1 Enqueuing price fetch task for pair: chain=%s, token=%s\n", pair.Chain, pair.Token)
		err := EnqueuePriceFetchTask(asynqClient, pair.Chain, pair.Token)
		if err != nil {
			return fmt.Errorf("failed to enqueue price fetch tasks: %v", err)
		}
		fmt.Printf("amogus1 Enqueued price fetch task for pair: chain=%s, token=%s\n", pair.Chain, pair.Token)
	}

	log.Printf("amogus1 Enqueued price fetch tasks for all active pairs")

	return nil
}
