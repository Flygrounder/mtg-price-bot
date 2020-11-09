package cardsinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	prices, err := GetPrices("Black lotus")
	assert.Nil(t, err)
	assert.NotEmpty(t, prices)
}
