package vk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"
)

type Handler struct {
	Scenario *scenario.Scenario
	SecretKey          string
	GroupId            int64
	ConfirmationString string
	DictPath           string
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
		go h.Scenario.HandleSearch(&scenario.UserMessage{
			Body:   req.Object.Body,
			UserId: req.Object.UserId,
		})
		c.String(http.StatusOK, "ok")
	}
}

func (h *Handler) handleConfirmation(c *gin.Context, req *messageRequest) {
	if (req.Type == "confirmation") && (req.GroupId == h.GroupId) {
		c.String(http.StatusOK, h.ConfirmationString)
	}
}
