package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	redis := asynq.RedisClientOpt{Addr: addr, DB: redisConfig.DB, Password: redisConfig.Password}

	server := asynq.NewServer(
		redis,
		asynq.Config{
			Concurrency: 1, // process 1 task concurrently, prevents overloading external APIs
			Queues: map[string]int{
				tasks.TypeBalanceFetch:             2,
				tasks.TypePointsCalculation:        1,
				tasks.TypeVaultBalanceFetch:        3,
				tasks.TypePriceFetch:               5,
				tasks.TypePriceFetchAllActivePairs: 1, // because this only schedules new jobs that have a lower priority and this is only 1x job every x period
				tasks.TypeBalanceFetchAll:          1, // same as above
			},
		},
	)

	// Create a context that is canceled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run tasks.Schedule in a separate goroutine
	go func() {
		err := tasks.Schedule(&redis)
		if err != nil {
			log.Fatalf("could not schedule tasks: %v", err)
		}
	}()

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeBalanceFetch, tasks.ProcessBalanceFetchTask)
	mux.HandleFunc(tasks.TypePointsCalculation, tasks.ProcessPointsCalculationTask)
	mux.HandleFunc(tasks.TypePriceFetch, tasks.ProcessPriceFetchTask)
	mux.HandleFunc(tasks.TypePriceFetchAllActivePairs, tasks.ProcessPriceFetchAllActivePairsTask)
	mux.HandleFunc(tasks.TypeBalanceFetchAll, tasks.ProcessBalanceFetchAllTask)

	// Run server in a separate goroutine
	go func() {
		if err := server.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// Set up channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan
	log.Println("Shutting down gracefully...")

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for the server to shutdown
	<-ctx.Done()

	// Close the database connection
	db.CloseDatabase()

	log.Println("Server stopped")
}
