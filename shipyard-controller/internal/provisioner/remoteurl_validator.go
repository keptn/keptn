package provisioner

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/filereader"
	"regexp"
)

type remoteURLValidator struct {
	denyListProvider filereader.FileReader
}

type RemoteURLValidator interface {
	Validate(url string) error
}

func NewRemoteURLValidator(denyListProvider filereader.FileReader) RemoteURLValidator {
	validator := remoteURLValidator{
		denyListProvider: denyListProvider,
	}
	return validator
}

func (c remoteURLValidator) Validate(url string) error {
	denyList := c.denyListProvider.GetLines(filereader.RemoteURLDenyListPath)

	for _, item := range denyList {
		if res, _ := regexp.MatchString(item, url); res {
			return fmt.Errorf("invalid RemoteURL: %s", url)
		}
	}
	return nil
}
