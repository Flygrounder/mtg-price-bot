package vk

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type Handler struct {
	Sender             Sender
	Logger             *log.Logger
	SecretKey          string
	GroupId            int64
	ConfirmationString string
	DictPath           string
	Cache              CardCache
	InfoFetcher        CardInfoFetcher
}

type CardInfoFetcher interface {
	GetPrices(name string) ([]cardsinfo.CardPrice, error)
	FormatCardPrices(name string, prices []cardsinfo.CardPrice) string
	GetNameByCardId(set string, number string) string
	GetOriginalName(name string, dict io.Reader) string
}

type CardCache interface {
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
		h.Sender.Send(req.Object.UserId, incorrectMessage)
		h.Logger.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), req.Object.Body)
	} else if cardName == "" {
		h.Sender.Send(req.Object.UserId, cardNotFoundMessage)
		h.Logger.Printf("[info] Could not find card. User input: %s", req.Object.Body)
	} else {
		message, err := h.getMessage(cardName)
		if err != nil {
			h.Sender.Send(req.Object.UserId, pricesUnavailableMessage)
			h.Logger.Printf("[error] Could not find SCG prices. Message: %s card name: %s", err.Error(), cardName)
			return
		}
		h.Sender.Send(req.Object.UserId, message)
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
		prices, err := h.InfoFetcher.GetPrices(cardName)
		if err != nil {
			return "", err
		}
		message := h.InfoFetcher.FormatCardPrices(cardName, prices)
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
		dict, _ := os.Open(h.DictPath)
		name = h.InfoFetcher.GetOriginalName(command, dict)
	}
	return name, nil
}
