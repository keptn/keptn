package utils

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
	"os"
	"time"
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
func SendEvent(keptnContext, triggeredID, eventType string, data interface{}, l keptnutils.LoggerInterface) error {
	source, _ := url.Parse("https://github.com/keptn/keptn/api")
	ev := cloudevents.NewEvent()
	ev.SetType(keptnevents.InternalProjectCreateEventType)
	ev.SetID(uuid.New().String())
	ev.SetTime(time.Now())
	ev.SetSource(source.String())
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
