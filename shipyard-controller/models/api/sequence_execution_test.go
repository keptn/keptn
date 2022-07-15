package api

import (
	"github.com/keptn/keptn/shipyard-controller/models"
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
		// TODO: Add test cases.
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
