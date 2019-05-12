package cardsinfo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const SCRYFALL_URL = "https://api.scryfall.com"

func GetNameByCardId(set string, number string) string {
	/*
		Note: number is string because some cards contain letters in their numbers.
	*/
	path := SCRYFALL_URL + "/cards/" + strings.ToLower(set) + "/" + number
	return GetCardByUrl(path)
}

func GetOriginalName(name string) string {
	path := SCRYFALL_URL + "/cards/named?fuzzy=" + url.QueryEscape(name)
	return GetCardByUrl(path)
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
	json.Unmarshal(data, &v)
	return v.getName()
}
