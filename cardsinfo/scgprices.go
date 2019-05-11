package cardsinfo

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"strconv"
	"strings"
)

const SCGURL = "http://www.starcitygames.com/results?name="

func GetSCGPrices(name string) ([]CardPrice, error) {
	splitted := strings.Split(name, " ")
	scgName := strings.Join(splitted, "+")
	url := SCGURL + scgName
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}
	nodesOdd := htmlquery.Find(doc, "//tr[contains(@class, 'deckdbbody_row')]")
	nodesEven := htmlquery.Find(doc, "//tr[contains(@class, 'deckdbbody2_row')]")
	nodes := append(nodesOdd, nodesEven...)
	prices := make([]CardPrice, 0)
	for _, node := range nodes {
		nameNode := htmlquery.FindOne(node, "//td[contains(@class, 'search_results_1')]")
		if nameNode == nil {
			continue
		}
		name := htmlquery.InnerText(nameNode)
		editionNode := htmlquery.FindOne(node, "//td[contains(@class, 'search_results_2')]")
		if editionNode == nil {
			continue
		}
		edition := strings.Trim(htmlquery.InnerText(editionNode), "\n ")
		priceNode := htmlquery.FindOne(node, "//td[contains(@class, 'search_results_9')]")
		if priceNode == nil {
			continue
		}
		priceS := htmlquery.InnerText(priceNode)
		price, err := fetchPrice(priceS)
		if err != nil {
			continue
		}
		linkNodes := htmlquery.Find(node, "//td[contains(@class, 'search_results_2')]/a")
		if len(linkNodes) == 0 {
			continue
		}
		linkNode := linkNodes[0]
		var link string
		for _, v := range linkNode.Attr {
			if v.Key == "href" {
				link = v.Val
			}
		}
		obj := CardPrice{
			Name:    name,
			Price:   price,
			Edition: edition,
			Link:    link,
		}
		prices = append(prices, obj)
	}
	return prices, nil
}

func fetchPrice(price string) (float64, error) {
	split := strings.Split(price, "$")
	if len(split) < 2 {
		return 0, errors.New("Not enough values")
	}
	p := split[1]
	v, err := strconv.ParseFloat(p, 64)
	return v, err
}
