package tasks

import (
	"log"

	"time"

	"github.com/hibiken/asynq"
)

func Schedule(redisConnOpt *asynq.RedisClientOpt) error {
	scheduler := asynq.NewScheduler(redisConnOpt, nil)

	// Price

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

	// Balance

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

	// Points calculation
	task, err = NewPointsCalculationAll()
	if err != nil {
		return err
	}

	// Runs at :15 every hour
	entryId, err = scheduler.Register("15 * * * * ", task, asynq.Queue(TypePointsCalculation), asynq.Retention(24*time.Hour))
	if err != nil {
		return err
	}

	log.Printf("Registered a scheduler entry for PointsCalculationAll: %v", entryId)

	if err := scheduler.Run(); err != nil {
		return err
	}

	log.Printf("Scheduler is running\n")

	return nil
}
