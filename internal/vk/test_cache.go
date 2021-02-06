package vk

import "errors"

type testCache struct {
	table map[string]string
}

func (t *testCache) Get(cardName string) (string, error) {
	msg, ok := t.table[cardName]
	if !ok {
		return "", errors.New("test")
	}
	return msg, nil
}

func (t *testCache) Set(cardName string, message string) {
	t.table[cardName] = message
}
