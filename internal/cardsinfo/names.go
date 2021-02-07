package cardsinfo

import (
	"encoding/json"
	"gitlab.com/flygrounder/go-mtg-vk/internal/dicttranslate"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const ScryfallUrl = "https://api.scryfall.com"

func (f *Fetcher) GetNameByCardId(set string, number string) string {
	/*
		Note: number is string because some cards contain letters in their numbers.
	*/
	path := ScryfallUrl + "/cards/" + strings.ToLower(set) + "/" + number
	return GetCardByUrl(path)
}

func (f *Fetcher) GetOriginalName(name string) string {
	path := ScryfallUrl + "/cards/named?fuzzy=" + ApplyFilters(name)
	result := GetCardByUrl(path)
	if result == "" && f.Dict != nil {
		result, _ = dicttranslate.FindFromReader(name, f.Dict, 5)
	}
	return result
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
	response, _ := http.Get(path)
	defer func() {
		_ = response.Body.Close()
	}()
	data, _ := ioutil.ReadAll(response.Body)
	var v Card
	err := json.Unmarshal(data, &v)
	if err != nil {
		return ""
	}
	return v.getName()
}
