package main

import (
	"context"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/keptn/go-utils/pkg/lib"
	"os"
	"reflect"
	"testing"
)

func Test_getEnableMeshCommandArgs(t *testing.T) {
	type args struct {
		project string
		stage   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Enable mesh command",
			args: args{
				project: "sockshop",
				stage:   "dev",
			},
			want: []string{
				"adm",
				"policy",
				"add-scc-to-group",
				"anyuid",
				"system:serviceaccounts",
				"-n",
				"sockshop-dev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnableMeshCommandArgs(tt.args.project, tt.args.stage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnableMeshCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCreateRouteCommandArgs(t *testing.T) {
	type args struct {
		project   string
		stage     string
		appDomain string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Create route command",
			args: args{
				project:   "sockshop",
				stage:     "dev",
				appDomain: "keptn.com",
			},
			want: []string{
				"create",
				"route",
				"edge",
				"sockshop-dev",
				"--service=istio-ingressgateway",
				"--hostname=www.sockshop-dev.keptn.com",
				"--port=http2",
				"--wildcard-policy=Subdomain",
				"--insecure-policy=Allow",
				"-n",
				"istio-system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCreateRouteCommandArgs(tt.args.project, tt.args.stage, tt.args.appDomain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCreateRouteCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gotEvent(t *testing.T) {
	type args struct {
		ctx   context.Context
		event cloudevents.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid event",
			args: args{
				ctx: nil,
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						SpecVersion: "0.2",
						Type:        "some-event",
						Source:      types.URLRef{},
						ID:          "123",
						Time:        nil,
						SchemaURL:   nil,
						ContentType: nil,
						Extensions:  nil,
					},
					Data:        "",
					DataEncoded: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := gotEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("gotEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createRoutes(t *testing.T) {
	type args struct {
		data *keptn.ProjectCreateEventData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "invalid shipyard should throw error",
			args: args{
				data: &keptn.ProjectCreateEventData{
					Project:  "sockshop",
					Shipyard: "this is not base 64 encoded",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_createRoutes1(t *testing.T) {
	type args struct {
		data *keptn.ProjectCreateEventData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid shipyard should throw error",
			args: args{
				data: &keptn.ProjectCreateEventData{
					Project:  "sockshop",
					Shipyard: "this is not base 64 encoded",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid shipyard should throw error",
			args: args{
				data: &keptn.ProjectCreateEventData{
					Project:  "sockshop",
					Shipyard: "aW52YWxpZAo=",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createRoutes(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("createRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_exposeRoute(t *testing.T) {
	type args struct {
		project string
		stage   string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		appDomainSet bool
	}{
		{
			name: "don't create anything if APP_DOMAIN env var is not set",
			args: args{
				project: "sockshop",
				stage:   "dev",
			},
			wantErr:      true,
			appDomainSet: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.appDomainSet {
				_ = os.Setenv("APP_DOMAIN", "")
			}
			if err := exposeRoute(tt.args.project, tt.args.stage); (err != nil) != tt.wantErr {
				t.Errorf("exposeRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__main(t *testing.T) {
	env := envConfig{
		Port: 8080,
		Path: "/",
	}
	type args struct {
		args []string
		env  envConfig
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Check if service starts correctly",
			args: args{
				args: nil,
				env:  env,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				if got := _main(tt.args.args, tt.args.env); got != tt.want {
					t.Errorf("_main() = %v, want %v", got, tt.want)
				}
			}()
			ctx.Done()
		})
	}
}
