package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type Sender struct {
	API *tgbotapi.BotAPI
}

func (s *Sender) SendPrices(userId int64, cardName string, prices []cardsinfo.ScgCardPrice) {
	msg := formatCardPrices(cardName, prices)
	s.Send(userId, msg)
}

func (h *Sender) Send(userId int64, message string) {
	msg := tgbotapi.NewMessage(userId, message)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "MarkdownV2"
	h.API.Send(msg)
}

func formatCardPrices(name string, prices []cardsinfo.ScgCardPrice) string {
	message := fmt.Sprintf("Оригинальное название: %v\n\n", name)
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v", i+1, formatPrice(v))
	}
	if len(prices) == 0 {
		message += "Цен не найдено\n"
	}
	return message
}

func formatPrice(s cardsinfo.ScgCardPrice) string {
	return fmt.Sprintf("[%v](%v): %v\n", s.Edition, s.Link, s.Price)
}
