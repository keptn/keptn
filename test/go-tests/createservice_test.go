package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)


const simpleShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"

    - name: "staging"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment
"`


func Test_createServiceWithTooLongName(t *testing.T) {
	PrepareEnvVars()
	projectName := "my-super-long-project"
	serviceName := "my-very-very-very-very-very-super-long-service-name"
	file, err := CreateTmpShipyardFile(simpleShipyard)
	require.Nil(t, err)
	defer os.Remove(file)

	// check if the project 'state' is already available - if not, delete it before creating it again
	resp, err := ApiGETRequest("/controlPlane/v1/project/" + projectName)

	if resp.Response().StatusCode != http.StatusNotFound {
		// delete project if it exists
		_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", projectName))
		require.NotNil(t, err)
	}

	output, err := ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=./%s", projectName, file))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	output, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	// ToDo: Should fail because of shipyard-controller
	require.Nil(t, err)
	require.Contains(t, output, "created successfully")
}


func Test_createServiceWithLongName(t *testing.T) {
	PrepareEnvVars()
	projectName := "my-super-long-project"
	serviceName := "my-very-very-very-super-long-service-name"
	file, err := CreateTmpShipyardFile(simpleShipyard)
	require.Nil(t, err)
	defer os.Remove(file)

	// check if the project 'state' is already available - if not, delete it before creating it again
	resp, err := ApiGETRequest("/controlPlane/v1/project/" + projectName)

	if resp.Response().StatusCode != http.StatusNotFound {
		// delete project if it exists
		_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", projectName))
		require.NotNil(t, err)
	}

	output, err := ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=./%s", projectName, file))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	output, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")
}
