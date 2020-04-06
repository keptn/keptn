package utils

import (
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const eventbroker = "EVENTBROKER"

// CreateAndSendConfigurationChangedEvent creates ConfigurationChangeEvent and sends it
func CreateAndSendConfigurationChangedEvent(problem *keptnevents.ProblemEventData,
	keptnHandler *keptnevents.Keptn, configChangedEvent keptnevents.ConfigurationChangeEventData) error {

	source, _ := url.Parse("https://github.com/keptn/keptn/remediation-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return keptnHandler.SendCloudEvent(event)
}
