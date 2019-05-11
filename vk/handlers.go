package vk

import (
	"errors"
	"github.com/flygrounder/mtg-price-vk/cardsinfo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const CARDSLIMIT = 8

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func HandleMessage(c *gin.Context) {
	var req MessageRequest
	c.BindJSON(&req)
	if req.Secret != SECRET_KEY {
		return
	}
	switch req.Type {
	case "confirmation":
		handleConfirmation(c, &req)
	case "message_new":
		handleSearch(c, &req)
	}
}

func handleSearch(c *gin.Context, req *MessageRequest) {
	defer c.String(http.StatusOK, "ok")
	cardName, err := getCardNameByCommand(req.Object.Body)
	if err != nil {
		Message(req.Object.UserId, "Некорректная команда")
	} else if cardName == "" {
		Message(req.Object.UserId, "Карта не найдена")
	} else {
		prices, _ := cardsinfo.GetSCGPrices(cardName)
		elements := min(CARDSLIMIT, len(prices))
		prices = prices[:elements]
		priceInfo := cardsinfo.FormatCardPrices(cardName, prices)
		Message(req.Object.UserId, priceInfo)
	}
}

func getCardNameByCommand(command string) (string, error) {
	var name string
	switch {
	case strings.HasPrefix(command, "!s"):
		split := strings.Split(command, " ")
		if len(split) < 3 {
			return "", errors.New("wrong command")
		}
		set := split[1]
		number := split[2]
		name = cardsinfo.GetNameByCardId(set, number)
	default:
		name = cardsinfo.GetOriginalName(command)
	}
	return name, nil
}

func handleConfirmation(c *gin.Context, req *MessageRequest) {
	if (req.Type == "confirmation") && (req.GroupId == GROUPID) {
		c.String(http.StatusOK, CONFIRMATION_STRING)
	}
}
