package scenario

import (
	"errors"
)

type testInfoFetcher struct{}

func (t *testInfoFetcher) GetFormattedCardPrices(name string) (string, error) {
	if name == "good" || name == "uncached" {
		return name, nil
	}
	return "", errors.New("test")
}

func (t *testInfoFetcher) GetNameByCardId(_ string, _ string) string {
	return "good"
}

func (t *testInfoFetcher) GetOriginalName(name string) string {
	if name == "good" || name == "bad" || name == "uncached" {
		return name
	}
	return ""
}
