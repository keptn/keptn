package common

import (
	"bufio"
)

type DenyListProvider interface {
	Get() []string
}

type denyListProvider struct {
	Scanner FileScanner
}

type FileScanner interface {
	Split(split bufio.SplitFunc)
	Scan() bool
	Text() string
}

func NewDenyListProvider(scanner FileScanner) DenyListProvider {
	return denyListProvider{
		Scanner: scanner,
	}
}

func (d denyListProvider) Get() []string {
	fileLines := []string{}
	d.Scanner.Split(bufio.ScanLines)

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
