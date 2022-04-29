package lib_test

import (
	"fmt"
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/require"
)

func TestRequestValidator_Validate(t *testing.T) {
	tests := []struct {
		name             string
		data             lib.Request
		want             error
		ipResolver       lib.IPResolver
		denyListProvider lib.IDenyListProvider
		wantErr          bool
	}{
		{
			name: "valid input",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method:  "POST",
				Options: "--some-options",
				Payload: "some payload",
				URL:     "http://some-valid-url",
			},
			ipResolver: fake.IPResolverMock{
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
				},
			},
			denyListProvider: fake.IDenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"1.1.1.2"}
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid input",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method:  "POST",
				Options: "--some-options",
				Payload: "some payload",
				URL:     "",
			},
			ipResolver: fake.IPResolverMock{
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
				},
			},
			denyListProvider: fake.IDenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"1.1.1.1"}
				},
			},
			want:    fmt.Errorf("Invalid curl URL: ''"),
			wantErr: true,
		},
		{
			name: "denied URL input",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method:  "POST",
				Options: "--some-options",
				Payload: "some payload",
				URL:     "http://some-denied-url",
			},
			ipResolver: fake.IPResolverMock{
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
				},
			},
			denyListProvider: fake.IDenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"some-denied"}
				},
			},
			want:    fmt.Errorf("curl command contains denied URL 'some-denied'"),
			wantErr: true,
		},
		{
			name: "denied IP input",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method:  "POST",
				Options: "--some-options",
				Payload: "some payload",
				URL:     "http://som-url",
			},
			ipResolver: fake.IPResolverMock{
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
				},
			},
			denyListProvider: fake.IDenyListProviderMock{
				GetDenyListFunc: func() []string {
					return []string{"1.1.1.1"}
				},
			},
			want:    fmt.Errorf("curl command contains denied IP address '1.1.1.1'"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestValidator := lib.NewRequestValidator(tt.denyListProvider, tt.ipResolver)
			err := requestValidator.Validate(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, err)
		})
	}
}
