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
            - name: "deployment"`

// Test_ImportCorrectManifest uploads a valid zip manifest which creates a Keptn service, secret and webhook and validates the result
func Test_ImportCorrectManifest(t *testing.T) {
	projectName := "keptn-importer-test"
	serviceName := "my-service-name"

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = createZipFileFromDirectory("../assets/import/sample-package/", "./sample-package.zip", false)
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject("v1/import", projectName, "./sample-package.zip")
	require.Nil(t, err)
	require.Equal(t, 200, responseCode, fmt.Sprintf("Expected response status 200 but got %d", responseCode))

	// Check if service exists
	getServiceResponse, err := ExecuteCommand(fmt.Sprintf("keptn get service --project %s %s", projectName, serviceName))
	require.NoError(t, err)
	require.NotContains(t, getServiceResponse, fmt.Sprintf("No services %s found in project", serviceName))

	// Check if secret exists
	err = checkIfSecretExists("slack-webhook", GetKeptnNameSpaceFromEnv())
	require.NoError(t, err)

	// Validate webhook
	exists, err := CheckIfWebhookSubscriptionExists(projectName, "sh.keptn.event.evaluation.triggered")
	require.NoError(t, err)
	require.True(t, exists, "Webhook service subscription does not exist")

}

// Test_ImportCorrectManifestNonExistingProject uploads a valid manifest with a non-existing project which throws an error when uploading
func Test_ImportCorrectManifestNonExistingProject(t *testing.T) {
	projectName := "keptn-importer-test-non-existing"
	wrongProjectName := "ketpn-importer-test-non-existing"
	errorMessage := fmt.Sprintf("project %s does not exist", projectName)
	expectedErrorCode := 404

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = createZipFileFromDirectory("../assets/import/sample-package/", "./sample-package.zip", false)
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject("v1/import", wrongProjectName, "./sample-package.zip")
	require.Equal(t, expectedErrorCode, responseCode, fmt.Sprintf("Expected response status %d but got %d", expectedErrorCode, responseCode))
	require.ErrorContains(t, err, errorMessage, fmt.Sprintf("Could not find expected error message: %s", errorMessage))
	require.Error(t, err)
}

// Test_ImportMalformedZipFileCorrectName uploads an invalid zip with incorrect structure and mapping that throws an error while uploading
func Test_ImportMalformedZipFileCorrectName(t *testing.T) {
	projectName := "keptn-importer-test-malformed-zip"
	errorMessage := "Error opening import archive"
	expectedErrorCode := 415

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(importerShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// Convert folder to Zip
	err = createZipFileFromDirectory("../assets/import/invalid-package/", "./invalid-package.zip", false)
	require.Nil(t, err)

	// Make API call with ZIP file
	responseCode, err := ImportUploadZipToProject("v1/import", projectName, "./invalid-package.zip")
	require.Equal(t, expectedErrorCode, responseCode, fmt.Sprintf("Expected response status %d but got %d", expectedErrorCode, responseCode))
	require.ErrorContains(t, err, errorMessage, fmt.Sprintf("Could not find expected error message: %s", errorMessage))
	require.Error(t, err)
}

func ImportUploadZipToProject(urlPath, projectName, filePath string) (int, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Get Keptn credentials
	token, keptnApiUrl, err := GetApiCredentials()
	if err != nil {
		return 500, err
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?project=%s", fmt.Sprintf("%s/%s", keptnApiUrl, urlPath), projectName), bytes.NewReader(body.Bytes()))
	if err != nil {
		return 400, err
	}

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
