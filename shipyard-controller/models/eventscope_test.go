package models

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"reflect"
	"testing"
)

func Test_NewEventScope(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name    string
		args    args
		want    *EventScope
		wantErr bool
	}{
		{
			name: "get event scope",
			args: args{
				event: Event{
					Data: keptnv2.EventData{Project: "sockshop", Stage: "dev", Service: "carts"},
					Type: common.Stringp("my-type"),
				},
			},
			want: &EventScope{EventData: keptnv2.EventData{Project: "sockshop", Stage: "dev", Service: "carts"}, EventType: "my-type", WrappedEvent: Event{
				Data: keptnv2.EventData{Project: "sockshop", Stage: "dev", Service: "carts"},
				Type: common.Stringp("my-type"),
			}},
			wantErr: false,
		},
		{
			name: "only project available, stage and service missing",
			args: args{
				event: Event{
					Data: keptnv2.EventData{Project: "sockshop"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty data",
			args: args{
				event: Event{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nonsense data",
			args: args{
				event: Event{Data: "invalid"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEventScope(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEventScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEventScope() got = %v, want %v", got, tt.want)
			}
		})
	}
}
