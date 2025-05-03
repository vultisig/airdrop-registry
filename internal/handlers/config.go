package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Api) getSessionAchievementsHandler(c *gin.Context) {
	sessionAchievements := a.cfg.SessionAchievements
	c.JSON(http.StatusOK, sessionAchievements)
}
