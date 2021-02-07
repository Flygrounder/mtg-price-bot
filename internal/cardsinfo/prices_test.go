package cardsinfo

import (
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrices_Ok(t *testing.T) {
	defer gock.Off()

	file, _ := os.Open("test_data/AcademyRuinsTest.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(file)
	f := &Fetcher{}
	prices, err := f.getPrices("card")
	assert.Nil(t, err)
	assert.Equal(t, []scgCardPrice{
		{
			price:   "$6.99",
			edition: "Double Masters",
			link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enn/?sku=SGL-MTG-2XM-309-ENN1",
		},
		{
			price:   "$9.99",
			edition: "Double Masters (Foil)",
			link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enf/?sku=SGL-MTG-2XM-309-ENF1",
		},
		{
			price:   "$11.99",
			edition: "Double Masters - Variants",
			link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enn/?sku=SGL-MTG-2XM2-369-ENN1",
		},
		{
			price:   "$14.99",
			edition: "Double Masters - Variants (Foil)",
			link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enf/?sku=SGL-MTG-2XM2-369-ENF1",
		},
		{
			price:   "$7.99",
			edition: "Modern Masters: 2013 Edition",
			link:    "https://starcitygames.com/academy-ruins-sgl-mtg-mma-219-enn/?sku=SGL-MTG-MMA-219-ENN1",
		},
	}, prices)
}

func TestGetPrices_Unavailable(t *testing.T) {
	defer gock.Off()

	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusBadGateway)
	f := &Fetcher{}
	_, err := f.getPrices("card")
	assert.NotNil(t, err)
}

func TestGetPrices_Empty(t *testing.T) {
	defer gock.Off()

	file, _ := os.Open("test_data/EmptyTest.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(file)
	f := &Fetcher{}
	prices, err := f.getPrices("card")
	assert.Nil(t, err)
	assert.Nil(t, prices)
}
