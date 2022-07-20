package api

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestGetSequenceExecutionParams_GetSequenceExecutionFilter(t *testing.T) {
	type fields struct {
		Project          string
		Service          string
		Stage            string
		KeptnContext     string
		Name             string
		Status           string
		PaginationParams models.PaginationParams
	}
	tests := []struct {
		name   string
		fields fields
		want   models.SequenceExecutionFilter
	}{
		{
			name: "with added state",
			fields: fields{
				Project:      "my-project",
				Service:      "my-service",
				Stage:        "my-stage",
				KeptnContext: "my-context",
				Name:         "my-sequence",
				Status:       "started",
			},
			want: models.SequenceExecutionFilter{
				Scope: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Service: "my-service",
						Stage:   "my-stage",
					},
					KeptnContext: "my-context",
				},
				Status: []string{"started"},
				Name:   "my-sequence",
			},
		},
		{
			name: "without state",
			fields: fields{
				Project:      "my-project",
				Service:      "my-service",
				Stage:        "my-stage",
				KeptnContext: "my-context",
				Name:         "my-sequence",
				Status:       "",
			},
			want: models.SequenceExecutionFilter{
				Scope: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Service: "my-service",
						Stage:   "my-stage",
					},
					KeptnContext: "my-context",
				},
				Status: nil,
				Name:   "my-sequence",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GetSequenceExecutionParams{
				Project:          tt.fields.Project,
				Service:          tt.fields.Service,
				Stage:            tt.fields.Stage,
				KeptnContext:     tt.fields.KeptnContext,
				Name:             tt.fields.Name,
				Status:           tt.fields.Status,
				PaginationParams: tt.fields.PaginationParams,
			}
			if got := p.GetSequenceExecutionFilter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSequenceExecutionFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSequenceExecutionParams_Validate(t *testing.T) {
	type fields struct {
		Project          string
		Service          string
		Stage            string
		KeptnContext     string
		Name             string
		Status           string
		PaginationParams models.PaginationParams
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "Project name not set",
			fields:  fields{},
			wantErr: ErrProjectNameMustNotBeEmpty,
		},
		{
			name: "Project name set",
			fields: fields{
				Project: "my-project",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GetSequenceExecutionParams{
				Project:          tt.fields.Project,
				Service:          tt.fields.Service,
				Stage:            tt.fields.Stage,
				KeptnContext:     tt.fields.KeptnContext,
				Name:             tt.fields.Name,
				Status:           tt.fields.Status,
				PaginationParams: tt.fields.PaginationParams,
			}
			err := p.Validate()

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
