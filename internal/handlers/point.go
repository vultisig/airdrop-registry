package handlers

import (
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"

	"github.com/gin-gonic/gin"
)

func StartPointsFetchHandler(c *gin.Context) {
	err := tasks.EnqueuePointsCalculationParentTask(asynqClient.AsynqClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points calculation task enqueued"})
}

func GetPointsHandler(c *gin.Context) {
	points, err := services.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, points)
}

func GetVaultPointsHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")

	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	balances, err := services.GetPointsByVault(vault.ECDSA, vault.EDDSA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPoints := 0.0
	for _, balance := range balances {
		totalPoints += balance.Share
	}

	c.JSON(http.StatusOK, gin.H{"total": totalPoints, "points": balances})
}
