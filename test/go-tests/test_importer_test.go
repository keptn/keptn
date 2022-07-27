package go_tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

	// Make API call with ZIP file
	err = ImportUploadZipToProject("http://35.192.145.202.nip.io/api/v1/import", "POST", projectName)
	require.Nil(t, err)

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

func ImportUploadZipToProject(urlPath, method, projectName string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("configPackage", "sample-package.zip")
	if err != nil {
		return err
	}
	file, err := os.Open("./sample-package.zip")
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	writer.Close()
	req, err := http.NewRequest(method, fmt.Sprintf("%s?project=%s", urlPath, projectName), bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	token, _, _ := GetApiCredentials()

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-token", token)

	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with response code: %d", rsp.StatusCode)
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
	return nil
}
