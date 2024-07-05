package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Balance fetch

func NewBalanceFetch(
	vaultId string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(BalanceFetchPayload{VaultID: vaultId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeBalanceFetch, payload), nil
}

func EnqueueBalanceFetchTask(
	asynqClient *asynq.Client,
	vaultId string,
) error {
	task, err := NewBalanceFetch(vaultId)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.Queue("balance"))
	return err
}

// Point calculation

func NewPointsCalculationPayload(
	vaultId string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(PointsCalculationPayload{VaultID: vaultId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePointsCalculation, payload), nil
}

func EnqueuePointsCalculationTask(
	asynqClient *asynq.Client,
	vaultId string,
) error {
	task, err := NewPointsCalculationPayload(vaultId)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.Queue("points"))
	return err
}
