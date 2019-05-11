package vk

import (
	"github.com/flygrounder/mtg-price-vk/cardsinfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleMessage(c *gin.Context) {
	var req MessageRequest
	c.BindJSON(&req)
	if (req.Type == "confirmation") && (req.GroupId == GROUPID) {
		c.String(http.StatusOK, CONFIRMATION_STRING)
		return
	}
	if req.Secret != SECRET_KEY {
		return
	}
	cardName := cardsinfo.GetOriginalName(req.Object.Body)
	Message(req.Object.UserId, cardName)
	c.String(http.StatusOK, "ok")
}
