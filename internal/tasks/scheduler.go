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
	entryID, err := scheduler.Register("30 * * * * ", task, asynq.Queue(TypePriceFetchAllActivePairs), asynq.Retention(24*time.Hour))
	if err != nil {
		return err
	}

	log.Printf("Registered a scheduler entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		return err
	}

	log.Printf("Scheduler is running\n")

	return nil
}
