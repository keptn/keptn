package types

import (
	"github.com/keptn/go-utils/pkg/api/models"
)

type RegistrationData models.Integration

// AdditionalSubscriptionData is the data the cp-connector
// will add as temporary data to the keptn events forwarded
// to the keptn integration
type AdditionalSubscriptionData struct {
	SubscriptionID string `json:"subscriptionID"`
}

// EventUpdate wraps a new Keptn event received from the Event source
type EventUpdate struct {
	KeptnEvent models.KeptnContextExtendedCE
	MetaData   EventUpdateMetaData
}

type EventSenderKeyType struct{}

var EventSenderKey = EventSenderKeyType{}

type EventSender func(ce models.KeptnContextExtendedCE) error

// EventUpdateMetaData is additional metadata for bound to the
// event received from the event source
type EventUpdateMetaData struct {
	Subject string
}
