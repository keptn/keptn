package fake

import "github.com/keptn/keptn/webhook-service/lib"

type IPResolverMock struct {
	ResolveIPAdressesFunc func(curlURL string) (lib.AdrDomainNameMapping, error)
}

func (r IPResolverMock) Resolve(curlURL string) (lib.AdrDomainNameMapping, error) {
	if r.ResolveIPAdressesFunc != nil {
		return r.ResolveIPAdressesFunc(curlURL)
	}
	panic("implement me")
}
