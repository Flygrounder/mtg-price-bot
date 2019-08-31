package tests

import (
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	prices, err := cardsinfo.GetSCGPrices("shock")
	assert.Nil(t, err)
	assert.NotEmpty(t, prices)
}
