package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Api) getAllSeasonInfo(c *gin.Context) {
	allSeasons := a.cfg.Seasons
	c.JSON(http.StatusOK, allSeasons)
}
