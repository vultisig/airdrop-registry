package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// Balance fetch

func newBalanceFetch(
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
	task, err := newBalanceFetch(eccdsa, eddsa, chain, address)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(10*time.Second), asynq.Retention(24*time.Hour), asynq.Queue(TypeBalanceFetch))
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

func newPointsCalculationPayload(
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
	task, err := newPointsCalculationPayload(eccdsa, eddsa)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(1), asynq.Unique(time.Hour), asynq.Retention(24*time.Hour), asynq.Queue(TypePointsCalculation))
	return err
}

// Price fetch

func newPriceFetch(
	chain string,
	token string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(PriceFetchPayload{Chain: chain, Token: token})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePriceFetch, payload), nil
}

func EnqueuePriceFetchTask(
	asynqClient *asynq.Client,
	chain string,
	token string,
) error {
	task, err := newPriceFetch(chain, token)
	if err != nil {
		return err
	}

	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(2), asynq.Unique(time.Minute*5), asynq.Retention(24*time.Hour), asynq.Queue(TypePriceFetch))

	return err
}

// Price fetch for all active pairs

func NewPriceFetchForAllActivePairs() (*asynq.Task, error) {
	return asynq.NewTask(TypePriceFetchAllActivePairs, nil), nil
}

// func EnqueuePriceFetchTasksForActivePairs(asynqClient *asynq.Client) error {
// 	task, err := newPriceFetchForAllActivePairs()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = asynqClient.Enqueue(task, asynq.Retention(24*time.Hour), asynq.Queue(TypePriceFetchAllActivePairs))

// 	return err
// }
