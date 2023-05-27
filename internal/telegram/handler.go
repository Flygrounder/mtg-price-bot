package telegram

import (
	"context"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"
)

const welcomeMessage = "Здравствуйте, вас приветствует бот для поиска цен на карты MTG, введите название карты, которая вас интересует."

type Handler struct {
	Scenario *scenario.Scenario
}

func (h *Handler) HandleMessage(c *gin.Context) {
	var upd tgbotapi.Update
	err := c.Bind(&upd)
	if err != nil || upd.Message == nil {
		return
	}

	if upd.Message.Text == "/start" {
		h.Scenario.Sender.Send(upd.Message.Chat.ID, welcomeMessage)
		return
	}

	h.Scenario.HandleSearch(context.Background(), &scenario.UserMessage{
		Body:   upd.Message.Text,
		UserId: upd.Message.Chat.ID,
	})
}
