package dicttranslate

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFindEmpty(t *testing.T) {
	dict := map[string]string{}
	_, f := find("", dict, 0)
	assert.False(t, f)
}

func TestFindEntry(t *testing.T) {
	dict := map[string]string{
		"entry": "value",
	}
	val, f := find("entry", dict, 0)
	assert.True(t, f)
	assert.Equal(t, "value", val)
}

func TestFindFromReaderFail(t *testing.T) {
	_, f := FindFromReader("entry", strings.NewReader("{}"), 0)
	assert.False(t, f)
}

func TestFindFromReaderSuccess(t *testing.T) {
	value, f := FindFromReader("entry", strings.NewReader("{\"entry\":\"value\"}"), 0)
	assert.True(t, f)
	assert.Equal(t, "value", value)
}
