package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

func (a *Api) getAllSeasonInfo(c *gin.Context) {
	allSeasons := a.cfg.Seasons
	c.JSON(http.StatusOK, allSeasons)
}

type SeassonPoints struct {
	Points float64 `json:"points"`
}

func (a *Api) getTotalPointsBySeasonHandler(c *gin.Context) {
	startId := uint(0)
	totalPoints := 0.0
	for {
		allVaults, err := a.s.GetVaultsWithPage(startId, 1_000)
		if err != nil {
			a.logger.Error("failed to get vaults: ", err)
			_ = c.Error(errAddressNotMatch)
			return
		}
		if len(allVaults) == 0 {
			break
		}
		for _, vault := range allVaults {
			referralMultiplier := utils.GetReferralMultiplier(vault.ReferralCount)
			swapMultiplier := utils.GetSwapVolumeMultiplier(vault.SwapVolume)
			totalPoints += vault.TotalPoints * referralMultiplier * swapMultiplier
			startId = vault.ID
		}

	}
	points := SeassonPoints{
		Points: totalPoints,
	}

	c.JSON(http.StatusOK, points)
}
