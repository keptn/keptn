package controlplane

import (
	"github.com/keptn/go-utils/pkg/api/models"
)

type RegistrationData models.Integration

type Integration interface {
	OnEvent(event models.KeptnContextExtendedCE)
	RegistrationData() RegistrationData
}
