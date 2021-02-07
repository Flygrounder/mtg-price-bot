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
	message := fmt.Sprintf("Оригинальное название: %v\n", name)
	message += fmt.Sprintf("Результатов: %v\n", len(prices))
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v", i+1, v.format())
	}
	return message
}
