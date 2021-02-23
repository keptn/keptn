package common

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
	"os"
)

const defaultKeptnNamespace = "keptn"

const keptnSpecVersionEnvVar = "KEPTN_SPEC_VERSION"

// GetKeptnNamespace godoc
func GetKeptnNamespace() string {
	ns := os.Getenv("POD_NAMESPACE")
	if ns != "" {
		return ns
	}
	return defaultKeptnNamespace
}

// GetShipyard godoc
func GetShipyard(eventScope *keptnv2.EventData) (*keptnv2.Shipyard, error) {
	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		return nil, errors.New("Could not get configuration-service URL: " + err.Error())
	}
	resourceHandler := keptnapi.NewResourceHandler(csEndpoint.String())
	resource, err := resourceHandler.GetProjectResource(eventScope.Project, "shipyard.yaml")
	if err != nil {
		return nil, errors.New("Could not retrieve shipyard.yaml for project " + eventScope.Project + ": " + err.Error())
	}

	shipyard, err := UnmarshalShipyard(resource.ResourceContent)
	if err != nil {
		return nil, err
	}

	return shipyard, err
}

// UnmarshalShipyard godoc
func UnmarshalShipyard(shipyardString string) (*keptnv2.Shipyard, error) {
	shipyard := &keptnv2.Shipyard{}
	err := yaml.Unmarshal([]byte(shipyardString), shipyard)
	if err != nil {
		return nil, errors.New("Could not decode shipyard file: " + err.Error())
	}
	return shipyard, nil
}

// ValidateShipyardVersion godoc
func ValidateShipyardVersion(shipyard *keptnv2.Shipyard) error {
	if shipyard.ApiVersion != "0.2.0" && shipyard.ApiVersion != "spec.keptn.sh/0.2.0" {
		return errors.New("Invalid shipyard APIVersion " + shipyard.ApiVersion)
	}
	return nil
}

// ValidateShipyardStages godoc
func ValidateShipyardStages(shipyard *keptnv2.Shipyard) error {
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name == "" {
			return errors.New("all stages within the shipyard must have a name")
		}
		if !keptncommon.ValidateKeptnEntityName(stage.Name) {
			errorMsg := "Stage " + stage.Name + " contains upper case letter(s) or special character(s).\n"
			errorMsg += "Keptn relies on the following conventions: "
			errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
			errorMsg += "Please update stage name in your shipyard and try again."
			return errors.New(errorMsg)
		}
	}
	return nil
}

// GetKeptnSpecVersion returns the Keptn Spec version the shipyard controller is based on
func GetKeptnSpecVersion() string {
	return os.Getenv(keptnSpecVersionEnvVar)
}

// SendEvent godoc
func SendEvent(event cloudevents.Event) error {
	ebEndpoint, err := keptncommon.GetServiceEndpoint("EVENTBROKER")
	if err != nil {
		return errors.New("Could not get eventbroker endpoint: " + err.Error())
	}
	k, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		EventBrokerURL: ebEndpoint.String(),
	})
	if err != nil {
		return errors.New("Could not initialize Keptn handler: " + err.Error())
	}
	if specVersion := GetKeptnSpecVersion(); specVersion != "" {
		event.SetExtension("shkeptnspecversion", specVersion)
	}
	err = k.SendCloudEvent(event)
	if err != nil {
		return errors.New("Could not send CloudEvent: " + err.Error())
	}
	return nil
}

// SendEventWithPayload godoc
// Deprecated will be removed, use functionality from go-utils instead
func SendEventWithPayload(keptnContext, triggeredID, eventType string, payload interface{}) error {
	event := CreateEventWithPayload(keptnContext, triggeredID, eventType, payload)

	ebEndpoint, err := keptncommon.GetServiceEndpoint("EVENTBROKER")
	if err != nil {
		return errors.New("Could not get eventbroker endpoint: " + err.Error())
	}
	k, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		EventBrokerURL: ebEndpoint.String(),
	})
	if err != nil {
		return errors.New("Could not initialize Keptn handler: " + err.Error())
	}

	err = k.SendCloudEvent(event)
	if err != nil {
		return errors.New("Could not send CloudEvent: " + err.Error())
	}
	return nil
}

func CreateEventWithPayload(keptnContext, triggeredID, eventType string, payload interface{}) cloudevents.Event {
	source, _ := url.Parse("shipyard-controller")
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	if keptnContext == "" {
		event.SetExtension("shkeptncontext", uuid.New().String())
	} else {
		event.SetExtension("shkeptncontext", keptnContext)
	}
	if triggeredID != "" {
		event.SetExtension("triggeredid", triggeredID)
	}
	if specVersion := GetKeptnSpecVersion(); specVersion != "" {
		event.SetExtension("shkeptnspecversion", specVersion)
	}
	event.SetData(cloudevents.ApplicationJSON, payload)
	return event
}
