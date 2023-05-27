package scenario

import (
	"context"
	"errors"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type testCache struct {
	table map[string][]cardsinfo.ScgCardPrice
}

func (t *testCache) Init(ctx context.Context) error {
	return nil
}

func (t *testCache) Get(ctx context.Context, cardName string) ([]cardsinfo.ScgCardPrice, error) {
	msg, ok := t.table[cardName]
	if !ok {
		return nil, errors.New("test")
	}
	return msg, nil
}

func (t *testCache) Set(ctx context.Context, cardName string, prices []cardsinfo.ScgCardPrice) error {
	t.table[cardName] = prices
	return nil
}
