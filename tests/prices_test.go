package tests

import (
	"github.com/flygrounder/go-mtg-vk/cardsinfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	prices, err := cardsinfo.GetPrices("Scavenging Ooze")
	assert.Nil(t, err)
	assert.NotEmpty(t, prices)
}
