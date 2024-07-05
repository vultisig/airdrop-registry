package asynq

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/config"

	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func Initialize() {
	redisConfig := config.Cfg.Redis
	redisAddr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	clientOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		DB:       redisConfig.DB,
		Password: redisConfig.Password,
	}

	AsynqClient = asynq.NewClient(clientOpt)
}
