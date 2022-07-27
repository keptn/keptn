package common

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_FileReader(t *testing.T) {
	denyListTestFileName := "test-file"
	tests := []struct {
		name     string
		provider FileReader
		result   []string
	}{
		{
			name: "valid input",
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\nstring")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			name: "valid input with many empty lines at the end",
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\nstring\n\n\n")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			name: "valid input with many empty lines",
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\n\n\nstring\n\n\n")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			name: "valid input with one line",
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some")},
				},
			},
			result: []string{"some"},
		},
		{
			name: "valid input with empty file",
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("")},
				},
			},
			result: []string{},
		},
		{
			name: "error: cannot open file",
			provider: &fileReader{
				FileSystem: fstest.MapFS{},
			},
			result: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.provider.GetLines(denyListTestFileName)
			require.Equal(t, tt.result, res)
		})
	}
}

func Test_FileReader_removeEmptyStrings(t *testing.T) {
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
