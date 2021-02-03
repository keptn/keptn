package utils

import (
	"fmt"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const eventBroker = "EVENTBROKER_URI"
const dataStore = "DATASTORE_URI"

// GetEventBrokerURL godoc
func GetEventBrokerURL() string {
	return SanitizeURL(os.Getenv(eventBroker))
}

// GetDatastoreURL godoc
func GetDatastoreURL() string {
	return SanitizeURL(os.Getenv(dataStore))
}

// SendEvent godoc
func SendEvent(keptnContext, triggeredID, eventType, source string, data interface{}, l keptnutils.LoggerInterface) error {
	if source == "" {
		source = "https://github.com/keptn/keptn/api"
	}
	ev := cloudevents.NewEvent()
	ev.SetType(eventType)
	ev.SetID(uuid.New().String())
	ev.SetTime(time.Now())
	ev.SetSource(source)
	ev.SetDataContentType(cloudevents.ApplicationJSON)
	ev.SetExtension("shkeptncontext", keptnContext)
	ev.SetExtension("triggeredid", triggeredID)
	ev.SetData(cloudevents.ApplicationJSON, data)

	k, err := keptnv2.NewKeptn(&ev, keptnutils.KeptnOpts{
		EventBrokerURL: GetEventBrokerURL(),
	})
	if err != nil {
		l.Error(fmt.Sprintf("Error initializing Keptn Handler %s", err.Error()))
		return err
	}
	err = k.SendCloudEvent(ev)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return err
	}
	return nil
}
