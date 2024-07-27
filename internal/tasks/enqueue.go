package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Helper function to create a new task
func newTask(taskType string, payload interface{}) (*asynq.Task, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	return asynq.NewTask(taskType, jsonPayload), nil
}

// Balance Fetch Tasks

func EnqueueBalanceFetchTask(client *asynq.Client, ecdsa, eddsa, chain, address string) error {
	task, err := newTask(TypeBalanceFetch, BalanceFetchPayload{ECDSA: ecdsa, EDDSA: eddsa, Chain: chain, Address: address})
	if err != nil {
		return err
	}

	_, err = client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Unique(time.Minute),
		asynq.Timeout(10*time.Second),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypeBalanceFetch),
	)
	return err
}

func EnqueueBalanceFetchParentTask(client *asynq.Client) error {
	task := asynq.NewTask(TypeBalanceFetchParent, nil)
	_, err := client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Unique(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypeBalanceFetchParent),
	)
	return err
}

// Points Calculation Tasks

func EnqueuePointsCalculationTask(client *asynq.Client, ecdsa, eddsa string) error {
	task, err := newTask(TypePointsCalculation, PointsCalculationPayload{ECDSA: ecdsa, EDDSA: eddsa})
	if err != nil {
		return err
	}

	_, err = client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Unique(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypePointsCalculation),
	)
	return err
}

func EnqueuePointsCalculationParentTask(client *asynq.Client) error {
	task := asynq.NewTask(TypePointsCalculationParent, nil)
	_, err := client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Unique(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypePointsCalculationParent),
	)
	return err
}

// Price Fetch Tasks

func EnqueuePriceFetchTask(client *asynq.Client, chain, token string) error {
	task, err := newTask(TypePriceFetch, PriceFetchPayload{Chain: chain, Token: token})
	if err != nil {
		return err
	}

	_, err = client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Unique(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypePriceFetch),
	)
	return err
}

func EnqueuePriceFetchParentTask(client *asynq.Client) error {
	task := asynq.NewTask(TypePriceFetchParent, nil)
	_, err := client.Enqueue(task,
		asynq.MaxRetry(2),
		asynq.Unique(time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.Queue(TypePriceFetchParent),
	)
	return err
}
