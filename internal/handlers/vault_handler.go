package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vultisig/airdrop-registry/internal/models"
)

func (a *Api) registerVaultHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	if _, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA); err == nil {
		a.logger.Error(err)
		c.Error(errVaultAlreadyRegist)
		return
	}
	vaultModel := models.Vault{
		Name:         vault.Name,
		ECDSA:        vault.PublicKeyECDSA,
		EDDSA:        vault.PublicKeyEDDSA,
		Uid:          vault.Uid,
		HexChainCode: vault.HexChainCode,
		TotalPoints:  0,
		JoinAirdrop:  false,
	}

	if err := a.s.RegisterVault(&vaultModel); err != nil {
		if errors.Is(err, models.ErrAlreadyExist) {
			c.Error(errVaultAlreadyRegist)
			return
		}
		a.logger.Error(err)
		c.Error(errFailedToRegisterVault)
		return
	}
	c.Status(http.StatusCreated)
}

func (a *Api) getVaultHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	vault, err := a.s.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	coins, err := a.s.GetCoins(vault.ID)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetCoin)
		return
	}
	vaultResp := models.VaultResponse{
		UId:            vault.Uid,
		Name:           vault.Name,
		Alias:          vault.Alias,
		PublicKeyECDSA: vault.ECDSA,
		PublicKeyEDDSA: vault.EDDSA,
		TotalPoints:    vault.TotalPoints,
		JoinAirdrop:    vault.JoinAirdrop,
		Coins:          []models.ChainCoins{},
	}
	for _, coin := range coins {
		found := false
		for i, _ := range vaultResp.Coins {
			if vaultResp.Coins[i].Name == coin.Chain {
				vaultResp.Coins[i].Coins = append(vaultResp.Coins[i].Coins, models.NewCoin(coin))
				found = true
			}
		}
		if !found {
			vaultResp.Coins = append(vaultResp.Coins, models.ChainCoins{
				Name:         coin.Chain,
				Address:      coin.Address,
				HexPublicKey: coin.HexPublicKey,
				Coins:        []models.Coin{models.NewCoin(coin)},
			})
		}
	}
	c.JSON(http.StatusOK, vaultResp)
}

func (a *Api) getVaultByUIDHandler(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.Error(errInvalidRequest)
		return
	}
	vault, err := a.s.GetVaultByUID(uid)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	if vault == nil {
		c.Error(errVaultNotFound)
		return
	}
	coins, err := a.s.GetCoins(vault.ID)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetCoin)
		return
	}
	if vault.Alias == "" {
		vault.Alias = vault.Name
	}
	vaultResp := models.VaultResponse{
		UId:            vault.Uid,
		Name:           vault.Alias,
		PublicKeyECDSA: "",
		PublicKeyEDDSA: "",
		TotalPoints:    vault.TotalPoints,
		JoinAirdrop:    vault.JoinAirdrop,
		Coins:          []models.ChainCoins{},
	}
	for i, _ := range coins {
		coin := coins[i]
		coin.VaultID = 0
		coin.HexPublicKey = ""
		found := false
		for j, _ := range vaultResp.Coins {
			if vaultResp.Coins[j].Name == coin.Chain {
				vaultResp.Coins[j].Coins = append(vaultResp.Coins[j].Coins, models.NewCoin(coin))
				found = true
			}
		}
		if !found {
			vaultResp.Coins = append(vaultResp.Coins, models.ChainCoins{
				Name:         coin.Chain,
				Address:      coin.Address,
				HexPublicKey: coin.HexPublicKey,
				Coins:        []models.Coin{models.NewCoin(coin)},
			})
		}
	}
	c.JSON(http.StatusOK, vaultResp)
}
func (a *Api) joinAirdrop(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.JoinAirdrop = true
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			c.Error(errFailedToJoinRegistry)
			return
		}
	} else {
		c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}
func (a *Api) exitAirdrop(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.JoinAirdrop = false
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			c.Error(errFailedToExitRegistry)
			return
		}
	}
	c.Status(http.StatusOK)
}
func (a *Api) deleteVaultHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		if err := a.s.DeleteVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA); err != nil {
			a.logger.Error(err)
			c.Error(errFailedToDeleteVault)
			return
		}
	}
	c.Status(http.StatusOK)
}

func (a *Api) updateAliasHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}

	if v == nil {
		c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.Alias = vault.Name
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			c.Error(errFailedToJoinRegistry)
			return
		}
	} else {
		fmt.Println(v.HexChainCode, vault.Uid)
		c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}
