package lib

import (
	"fmt"
	"strings"
)

type CurlValidator struct {
	denyListProvider IDenyListProvider
	ipResolver       IPResolver
}

type ICurlValidator interface {
	Validate(request Request) error
}

func NewCurlValidator(denyListProvider IDenyListProvider, ipResolver IPResolver) CurlValidator {
	validator := CurlValidator{
		denyListProvider: denyListProvider,
		ipResolver:       ipResolver,
	}
	return validator
}

func (c CurlValidator) Validate(request Request) error {
	if request.URL == "" {
		return fmt.Errorf("Invalid curl URL: '%s'", request.URL)
	}

	denyList := c.denyListProvider.GetDenyList()
	ipAddresses := c.ipResolver.ResolveIPAdresses(request.URL)

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
