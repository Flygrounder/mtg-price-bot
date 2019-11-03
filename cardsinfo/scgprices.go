package cardsinfo

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

const Scgurl = "http://www.starcitygames.com/results?name="

func GetSCGPrices(name string) ([]CardPrice, error) {
	preprocessedName := preprocessNameForSearch(name)
	url := getSCGUrl(preprocessedName)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}
	return fetchPrices(doc)
}

func preprocessNameForSearch(name string) string {
	return strings.Split(name, "|")[0]
}

func fetchPrices(doc *html.Node) ([]CardPrice, error) {
	priceContainers := getPriceContainers(doc)
	prices := make([]CardPrice, 0)
	for _, container := range priceContainers {
		name := parseName(container)
		edition := parseEdition(container)
		price := parsePrice(container)
		link := parseLink(container)
		cardPrice := buildCardPrice(name, edition, price, link)
		if isValidPrice(&cardPrice) {
			prices = append(prices, cardPrice)
		}
	}
	return prices, nil
}

func isValidPrice(price *CardPrice) bool {
	isValid := true
	isValid = isValid && (price.Name != "")
	isValid = isValid && (price.Edition != "")
	isValid = isValid && (price.Price != 0.0)
	isValid = isValid && (price.Link != "")
	return isValid
}

func buildCardPrice(name, edition string, price float64, link string) CardPrice {
	cardPrice := CardPrice{
		Name:    name,
		Edition: edition,
		Price:   price,
		Link:    link,
	}
	return cardPrice
}

func getPriceContainers(doc *html.Node) []*html.Node {
	nodesOdd := htmlquery.Find(doc, "//tr[contains(@class, 'deckdbbody_row')]")
	nodesEven := htmlquery.Find(doc, "//tr[contains(@class, 'deckdbbody2_row')]")
	nodes := append(nodesOdd, nodesEven...)
	return nodes
}

func fetchPrice(price string) (float64, error) {
	split := strings.Split(price, "$")
	if len(split) < 2 {
		return 0, errors.New("not enough values")
	}
	p := split[1]
	v, err := strconv.ParseFloat(p, 64)
	return v, err
}

func getSCGUrl(name string) string {
	words := strings.Split(name, " ")
	scgName := strings.Join(words, "+")
	url := Scgurl + scgName
	return url
}

func parseName(container *html.Node) string {
	nameNode := htmlquery.FindOne(container, "//td[contains(@class, 'search_results_1')]")
	if nameNode == nil {
		return ""
	}
	name := htmlquery.InnerText(nameNode)
	return name
}

func parseEdition(container *html.Node) string {
	editionNode := htmlquery.FindOne(container, "//td[contains(@class, 'search_results_2')]")
	if editionNode == nil {
		return ""
	}
	edition := strings.Trim(htmlquery.InnerText(editionNode), "\n ")
	return edition
}

func parsePrice(container *html.Node) float64 {
	priceNode := htmlquery.FindOne(container, "//td[contains(@class, 'search_results_9')]")
	if priceNode == nil {
		return 0.0
	}
	priceString := htmlquery.InnerText(priceNode)
	price, err := fetchPrice(priceString)
	if err != nil {
		return 0.0
	}
	return price
}

func parseLink(container *html.Node) string {
	linkNodes := htmlquery.Find(container, "//td[contains(@class, 'search_results_1')]/b/a")
	if len(linkNodes) == 0 {
		return ""
	}
	linkNode := linkNodes[0]
	link := fetchLink(linkNode)
	return link
}

func fetchLink(linkContainer *html.Node) string {
	for _, v := range linkContainer.Attr {
		if v.Key == "href" {
			return v.Val
		}
	}
	return ""
}
