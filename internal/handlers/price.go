package handlers

import (
	"log"
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"

	"github.com/gin-gonic/gin"
)

func StartPricesFetchHandler(c *gin.Context) {
	pairs, err := services.GetUniqueActiveChainTokenPairs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, pair := range pairs {
		err := tasks.EnqueuePriceFetchTask(asynqClient.AsynqClient, pair.Chain, pair.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Enqueued task: PriceFetch for Chain: %s, Token: %s", pair.Chain, pair.Token)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prices fetch task enqueued", "pairs": pairs})
}

func GetPricesHandler(c *gin.Context) {
	prices, err := services.GetLatestPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prices)
}
