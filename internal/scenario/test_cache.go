package scenario

import (
	"errors"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type testCache struct {
	table map[string][]cardsinfo.ScgCardPrice
}

func (t *testCache) Get(cardName string) ([]cardsinfo.ScgCardPrice, error) {
	msg, ok := t.table[cardName]
	if !ok {
		return nil, errors.New("test")
	}
	return msg, nil
}

func (t *testCache) Set(cardName string, prices []cardsinfo.ScgCardPrice) {
	t.table[cardName] = prices
}
