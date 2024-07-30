package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	"github.com/vultisig/airdrop-registry/pkg/address"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"
)

func FetchVaultBalancesHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	addresses, err := address.GenerateSupportedChainAddresses(vault.ECDSA, vault.HexChainCode, vault.EDDSA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for chain, addr := range addresses {
		err := tasks.EnqueueBalanceFetchTask(asynqClient.AsynqClient, ecdsaPublicKey, eddsaPublicKey, chain, addr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Enqueued task: BalanceFetch for Vault: ecdsa=%s, eddsa=%s, chain=%s, address=%s", ecdsaPublicKey, eddsaPublicKey, chain, addr)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Balance fetch task enqueued", "vault": vault, "addresses": addresses})
}

func GetVaultBalancesHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	balances, err := services.GetBalancesByVaultKeys(vault.ECDSA, vault.EDDSA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func GetVaultBalanceHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	balance, err := services.GetLatestBalancesByVaultKeys(vault.ECDSA, vault.EDDSA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func GetBalanceByIDHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid balance ID"})
		return
	}
	balance, err := services.GetBalanceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "balance not found"})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func GetAllBalancesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	balances, err := services.GetAllBalances(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func GetAverageBalanceSinceHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	sinceStr := c.Query("since")
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'since' parameter"})
		return
	}
	avgBalance, err := services.GetAverageBalanceSince(ecdsaPublicKey, eddsaPublicKey, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"average_balance": avgBalance})
}

func GetAverageBalanceForTimeRangeHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'start' parameter"})
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'end' parameter"})
		return
	}
	avgBalance, err := services.GetAverageBalanceForTimeRange(ecdsaPublicKey, eddsaPublicKey, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"average_balance": avgBalance})
}
