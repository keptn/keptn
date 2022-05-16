package fake

type IPResolverMock struct {
	ResolveIPAdressesFunc func(curlURL string) []string
}

func (r IPResolverMock) Resolve(curlURL string) []string {
	if r.ResolveIPAdressesFunc != nil {
		return r.ResolveIPAdressesFunc(curlURL)
	}
	panic("implement me")
}
