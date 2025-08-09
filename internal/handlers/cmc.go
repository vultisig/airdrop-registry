package handlers

import "github.com/gin-gonic/gin"

type cmcQuestResponse struct {
	Result struct {
		IsValid bool `json:"isValid"`
	} `json:"result"`
}

func (a *Api) verifyCoinMarketCapQuest(c *gin.Context) {
	//TODO: add whitelist ip address check for cmc
	address := c.Query("address")
	isValid := a.questService.Exists(address)

	result := cmcQuestResponse{
		Result: struct {
			IsValid bool `json:"isValid"`
		}{
			IsValid: isValid,
		},
	}
	c.JSON(200, result)
}
