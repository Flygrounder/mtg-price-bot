package cardsinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	data := []CardPrice{
		&TcgCardPrice{
			Name:    "Green lotus",
			PriceFoil:   "22.8",
			Link:    "scg.com/1",
			Edition: "alpha",
		},
		&TcgCardPrice{
			Name:    "White lotus",
			Price:   "3.22",
			Link:    "scg.com/2",
			Edition: "gamma",
		},
	}
	res := FormatCardPrices("Black Lotus", data)
	ans := "Оригинальное название: Black Lotus\nРезультатов: 2\n1. alpha\nRegular: -\nFoil: $22.8\nscg.com/1\n2. gamma\nRegular: $3.22\nFoil: -\nscg.com/2\n"
	assert.Equal(t, ans, res)
}
