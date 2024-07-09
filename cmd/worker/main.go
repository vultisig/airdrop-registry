package main

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func main() {
	config.LoadConfig()

	db.ConnectDatabase()

	redisConfig := config.Cfg.Redis

	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: addr, DB: redisConfig.DB, Password: redisConfig.Password},
		asynq.Config{
			Concurrency: 1, // process 1 tasks concurrently, prevents overloading external APIs
			Queues: map[string]int{
				// "critical":                  6,
				// "default":                   3,
				// "low":                       1,
				tasks.TypeBalanceFetch:      2,
				tasks.TypePointsCalculation: 1,
				tasks.TypeVaultBalanceFetch: 3,
				tasks.TypePriceFetch:        5,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeBalanceFetch, tasks.ProcessBalanceFetchTask)
	// mux.HandleFunc(tasks.TypeVaultBalanceFetch, tasks.ProcessVaultBalanceFetchTask)
	mux.HandleFunc(tasks.TypePointsCalculation, tasks.ProcessPointsCalculationTask)
	mux.HandleFunc(tasks.TypePriceFetch, tasks.ProcessPriceFetchTask)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	defer db.CloseDatabase()
}
