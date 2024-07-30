package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	points, err := services.GetAllPoints(page, pageSize)
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
	points, err := services.GetPointsByVault(vault.ECDSA, vault.EDDSA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPoints := 0.0
	for _, point := range points {
		totalPoints += point.Share
	}
	c.JSON(http.StatusOK, gin.H{"total": totalPoints, "points": points})
}

func GetPointByIDHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point ID"})
		return
	}
	point, err := services.GetPointByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
		return
	}
	c.JSON(http.StatusOK, point)
}

func GetPointsForVaultByCycleHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	cycleID, err := strconv.ParseUint(c.Param("cycleID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cycle ID"})
		return
	}
	points, err := services.GetPointsForVaultByCycle(ecdsaPublicKey, eddsaPublicKey, uint(cycleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, points)
}
