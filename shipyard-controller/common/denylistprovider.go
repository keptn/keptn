package common

import (
	"bufio"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type DenyListProvider interface {
	Get() []string
}

type denyListProvider struct {
	Scanner FileScanner
	File    *os.File
}

type FileScanner interface {
	Scan() bool
	Text() string
}

const denyListFileName = "/keptn-git-config/git-remote-url-denylist"

func NewDenyListProvider() DenyListProvider {
	gitConfigFile, err := os.Open(denyListFileName)
	if err != nil {
		logrus.Errorf("cannot open %s file: %s", denyListFileName, err.Error())
	}
	return denyListProvider{
		Scanner: bufio.NewScanner(gitConfigFile),
		File:    gitConfigFile,
	}
}

func (d denyListProvider) Get() []string {
	_, err := d.File.Seek(0, io.SeekStart)
	if err != nil {
		logrus.Errorf("cannot seek %s file: %s", denyListFileName, err.Error())
	}

	fileLines := []string{}

	for d.Scanner.Scan() {
		fileLines = append(fileLines, d.Scanner.Text())
	}

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
