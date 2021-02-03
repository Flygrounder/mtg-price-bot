package dicttranslate

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func find(query string, dict map[string]string, maxDist int) (string, bool) {
	var keys []string
	for i := range dict {
		keys = append(keys, i)
	}
	key, f := match(query, keys, maxDist)
	return dict[key], f
}

func FindFromReader(query string, reader io.Reader, maxDist int) (string, bool) {
	content, _ := ioutil.ReadAll(reader)
	dict := map[string]string{}
	_ = json.Unmarshal(content, &dict)
	return find(query, dict, maxDist)
}
