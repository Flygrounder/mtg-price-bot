package cardsinfo

import (
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
)

const scgDomain = "https://starcitygames.com"
const scgSearchUrlTemplate = "https://starcitygames.hawksearch.com/sites/starcitygames/?search_query="

func (f *Fetcher) GetPrices(name string) ([]ScgCardPrice, error) {
	prices, err := getPricesScg(name)
	if err != nil {
		return nil, err
	}
	if len(prices) > 5 {
		return prices[:5], nil
	}
	return prices, nil
}

func getPricesScg(name string) ([]ScgCardPrice, error) {
	escapedName := url.QueryEscape(name)
	searchUrl := scgSearchUrlTemplate + escapedName
	node, err := htmlquery.LoadURL(searchUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load url")
	}
	blocks := htmlquery.Find(node, "//div[@class=\"hawk-results-item\"]")
	var results []ScgCardPrice
	for _, block := range blocks {
		price := ScgCardPrice{}
		linkNode := htmlquery.FindOne(block, "//h2/a")
		price.Link = scgDomain + htmlquery.SelectAttr(linkNode, "href")
		editionNode := htmlquery.FindOne(block, "//p[@class=\"hawk-results-item__category\"]/a")
		if !strings.HasPrefix(htmlquery.SelectAttr(editionNode, "href"), "/shop/singles/") {
			continue
		}
		price.Edition = editionNode.FirstChild.Data
		priceNode := htmlquery.FindOne(block, "//span[@class='hawk-old-price']|//div[contains(concat(' ',normalize-space(@class),' '),' hawk-results-item__options-table-cell--price ')]")
		price.Price = priceNode.FirstChild.Data
		results = append(results, price)
	}
	return results, nil
}
