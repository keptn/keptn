package models

import (
	"errors"
	"strings"
)

func validateEntityName(s string) error {
	if strings.Contains(s, " ") {
		return errors.New("stage name must not contain whitespaces")
	}
	if strings.Contains(s, "/") {
		return errors.New("stage name must not contain '/'")
	}
	if strings.ReplaceAll(s, " ", "") == "" {
		return errors.New("stage name must not be empty")
	}
	return nil
}
