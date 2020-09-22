package cardsinfo

import (
	"context"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/pkg/errors"
)

func GetPrices(name string) ([]CardPrice, error) {
	client, err := scryfall.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot fetch prices")
	}
	ctx := context.Background()
	opts := scryfall.SearchCardsOptions{
		Unique: scryfall.UniqueModePrints,
	}
	resp, err := client.SearchCards(ctx, name, opts)
	var prices []CardPrice
	for _, card := range resp.Cards {
		foilString := ""
		price := card.Prices.USD
		if card.Foil {
			foilString = "(Foil)"
			price = card.Prices.USDFoil
		}
		edition := card.SetName + foilString
		cardPrice := CardPrice {
			Edition: edition,
			Price: price,
			Name: card.Name,
			Link: card.PurchaseURIs.TCGPlayer,
		}
		prices = append(prices, cardPrice)
	}
	return prices, nil
}
