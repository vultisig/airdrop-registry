package tasks

import (
	"log"

	"time"

	"github.com/hibiken/asynq"
)

func Schedule(redisConnOpt *asynq.RedisClientOpt) error {
	scheduler := asynq.NewScheduler(redisConnOpt, nil)

	task, err := NewPriceFetchForAllActivePairs()
	if err != nil {
		return err
	}

	// Runs at :30 every hour
	entryId, err := scheduler.Register("30 * * * * ", task, asynq.Queue(TypePriceFetchAllActivePairs), asynq.Retention(24*time.Hour))
	if err != nil {
		return err
	}

	log.Printf("Registered a scheduler entry for PriceFetchAllActivePairs: %v", entryId)

	// Runs at :00 every hour
	task, err = NewBalanceFetchAll()
	if err != nil {
		return err
	}

	// Runs at :00 every hour
	entryId, err = scheduler.Register("0 * * * * ", task, asynq.Queue(TypeBalanceFetchAll), asynq.Retention(24*time.Hour))
	if err != nil {
		return err
	}

	log.Printf("Registered a scheduler entry for BalanceFetchAll: %v", entryId)

	if err := scheduler.Run(); err != nil {
		return err
	}

	log.Printf("Scheduler is running\n")

	return nil
}
