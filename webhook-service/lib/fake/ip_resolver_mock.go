package fake

type IPResolverMock struct {
	ResolveIPAdressesFunc func(curlURL string) []string
}

func (r IPResolverMock) ResolveIPAdresses(curlURL string) []string {
	if r.ResolveIPAdressesFunc != nil {
		return r.ResolveIPAdressesFunc(curlURL)
	}
	panic("implement me")
}
