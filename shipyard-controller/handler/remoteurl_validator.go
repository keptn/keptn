package handler

import (
	"fmt"
	"regexp"

	"github.com/keptn/keptn/shipyard-controller/common"
)

type remoteURLValidator struct {
	denyListProvider common.FileReader
}

type RemoteURLValidator interface {
	Validate(url string) error
}

func NewRemoteURLValidator(denyListProvider common.FileReader) RemoteURLValidator {
	validator := remoteURLValidator{
		denyListProvider: denyListProvider,
	}
	return validator
}

func (c remoteURLValidator) Validate(url string) error {
	denyList := c.denyListProvider.Get(common.RemoteURLDenyListPath)

	for _, item := range denyList {
		if res, _ := regexp.MatchString(item, url); res {
			return fmt.Errorf("invalid RemoteURL: %s", url)
		}
	}
	return nil
}
