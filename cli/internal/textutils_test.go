package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StrtoJSON(t *testing.T) {

	var tests = []struct {
		in  string
		out string
		err bool
	}{
		{"", "", true},
		{"a.b.c=v1", `{"a":{"b":{"c":"v1"}}}`, false},
		{"a.b.c=v, a.b.d=v2", `{"a":{"b":{"c":"v","d":"v2"}}}`, false},
		{"a.b.c=v, a.b.d=v2, b.c.d.e=v3", `{"a":{"b":{"c":"v","d":"v2"}},"b":{"c":{"d":{"e":"v3"}}}}`, false}, //TODO this will eventually fail (unordered map)
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := JSONPathToJSONObj(tt.in)
			assert.Equal(t, tt.out, out)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}
