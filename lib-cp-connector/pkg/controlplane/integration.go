package controlplane

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
)

type RegistrationData models.Integration

// Integration represents a Keptn Service that wants to receive events from the Keptn Control plane
type Integration interface {
	// OnEvent is called when a new event was received
	OnEvent(context.Context, models.KeptnContextExtendedCE) error

	// RegistrationData is called to get the initial registration data
	RegistrationData() RegistrationData
}
