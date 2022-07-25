package common

import (
	"fmt"
	"regexp"
)

type remoteURLValidator struct {
	denyListProvider DenyListProvider
}

type RemoteURLValidator interface {
	Validate(url string) error
}

func NewRemoteURLValidator(denyListProvider DenyListProvider) RemoteURLValidator {
	validator := remoteURLValidator{
		denyListProvider: denyListProvider,
	}
	return validator
}

func (c remoteURLValidator) Validate(url string) error {
	denyList := c.denyListProvider.Get()

	for _, item := range denyList {
		if res, _ := regexp.MatchString(item, url); res {
			return fmt.Errorf("invalid RemoteURL: %s", url)
		}
	}
	return nil
}
