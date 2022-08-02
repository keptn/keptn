package filereader

import (
	"bufio"
	"io/fs"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type FileReader interface {
	GetLines(filePath string) []string
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

func New() *fileReader {
	return &fileReader{
		FileSystem: fileOpener{},
	}
}

func (d *fileReader) GetLines(filePath string) []string {
	gitConfigFile, err := d.FileSystem.Open(filePath)
	if err != nil {
		logrus.Errorf("Cannot open %s: %s", filePath, err.Error())
		return []string{}
	}
	defer gitConfigFile.Close()

	lines := []string{}
	scanner := bufio.NewScanner(gitConfigFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return removeEmptyStrings(lines)
}

func removeEmptyStrings(s []string) []string {
	r := []string{}
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}
