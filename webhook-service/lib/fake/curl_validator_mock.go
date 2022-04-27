package fake

import "github.com/keptn/keptn/webhook-service/lib"

type ICurlValidatorMock struct {
	ValidateFunc func(request lib.Request) error
}

func (r ICurlValidatorMock) Validate(request lib.Request) error {
	if r.ValidateFunc != nil {
		return r.ValidateFunc(request)
	}
	panic("implement me")
}
