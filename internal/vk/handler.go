package vk

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Sender             sender
	Logger             *log.Logger
	SecretKey          string
	GroupId            int64
	ConfirmationString string
	DictPath           string
	Cache              cardCache
	InfoFetcher        cardInfoFetcher
}

type cardInfoFetcher interface {
	GetFormattedCardPrices(name string) (string, error)
	GetNameByCardId(set string, number string) string
	GetOriginalName(name string) string
}

type cardCache interface {
	Get(cardName string) (string, error)
	Set(cardName string, message string)
}

type messageRequest struct {
	Type    string      `json:"type"`
	GroupId int64       `json:"group_id"`
	Object  userMessage `json:"object"`
	Secret  string      `json:"secret"`
}

type userMessage struct {
	Body   string `json:"text"`
	UserId int64  `json:"peer_id"`
}

const (
	incorrectMessage         = "Некорректная команда"
	cardNotFoundMessage      = "Карта не найдена"
	pricesUnavailableMessage = "Цены временно недоступны, попробуйте позже"
)

func (h *Handler) HandleMessage(c *gin.Context) {
	var req messageRequest
	_ = c.BindJSON(&req)
	if req.Secret != h.SecretKey {
		return
	}
	switch req.Type {
	case "confirmation":
		h.handleConfirmation(c, &req)
	case "message_new":
		go h.handleSearch(&req)
		c.String(http.StatusOK, "ok")
	}
}

func (h *Handler) handleSearch(req *messageRequest) {
	cardName, err := h.getCardNameByCommand(req.Object.Body)
	if err != nil {
		h.Sender.send(req.Object.UserId, incorrectMessage)
		h.Logger.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), req.Object.Body)
	} else if cardName == "" {
		h.Sender.send(req.Object.UserId, cardNotFoundMessage)
		h.Logger.Printf("[info] Could not find card. User input: %s", req.Object.Body)
	} else {
		message, err := h.getMessage(cardName)
		if err != nil {
			h.Sender.send(req.Object.UserId, pricesUnavailableMessage)
			h.Logger.Printf("[error] Could not find SCG prices. Message: %s card name: %s", err.Error(), cardName)
			return
		}
		h.Sender.send(req.Object.UserId, message)
	}
}

func (h *Handler) handleConfirmation(c *gin.Context, req *messageRequest) {
	if (req.Type == "confirmation") && (req.GroupId == h.GroupId) {
		c.String(http.StatusOK, h.ConfirmationString)
	}
}

func (h *Handler) getMessage(cardName string) (string, error) {
	val, err := h.Cache.Get(cardName)
	if err != nil {
		message, err := h.InfoFetcher.GetFormattedCardPrices(cardName)
		if err != nil {
			return "", err
		}
		h.Cache.Set(cardName, message)
		return message, nil
	}
	return val, nil
}

func (h *Handler) getCardNameByCommand(command string) (string, error) {
	var name string
	switch {
	case strings.HasPrefix(command, "!s"):
		split := strings.Split(command, " ")
		if len(split) < 3 {
			return "", errors.New("wrong command")
		}
		set := split[1]
		number := split[2]
		name = h.InfoFetcher.GetNameByCardId(set, number)
	default:
		name = h.InfoFetcher.GetOriginalName(command)
	}
	return name, nil
}
