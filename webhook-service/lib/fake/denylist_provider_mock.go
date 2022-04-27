package fake

type IDenyListProviderMock struct {
	GetDenyListFunc func() []string
}

func (r IDenyListProviderMock) GetDenyList() []string {
	if r.GetDenyListFunc != nil {
		return r.GetDenyListFunc()
	}
	panic("implement me")
}
