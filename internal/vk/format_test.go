package vk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

func TestFormatCardPrices(t *testing.T) {
	formatted := formatCardPrices("card", []cardsinfo.ScgCardPrice{
		{
			Price:   "1.5$",
			Edition: "ED",
			Link:    "scg.com",
		},
	})
	assert.Equal(t, "Оригинальное название: card\n\n1. ED: 1.5$\nscg.com\n", formatted)
}
