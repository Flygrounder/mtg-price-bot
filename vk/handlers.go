package vk

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/flygrounder/go-mtg-vk/caching"
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
	"github.com/gin-gonic/gin"
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
		prices, err := GetPrices(cardName)
		if err != nil {
			Message(req.Object.UserId, "Цены временно недоступны, попробуйте позже")
			return
		}
		elements := min(CARDSLIMIT, len(prices))
		prices = prices[:elements]
		priceInfo := cardsinfo.FormatCardPrices(cardName, prices)
		Message(req.Object.UserId, priceInfo)
	}
}

func GetPrices(cardName string) ([]cardsinfo.CardPrice, error) {
	client := caching.GetClient()
	val, err := client.Get(cardName)
	var prices []cardsinfo.CardPrice
	if err != nil {
		prices, err = cardsinfo.GetSCGPrices(cardName)
		if err != nil {
			return nil, err
		}
		serialized, err := json.Marshal(prices)
		if err != nil {
			return nil, err
		}
		client.Set(cardName, string(serialized))
		return prices, nil
	}
	json.Unmarshal([]byte(val), &prices)
	return prices, nil
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
