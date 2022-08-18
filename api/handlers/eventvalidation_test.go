package handlers

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/keptn/keptn/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	type args struct {
		e models.KeptnContextExtendedCE
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "invalid event type",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("garbage"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{Msg: "unknown event type"})
			},
		},
		{
			name: "sequence .triggered event - valid",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.st.sequence.triggered"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "sequence event with action other than .triggered is not allowed",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - common event data missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - project missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"service": "svc", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - service missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - stage missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - event data invalid",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        -1,
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "sequence .triggered event - stage mismatch between event data info and event type",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc", "stage": "anotherStage"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.stage.sequence.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - valid",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
					Triggeredid: "triggeredid",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "task .started event - project missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"service": "svc", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - service missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - stage missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - project, stage and service missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - triggeredid missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .started event - common data un-parsable",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        -1,
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.started"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - valid",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc", "result": "succeeded", "status": "pass"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
					Triggeredid: "triggeredid",
				},
			},
			wantErr: assert.NoError,
		},

		{
			name: "task .finished event - project missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"service": "svc", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - service missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - stage missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - project, stage and service missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - triggeredid missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc", "result": "succeeded", "status": "pass"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - common data un-parsable",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        -1,
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .triggered event - not allowed",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc", "result": "succeeded", "status": "pass"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.triggered"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - data result field missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc", "status": "pass"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task .finished event - data status field missing",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "stage": "st", "service": "svc", "result": "succeeded"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.task.finished"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "monitoring configure - valid",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc", "type": "prometheus"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.monitoring.configure"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "monitoring configure - missing type",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"project": "pr", "service": "svc"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.monitoring.configure"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "monitoring configure - missing project",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{"service": "svc", "type": "prometheus"},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Type:        stringp("sh.keptn.event.monitoring.configure"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &EventValidationError{})
			},
		},
		{
			name: "task invalidated event -generally allowed",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Time:        strfmt.DateTime{},
					Type:        stringp("sh.keptn.event.task.invalidated"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "problem event - generally allowed",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Time:        strfmt.DateTime{},
					Type:        stringp("sh.keptn.events.problem"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error event - generally allowed",
			args: args{
				e: models.KeptnContextExtendedCE{
					Contenttype: "application/json",
					Data:        map[string]interface{}{},
					Source:      stringp("test-source"),
					Specversion: "1.0",
					Time:        strfmt.DateTime{},
					Type:        stringp("sh.keptn.log.error"),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, Validate(tt.args.e), fmt.Sprintf("Validate(%v)", tt.args.e))
		})
	}
}
