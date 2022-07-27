package handler

import (
	"testing"

	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
)

func Test_RemoteURLValidator(t *testing.T) {
	tests := []struct {
		url        string
		fileReader common.FileReader
		expectErr  bool
	}{
		{
			url: "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"some", "list"}
				},
			},
			expectErr: true,
		},
		{
			url: "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{}
				},
			},
			expectErr: false,
		},
		{
			url: "some",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"something"}
				},
			},
			expectErr: false,
		},
		{
			url: "something",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{"some"}
				},
			},
			expectErr: true,
		},
		{
			url: "something",
			fileReader: common_mock.FileReaderMock{
				GetLinesFunc: func(path string) []string {
					return []string{""}
				},
			},
			expectErr: true,
		},
		{
			url: "something",
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
