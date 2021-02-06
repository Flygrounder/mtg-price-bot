package cardsinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatCardPrices(t *testing.T) {
	formatted := FormatCardPrices("card", []CardPrice{
		&ScgCardPrice{
			Price:   "1.5$",
			Edition: "ED",
			Link:    "scg.com",
		},
	})
	assert.Equal(t, "Оригинальное название: card\nРезультатов: 1\n1. ED: 1.5$\nscg.com\n", formatted)
}
