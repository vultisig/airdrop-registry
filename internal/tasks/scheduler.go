package tasks

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func Schedule(redisConnOpt *asynq.RedisClientOpt) error {
	scheduler := asynq.NewScheduler(redisConnOpt, nil)

	// Price fetch parent task
	priceFetchTask := asynq.NewTask(TypePriceFetchParent, nil)
	entryID, err := scheduler.Register("30 * * * *", priceFetchTask,
		asynq.Queue(TypePriceFetchParent),
		asynq.Retention(24*time.Hour),
	)
	if err != nil {
		return err
	}
	log.Printf("Registered a scheduler entry for PriceFetchParent: %v", entryID)

	// Balance fetch parent task
	balanceFetchTask := asynq.NewTask(TypeBalanceFetchParent, nil)
	entryID, err = scheduler.Register("0 * * * *", balanceFetchTask,
		asynq.Queue(TypeBalanceFetchParent),
		asynq.Retention(24*time.Hour),
	)
	if err != nil {
		return err
	}
	log.Printf("Registered a scheduler entry for BalanceFetchParent: %v", entryID)

	// Points calculation parent task
	pointsCalcTask := asynq.NewTask(TypePointsCalculationParent, nil)
	entryID, err = scheduler.Register("15 * * * *", pointsCalcTask,
		asynq.Queue(TypePointsCalculationParent),
		asynq.Retention(24*time.Hour),
	)
	if err != nil {
		return err
	}
	log.Printf("Registered a scheduler entry for PointsCalculationParent: %v", entryID)

	if err := scheduler.Run(); err != nil {
		return err
	}

	log.Printf("Scheduler is running")
	return nil
}
