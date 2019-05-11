package tests

import (
	"github.com/flygrounder/mtg-price-vk/cardsinfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	data := []cardsinfo.CardPrice{
		{
			Name:    "Green lotus",
			Price:   22.8,
			Link:    "scg.com/1",
			Edition: "alpha",
		},
		{
			Name:    "White lotus",
			Price:   3.22,
			Link:    "scg.com/2",
			Edition: "gamma",
		},
	}
	res := cardsinfo.FormatCardPrices("Black Lotus", data)
	ans := "Оригинальное название: Black Lotus\nРезультатов: 2\n1. alpha: $22.8\nscg.com/1\n2. gamma: $3.22\nscg.com/2\n"
	assert.Equal(t, res, ans)
}
