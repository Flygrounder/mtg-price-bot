package dicttranslate

import (
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"strings"
)

func match(query string, opts []string, maxDist int) (string, bool) {
	bestInd := -1
	bestDist := 0
	for i, s := range opts {
		cfg := levenshtein.DefaultOptions
		cfg.SubCost = 1
		dist := levenshtein.DistanceForStrings([]rune(strings.ToLower(s)), []rune(strings.ToLower(query)), cfg)
		if dist <= maxDist && (bestInd == -1 || dist < bestDist) {
			bestInd = i
			bestDist = dist
		}
	}
	if bestInd == -1 {
		return "", false
	}
	return opts[bestInd], true
}
