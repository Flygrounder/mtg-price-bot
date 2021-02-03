package cardsinfo

import (
	"fmt"
)

func FormatCardPrices(name string, prices []CardPrice) string {
	message := fmt.Sprintf("Оригинальное название: %v\n", name)
	message += fmt.Sprintf("Результатов: %v\n", len(prices))
	for i, v := range prices {
		message += fmt.Sprintf("%v. %v", i+1, v.Format())
	}
	return message
}
