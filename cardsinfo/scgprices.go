package cardsinfo

import (
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

const Scgurl = "https://www.starcitygames.com/search.php?search_query="
const Scgapi = "https://newstarcityconnector.herokuapp.com/eyApi/products/"
const MaxCards = 4

func GetSCGPrices(name string) ([]CardPrice, error) {
	preprocessedName := preprocessNameForSearch(name)
	url := getSCGUrl(preprocessedName)
	doc, err := getScgHTML(url)
	if err != nil {
		return nil, err
	}
	return fetchPrices(doc)
}

func getScgHTML(url string) (*html.Node, error) {
	response, err := http.Get(url)
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("not ok status")
	}

	r, err := charset.NewReader(response.Body, response.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	return html.Parse(r)
}

func preprocessNameForSearch(name string) string {
	return strings.Split(name, "|")[0]
}

func fetchPrices(doc *html.Node) ([]CardPrice, error) {
	priceContainers := getPriceContainers(doc)
	if MaxCards < len(priceContainers) {
		priceContainers = priceContainers[:MaxCards]
	}
	length := len(priceContainers)
	prices := make(chan CardPrice, length)
	finished := make(chan bool, length)
	cardPrices := make([]CardPrice, 0)
	for _, container := range priceContainers {
		go processCard(container, prices, finished)
	}
	processed := 0
	for {
		c := <-finished
		processed++
		if c {
			cardPrices = append(cardPrices, <-prices)
		}
		if processed == length {
			break
		}
	}
	return cardPrices, nil
}

func processCard(container *html.Node, prices chan CardPrice, finished chan bool) {
	name := parseName(container)
	edition := parseEdition(container)
	price := parsePrice(container)
	link := parseLink(container)
	cardPrice := buildCardPrice(name, edition, price, link)
	if isValidPrice(&cardPrice) {
		prices <- cardPrice
		finished <- true
	} else {
		finished <- false
	}
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
	nodes := htmlquery.Find(doc, "//tr[@class='product']")
	return nodes
}

func getItemId(item *html.Node) string {
	return htmlquery.SelectAttr(item, "data-id")
}

func getPriceById(id string) float64 {
	path := Scgapi + id + "/variants"
	resp, err := http.Get(path)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0.0
	}
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.0
	}
	var scgResponse ScgResponse
	err = json.Unmarshal(respString, &scgResponse)
	if err != nil {
		return 0.0
	}
	for _, v := range scgResponse.Response.Data {
		if len(v.OptionValues) > 0 && v.OptionValues[0].Label == "Near Mint" {
			return v.Price
		}
	}
	return 0.0
}

func getSCGUrl(name string) string {
	scgName := strings.Replace(name, " ", "+", -1)
	url := Scgurl + scgName
	return url
}

func parseName(container *html.Node) string {
	nameNode := htmlquery.FindOne(container, "//h4[@class='listItem-title']")
	if nameNode == nil {
		return ""
	}
	name := htmlquery.InnerText(nameNode)
	name = strings.Trim(name, "\n ")
	return name
}

func parseEdition(container *html.Node) string {
	editionNode := htmlquery.FindOne(container, "//span[@class='category-row-name-search']")
	if editionNode == nil {
		return ""
	}
	edition := strings.Trim(htmlquery.InnerText(editionNode), "\n ")
	parts := strings.Split(edition, "/")
	if len(parts) == 0 {
		return ""
	}
	last := len(parts) - 1
	return parts[last]
}

func parsePrice(container *html.Node) float64 {
	id := getItemId(container)
	price := getPriceById(id)
	return price
}

func parseLink(container *html.Node) string {
	linkNodes := htmlquery.Find(container, "//h4[@class='listItem-title']/a")
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
