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

type Handler struct {
	Sender             Sender
	Logger             *log.Logger
	SecretKey          string
	GroupId            int64
	ConfirmationString string
	DictPath           string
	CachingClient      *caching.CacheClient
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
	defer func() {
		if r := recover(); r != nil {
			h.Logger.Printf("[error] Search panicked. Exception info: %s", r)
		}
	}()

	cardName, err := h.getCardNameByCommand(req.Object.Body)
	if err != nil {
		h.Sender.Send(req.Object.UserId, "Некорректная команда")
		h.Logger.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), req.Object.Body)
	} else if cardName == "" {
		h.Sender.Send(req.Object.UserId, "Карта не найдена")
		h.Logger.Printf("[info] Could not find card. User input: %s", req.Object.Body)
	} else {
		message, err := h.getMessage(cardName)
		if err != nil {
			h.Sender.Send(req.Object.UserId, "Цены временно недоступны, попробуйте позже")
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
	val, err := h.CachingClient.Get(cardName)
	if err != nil {
		prices, err := cardsinfo.GetPrices(cardName)
		if err != nil {
			return "", err
		}
		message := cardsinfo.FormatCardPrices(cardName, prices)
		h.CachingClient.Set(cardName, message)
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
		name = cardsinfo.GetNameByCardId(set, number)
	default:
		dict, _ := os.Open(h.DictPath)
		name = cardsinfo.GetOriginalName(command, dict)
	}
	return name, nil
}
