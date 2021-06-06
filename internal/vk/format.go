package vk

import (
	"fmt"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

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
	return fmt.Sprintf("%v: %v\n%v\n", s.Edition, s.Price, s.Link)
}
