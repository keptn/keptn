package fake

type FileReaderMock struct {
	GetLinesFunc func(path string) []string
}

func (r FileReaderMock) GetLines(path string) []string {
	if r.GetLinesFunc != nil {
		return r.GetLinesFunc(path)
	}
	panic("implement me")
}
