package handlers

import (
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterVaultHandler(c *gin.Context) {
	var vault models.Vault
	if err := c.ShouldBindJSON(&vault); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.RegisterVault(&vault); err != nil {
		if err.Error() == "vault with given ECDSA and EDDSA already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "already registered"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vault)
}

func QueryVaultHandler(c *gin.Context) {
	eccdsaPublicKey := c.Param("eccdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")

	vault, err := services.GetVault(eccdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	c.JSON(http.StatusOK, vault)
}

func ListVaultsHandler(c *gin.Context) {
	vaults, err := services.GetAllVaults()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vaults)
}
