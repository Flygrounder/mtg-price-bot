package cardsinfo

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetCardByStringFull(t *testing.T) {
	name := GetOriginalName("Шок", nil)
	assert.Equal(t, "Shock", name)
}

func TestGetCardByStringSplit(t *testing.T) {
	name := GetOriginalName("commit", nil)
	assert.Equal(t, "Commit // Memory", name)
}

func TestGetCardByStringDouble(t *testing.T) {
	name := GetOriginalName("Legion's landing", nil)
	assert.Equal(t, "Legion's Landing | Adanto, the First Fort", name)
}

func TestGetCardByStringPrefix(t *testing.T) {
	name := GetOriginalName("Тефери, герой", nil)
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardByStringEnglish(t *testing.T) {
	name := GetOriginalName("Teferi, Hero of Dominaria", nil)
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardByStringWrong(t *testing.T) {
	name := GetOriginalName("fwijefiwjfew", nil)
	assert.Equal(t, "", name)
}

func TestGetCardBySetId(t *testing.T) {
	name := GetNameByCardId("DOM", "207")
	assert.Equal(t, "Teferi, Hero of Dominaria", name)
}

func TestGetCardBySetIdWrong(t *testing.T) {
	name := GetNameByCardId("DOM", "1207")
	assert.Equal(t, "", name)
}

func TestGetCardByStringDict(t *testing.T) {
	dictContent := "{\"n0suchc8rdc8n3x1s1\":\"Success\"}"
	name := GetOriginalName("n0suchc8rdc8n3x1s1", strings.NewReader(dictContent))
	assert.Equal(t, "Success", name)
}
