package mv

import (
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	dbcommon "github.com/keptn/keptn/shipyard-controller/internal/db/common"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrEmptyUpdate = errors.New("update object does not contain any changes")

type EventUpdate struct {
	EventType string
	EventInfo models.EventContextInfo
}

// ServiceUpdate defines a set of properties that can be updated for a specific service within a project, without needing to overwrite the complete project
// this will help to reduce concurrency issues due to simultaneous updates of a service within a project
type ServiceUpdate struct {
	deployedImage   string
	eventTypeUpdate *EventUpdate
}

func (s *ServiceUpdate) DeployedImage() string {
	return s.deployedImage
}

func (s *ServiceUpdate) EventTypeUpdate() *EventUpdate {
	return s.eventTypeUpdate
}

func (s *ServiceUpdate) SetDeployedImage(deployedImage string) {
	s.deployedImage = deployedImage
}

func (s *ServiceUpdate) SetEventTypeUpdate(update *EventUpdate) {
	s.eventTypeUpdate = update
}

func (s *ServiceUpdate) GetBSONUpdate() (bson.D, error) {

	if s.deployedImage == "" && s.eventTypeUpdate == nil {
		return nil, ErrEmptyUpdate
	}
	changeSet := bson.M{}

	if s.deployedImage != "" {
		changeSet["stages.$.services.$[service].deployedImage"] = s.DeployedImage()
	}

	if s.eventTypeUpdate != nil {
		encodedEventInfo, err := dbcommon.ToInterface(s.eventTypeUpdate.EventInfo)
		if err != nil {
			return nil, err
		}
		encodedEventType := dbcommon.EncodeKey(s.eventTypeUpdate.EventType)
		changeSet["stages.$.services.$[service].lastEventTypes."+encodedEventType] = encodedEventInfo
	}

	update := bson.D{
		{"$set", changeSet},
	}

	return update, nil
}
