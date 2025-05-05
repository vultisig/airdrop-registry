package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Api) getCurrentSeasonInfo(c *gin.Context) {
	currentSeason := a.cfg.Season
	c.JSON(http.StatusOK, currentSeason)
}
