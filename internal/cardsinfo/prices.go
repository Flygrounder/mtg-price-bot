package cardsinfo

import (
	"net/url"

	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
)

const scgDomain = "https://starcitygames.com"
const scgSearchUrlTemplate = "https://starcitygames.hawksearch.com/sites/starcitygames/?search_query="

func (f *Fetcher) getPrices(name string) ([]scgCardPrice, error) {
	prices, err := getPricesScg(name)
	if err != nil {
		return nil, err
	}
	if len(prices) > 5 {
		return prices[:5], nil
	}
	return prices, nil
}

func getPricesScg(name string) ([]scgCardPrice, error) {
	escapedName := url.QueryEscape(name)
	searchUrl := scgSearchUrlTemplate + escapedName
	node, err := htmlquery.LoadURL(searchUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load url")
	}
	blocks := htmlquery.Find(node, "//div[@class=\"hawk-results-item\"]")
	var results []scgCardPrice
	for _, block := range blocks {
		price := scgCardPrice{}
		linkNode := htmlquery.FindOne(block, "//h2/a")
		for _, attr := range linkNode.Attr {
			if attr.Key == "href" {
				price.link = scgDomain + attr.Val
				break
			}
		}
		editionNode := htmlquery.FindOne(block, "//p[@class=\"hawk-results-item__category\"]/a")
		if editionNode.FirstChild != nil {
			price.edition = editionNode.FirstChild.Data
		}
		priceNode := htmlquery.FindOne(block, "//span[@class='hawk-old-price']|//div[contains(concat(' ',normalize-space(@class),' '),' hawk-results-item__options-table-cell--price ')]")
		if priceNode.FirstChild != nil {
			price.price = priceNode.FirstChild.Data
		}
		results = append(results, price)
	}
	return results, nil
}
