package common

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_FileReader(t *testing.T) {
	denyListTestFileName := "test-file"
	tests := []struct {
		provider FileReader
		result   []string
	}{
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\nstring")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\nstring\n\n\n")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some\n\n\nstring\n\n\n")},
				},
			},
			result: []string{"some", "string"},
		},
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("some")},
				},
			},
			result: []string{"some"},
		},
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{
					denyListTestFileName: {Data: []byte("")},
				},
			},
			result: []string{},
		},
		{
			provider: &fileReader{
				FileSystem: fstest.MapFS{},
			},
			result: []string{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res := tt.provider.Get(denyListTestFileName)
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
