package cardsinfo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const scryfallUrl = "https://api.scryfall.com"

func (f *Fetcher) GetNameByCardId(set string, number string) string {
	/*
		Note: number is string because some cards contain letters in their numbers.
	*/
	path := scryfallUrl + "/cards/" + strings.ToLower(set) + "/" + number
	return getCardByUrl(path)
}

func (f *Fetcher) GetOriginalName(name string) string {
	path := scryfallUrl + "/cards/named?fuzzy=" + applyFilters(name)
	result := getCardByUrl(path)
	return result
}

func applyFilters(name string) string {
	/*
		Despite of the rules of Russian language, letter ё is replaced with e on cards
		Sometimes it leads to wrong search results
	*/
	name = strings.ReplaceAll(name, "ё", "е")
	return url.QueryEscape(name)
}

func getCardByUrl(path string) string {
	response, err := http.Get(path)
	if err != nil {
		return ""
	}
	defer func() {
		_ = response.Body.Close()
	}()
	data, _ := ioutil.ReadAll(response.Body)
	var v card
	err = json.Unmarshal(data, &v)
	if err != nil {
		return ""
	}
	return v.getName()
}
