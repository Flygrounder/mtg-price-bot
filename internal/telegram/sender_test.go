package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

func Test_formatCardPrices(t *testing.T) {
	prices := []cardsinfo.ScgCardPrice{
		{
			Price:   "1",
			Edition: "Alpha",
			Link:    "scg1",
		},
		{
			Price:   "2",
			Edition: "Beta",
			Link:    "scg2",
		},
	}
	result := formatCardPrices("card", prices)
	assert.Equal(t, "Оригинальное название: card\n\n1. [Alpha](scg1): 1\n2. [Beta](scg2): 2\n", result)
}

func Test_formatCardPricesEscapeUnderscore(t *testing.T) {
	prices := []cardsinfo.ScgCardPrice{}
	result := formatCardPrices("_____", prices)
	assert.Equal(t, "Оригинальное название: \\_\\_\\_\\_\\_\n\nЦен не найдено\n", result)
}
