package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsSLORequired(t *testing.T) {
	logger := keptnutils.NewLogger("1234-9876-abcd-zyxw", "1234-9876-abcd-zyxw", "lighthouse-service")
 	event := cloudevents.New()

	handler:= &StartEvaluationHandler{Logger: logger, Event: event}
	assert.EqualValues(t, handler.isSLORequired(), true)

	os.Setenv("SLO_REQUIRED", "false")
	assert.EqualValues(t, handler.isSLORequired(), false)

	os.Setenv("SLO_REQUIRED", "true")
	assert.EqualValues(t, handler.isSLORequired(), true)

	os.Setenv("SLO_REQUIRED", "yes")
	assert.EqualValues(t, handler.isSLORequired(), true)
}
