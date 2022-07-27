package common

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type FileReader interface {
	Get(filePath string) []string
}

type fileReader struct {
	FileSystem fs.FS
}

type fileOpener struct {
}

func (f fileOpener) Open(name string) (fs.File, error) {
	return os.Open(name)
}

const RemoteURLDenyListPath = "/keptn-git-config/git-remote-url-denylist"

func NewFileReader() *fileReader {
	return &fileReader{
		FileSystem: fileOpener{},
	}
}

func (d *fileReader) Get(filePath string) []string {
	gitConfigFile, err := d.FileSystem.Open(filePath)
	if err != nil {
		fmt.Printf(err.Error())
		logrus.Errorf("cannot open %s: %s", filePath, err.Error())
		return []string{}
	}
	defer gitConfigFile.Close()

	configFileContent, err := io.ReadAll(gitConfigFile)
	if err != nil {
		logrus.Errorf("cannot read %s: %s", filePath, err.Error())
		return []string{}
	}

	fileLines := strings.Split(strings.ReplaceAll(string(configFileContent), "\r\n", "\n"), "\n")

	return removeEmptyStrings(fileLines)
}

func removeEmptyStrings(s []string) []string {
	r := []string{}
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
