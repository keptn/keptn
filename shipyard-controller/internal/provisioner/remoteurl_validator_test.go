package provisioner

import (
	"github.com/keptn/keptn/shipyard-controller/internal/filereader"
	common_mock "github.com/keptn/keptn/shipyard-controller/internal/filereader/fake"
	"testing"
)

func Test_RemoteURLValidator(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		fileReader filereader.FileReader
		expectErr  bool
	}{
		{
			name: "invalid url",
			url:  "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"some", "list"}
				},
			},
			expectErr: true,
		},
		{

			name: "valid url",
			url:  "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{}
				},
			},
			expectErr: false,
		},
		{

			name: "valid url",
			url:  "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"something"}
				},
			},
			expectErr: false,
		},
		{

			name: "invalid url regex match",
			url:  "something",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"some"}
				},
			},
			expectErr: true,
		},
		{
			name: "invalid url regex match",
			url:  "something",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"."}
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			validator := NewRemoteURLValidator(tt.fileReader)
			res := validator.Validate(tt.url)
			if (res != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, wantErr %v", res, tt.expectErr)
			}
		})
	}
}
