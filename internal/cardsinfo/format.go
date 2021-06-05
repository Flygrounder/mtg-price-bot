package cardsinfo

import (
	"fmt"
)

func (f *Fetcher) GetFormattedCardPrices(name string) (string, error) {
	prices, err := f.getPrices(name)
	if err != nil {
		return "", err
	}
	return f.formatCardPrices(name, prices), nil
}

func (f *Fetcher) formatCardPrices(name string, prices []scgCardPrice) string {
	message := fmt.Sprintf("Оригинальное название: %v\n\n", name)
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v", i+1, v.format())
	}
	if len(prices) == 0 {
	    message += "Цен не найдено\n"
	}
	return message
}
