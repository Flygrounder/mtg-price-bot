package vk

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

var dictPath = "./assets/additional_cards.json"

func HandleMessage(c *gin.Context) {
	var req MessageRequest
	_ = c.BindJSON(&req)
	if req.Secret != SecretKey {
		return
	}
	switch req.Type {
	case "confirmation":
		handleConfirmation(c, &req)
	case "message_new":
		go handleSearch(&req)
		c.String(http.StatusOK, "ok")
	}
}

func handleSearch(req *MessageRequest) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[error] Search panicked. Exception info: %s", r)
		}
	}()

	cardName, err := getCardNameByCommand(req.Object.Body)
	if err != nil {
		Message(req.Object.UserId, "Некорректная команда")
		log.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), req.Object.Body)
	} else if cardName == "" {
		Message(req.Object.UserId, "Карта не найдена")
		log.Printf("[info] Could not find card. User input: %s", req.Object.Body)
	} else {
		message, err := GetMessage(cardName)
		if err != nil {
			Message(req.Object.UserId, "Цены временно недоступны, попробуйте позже")
			log.Printf("[error] Could not find SCG prices. Message: %s card name: %s", err.Error(), cardName)
			return
		}
		Message(req.Object.UserId, message)
	}
}

func GetMessage(cardName string) (string, error) {
	client := caching.GetClient()
	val, err := client.Get(cardName)
	if err != nil {
		prices, err := cardsinfo.GetPrices(cardName)
		if err != nil {
			return "", err
		}
		message := cardsinfo.FormatCardPrices(cardName, prices)
		client.Set(cardName, message)
		return message, nil
	}
	return val, nil
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
		dict, _ := os.Open(dictPath)
		name = cardsinfo.GetOriginalName(command, dict)
	}
	return name, nil
}

func handleConfirmation(c *gin.Context, req *MessageRequest) {
	if (req.Type == "confirmation") && (req.GroupId == GroupId) {
		c.String(http.StatusOK, ConfirmationString)
	}
}
