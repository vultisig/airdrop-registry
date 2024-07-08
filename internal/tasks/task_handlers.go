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
)

// func ProcessVaultBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
// 	var p VaultBalanceFetchPayload
// 	if err := json.Unmarshal(t.Payload(), &p); err != nil {
// 		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
// 	}
// 	log.Printf("Fetching addresses for Vault: eccdsa=%s, eddsa=%s and spinning up new jobs to fetch balance", p.ECCDSA, p.EDDSA)

// 	vault, err := services.GetVault(p.ECCDSA, p.EDDSA)
// 	if err != nil {
// 		return fmt.Errorf("services.GetVault failed: %v: %w", err, asynq.SkipRetry)
// 	}

// 	addresses, err := address.GenerateSupportedChainAddresses(vault.ECDSA, vault.HexChainCode)
// 	if err != nil {
// 		return fmt.Errorf("address.GenerateSupportedChainAddresses failed: %v: %w", err, asynq.SkipRetry)
// 	}

// 	for chain, addr := range addresses {
// 		// err := EnqueueBalanceFetchTask(asynqClient.AsynqClient, p.ECCDSA, p.EDDSA, chain, addr)
// 		// if err != nil {
// 		// 	return fmt.Errorf("EnqueueBalanceFetchTask failed: %v: %w", err, asynq.SkipRetry)
// 		// }
// 		log.Printf("Enqueued task: BalanceFetch for Vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s", p.ECCDSA, p.EDDSA, chain, addr)
// 	}

// 	return nil
// }

func ProcessBalanceFetchTask(ctx context.Context, t *asynq.Task) error {
	var p BalanceFetchPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Fetching balance for vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s", p.ECCDSA, p.EDDSA, p.Chain, p.Address)

	balance, err := services.FetchBalanceOfAddress(p.Chain, p.Address)
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
	// Points calculation logic ...
	return nil
}
