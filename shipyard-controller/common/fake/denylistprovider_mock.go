package common_mock

type DenyListProviderMock struct {
	GetDenyListFunc func() []string
}

func (r DenyListProviderMock) Get() []string {
	if r.GetDenyListFunc != nil {
		return r.GetDenyListFunc()
	}
	panic("implement me")
}
