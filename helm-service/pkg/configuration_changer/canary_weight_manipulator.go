package configuration_changer

import (
	"fmt"
	"strings"

	"github.com/keptn/keptn/helm-service/pkg/mesh"

	"helm.sh/helm/v3/pkg/chart"
)

// CanaryWeightManipulator allows to manipulate the traffic weight for the canary in the VirtualService
type CanaryWeightManipulator struct {
	mesh         mesh.Mesh
	canaryWeight int32
}

// NewCanaryWeightManipulator creates a CanaryWeightManipulator
func NewCanaryWeightManipulator(mesh mesh.Mesh, canaryWeight int32) *CanaryWeightManipulator {
	return &CanaryWeightManipulator{
		mesh:         mesh,
		canaryWeight: canaryWeight,
	}
}

// Update updates the provided traffic weight in the VirtualService contained in the chart
func (c *CanaryWeightManipulator) Manipulate(ch *chart.Chart) error {

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
