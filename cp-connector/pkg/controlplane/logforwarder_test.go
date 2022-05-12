package controlplane

import (
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane/fake"
	"github.com/stretchr/testify/require"
)

func TestLogForwarderNoForward(t *testing.T) {
	logHandler := &fake.LogAPIMock{}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.Nil(t, err)
	require.Len(t, logHandler.LogCalls(), 0)
}

func TestLogForwarderFinishedNoForward(t *testing.T) {
	logHandler := &fake.LogAPIMock{}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.finished"), Data: keptnv2.EventData{Status: keptnv2.StatusSucceeded}}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.Nil(t, err)
	require.Len(t, logHandler.LogCalls(), 0)
}

func TestLogForwarderFinishedInvalidEventType(t *testing.T) {
	logHandler := &fake.LogAPIMock{}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.finished"), Data: "some invalid data"}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.NotNil(t, err)
	require.Len(t, logHandler.LogCalls(), 0)
}

func TestLogForwarderFinishedForward(t *testing.T) {
	logHandler := &fake.LogAPIMock{
		LogFunc:   func(logs []models.LogEntry) {},
		FlushFunc: func() error { return nil },
	}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.finished"), Data: keptnv2.EventData{Status: keptnv2.StatusErrored}}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.Nil(t, err)
	require.Len(t, logHandler.LogCalls(), 1)
}

func TestLogForwarderErrorInvalidEventType(t *testing.T) {
	logHandler := &fake.LogAPIMock{}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.log.error"), Data: "some invalid data"}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.NotNil(t, err)
	require.Len(t, logHandler.LogCalls(), 0)
}

func TestLogForwarderErrorForward(t *testing.T) {
	logHandler := &fake.LogAPIMock{
		LogFunc:   func(logs []models.LogEntry) {},
		FlushFunc: func() error { return nil },
	}
	logForwarder := NewLogForwarder(logHandler)
	keptnEvent := models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.log.error")}
	err := logForwarder.Forward(keptnEvent, "some-other-id")
	require.Nil(t, err)
	require.Len(t, logHandler.LogCalls(), 1)
}
