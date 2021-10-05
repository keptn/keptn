package lib_test

import (
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/lib"
	"testing"
)

func TestEventDataAdapter_SubscriptionID(t *testing.T) {
	type fields struct {
		event sdk.KeptnEvent
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "get subscription id",
			fields: fields{
				event: sdk.KeptnEvent{
					Data: map[string]interface{}{
						"project": "my-project",
						"stage":   "my-stage",
						"service": "my-service",
						"temporaryData": map[string]interface{}{
							"distributor": map[string]string{
								"subscriptionID": "sub-id",
							},
						},
					},
				},
			},
			want:    "sub-id",
			wantErr: false,
		},
		{
			name: "empty subscription id",
			fields: fields{
				event: sdk.KeptnEvent{
					Data: map[string]interface{}{
						"project": "my-project",
						"stage":   "my-stage",
						"service": "my-service",
						"temporaryData": map[string]interface{}{
							"distributor": map[string]string{
								"subscriptionID": "",
							},
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no subscription id",
			fields: fields{
				event: sdk.KeptnEvent{
					Data: map[string]interface{}{
						"project": "my-project",
						"stage":   "my-stage",
						"service": "my-service",
						"temporaryData": map[string]interface{}{
							"distributor": map[string]string{},
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no distributor info",
			fields: fields{
				event: sdk.KeptnEvent{
					Data: map[string]interface{}{
						"project":       "my-project",
						"stage":         "my-stage",
						"service":       "my-service",
						"temporaryData": map[string]interface{}{},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := lib.NewEventDataAdapter(tt.fields.event)
			got, err := e.SubscriptionID()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscriptionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SubscriptionID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
