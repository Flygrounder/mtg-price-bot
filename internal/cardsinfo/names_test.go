package cardsinfo

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestGetNameByCardId(t *testing.T) {
	defer gock.Off()

	gock.New(scryfallUrl + "/set/1").Reply(http.StatusOK).JSON(card{
		Name: "card",
	})
	f := &Fetcher{}
	name := f.GetNameByCardId("set", "1")
	assert.Equal(t, "card", name)
}

func TestGetOriginalName_Scryfall(t *testing.T) {
	defer gock.Off()

	gock.New(scryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).JSON(card{
		Name: "Result Card",
	})
	f := &Fetcher{}
	name := f.GetOriginalName("card")
	assert.Equal(t, "Result Card", name)
}

func TestGetOriginalName_DictTwice(t *testing.T) {
	defer gock.Off()

	gock.New(scryfallUrl + "/cards/named?fuzzy=card").Persist().Reply(http.StatusOK).JSON(card{})
	f := &Fetcher{
		Dict: map[string]string{
			"card": "Card",
		},
	}
	name := f.GetOriginalName("card")
	assert.Equal(t, "Card", name)
	name = f.GetOriginalName("card")
	assert.Equal(t, "Card", name)
}

func TestGetOriginalName_BadJson(t *testing.T) {
	defer gock.Off()

	gock.New(scryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).BodyString("}")
	f := &Fetcher{}
	name := f.GetOriginalName("card")
	assert.Equal(t, "", name)
}

func TestGetOriginalName_DoubleSide(t *testing.T) {
	defer gock.Off()

	gock.New(scryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).JSON(card{
		Name:   "Legion's Landing // Adanto, the First Fort",
		Layout: "transform",
	})
	f := &Fetcher{}
	name := f.GetOriginalName("card")
	assert.Equal(t, "Legion's Landing | Adanto, the First Fort", name)
}
