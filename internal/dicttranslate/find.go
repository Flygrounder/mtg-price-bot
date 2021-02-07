package dicttranslate

func Find(query string, dict map[string]string, maxDist int) (string, bool) {
	var keys []string
	for i := range dict {
		keys = append(keys, i)
	}
	key, f := match(query, keys, maxDist)
	return dict[key], f
}
