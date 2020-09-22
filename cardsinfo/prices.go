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
		fullArtString := ""
		if card.FullArt {
			fullArtString = " (Fullart)"
		}
		edition := card.SetName + fullArtString
		if card.Prices.USD == "" && card.Prices.USDFoil == "" {
			continue
		}
		cardPrice := CardPrice {
			Edition: edition,
			Price: card.Prices.USD,
			PriceFoil: card.Prices.USDFoil,
			Name: card.Name,
			Link: card.PurchaseURIs.TCGPlayer,
		}
		prices = append(prices, cardPrice)
	}
	if len(prices) > 5 {
		return prices[:5], nil
	}
	return prices, nil
}
