package cardsinfo

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestGetNameByCardId(t *testing.T) {
	defer gock.Off()

	gock.New(ScryfallUrl + "/set/1").Reply(http.StatusOK).JSON(Card{
		Name: "card",
	})
	name := GetNameByCardId("set", "1")
	assert.Equal(t, "card", name)
}

func TestGetOriginalName_Scryfall(t *testing.T) {
	defer gock.Off()

	gock.New(ScryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).JSON(Card{
		Name: "Result Card",
	})
	name := GetOriginalName("card", nil)
	assert.Equal(t, "Result Card", name)
}

func TestGetOriginalName_Dict(t *testing.T) {
	defer gock.Off()

	gock.New(ScryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).JSON(Card{})
	serialized, _ := json.Marshal(map[string]string{
		"card": "Card",
	})
	dict := strings.NewReader(string(serialized))
	name := GetOriginalName("card", dict)
	assert.Equal(t, "Card", name)
}

func TestGetOriginalName_BadJson(t *testing.T) {
	defer gock.Off()

	gock.New(ScryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).BodyString("}")
	name := GetOriginalName("card", nil)
	assert.Equal(t, "", name)
}

func TestGetOriginalName_DoubleSide(t *testing.T) {
	defer gock.Off()

	gock.New(ScryfallUrl + "/cards/named?fuzzy=card").Reply(http.StatusOK).JSON(Card{
		Name:   "Legion's Landing // Adanto, the First Fort",
		Layout: "transform",
	})
	name := GetOriginalName("card", nil)
	assert.Equal(t, "Legion's Landing | Adanto, the First Fort", name)
}
