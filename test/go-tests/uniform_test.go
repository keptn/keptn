package go_tests

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_UniformRegistration(t *testing.T) {

	uniformIntegration := &models.Integration{
		Name: "my-uniform-service",
		MetaData: models.MetaData{
			DeploymentName:     "my-uniform-service",
			DistributorVersion: "0.8.3",
			Status:             "active",
		},
		Subscription: models.Subscription{
			Topics: []string{keptnv2.GetTriggeredEventType(keptnv2.TestTaskName)},
		},
	}

	// register the integration at the shipyard controller
	resp, err := ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// retrieve the integration
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=my-uniform-id")

	integrations := []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Equal(t, uniformIntegration.Name, integrations[0].Name)
	require.Equal(t, uniformIntegration.MetaData, integrations[0].MetaData)
	require.Equal(t, uniformIntegration.Subscription, integrations[0].Subscription)
}
