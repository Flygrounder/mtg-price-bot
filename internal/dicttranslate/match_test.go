package dicttranslate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	type testCase struct {
		name       string
		query      string
		opts       []string
		shouldFind bool
		match      string
	}
	tests := []testCase{
		{
			name:       "No options",
			query:      "opt",
			opts:       []string{},
			shouldFind: false,
		},
		{
			name:       "Match one",
			query:      "option",
			opts:       []string{"opt1on"},
			shouldFind: true,
			match:      "opt1on",
		},
		{
			name:       "Match exact",
			query:      "opt1on",
			opts:       []string{"option", "opt1on"},
			shouldFind: true,
			match:      "opt1on",
		},
		{
			name:       "Do not match bad options",
			query:      "random",
			opts:       []string{"option", "opt1on"},
			shouldFind: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, f := match(test.query, test.opts, 1)
			assert.Equal(t, test.shouldFind, f)
			assert.Equal(t, test.match, val)
		})
	}
}
