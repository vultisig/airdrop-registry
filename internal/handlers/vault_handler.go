package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vultisig/airdrop-registry/internal/models"
)

const MaxPageSize = 100
const MinBalanceForValidReferral = 50

func (a *Api) registerVaultHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	if _, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA); err == nil {
		a.logger.Error(err)
		_ = c.Error(errVaultAlreadyRegist)
		return
	}
	vaultModel := models.Vault{
		Name:         vault.Name,
		Alias:        vault.Name,
		ECDSA:        vault.PublicKeyECDSA,
		EDDSA:        vault.PublicKeyEDDSA,
		Uid:          vault.Uid,
		HexChainCode: vault.HexChainCode,
		TotalPoints:  0,
		JoinAirdrop:  false,
	}

	if err := a.s.RegisterVault(&vaultModel); err != nil {
		if errors.Is(err, models.ErrAlreadyExist) {
			_ = c.Error(errVaultAlreadyRegist)
			return
		}
		a.logger.Error(err)
		_ = c.Error(errFailedToRegisterVault)
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
		_ = c.Error(errFailedToGetVault)
		return
	}
	coins, err := a.s.GetCoins(vault.ID)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetCoin)
		return
	}
	vaultResp := models.VaultResponse{
		UId:                   vault.Uid,
		Name:                  vault.Name,
		Alias:                 vault.Alias,
		PublicKeyECDSA:        vault.ECDSA,
		PublicKeyEDDSA:        vault.EDDSA,
		TotalPoints:           vault.TotalPoints,
		JoinAirdrop:           vault.JoinAirdrop,
		Rank:                  vault.Rank,
		Balance:               vault.Balance,
		LPValue:               vault.LPValue,
		NFTValue:              vault.NFTValue,
		RegisteredAt:          vault.Model.CreatedAt.UTC().Unix(),
		Coins:                 []models.ChainCoins{},
		AvatarURL:             vault.AvatarURL,
		ShowNameInLeaderboard: vault.ShowNameInLeaderboard,
		ReferralCode:          vault.ReferralCode,
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
		_ = c.Error(errInvalidRequest)
		return
	}
	vault, err := a.s.GetVaultByUID(uid)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	if vault == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	coins, err := a.s.GetCoins(vault.ID)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetCoin)
		return
	}
	if vault.Alias == "" {
		vault.Alias = vault.Name
	}
	vaultResp := models.VaultResponse{
		UId:            vault.Uid,
		Name:           vault.Alias,
		Alias:          vault.Alias,
		PublicKeyECDSA: "",
		PublicKeyEDDSA: "",
		TotalPoints:    vault.TotalPoints,
		JoinAirdrop:    vault.JoinAirdrop,
		Balance:        vault.Balance,
		LPValue:        vault.LPValue,
		NFTValue:       vault.NFTValue,
		Rank:           vault.Rank,
		RegisteredAt:   vault.Model.CreatedAt.UTC().Unix(),
		Coins:          []models.ChainCoins{},
		AvatarURL:      vault.AvatarURL,
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
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.JoinAirdrop = true
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			_ = c.Error(errFailedToJoinRegistry)
			return
		}
	} else {
		_ = c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}
