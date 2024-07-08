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
	"github.com/vultisig/airdrop-registry/pkg/balance"
)

func ProcessBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p BalanceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching balance for vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s", p.ECCDSA, p.EDDSA, p.Chain, p.Address)

	balance, err := balance.FetchBalanceOfAddress(p.Chain, p.Address)
	if err != nil {
		return fmt.Errorf("services.GetBalanceOfAddress failed: %v: %w", err, asynq.SkipRetry)
	}

	b := &models.Balance{
		ECDSA:   p.ECCDSA,
		EDDSA:   p.EDDSA,
		Chain:   p.Chain,
		Address: p.Address,
		Balance: balance,
		Date:    time.Now().Unix(),
	}

	err = services.SaveBalance(b)
	if err != nil {
		return fmt.Errorf("services.SaveBalance failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Balance for vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s is %f", p.ECCDSA, p.EDDSA, p.Chain, p.Address, balance)

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
