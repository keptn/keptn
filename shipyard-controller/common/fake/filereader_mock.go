package common_mock

type FileReaderMock struct {
	GetFunc func(path string) []string
}

func (r FileReaderMock) Get(path string) []string {
	if r.GetFunc != nil {
		return r.GetFunc(path)
	}
	panic("implement me")
}
