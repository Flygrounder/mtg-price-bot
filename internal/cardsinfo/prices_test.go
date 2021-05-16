package cardsinfo

import (
	"net/http"
	"os"
	"testing"

	"gopkg.in/h2non/gock.v1"

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

func TestGetPrices_FilterNonCards(t *testing.T) {
	defer gock.Off()

	file, _ := os.Open("test_data/NonCards.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(file)
	f := &Fetcher{}
	prices, err := f.getPrices("card")
	assert.Nil(t, err)
	expected := []scgCardPrice{
		{
			price:   "$72.99",
			edition: "3rd Edition - Black Border",
			link:    "https://starcitygames.com/sol-ring-sgl-mtg-3bb-274-frn/?sku=SGL-MTG-3BB-274-FRN3",
		},
		{
			price:   "$24.99",
			edition: "3rd Edition/Revised",
			link:    "https://starcitygames.com/sol-ring-sgl-mtg-3ed-274-enn/?sku=SGL-MTG-3ED-274-ENN1",
		},
		{
			price:   "$1,999.99",
			edition: "Alpha",
			link:    "https://starcitygames.com/sol-ring-sgl-mtg-lea-269-enn/?sku=SGL-MTG-LEA-269-ENN1",
		},
		{
			price:   "$1,199.99",
			edition: "Beta",
			link:    "https://starcitygames.com/sol-ring-sgl-mtg-leb-270-enn/?sku=SGL-MTG-LEB-270-ENN1",
		},
		{
			price:   "$99.99",
			edition: "Collectors' Edition",
			link:    "https://starcitygames.com/sol-ring-sgl-mtg-ced-270-enn/?sku=SGL-MTG-CED-270-ENN1",
		},
	}
	assert.Equal(t, expected, prices)
}
