// FetchVaultBalancesHandler

package handlers

import (
	"log"
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tasks"
	"github.com/vultisig/airdrop-registry/pkg/address"
	asynqClient "github.com/vultisig/airdrop-registry/pkg/asynq"

	"github.com/gin-gonic/gin"
)

func FetchVaultBalancesHandler(c *gin.Context) {
	eccdsaPublicKey := c.Param("eccdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")

	vault, err := services.GetVault(eccdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	addresses, err := address.GenerateSupportedChainAddresses(vault.ECDSA, vault.HexChainCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for chain, addr := range addresses {
		err := tasks.EnqueueBalanceFetchTask(asynqClient.AsynqClient, eccdsaPublicKey, eddsaPublicKey, chain, addr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Enqueued task: BalanceFetch for Vault: eccdsa=%s, eddsa=%s, chain=%s, address=%s", eccdsaPublicKey, eddsaPublicKey, chain, addr)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance fetch task enqueued", "vault": vault, "addresses": addresses})
}
