package vk

import (
	"github.com/flygrounder/mtg-price-vk/cardsinfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleMessage(c *gin.Context) {
	defer c.String(http.StatusOK, "ok")
	var req MessageRequest
	c.BindJSON(&req)
	if req.Secret != SECRET_KEY {
		return
	}
	cardName := cardsinfo.GetOriginalName(req.Object.Body)
	Message(req.Object.UserId, cardName)
}
