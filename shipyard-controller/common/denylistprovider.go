package common

import (
	"bufio"
	"os"

	"github.com/sirupsen/logrus"
)

type DenyListProvider interface {
	Get() []string
}

type denyListProvider struct {
	Scanner FileScanner
}

type FileScanner interface {
	Scan() bool
	Text() string
}

const denyListFileName = "/keptn-git-config/git-remote-url-denylist"

func NewDenyListProvider() DenyListProvider {
	gitConfigFile, err := os.Open(denyListFileName)
	if err != nil {
		logrus.Errorf("cannot open keptn-git-config file %s", err.Error())
	}
	defer gitConfigFile.Close()
	return denyListProvider{
		Scanner: bufio.NewScanner(gitConfigFile),
	}
}

func (d denyListProvider) Get() []string {
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
