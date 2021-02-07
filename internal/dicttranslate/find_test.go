package dicttranslate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindEmpty(t *testing.T) {
	dict := map[string]string{}
	_, f := Find("", dict, 0)
	assert.False(t, f)
}

func TestFindEntry(t *testing.T) {
	dict := map[string]string{
		"entry": "value",
	}
	val, f := Find("entry", dict, 0)
	assert.True(t, f)
	assert.Equal(t, "value", val)
}
