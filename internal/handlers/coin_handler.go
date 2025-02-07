package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
)

func addCoin(store *services.Storage, coin models.CoinBase, ecdsaPublicKey, eddsaPublicKey, hexChainCode string) (uint, error) {
	coin.Balance = ""
	coin.USDValue = ""
	coin.PriceUSD = ""
	if hexChainCode == "" {
		return 0, errForbiddenAccess
	}
	// Ensure the relevant vault exist
	vault, err := store.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		return 0, errVaultNotFound
	}
	if vault.HexChainCode != hexChainCode {
		return 0, errForbiddenAccess
	}
	addr, err := vault.GetAddress(coin.Chain)
	if err != nil {
		return 0, errFailedToGetAddress
	}
	if coin.Address != addr {
		return 0, errAddressNotMatch
	}
	coinDB := models.CoinDBModel{
		CoinBase: coin,
		VaultID:  vault.ID,
	}
	if err := store.AddCoin(&coinDB); err != nil {
		return 0, errFailedToAddCoin
	}
	return coinDB.ID, nil
}

func (a *Api) addCoin(c *gin.Context) {
	var coin models.CoinBase
	if err := c.ShouldBindJSON(&coin); err != nil {
		a.logger.Errorf("failed to bind json: %v", err)
		_ = c.Error(errInvalidRequest)
		return
	}
	coin.Balance = ""
	coin.USDValue = ""
	coin.PriceUSD = ""
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	hexChainCode := c.GetHeader("x-hex-chain-code")
	id, err := addCoin(a.s, coin, ecdsaPublicKey, eddsaPublicKey, hexChainCode)
	if err != nil {
		a.logger.Errorf("failed to add coin: %v", err)
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"coin_id": id})
}

func (a *Api) addCoins(c *gin.Context) {
	var coins []models.CoinBase
	if err := c.ShouldBindJSON(&coins); err != nil {
		a.logger.Errorf("failed to bind json: %v", err)
		_ = c.Error(errInvalidRequest)
		return
	}
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	hexChainCode := c.GetHeader("x-hex-chain-code")
	ids := make([]uint, 0)
	for i := range coins {
		id, err := addCoin(a.s, coins[i], ecdsaPublicKey, eddsaPublicKey, hexChainCode)
		if err != nil {
			a.logger.Errorf("failed to add coin: %v", err)
			_ = c.Error(err)
			return
		}
		ids = append(ids, id)
	}
	c.JSON(http.StatusCreated, gin.H{"coin_ids": ids})
}

func (a *Api) deleteCoin(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	strCoinID := c.Param("coinID")
	hexChainCode := c.GetHeader("x-hex-chain-code")
	if hexChainCode == "" {
		_ = c.Error(errForbiddenAccess)
		return
	}

	// Ensure the relevant vault exist
	vault, err := a.s.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		a.logger.Errorf("failed to get vault: %v", err)
		_ = c.Error(errVaultNotFound)
		return
	}

	if vault.HexChainCode != hexChainCode {
		_ = c.Error(errForbiddenAccess)
		return
	}
	coin, err := a.s.GetCoin(strCoinID)
	if err != nil {
		a.logger.Errorf("failed to get coin: %v", err)
		_ = c.Error(errFailedToGetCoin)
		return
	}
	// If the coin is native token, delete all coins with the same chain
	if coin.IsNative {
		coins, err := a.s.GetCoins(vault.ID)
		if err != nil {
			a.logger.Errorf("failed to get coins: %v", err)
			_ = c.Error(errFailedToGetCoin)
			return
		}
		coinIds := make([]uint, 0)
		for i := range coins {
			if coins[i].Chain == coin.Chain {
				coinIds = append(coinIds, coins[i].ID)
			}
		}
		if err := a.s.DeleteCoins(coinIds, vault.ID); err != nil {
			a.logger.Errorf("failed to delete coins: %v", err)
			_ = c.Error(errFailedToDeleteCoin)
			return
		}
	} else if err := a.s.DeleteCoin(strCoinID, vault.ID); err != nil {
		a.logger.Errorf("failed to delete coin: %v", err)
		_ = c.Error(errFailedToDeleteCoin)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *Api) getCoin(c *gin.Context) {
	strCoinID := c.Param("coinID")
	coin, err := a.s.GetCoin(strCoinID)
	if err != nil {
		a.logger.Errorf("failed to get coin: %v", err)
		_ = c.Error(errFailedToGetCoin)
		return
	}
	c.JSON(http.StatusOK, coin.CoinBase)
}
