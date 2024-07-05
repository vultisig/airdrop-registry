package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func HandleBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p BalanceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching balance for Vault: vault_id=%s", p.VaultID)
	// Fetch balance logic ...
	return nil
}

func HandlePointsCalculationTask(ctx context.Context, t *asynq.Task) error {
	var p PointsCalculationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Calculating points for Vault: vault_id=%s", p.VaultID)
	// Points calculation logic ...
	return nil
}
