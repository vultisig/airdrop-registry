package handlers

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func (a *Api) getVaultShareAppearanceHandler(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.Error(errInvalidRequest)
		return
	}
	vault, err := a.s.GetVaultByUID(uid)
	if err != nil {
		c.Error(errFailedToGetVault)
		return
	}
	if vault == nil {
		c.Error(errVaultNotFound)
		return
	}
	appearance := a.s.GetTheme(vault.ID)
	v := models.SharedVaultRequest{
		Theme: appearance.Theme,
		Logo:  appearance.Logo,
		Uid:   vault.Uid,
	}
	c.JSON(http.StatusOK, v)
}

const MaxLogoSize = 100 * 1024 // 100KB in bytes

func (a *Api) updateVaultShareAppearanceHandler(c *gin.Context) {
	var app models.SharedVaultRequest
	if err := c.ShouldBindJSON(&app); err != nil {
		c.Error(errInvalidRequest)
		return
	}
	base64Logo := app.Logo
	//Clean the base64 string (if there's a prefix like "data:image/png;base64,")
	if strings.Contains(base64Logo, "base64,") {
		base64Logo = strings.Split(base64Logo, ",")[1]
	}
	//Decode the base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64Logo)
	if err != nil {
		c.Error(errInvalidRequest)
		return
	}
	// Check logo size
	if len(imageData) > MaxLogoSize {
		c.Error(errLogoTooLarge)
		return
	}

	// check vault already exists , should we tell front-end that vault already registered?
	v, err := a.s.GetVault(app.PublicKeyECDSA, app.PublicKeyEDDSA)
	if err != nil {
		a.logger.Error(err)
		c.Error(errFailedToGetVault)
		return
	}
	if v == nil {
		c.Error(errVaultNotFound)
		return
	}
	if v.HexChainCode == app.HexChainCode && v.Uid == app.Uid {
		err = a.s.UpdateTheme(models.VaultShareAppearance{
			VaultID: v.ID,
			Theme:   app.Theme,
			Logo:    app.Logo,
		})
		if err != nil {
			a.logger.Error(err)
			c.Error(errFailedToSetTheme)
			return
		}
	} else {
		c.Error(errForbiddenAccess)
		return
	}
	c.Status(http.StatusOK)
}
