package tasks

// func SchedulePeriodicTasks() {
// 	for {
// 		randomDuration := time.Duration(rand.Intn(24)) * time.Hour
// 		time.Sleep(randomDuration)

// 		// Enqueue balance fetch tasks for each vault
// 		// Assuming vault IDs are from 1 to N
// 		for vaultID := 1; vaultID <= 100; vaultID++ {
// 			payload, err := json.Marshal(BalanceFetchPayload{VaultID: uint(vaultID)})
// 			if err != nil {
// 				log.Printf("Could not marshal payload: %v", err)
// 				continue
// 			}
// 			task := asynq.NewTask(TypeBalanceFetch, payload)
// 			_, err = asynqClient.AsynqClient(task, asynq.Queue("balance"))
// 			if err != nil {
// 				log.Printf("Could not enqueue task: %v", err)
// 			}
// 		}

// 		// Enqueue points calculation tasks for each vault
// 		for vaultID := 1; vaultID <= 100; vaultID++ {
// 			payload, err := json.Marshal(PointsCalculationPayload{VaultID: uint(vaultID)})
// 			if err != nil {
// 				log.Printf("Could not marshal payload: %v", err)
// 				continue
// 			}
// 			task := asynq.NewTask(TypePointsCalculation, payload)
// 			_, err = asynq.Client.Enqueue(task, asynq.Queue("points"))
// 			if err != nil {
// 				log.Printf("Could not enqueue task: %v", err)
// 			}
// 		}
// 	}
// }
