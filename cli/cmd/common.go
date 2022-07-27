package cmd

import (
	"io/fs"
	"os"
)

func areStringFlagsSet(el ...*string) bool {
	for _, e := range el {
		if !isStringFlagSet(e) {
			return false
		}
	}
	return true
}

func isStringFlagSet(s *string) bool {
	return s != nil && *s != ""
}

func isStringFlagNotSet(s *string) bool {
	return s == nil || (s != nil && *s == "")
}

func areBoolFlagsSet(el ...*bool) bool {
	for _, e := range el {
		if !isBoolFlagSet(e) {
			return false
		}
	}
	return true
}

func isBoolFlagSet(b *bool) bool {
	return b != nil && *b
}

type fileOpener struct {
}

func (f fileOpener) Open(name string) (fs.File, error) {
	return os.Open(name)
}
