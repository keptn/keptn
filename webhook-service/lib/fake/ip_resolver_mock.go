package fake

type IIPResolverMock struct {
	ResolveIPAdressesFunc func(curlURL string) []string
}

func (r IIPResolverMock) ResolveIPAdresses(curlURL string) []string {
	if r.ResolveIPAdressesFunc != nil {
		return r.ResolveIPAdressesFunc(curlURL)
	}
	panic("implement me")
}
