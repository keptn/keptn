package sequencehooks_test

import (
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	"testing"
)

type SequenceStateMVTestFields struct {
	SequenceStateRepo *db_mock.StateRepoMock
}

func TestSequenceStateMaterializedView_OnSequenceFinished(t *testing.T) {

	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			if err := smv.OnSequenceFinished(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("OnSequenceFinished() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTaskFinished(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			if err := smv.OnSequenceTaskFinished(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("OnSequenceTaskFinished() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTaskStarted(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			if err := smv.OnSequenceTaskStarted(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("OnSequenceTaskStarted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTaskTriggered(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			if err := smv.OnSequenceTaskTriggered(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("OnSequenceTaskTriggered() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTriggered(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			if err := smv.OnSequenceTriggered(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("OnSequenceTriggered() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
