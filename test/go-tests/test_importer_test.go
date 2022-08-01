package go_tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const importerShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "user_managed"`

func Test_ImportCorrectManifest(t *testing.T) {
	projectName := "keptn-importer-test"
	serviceName := "my-service-name"

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = recursiveZip("../assets/import/sample-package/", "./sample-package.zip", false)
	require.Nil(t, err)

	// Getting Keptn API credentials
	_, keptnApiUrl, err := GetApiCredentials()
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject(fmt.Sprintf("%s/v1/import", keptnApiUrl), projectName, "./sample-package.zip")
	require.Nil(t, err)
	require.Equal(t, 200, responseCode, fmt.Sprintf("Expected response status 200 but got %d", responseCode))

	// Check if service exists
	approvalCLIOutput, err := ExecuteCommand(fmt.Sprintf("keptn get service --project %s %s", projectName, serviceName))
	require.NoError(t, err)
	require.NotContains(t, approvalCLIOutput, fmt.Sprintf("No services %s found in project", serviceName))

	// Check if secret exists
	err = checkIfSecretExists("slack-webhook")
	require.NoError(t, err)

	// Validate webhook
	exists, err := CheckIfWebhookSubscriptionExists(projectName, "sh.keptn.event.evaluation.triggered")
	require.NoError(t, err)
	require.True(t, exists, "Webhook service subscription does not exist")

}

func Test_ImportCorrectManifestNonExistingProject(t *testing.T) {
	projectName := "keptn-importer-test-non-existing"
	wrongProjectName := "ketpn-importer-test-non-existing"
	errorMessage := "project not found"

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = recursiveZip("../assets/import/sample-package/", "./sample-package.zip", false)
	require.Nil(t, err)

	// Getting Keptn API credentials
	_, keptnApiUrl, err := GetApiCredentials()
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject(fmt.Sprintf("%s/v1/import", keptnApiUrl), wrongProjectName, "./sample-package.zip")
	require.Equal(t, 424, responseCode, fmt.Sprintf("Expected response status 424 but got %d", responseCode))
	require.ErrorContains(t, err, errorMessage, fmt.Sprintf("Could not find expected error message: %s", errorMessage))
	require.Error(t, err)
}

func Test_ImportMalformedZipFileCorrectName(t *testing.T) {
	projectName := "keptn-importer-test-malformed-zip"
	errorMessage := "Error opening import archive"

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = recursiveZip("../assets/import/invalid-package/", "./invalid-package.zip", false)
	require.Nil(t, err)

	// Getting Keptn API credentials
	_, keptnApiUrl, err := GetApiCredentials()
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject(fmt.Sprintf("%s/v1/import", keptnApiUrl), projectName, "./invalid-package.zip")
	require.Equal(t, 415, responseCode, fmt.Sprintf("Expected response status 424 but got %d", responseCode))
	require.ErrorContains(t, err, errorMessage, fmt.Sprintf("Could not find expected error message: %s", errorMessage))
	require.Error(t, err)
}

func ImportUploadZipToProject(urlPath, projectName, filePath string) (int, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("configPackage", filepath.Base(filePath))
	if err != nil {
		return 400, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return 400, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return 400, err
	}
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?project=%s", urlPath, projectName), bytes.NewReader(body.Bytes()))
	if err != nil {
		return 400, err
	}
	token, _, _ := GetApiCredentials()

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-token", token)

	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return rsp.StatusCode, err
		}
		bodyString := string(bodyBytes)
		return rsp.StatusCode, fmt.Errorf(bodyString)
	}
	return rsp.StatusCode, nil
}
