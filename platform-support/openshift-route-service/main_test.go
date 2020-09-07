package main

import (
	"context"
	"reflect"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/keptn/go-utils/pkg/lib"
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

func Test_gotEvent(t *testing.T) {

	testEvent := cloudevents.NewEvent()
	testEvent.SetID("123")
	testEvent.SetType("some-event")

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
				ctx:   nil,
				event: testEvent,
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
