package cardsinfo

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestFetcher_GetFormattedCardPrices_Error(t *testing.T) {
	defer gock.Off()

	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusInternalServerError)
	f := &Fetcher{}
	_, err := f.GetFormattedCardPrices("card")
	assert.NotNil(t, err)
}

func TestFetcher_GetFormattedCardPrices_Empty(t *testing.T) {
	defer gock.Off()

	resp, _ := os.Open("test_data/EmptyTest.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(resp)
	f := &Fetcher{}
	msg, err := f.GetFormattedCardPrices("card")
	assert.Nil(t, err)
	assert.Equal(t, "Оригинальное название: card\n\nЦен не найдено\n", msg)
}

func TestFormatCardPrices(t *testing.T) {
	f := &Fetcher{}
	formatted := f.formatCardPrices("card", []scgCardPrice{
		{
			price:   "1.5$",
			edition: "ED",
			link:    "scg.com",
		},
	})
	assert.Equal(t, "Оригинальное название: card\n\n1. ED: 1.5$\nscg.com\n", formatted)
}
