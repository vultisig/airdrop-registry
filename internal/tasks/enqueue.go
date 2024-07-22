package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// Balance fetch

func newBalanceFetch(
	ecdsa string,
	eddsa string,
	chain string,
	address string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(BalanceFetchPayload{ECDSA: ecdsa, EDDSA: eddsa, Chain: chain, Address: address})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeBalanceFetch, payload), nil
}

func EnqueueBalanceFetchTask(
	asynqClient *asynq.Client,
	ecdsa string,
	eddsa string,
	chain string,
	address string,
) error {
	task, err := newBalanceFetch(ecdsa, eddsa, chain, address)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Unique(time.Minute*1), asynq.Timeout(10*time.Second), asynq.Retention(24*time.Hour), asynq.Queue(TypeBalanceFetch))
	return err
}

// Point calculation

func newPointsCalculationPayload(
	ecdsa string,
	eddsa string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(PointsCalculationPayload{ECDSA: ecdsa, EDDSA: eddsa})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePointsCalculation, payload), nil
}

func EnqueuePointsCalculationTask(
	asynqClient *asynq.Client,
	ecdsa string,
	eddsa string,
) error {
	task, err := newPointsCalculationPayload(ecdsa, eddsa)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(2), asynq.Unique(time.Minute*1), asynq.Retention(24*time.Hour), asynq.Queue(TypePointsCalculation))
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

	_, err = asynqClient.Enqueue(task, asynq.MaxRetry(2), asynq.Unique(time.Minute*1), asynq.Retention(24*time.Hour), asynq.Queue(TypePriceFetch))

	return err
}

// Price fetch for all active pairs
func NewPriceFetchForAllActivePairs() (*asynq.Task, error) {
	return asynq.NewTask(TypePriceFetchAllActivePairs, nil), nil
}

// Balance fetch all
func NewBalanceFetchAll() (*asynq.Task, error) {
	return asynq.NewTask(TypeBalanceFetchAll, nil), nil
}
