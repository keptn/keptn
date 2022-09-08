package common

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"reflect"
	"testing"
)

func TestEncodeKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode dots",
			args: args{
				key: "sh.keptn.event",
			},
			want: "sh~pkeptn~pevent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeKey(tt.args.key); got != tt.want {
				t.Errorf("EncodeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode dots",
			args: args{
				key: "sh~pkeptn~pevent",
			},
			want: "sh.keptn.event",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeKey(tt.args.key); got != tt.want {
				t.Errorf("DecodeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInterface(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "convert event scope",
			args: args{
				item: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
					},
					KeptnContext: "my-context",
					GitCommitID:  "my-commit-id",
					EventType:    "my-event-type",
					TriggeredID:  "my-triggered-id",
				},
			},
			want: map[string]interface{}{
				"project":      "my-project",
				"keptnContext": "my-context",
				"triggeredId":  "my-triggered-id",
				"gitcommitid":  "my-commit-id",
				"eventType":    "my-event-type",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInterface(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInterface() got = %v, want %v", got, tt.want)
			}
		})
	}
}
