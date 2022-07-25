package handler

import (
	"testing"

	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
)

func Test_RemoteURLValidator(t *testing.T) {
	tests := []struct {
		url              string
		denyListProvider common.DenyListProvider
		expectErr        bool
	}{
		{
			url: "some",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"some", "list"}
				},
			},
			expectErr: true,
		},
		{
			url: "some",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{}
				},
			},
			expectErr: false,
		},
		{
			url: "some",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"something"}
				},
			},
			expectErr: false,
		},
		{
			url: "something",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"some"}
				},
			},
			expectErr: true,
		},
		{
			url: "something",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{""}
				},
			},
			expectErr: true,
		},
		{
			url: "something",
			denyListProvider: common_mock.DenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"."}
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			validator := NewRemoteURLValidator(tt.denyListProvider)
			res := validator.Validate(tt.url)
			if (res != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, wantErr %v", res, tt.expectErr)
			}
		})
	}
}
