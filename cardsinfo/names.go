package cardsinfo

import (
	mtg "github.com/MagicTheGathering/mtg-sdk-go"
	"strings"
)

func GetNameByCardId(set string, number string) string {
	/*
		From https://docs.magicthegathering.io/#api_v1cards_list

		Number: This is a string, not an integer,
		because some cards have letters in their numbers.
	*/
	cards, _, _ := mtg.NewQuery().Where(mtg.CardSet, set).Where(mtg.CardNumber, number).PageS(1, 1)
	name := fetchCardNameFromSlice(cards)
	return name
}

func GetOriginalName(name string) string {
	langs := []string{"Russian", ""}
	channel := make(chan string)
	for i := range langs {
		go getOriginalNameFromLang(name, langs[i], channel)
	}
	for i := 0; i < len(langs); i++ {
		name := <-channel
		if name != "" {
			return name
		}
	}
	return ""
}

func getOriginalNameFromLang(name, lang string, channel chan string) {
	cards, _, _ := mtg.NewQuery().Where(mtg.CardLanguage, lang).Where(mtg.CardName, name).PageS(1, 1)
	originalName := fetchCardNameFromSlice(cards)
	channel <- originalName
}

func fetchCardNameFromSlice(cards []*mtg.Card) string {
	if len(cards) > 0 {
		name := getCardName(cards[0])
		return name
	}
	return ""
}

func getCardName(card *mtg.Card) string {
	switch card.Layout {
	case "split":
		return strings.Join(card.Names, " // ")
	case "transform":
		return strings.Join(card.Names, " | ")
	default:
		return card.Name
	}
}
