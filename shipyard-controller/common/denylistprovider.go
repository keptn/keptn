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
	File *os.File
}

const denyListFileName = "/keptn-git-config/git-remote-url-denylist"

func NewDenyListProvider() DenyListProvider {
	gitConfigFile, err := os.Open(denyListFileName)
	if err != nil {
		logrus.Errorf("cannot open %s file: %s", denyListFileName, err.Error())
	}
	return denyListProvider{
		File: gitConfigFile,
	}
}

func (d denyListProvider) Get() []string {
	_, err := d.File.Seek(0, io.SeekStart)
	if err != nil {
		logrus.Errorf("cannot seek %s file: %s", denyListFileName, err.Error())
	}

	scanner := bufio.NewScanner(d.File)
	fileLines := []string{}

	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
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
