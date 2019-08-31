package tests

import (
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCardByStringFull(t *testing.T) {
	name := cardsinfo.GetOriginalName("Шок")
	assert.Equal(t, "Shock", name)
}

func TestGetCardByStringSplit(t *testing.T) {
	name := cardsinfo.GetOriginalName("commit")
	assert.Equal(t, "Commit // Memory", name)
}

func TestGetCardByStringDouble(t *testing.T) {
	name := cardsinfo.GetOriginalName("Legion's landing")
	assert.Equal(t, "Legion's Landing | Adanto, the First Fort", name)
}

func TestGetCardByStringPrefix(t *testing.T) {
	name := cardsinfo.GetOriginalName("Тефери, герой")
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardByStringEnglish(t *testing.T) {
	name := cardsinfo.GetOriginalName("Teferi, Hero of Dominaria")
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardByStringWrong(t *testing.T) {
	name := cardsinfo.GetOriginalName("fwijefiwjfew")
	assert.Equal(t, "", name)
}

func TestGetCardBySetId(t *testing.T) {
	name := cardsinfo.GetNameByCardId("DOM", "207")
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardBySetIdWrong(t *testing.T) {
	name := cardsinfo.GetNameByCardId("DOM", "1207")
	assert.Equal(t, "", name)
}
