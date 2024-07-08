package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Balance fetch

func NewBalanceFetch(
	eccdsa string,
	eddsa string,
	chain string,
	address string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(BalanceFetchPayload{ECCDSA: eccdsa, EDDSA: eddsa, Chain: chain, Address: address})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeBalanceFetch, payload), nil
}

func EnqueueBalanceFetchTask(
	asynqClient *asynq.Client,
	eccdsa string,
	eddsa string,
	chain string,
	address string,
) error {
	task, err := NewBalanceFetch(eccdsa, eddsa, chain, address)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.Queue(TypeBalanceFetch))
	return err
}

// Vault balance fetch

// func NewVaultBalanceFetch(
// 	eccdsa string,
// 	eddsa string,
// ) (*asynq.Task, error) {
// 	payload, err := json.Marshal(VaultBalanceFetchPayload{ECCDSA: eccdsa, EDDSA: eddsa})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return asynq.NewTask(TypeVaultBalanceFetch, payload), nil
// }

// func EnqueueVaultBalanceFetchTask(
// 	asynqClient *asynq.Client,
// 	eccdsa string,
// 	eddsa string,
// ) error {
// 	task, err := NewVaultBalanceFetch(eccdsa, eddsa)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = asynqClient.Enqueue(task, asynq.Queue(TypeVaultBalanceFetch))
// 	return err
// }

// Point calculation

func NewPointsCalculationPayload(
	eccdsa string,
	eddsa string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(PointsCalculationPayload{ECCDSA: eccdsa, EDDSA: eddsa})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePointsCalculation, payload), nil
}

func EnqueuePointsCalculationTask(
	asynqClient *asynq.Client,
	eccdsa string,
	eddsa string,
) error {
	task, err := NewPointsCalculationPayload(eccdsa, eddsa)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.Queue(TypePointsCalculation))
	return err
}
