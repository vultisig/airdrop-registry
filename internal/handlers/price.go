package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"
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

func GetPriceByIDHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price ID"})
		return
	}
	price, err := services.GetPriceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "price not found"})
		return
	}
	c.JSON(http.StatusOK, price)
}

func GetPricesByTokenHandler(c *gin.Context) {
	chain := c.Query("chain")
	token := c.Query("token")
	if chain == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chain and token parameters are required"})
		return
	}
	prices, err := services.GetPricesByToken(chain, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

func GetPricesByTokenAndTimeRangeHandler(c *gin.Context) {
	chain := c.Query("chain")
	token := c.Query("token")
	startStr := c.Query("start")
	endStr := c.Query("end")
	if chain == "" || token == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chain, token, start, and end parameters are required"})
		return
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}
	prices, err := services.GetPricesByTokenAndTimeRange(chain, token, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

func GetLatestPriceByTokenHandler(c *gin.Context) {
	chain := c.Query("chain")
	token := c.Query("token")
	if chain == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chain and token parameters are required"})
		return
	}
	price, err := services.GetLatestPriceByToken(chain, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, price)
}
