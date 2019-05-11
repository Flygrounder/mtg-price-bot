package cardsinfo

import (
	"fmt"
)

func FormatCardPrices(name string, prices []CardPrice) string {
	message := fmt.Sprintf("Оригинальное название: %v\n", name)
	message += fmt.Sprintf("Результатов: %v\n", len(prices))
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v: $%v\n", i+1, v.Edition, v.Price)
		message += fmt.Sprintf("%v\n", v.Link)
	}
	return message
}
