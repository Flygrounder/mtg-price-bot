package cardsinfo

import (
	"context"
	"fmt"
	"net/url"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
)

const scgDomain = "https://starcitygames.com"
const scgSearchUrlTemplate = "https://starcitygames.hawksearch.com/sites/starcitygames/?search_query=%v"

func GetPrices(name string) ([]CardPrice, error) {
	prices, err := GetPricesScg(name)
	if err != nil {
		return nil, err
	}
	if len(prices) > 5 {
		return prices[:5], nil
	}
	return prices, nil
}

func GetPricesScg(name string) ([]CardPrice, error) {
	escapedName := url.QueryEscape(name)
	searchUrl := fmt.Sprintf(scgSearchUrlTemplate, escapedName)
	node, err := htmlquery.LoadURL(searchUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load url")
	}
	blocks := htmlquery.Find(node, "//div[@class=\"hawk-results-item\"]")
	var results []CardPrice
	for _, block := range blocks {
		price := &ScgCardPrice{}
		linkNode := htmlquery.FindOne(block, "//h2/a")
		for _, attr := range linkNode.Attr {
			if attr.Key == "href" {
				price.Link = scgDomain + attr.Val
				break
			}
		}
		editionNode := htmlquery.FindOne(block, "//p[@class=\"hawk-results-item__category\"]/a")
		price.Edition = editionNode.FirstChild.Data
		priceNode := htmlquery.FindOne(block, "//div[contains(concat(' ',normalize-space(@class),' '),' hawk-results-item__options-table-cell--price ')]")
		price.Price = priceNode.FirstChild.Data
		results = append(results, price)
	}
	return results, nil
}

func GetPricesTcg(name string) ([]CardPrice, error) {
	client, err := scryfall.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot fetch prices")
	}
	ctx := context.Background()
	opts := scryfall.SearchCardsOptions{
		Unique: scryfall.UniqueModePrints,
	}
	resp, err := client.SearchCards(ctx, fmt.Sprintf("!\"%v\"", name), opts)
	var prices []CardPrice
	for _, card := range resp.Cards {
		edition := card.SetName + " #" + card.CollectorNumber
		if card.Prices.USD == "" && card.Prices.USDFoil == "" {
			continue
		}
		cardPrice := &TcgCardPrice {
			Edition: edition,
			Price: card.Prices.USD,
			PriceFoil: card.Prices.USDFoil,
			Name: card.Name,
			Link: card.PurchaseURIs.TCGPlayer,
		}
		prices = append(prices, cardPrice)
	}
	return prices, nil
}