func (a *Api) exitAirdrop(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.JoinAirdrop = false
		v.Rank = 0
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			_ = c.Error(errFailedToExitRegistry)
			return
		}
	}
	c.Status(http.StatusOK)
}
func (a *Api) deleteVaultHandler(c *gin.Context) {
	ecdsaPublicKey := c.Param("ecdsaPublicKey")
	eddsaPublicKey := c.Param("eddsaPublicKey")
	hexChainCode := c.GetHeader("x-hex-chain-code")
	if hexChainCode == "" {
		_ = c.Error(errForbiddenAccess)
		return
	}
	vault, err := a.s.GetVault(ecdsaPublicKey, eddsaPublicKey)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	if vault == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if hexChainCode == vault.HexChainCode {
		if err := a.s.DeleteVault(ecdsaPublicKey, eddsaPublicKey); err != nil {
			a.logger.Error(err)
			_ = c.Error(errFailedToDeleteVault)
			return
		}
	} else {
		_ = c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Api) updateAliasHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}

	if v == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.Alias = vault.Name
		v.ShowNameInLeaderboard = vault.ShowNameInLeaderboard
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			_ = c.Error(errFailedToUpdateVault)
			return
		}
	} else {
		_ = c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Api) updateReferralHandler(c *gin.Context) {
	var vault models.VaultRequest
	if err := c.ShouldBindJSON(&vault); err != nil {
		a.logger.Error(err)
		_ = c.Error(errInvalidRequest)
		return
	}
	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(vault.PublicKeyECDSA, vault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return
	}

	if v == nil {
		_ = c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == vault.HexChainCode && v.Uid == vault.Uid {
		v.ReferralCode = vault.ReferralCode
		if err := a.s.UpdateVault(v); err != nil {
			a.logger.Error(err)
			_ = c.Error(errFailedToUpdateVault)
			return
		}
	} else {
		_ = c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Api) getVaultsByRankHandler(c *gin.Context) {
	fromStr := c.DefaultQuery("from", "0")
	limitStr := c.DefaultQuery("limit", "10")
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		_ = c.Error(errInvalidRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		_ = c.Error(errInvalidRequest)
		return
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	vaultsResp := models.VaultsResponse{
		Vaults:          []models.VaultResponse{},
		TotalVaultCount: 0,
		TotalBalance:    0,
		TotalLP:         0,
		TotalNFT:        0,
	}
	vaultsResp.TotalVaultCount, err = a.s.GetLeaderVaultCount()
	if err != nil {
		a.logger.Errorf("failed to get leader vault count: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	vaultsResp.TotalBalance, err = a.s.GetLeaderVaultTotalBalance()
	if err != nil {
		a.logger.Errorf("failed to get leader vault total balance: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	vaultsResp.TotalLP, err = a.s.GetLeaderVaultTotalLP()
	if err != nil {
		a.logger.Errorf("failed to get leader vault total LP: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	vaultsResp.TotalNFT, err = a.s.GetLeaderVaultTotalNFT()
	if err != nil {
		a.logger.Errorf("failed to get leader vault total NFT: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	vaults, err := a.s.GetLeaderVaults(from, limit)
	if err != nil {
		a.logger.Errorf("failed to get leader vaults: %v", err)
		_ = c.Error(errFailedToGetVault)
		return
	}
	for _, vault := range vaults {
		vaultName := vault.Alias
		if !vault.ShowNameInLeaderboard {
			length := 10
			if len(vault.Uid) < 10 {
				length = len(vault.Uid)
			}
			vaultName = vault.Uid[:length]
		}
		vaultResp := models.VaultResponse{
			Name:         vaultName,
			Alias:        vaultName,
			TotalPoints:  vault.TotalPoints,
			Rank:         vault.Rank,
			Balance:      vault.Balance,
			LPValue:      vault.LPValue,
			NFTValue:     vault.NFTValue,
			RegisteredAt: vault.Model.CreatedAt.UTC().Unix(),
			AvatarURL:    vault.AvatarURL,
		}
		vaultsResp.Vaults = append(vaultsResp.Vaults, vaultResp)
	}
	c.JSON(http.StatusOK, vaultsResp)
}

func (a *Api) getUserReferrals(c *gin.Context) {
	var requestedVault models.VaultRequest
	if err := c.ShouldBindJSON(&requestedVault); err != nil {
		a.logger.Error(err)
		_ = c.Error(errInvalidRequest)
		return
	}
	refferer, err := a.s.GetVault(requestedVault.PublicKeyECDSA, requestedVault.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(errFailedToGetVault)
		return

	}

	if refferer == nil {
		c.Error(errVaultNotFound)
		return
	}

	var apiResponse models.ReferralsAPIResponse
	var summary models.ReferralsSummary
	if refferer.HexChainCode == requestedVault.HexChainCode && refferer.Uid == requestedVault.Uid {
		url := fmt.Sprintf("%s/user/referrals?eddsaKey=%s&ecdsaKey=%s&apiKey=%s",
			a.cfg.Vultiref.BaseAddress,
			refferer.EDDSA,
			refferer.ECDSA,
			a.cfg.Vultiref.APIKey)

		resp, err := http.Get(url)
		if err != nil {
			a.logger.WithError(err).Error("Failed to fetch referrals from API")
			c.Error(errFailedToFetchFromBotApi)
			return
		}
		if resp.StatusCode != http.StatusOK {
			a.logger.Errorf("API returned non-200 status code: %d", resp.StatusCode)
			_ = c.Error(errFailedToFetchFromBotApi)
			return
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			a.logger.WithError(err).Error("Failed to fetch referrals from API")
			c.Error(errUnknown)
			return
		}
		summary.TotalReferrals = apiResponse.Total
		for _, referral := range apiResponse.Items {
			if referral.WalletPublicKeyEcdsa == "" || referral.WalletPublicKeyEddsa == "" {
				continue
			}
			referralVault, err := a.s.GetVault(referral.WalletPublicKeyEcdsa, referral.WalletPublicKeyEddsa)
			if err != nil || referralVault == nil {
				a.logger.Warnf("Referral vault not found for ECDSA: %s, EDDSA: %s",
					referral.WalletPublicKeyEcdsa, referral.WalletPublicKeyEddsa)
				continue
			}

			//User can not refer himself
			if refferer.EDDSA == referralVault.EDDSA && refferer.ECDSA == referralVault.ECDSA {
				continue
			}

			if referralVault.Balance+referralVault.LPValue+referralVault.NFTValue >= MinBalanceForValidReferral {
				summary.ValidReferrals++
			}
		}
		c.JSON(http.StatusOK, summary)
	} else {
		c.Error(errForbiddenAccess)
	}
}
