package configuration_changer

import (
	"fmt"
	"strings"

	"github.com/keptn/keptn/helm-service/controller/mesh"

	"helm.sh/helm/v3/pkg/chart"
)

type CanaryWeightUpdater struct {
	mesh         mesh.Mesh
	canaryWeight int32
}

func NewCanaryWeightUpdater(mesh mesh.Mesh, canaryWeight int32) *CanaryWeightUpdater {
	return &CanaryWeightUpdater{
		mesh:         mesh,
		canaryWeight: canaryWeight,
	}
}

// Update updates the provided traffic weight in the VirtualService contained in the chart
func (c *CanaryWeightUpdater) Update(ch *chart.Chart) error {

	// Set weights in all virtualservices
	for _, template := range ch.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, c.mesh.GetVirtualServiceSuffix()) {

			vs, err := c.mesh.UpdateWeights(template.Data, c.canaryWeight)
			if err != nil {
				return fmt.Errorf("Error when setting new weights in VirtualService %s: %s",
					template.Name, err.Error())
			}
			template.Data = vs
			break
		}
	}
	return nil
}
