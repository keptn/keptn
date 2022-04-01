package controlplane

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
)

type RegistrationData models.Integration

type Integration interface {
	OnEvent(context.Context, models.KeptnContextExtendedCE) error
	RegistrationData() RegistrationData
}
