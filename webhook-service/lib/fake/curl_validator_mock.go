package fake

import "github.com/keptn/keptn/webhook-service/lib"

type ICurlValidatorMock struct {
	ValidateFunc          func(request lib.Request) error
	ResolveIPAdressesFunc func(curlURL string) []string
	GetConfigDenyListFunc func() ([]string, error)
}

func (r *ICurlValidatorMock) Validate(request lib.Request) error {
	if r.ValidateFunc != nil {
		return r.ValidateFunc(request)
	}
	panic("implement me")
}

func (r *ICurlValidatorMock) ResolveIPAdresses(curlURL string) []string {
	if r.ResolveIPAdressesFunc != nil {
		return r.ResolveIPAdressesFunc(curlURL)
	}
	panic("implement me")
}

func (r *ICurlValidatorMock) GetConfigDenyList() ([]string, error) {
	if r.GetConfigDenyListFunc != nil {
		return r.GetConfigDenyListFunc()
	}
	panic("implement me")
}
