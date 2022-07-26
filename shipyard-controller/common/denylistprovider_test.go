package common

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DenyListProvider(t *testing.T) {
	tests := []struct {
		input  string
		result []string
	}{
		{
			input:  "some\nstring",
			result: []string{"some", "string"},
		},
		{
			input:  "some\nstring\n\n\n",
			result: []string{"some", "string"},
		},
		{
			input:  "some",
			result: []string{"some"},
		},
		{
			input:  "",
			result: []string{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			provider := NewDenyListProvider(bufio.NewScanner(strings.NewReader(tt.input)))
			res := provider.Get()
			require.Equal(t, tt.result, res)
		})
	}
}

func Test_DenyListProvider_removeEmptyStrings(t *testing.T) {
	tests := []struct {
		input  []string
		result []string
	}{
		{
			input:  []string{"some", "string"},
			result: []string{"some", "string"},
		},
		{
			input:  []string{"some", "string", "", "", "something"},
			result: []string{"some", "string", "something"},
		},
		{
			input:  []string{"", ""},
			result: []string{},
		},
		{
			input:  []string{},
			result: []string{},
		},
		{
			input:  nil,
			result: []string{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res := removeEmptyStrings(tt.input)
			require.Equal(t, tt.result, res)
		})
	}
}
