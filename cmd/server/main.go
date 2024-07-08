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

	router.POST("/vault", handlers.RegisterVaultHandler)

	router.GET("/vaults", handlers.ListVaultsHandler)

	router.GET("/vault/:eccdsaPublicKey/:eddsaPublicKey", handlers.GetVaultHandler)

	router.POST("/vault/:eccdsaPublicKey/:eddsaPublicKey/address", handlers.FetchVaultBalancesHandler)
	router.GET("/vault/:eccdsaPublicKey/:eddsaPublicKey/address", handlers.GetVaultAddressesHandler)

	router.GET("/vault/:eccdsaPublicKey/:eddsaPublicKey/balances", handlers.GetVaultBalancesHandler)
	router.GET("/vault/:eccdsaPublicKey/:eddsaPublicKey/balance", handlers.GetVaultBalanceHandler)

	go asynq.Initialize()
	defer asynq.AsynqClient.Close()
	defer db.CloseDatabase()

	router.Run(fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port))
}
