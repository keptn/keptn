// This file is safe to edit. Once it exists it will not be overwritten

package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/keptn/api/restapi/operations/event"
	"github.com/keptn/api/utils"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const eventbrokerURL = "http://event-broker.keptn.svc.cluster.local/keptn"

// ForwardEvent sends the received event to the eventbroker
func ForwardEvent(e event.SendEventBody) error {

	if e.Shkeptncontext == nil || *e.Shkeptncontext == "" {
		uuidStr := uuid.New().String()
		e.Shkeptncontext = &uuidStr
	}

	keptnutils.Info(*e.Shkeptncontext, fmt.Sprintf("Sending keptn event with type %s", *e.Type))

	return utils.PostToEventBroker(e, *e.Shkeptncontext)
}
