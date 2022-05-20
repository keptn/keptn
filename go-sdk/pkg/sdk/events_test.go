package sdk

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_createFinishedEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing event type",
			args: args{
				source:      "source",
				parentEvent: models.KeptnContextExtendedCE{},
				eventData:   nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing event context",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.evaluation.triggered"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "event type cannot be replaced",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("something.weird"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "passed event data missing",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           v0_2_0.EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "defalut status succeeded",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           v0_2_0.EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: v0_2_0.EventData{},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["status"] == string(v0_2_0.StatusSucceeded) && ce.Data.(map[string]interface{})["result"] == string(v0_2_0.ResultPass)
			},
			wantErr: false,
		},
		{
			name: "defalut result pass",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           v0_2_0.EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: v0_2_0.EventData{},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["result"] == "pass"
			},
			wantErr: false,
		},
		{
			name: "passed status reused",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           v0_2_0.EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: v0_2_0.EventData{Status: v0_2_0.StatusErrored, Result: v0_2_0.ResultFailed},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				fmt.Println(ce.Data.(map[string]interface{})["status"])
				return ce.Data.(map[string]interface{})["status"] == string(v0_2_0.StatusErrored) && ce.Data.(map[string]interface{})["result"] == string(v0_2_0.ResultFailed)
			},
			wantErr: false,
		},
		{
			name: "correct event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           v0_2_0.EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: v0_2_0.EventData{Status: v0_2_0.StatusSucceeded, Result: v0_2_0.ResultPass},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				fmt.Println(ce.Data.(map[string]interface{})["status"])
				return *ce.Type == "sh.keptn.event.eval.finished"
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createFinishedEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_createStartedEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing keptn context",
			args: args{
				source:      "source",
				parentEvent: models.KeptnContextExtendedCE{},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "non-replacable event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
					Type:           strutils.Stringp("somethin.weird"),
				},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "ok",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == "sh.keptn.event.eval.started" },
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createStartedEvent(tt.args.source, tt.args.parentEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("createStartedEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_createErrorLogEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
		errVal      *Error
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing event type",
			args: args{
				source:      "soure",
				parentEvent: models.KeptnContextExtendedCE{},
				eventData:   nil,
				errVal:      nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing keptn context",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "creates finished event for events of type .triggered",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == "sh.keptn.event.eval.finished" },
			wantErr:     false,
		},
		{
			name: "creates error event for events other than .triggered",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == v0_2_0.ErrorLogEventName },
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createErrorLogEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData, tt.args.errVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("createErrorLogEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_createErrorEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
		err         *Error
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing keptn context",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "creates error event for events other than .triggered",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == v0_2_0.ErrorLogEventName },
			wantErr:     false,
		},
		{
			name: "creates finished event for events of type .triggered",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == "sh.keptn.event.eval.finished" },
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createErrorEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("createErrorEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}
