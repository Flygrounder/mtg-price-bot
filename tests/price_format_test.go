package tests

import (
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	data := []cardsinfo.CardPrice{
		{
			Name:    "Green lotus",
			PriceFoil:   "22.8",
			Link:    "scg.com/1",
			Edition: "alpha",
		},
		{
			Name:    "White lotus",
			Price:   "3.22",
			Link:    "scg.com/2",
			Edition: "gamma",
		},
	}
	res := cardsinfo.FormatCardPrices("Black Lotus", data)
	ans := "Оригинальное название: Black Lotus\nРезультатов: 2\n1. alpha\nRegular: -\nFoil: $22.8\nscg.com/1\n2. gamma\nRegular: $3.22\nFoil: -\nscg.com/2\n"
	assert.Equal(t, res, ans)
}
