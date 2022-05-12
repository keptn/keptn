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
		return fmt.Errorf("curl command contains empty URL")
	}
	denyList := c.denyListProvider.Get()
	ipAddresses := c.ipResolver.Resolve(request.URL)

	for _, url := range denyList {
		if strings.Contains(request.URL, url) {
			return fmt.Errorf("curl command contains denied URL '%s'", url)
		}
		err := validateIPDomain(ipAddresses, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateIPDomain(ipAddresses AdrDomainNameMapping, url string) error {
	for ip, hosts := range ipAddresses {
		if strings.Contains(ip, url) {
			return fmt.Errorf("curl command contains denied IP address '%s'", url)
		}
		for _, h := range hosts {
			if strings.Contains(trimEndDot(h), url) {
				return fmt.Errorf("curl command url resolves to denied host '%s'", url)
			}
		}
	}
	return nil
}

func trimEndDot(h string) string {
	lastIdx := len(h) - 1
	if h[lastIdx] == '.' {
		h = h[:lastIdx]
	}
	return h
}
