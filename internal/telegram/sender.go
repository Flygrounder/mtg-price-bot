package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Sender struct {
    API *tgbotapi.BotAPI 
}

func (h *Sender) Send(userId int64, message string) {
    msg := tgbotapi.NewMessage(userId, message)
    msg.DisableWebPagePreview = true
    h.API.Send(msg)
}

