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

	// Vaults
	router.POST("/vault", handlers.RegisterVaultHandler)
	router.GET("/vaults", handlers.ListVaultsHandler)
	router.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey", handlers.GetVaultHandler)

	// Address
	router.POST("/vault/:ecdsaPublicKey/:eddsaPublicKey/address", handlers.FetchVaultBalancesHandler)
	router.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey/address", handlers.GetVaultAddressesHandler)

	// Balances
	router.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey/balances", handlers.GetVaultBalancesHandler)
	router.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey/balance", handlers.GetVaultBalanceHandler)

	// Prices
	router.POST("/prices", handlers.StartPricesFetchHandler)
	router.GET("/prices", handlers.GetPricesHandler)

	// Points
	router.GET("/points", handlers.GetPointsHandler)
	router.POST("/points", handlers.StartPointsFetchHandler)
	router.GET("/vault/:ecdsaPublicKey/:eddsaPublicKey/points", handlers.GetVaultPointsHandler)

	go asynq.Initialize()
	defer asynq.AsynqClient.Close()
	defer db.CloseDatabase()

	router.Run(fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port))
}
