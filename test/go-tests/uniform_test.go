package go_tests

import (
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

// Test_UniformRegistration_TestAPI directly tests the API for (un)registering Keptn integrations
// to the Keptn control plane
func Test_UniformRegistration_TestAPI(t *testing.T) {

	uniformIntegration := &keptnmodels.Integration{
		Name: "my-uniform-service",
		MetaData: keptnmodels.MetaData{
			DistributorVersion: "0.8.3",
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscription: keptnmodels.Subscription{
			Topics: []string{keptnv2.GetTriggeredEventType(keptnv2.TestTaskName)},
		},
	}

	// register the integration at the shipyard controller
	resp, err := ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	registrationResponse := &models.RegisterResponse{}
	err = resp.ToJSON(registrationResponse)
	require.Nil(t, err)

	// retrieve the integration
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations := []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Equal(t, uniformIntegration.Name, integrations[0].Name)
	require.Equal(t, uniformIntegration.MetaData, integrations[0].MetaData)
	require.Equal(t, uniformIntegration.Subscription, integrations[0].Subscription)

	// delete the integration
	resp, err = ApiDELETERequest("/controlPlane/v1/uniform/registration/" + registrationResponse.ID)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// try to retrieve the integration again - should not be available anymore
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.Empty(t, integrations)
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane
func Test_UniformRegistration_RegistrationOfKeptnIntegration(t *testing.T) {
	// install echo integration
	deleteEchoIntegration, err := KubeCtlApplyFromURL("https://raw.githubusercontent.com/keptn-sandbox/echo-service/c7c97bb1b5affa938adb2f65260bcba8619a343f/deploy/service.yaml")
	require.Nil(t, err)

	// wait for echo integration registered
	var fetchedEchoIntegration models.Integration
	require.Eventually(t, func() bool {
		fetchedEchoIntegration, err = getIntegrationWithName("echo-service")
		return err == nil
	}, time.Second*20, time.Second*3)

	require.Nil(t, err)
	require.NotNil(t, fetchedEchoIntegration)
	require.Equal(t, "echo-service", fetchedEchoIntegration.Name)
	require.Equal(t, "echo-service", fetchedEchoIntegration.MetaData.KubernetesMetaData.DeploymentName)
	require.Equal(t, GetKeptnNameSpaceFromEnv(), fetchedEchoIntegration.MetaData.KubernetesMetaData.Namespace)
	require.Equal(t, "control-plane", fetchedEchoIntegration.MetaData.Location)
	require.Equal(t, StringArr("sh.keptn.>"), fetchedEchoIntegration.Subscription.Topics)

	// uninstall echo integration
	err = deleteEchoIntegration()
	require.Nil(t, err)

	// wait for echo integration unregistered
	require.Eventually(t, func() bool {
		fetchedEchoIntegration, err = getIntegrationWithName("echo-service")
		return err != nil
	}, time.Second*20, time.Second*3)
}

func getIntegrationWithName(name string) (models.Integration, error) {
	resp, _ := ApiGETRequest("/controlPlane/v1/uniform/registration")
	integrations := []models.Integration{}
	if err := resp.ToJSON(&integrations); err != nil {
		return models.Integration{}, err
	}
	for _, r := range integrations {
		if r.Name == "echo-service" {
			return r, nil
		}
	}
	return models.Integration{}, fmt.Errorf("No Keptn Inegration with name %s found", name)
}
