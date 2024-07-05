package main

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/handlers"
	"github.com/vultisig/airdrop-registry/pkg/asynq"
	"github.com/vultisig/airdrop-registry/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db.ConnectDatabase()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Vultisig Airdrop Registry",
		})
	})

	router.POST("/register", handlers.RegisterVaultHandler)
	router.GET("/vault/:eccdsaPublicKey/:eddsaPublicKey", handlers.QueryVaultHandler)
	router.GET("/vaults", handlers.ListVaultsHandler)

	go asynq.Initialize()
	defer asynq.AsynqClient.Close()
	defer db.CloseDatabase()

	router.Run(fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port))
}
