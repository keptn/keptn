package controlplane

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
)

type Integration struct{}

func (i Integration) OnEvent(event models.KeptnContextExtendedCE) {
	fmt.Println("Integration OnEvent " + event.ID)
}
