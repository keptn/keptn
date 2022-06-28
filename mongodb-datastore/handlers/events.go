package handlers

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
	"github.com/keptn/keptn/mongodb-datastore/db"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	log "github.com/sirupsen/logrus"
)

type ProjectEventData struct {
	Project *string `json:"project,omitempty"`
}

type EventRequestHandler struct {
	eventRepo db.EventRepo
	Env       envConfig
}

func (erh EventRequestHandler) OnEvent(ctx context.Context, event keptnapi.KeptnContextExtendedCE) error {
	return erh.ProcessEvent(event)
}

func (erh EventRequestHandler) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
		Name: erh.Env.K8SDeploymentName,
		MetaData: keptnapi.MetaData{
			Hostname:           erh.Env.K8SNodeName,
			IntegrationVersion: erh.Env.K8SDeploymentVersion,
			Location:           erh.Env.K8SDeploymentComponent,
			KubernetesMetaData: keptnapi.KubernetesMetaData{
				Namespace:      erh.Env.K8SNamespace,
				PodName:        erh.Env.K8SPodName,
				DeploymentName: erh.Env.K8SDeploymentName,
			},
		},
		Subscriptions: []keptnapi.EventSubscription{
			{
				Event:  "sh.keptn.event.>",
				Filter: keptnapi.EventSubscriptionFilter{},
			},
		},
	}
}
func NewEventRequestHandler(eventRepo db.EventRepo) *EventRequestHandler {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	return &EventRequestHandler{eventRepo: eventRepo, Env: env}
}

func (erh *EventRequestHandler) ProcessEvent(event keptnapi.KeptnContextExtendedCE) error {
	if *event.Type == keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName) {
		return erh.eventRepo.DropProjectCollections(event)
	}

	return erh.eventRepo.InsertEvent(event)
}

func (erh *EventRequestHandler) GetEvents(params event.GetEventsParams) (*event.GetEventsOKBody, error) {
	events, err := erh.eventRepo.GetEvents(params)
	if err != nil {
		return nil, err
	}
	if events.Events == nil {
		events.Events = []keptnapi.KeptnContextExtendedCE{}
	}
	return (*event.GetEventsOKBody)(events), nil
}

func (erh *EventRequestHandler) GetEventsByType(params event.GetEventsByTypeParams) (*event.GetEventsByTypeOKBody, error) {
	events, err := erh.eventRepo.GetEventsByType(params)
	if err != nil {
		return nil, err
	}
	if events.Events == nil {
		events.Events = []keptnapi.KeptnContextExtendedCE{}
	}
	return &event.GetEventsByTypeOKBody{Events: events.Events}, nil
}
