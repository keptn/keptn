package lib

import (
	"fmt"
	"strings"
)

type requestValidator struct {
	denyListProvider DenyListProvider
	ipResolver       IPResolver
}

type RequestValidator interface {
	Validate(request Request) error
}

func NewRequestValidator(denyListProvider DenyListProvider, ipResolver IPResolver) RequestValidator {
	validator := requestValidator{
		denyListProvider: denyListProvider,
		ipResolver:       ipResolver,
	}
	return validator
}

func (c requestValidator) Validate(request Request) error {
	if request.URL == "" {
		return fmt.Errorf("invalid curl URL: '%s'", request.URL)
	}

	denyList := c.denyListProvider.Get()
	ipAddresses := c.ipResolver.Resolve(request.URL)

	for _, url := range denyList {
		if strings.Contains(request.URL, url) {
			return fmt.Errorf("curl command contains denied URL '%s'", url)
		}
		for _, ip := range ipAddresses {
			if strings.Contains(ip, url) {
				return fmt.Errorf("curl command contains denied IP address '%s'", url)
			}
		}
	}
	return nil
}
