package cardsinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchPriceOne(t *testing.T) {
	price, _ := fetchPrice("$2.28")
	assert.Equal(t, 2.28, price)
}

func TestFetchPriceTwo(t *testing.T) {
	price, _ := fetchPrice("$2.28$1.14")
	assert.Equal(t, 2.28, price)
}

func TestFetchPriceNo(t *testing.T) {
	_, err := fetchPrice("")
	assert.NotNil(t, err)
}
