package common

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystem_WriteAndReadFile(t *testing.T) {
	// create a tmp directory in test/tmp
	dir, err := ioutil.TempDir("../test/tmp/", "project-")
	require.Nil(t, err)
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			t.Logf("could not delete tmp directory: %v", err)
		}
	}(dir)

	fs := FileSystem{}

	filePath := dir + "/my-file"

	fileContent := "content"

	err = fs.WriteFile(filePath, []byte(fileContent))
	require.Nil(t, err)

	fileExists := fs.FileExists(filePath)
	require.True(t, fileExists)

	res, err := fs.ReadFile(filePath)
	require.Nil(t, err)

	require.Equal(t, fileContent, string(res))

	err = fs.DeleteFile(filePath)
	require.Nil(t, err)

	fileExists = fs.FileExists(filePath)
	require.False(t, fileExists)
}
