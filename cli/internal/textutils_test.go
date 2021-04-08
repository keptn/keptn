package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JSONPathToJSONObj(t *testing.T) {

	var tests = []struct {
		name string
		in   []string
		out  string
		err  bool
	}{
		{"1", []string{""}, "", true},
		{"2", []string{"a.b.c=v1"}, `{"a":{"b":{"c":"v1"}}}`, false},
		{"3", []string{"a.b.c=v", "a.b.d=v2"}, `{"a":{"b":{"c":"v","d":"v2"}}}`, false},
		{"4", []string{"a.b.c=v", "a.b.d=v2", "b.c.d.e=v3"}, `{"a":{"b":{"c":"v","d":"v2"}},"b":{"c":{"d":{"e":"v3"}}}}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := JSONPathToJSONObj(tt.in)
			assert.Equal(t, tt.out, out)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}
