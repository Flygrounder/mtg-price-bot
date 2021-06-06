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
	prices, err := f.GetPrices("card")
	assert.Nil(t, err)
	assert.Equal(t, []ScgCardPrice{
		{
			Price:   "$6.99",
			Edition: "Double Masters",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enn/?sku=SGL-MTG-2XM-309-ENN1",
		},
		{
			Price:   "$9.99",
			Edition: "Double Masters (Foil)",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enf/?sku=SGL-MTG-2XM-309-ENF1",
		},
		{
			Price:   "$11.99",
			Edition: "Double Masters - Variants",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enn/?sku=SGL-MTG-2XM2-369-ENN1",
		},
		{
			Price:   "$14.99",
			Edition: "Double Masters - Variants (Foil)",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enf/?sku=SGL-MTG-2XM2-369-ENF1",
		},
		{
			Price:   "$7.99",
			Edition: "Modern Masters: 2013 Edition",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-mma-219-enn/?sku=SGL-MTG-MMA-219-ENN1",
		},
	}, prices)
}

func TestGetPrices_Unavailable(t *testing.T) {
	defer gock.Off()

	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusBadGateway)
	f := &Fetcher{}
	_, err := f.GetPrices("card")
	assert.NotNil(t, err)
}

func TestGetPrices_Empty(t *testing.T) {
	defer gock.Off()

	file, _ := os.Open("test_data/EmptyTest.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(file)
	f := &Fetcher{}
	prices, err := f.GetPrices("card")
	assert.Nil(t, err)
	assert.Nil(t, prices)
}

func TestGetPrices_FilterNonCards(t *testing.T) {
	defer gock.Off()

	file, _ := os.Open("test_data/NonCards.html")
	gock.New(scgSearchUrlTemplate + "card").Reply(http.StatusOK).Body(file)
	f := &Fetcher{}
	prices, err := f.GetPrices("card")
	assert.Nil(t, err)
	expected := []ScgCardPrice{
		{
			Price:   "$72.99",
			Edition: "3rd Edition - Black Border",
			Link:    "https://starcitygames.com/sol-ring-sgl-mtg-3bb-274-frn/?sku=SGL-MTG-3BB-274-FRN3",
		},
		{
			Price:   "$24.99",
			Edition: "3rd Edition/Revised",
			Link:    "https://starcitygames.com/sol-ring-sgl-mtg-3ed-274-enn/?sku=SGL-MTG-3ED-274-ENN1",
		},
		{
			Price:   "$1,999.99",
			Edition: "Alpha",
			Link:    "https://starcitygames.com/sol-ring-sgl-mtg-lea-269-enn/?sku=SGL-MTG-LEA-269-ENN1",
		},
		{
			Price:   "$1,199.99",
			Edition: "Beta",
			Link:    "https://starcitygames.com/sol-ring-sgl-mtg-leb-270-enn/?sku=SGL-MTG-LEB-270-ENN1",
		},
		{
			Price:   "$99.99",
			Edition: "Collectors' Edition",
			Link:    "https://starcitygames.com/sol-ring-sgl-mtg-ced-270-enn/?sku=SGL-MTG-CED-270-ENN1",
		},
	}
	assert.Equal(t, expected, prices)
}
