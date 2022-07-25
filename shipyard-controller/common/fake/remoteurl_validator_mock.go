package common_mock

type RequestValidatorMock struct {
	ValidateFunc func(url string) error
}

func (r RequestValidatorMock) Validate(url string) error {
	if r.ValidateFunc != nil {
		return r.ValidateFunc(url)
	}
	panic("implement me")
}
