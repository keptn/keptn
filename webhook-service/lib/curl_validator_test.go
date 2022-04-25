package lib_test

import (
	"fmt"
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/require"
)

func TestCurlValidator_Validate(t *testing.T) {
	curlValidator := &lib.CurlValidator{}
	tests := []struct {
		name              string
		data              lib.Request
		curlValidatorMock lib.ICurlValidator
		want              error
		wantErr           bool
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
			curlValidatorMock: lib.NewCurlValidator(curlValidator, []string{"some-invalid-url"}),
			want:              nil,
			wantErr:           false,
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
			curlValidatorMock: lib.NewCurlValidator(curlValidator, []string{"some-invalid-url"}),
			want:              fmt.Errorf("Invalid curl URL: ''"),
			wantErr:           true,
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
			curlValidatorMock: &fake.ICurlValidatorMock{
				ValidateFunc: func(request lib.Request, denyList []string, ipAddresses []string) error {
					return curlValidator.Validate(request, denyList, ipAddresses)
				},
				GetConfigDenyListFunc: func() []string {
					return []string{"some-denied"}
				},
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
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
			curlValidatorMock: &fake.ICurlValidatorMock{
				ValidateFunc: func(request lib.Request, denyList []string, ipAddresses []string) error {
					return curlValidator.Validate(request, denyList, ipAddresses)
				},
				GetConfigDenyListFunc: func() []string {
					return []string{"1.1.1.1"}
				},
				ResolveIPAdressesFunc: func(curlURL string) []string {
					return []string{"1.1.1.1"}
				},
			},
			want:    fmt.Errorf("curl command contains denied IP address '1.1.1.1'"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.curlValidatorMock.Validate(tt.data, tt.curlValidatorMock.GetConfigDenyList(), tt.curlValidatorMock.ResolveIPAdresses(tt.data.URL))
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, err)
		})
	}
}
