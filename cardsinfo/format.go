package cardsinfo

import (
	"fmt"
)

func FormatCardPrices(name string, prices []CardPrice) string {
	message := fmt.Sprintf("Оригинальное название: %v\n", name)
	message += fmt.Sprintf("Результатов: %v\n", len(prices))
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v\nRegular: %v\nFoil: %v\n", i+1, v.Edition, formatPrice(v.Price), formatPrice(v.PriceFoil))
		message += fmt.Sprintf("%v\n", v.Link)
	}
	return message
}

func formatPrice(price string) string {
	if price == "" {
		return "-"
	}
	return fmt.Sprintf("$%v", price)
}
