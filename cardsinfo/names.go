package cardsinfo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const ScryfallUrl = "https://api.scryfall.com"

func GetNameByCardId(set string, number string) string {
	/*
		Note: number is string because some cards contain letters in their numbers.
	*/
	path := ScryfallUrl + "/cards/" + strings.ToLower(set) + "/" + number
	return GetCardByUrl(path)
}

func GetOriginalName(name string) string {
	path := ScryfallUrl + "/cards/named?fuzzy=" + ApplyFilters(name)
	return GetCardByUrl(path)
}

func ApplyFilters(name string) string {
	/*
		Despite of the rules of Russian language, letter ё is replaced with e on cards
		Sometimes it leads to wrong search results
	*/
	name = strings.ReplaceAll(name, "ё", "е")
	return url.QueryEscape(name)
}

func GetCardByUrl(path string) string {
	response, err := http.Get(path)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	var v Card
	err = json.Unmarshal(data, &v)
	if err != nil {
		return ""
	}
	return v.getName()
}
