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
	prices, err := f.GetPrices("card")
	assert.Nil(t, err)
	assert.Equal(t, []CardPrice{
		&ScgCardPrice{
			Price:   "$6.99",
			Edition: "Double Masters",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enn/?sku=SGL-MTG-2XM-309-ENN1",
		},
		&ScgCardPrice{
			Price:   "$9.99",
			Edition: "Double Masters (Foil)",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm-309-enf/?sku=SGL-MTG-2XM-309-ENF1",
		},
		&ScgCardPrice{
			Price:   "$11.99",
			Edition: "Double Masters - Variants",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enn/?sku=SGL-MTG-2XM2-369-ENN1",
		},
		&ScgCardPrice{
			Price:   "$14.99",
			Edition: "Double Masters - Variants (Foil)",
			Link:    "https://starcitygames.com/academy-ruins-sgl-mtg-2xm2-369-enf/?sku=SGL-MTG-2XM2-369-ENF1",
		},
		&ScgCardPrice{
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
