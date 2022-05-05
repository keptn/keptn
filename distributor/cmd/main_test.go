package main

import (
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	dapi "github.com/keptn/keptn/distributor/pkg/api"
	"github.com/keptn/keptn/distributor/pkg/config"
	"net/http"
	"reflect"
	"testing"
)

func Test_createAPI(t *testing.T) {

	api := APIInitializer{
		internal: func(client *http.Client, apiMappings ...dapi.InClusterAPIMappings) (*dapi.InternalAPISet, error) {
			return &dapi.InternalAPISet{}, nil
		},
		external: func(baseURL string, options ...func(*keptnapi.APISet)) (*keptnapi.APISet, error) {
			return &keptnapi.APISet{}, nil
		},
	}

	tests := []struct {
		name         string
		env          config.EnvConfig
		wantInternal bool
		wantErr      bool
	}{
		{
			name:         "test no env internal NATS ",
			env:          config.EnvConfig{},
			wantInternal: true,
			wantErr:      false,
		},
		{
			name: "test FAIL for no http address",
			env: config.EnvConfig{
				KeptnAPIEndpoint: "ssh://mynotsogoodendpoint",
			},
			wantErr:      true,
			wantInternal: false,
		},
		{
			name: "test FAIL for no good address",
			env: config.EnvConfig{
				KeptnAPIEndpoint: ":///MALFORMEDendpoint",
			},
			wantErr:      true,
			wantInternal: false,
		},
		{
			name: "test PASS for http address",
			env: config.EnvConfig{
				KeptnAPIEndpoint: "http://endpoint",
			},
			wantErr:      false,
			wantInternal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAPI(nil, tt.env, api)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantInternal && !reflect.DeepEqual(got, &dapi.InternalAPISet{}) {
				t.Errorf("createAPI() got = %v, wanted internal API", got)
			} else if err == nil && !tt.wantInternal && !reflect.DeepEqual(got, &keptnapi.APISet{}) {
				t.Errorf("createAPI() got = %v, want remote execution plane", got)
			}

		})
	}
}
