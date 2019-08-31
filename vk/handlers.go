package vk

import (
	"encoding/json"
	"errors"
	"github.com/flygrounder/go-mtg-vk/caching"
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
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
		go handleSearch(c, &req)
		c.String(http.StatusOK, "ok")
	}
}

func handleSearch(c *gin.Context, req *MessageRequest) {
	cardName, err := getCardNameByCommand(req.Object.Body)
	if err != nil {
		Message(req.Object.UserId, "Некорректная команда")
	} else if cardName == "" {
		Message(req.Object.UserId, "Карта не найдена")
	} else {
		prices := GetPrices(cardName)
		elements := min(CARDSLIMIT, len(prices))
		prices = prices[:elements]
		priceInfo := cardsinfo.FormatCardPrices(cardName, prices)
		Message(req.Object.UserId, priceInfo)
	}
}

func GetPrices(cardName string) []cardsinfo.CardPrice {
	client := caching.GetClient()
	val, err := client.Get(cardName)
	var prices []cardsinfo.CardPrice
	if err != nil {
		prices, _ = cardsinfo.GetSCGPrices(cardName)
		serialized, _ := json.Marshal(prices)
		client.Set(cardName, string(serialized))
		return prices
	}
	json.Unmarshal([]byte(val), &prices)
	return prices
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
