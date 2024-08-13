package main

func main() {
	//config.LoadConfig()
	//
	//db.ConnectDatabase()
	//
	//redisConfig := config.Cfg.Redis
	//addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	//redis := asynq.RedisClientOpt{Addr: addr, DB: redisConfig.DB, Password: redisConfig.Password}
	//
	//server := asynq.NewServer(
	//	redis,
	//	asynq.Config{
	//		Concurrency: 5,
	//		Queues: map[string]int{
	//			tasks.TypeBalanceFetch:            2,
	//			tasks.TypeBalanceFetchParent:      1,
	//			tasks.TypePointsCalculation:       1,
	//			tasks.TypePointsCalculationParent: 1,
	//			tasks.TypePriceFetch:              5,
	//			tasks.TypePriceFetchParent:        1,
	//		},
	//	},
	//)
	//
	//// Create a context that is canceled on interrupt signal
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//// Run tasks.Schedule in a separate goroutine
	//go func() {
	//	err := tasks.Schedule(&redis)
	//	if err != nil {
	//		log.Fatalf("could not schedule tasks: %v", err)
	//	}
	//}()
	//
	//mux := asynq.NewServeMux()
	//mux.HandleFunc(tasks.TypeBalanceFetch, tasks.ProcessBalanceFetchTask)
	//mux.HandleFunc(tasks.TypeBalanceFetchParent, tasks.ProcessBalanceFetchParentTask)
	//mux.HandleFunc(tasks.TypePointsCalculation, tasks.ProcessPointsCalculationTask)
	//mux.HandleFunc(tasks.TypePointsCalculationParent, tasks.ProcessPointsCalculationParentTask)
	//mux.HandleFunc(tasks.TypePriceFetch, tasks.ProcessPriceFetchTask)
	//mux.HandleFunc(tasks.TypePriceFetchParent, tasks.ProcessPriceFetchParentTask)
	//
	//// Run server in a separate goroutine
	//go func() {
	//	if err := server.Run(mux); err != nil {
	//		log.Fatalf("could not run server: %v", err)
	//	}
	//}()
	//
	//// Set up channel to receive OS signals
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//
	//// Block until a signal is received
	//<-sigChan
	//log.Println("Shutting down gracefully...")
	//
	//// Cancel the context to signal the goroutines to stop
	//cancel()
	//
	//// Wait for the server to shutdown
	//<-ctx.Done()
	//
	//// Close the database connection
	//db.CloseDatabase()
	//
	//log.Println("Server stopped")
}
