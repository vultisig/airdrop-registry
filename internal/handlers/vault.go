package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/pkg/address"
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

func GetVaultHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	c.JSON(http.StatusOK, vault)
}

func GetVaultAddressesHandler(c *gin.Context) {
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
	c.JSON(http.StatusOK, addresses)
}

func ListVaultsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	vaults, err := services.GetAllVaults(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vaults)
}

// func UpdateVaultHandler(c *gin.Context) {
// 	ecdsaPublicKey := c.Param("ecdsaPublicKey")
// 	eddsaPublicKey := c.Param("eddsaPublicKey")
// 	var updatedVault models.Vault
// 	if err := c.ShouldBindJSON(&updatedVault); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	vault, err := services.GetVault(ecdsaPublicKey, eddsaPublicKey)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
// 		return
// 	}
// 	vault.HexChainCode = updatedVault.HexChainCode
// 	if err := services.UpdateVault(vault); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, vault)
// }

// func DeleteVaultHandler(c *gin.Context) {
// 	ecdsaPublicKey := c.Param("ecdsaPublicKey")
// 	eddsaPublicKey := c.Param("eddsaPublicKey")
// 	if err := services.DeleteVault(ecdsaPublicKey, eddsaPublicKey); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "vault deleted successfully"})
// }
