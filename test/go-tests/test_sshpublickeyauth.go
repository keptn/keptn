package go_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

const testingSSHShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-podtato-ohead"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "prod"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "test"
              properties:
                teststrategy: "performance"
            - name: "evaluation"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "prod.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"

        - name: "delivery-direct"
          triggeredOn:
            - event: "dev.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"
`
const baseSSHProjectPath = "/controlPlane/v1/project"

func Test_SSHPublicKeyAuth(t *testing.T) {
	projectName := "public-key-auth"
	serviceName := "helloservice"
	secondServiceName := "helloservice2"

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(testingSSHShipyard)
	require.Nil(t, err)
	projectName, err = CreateProjectWithSSH(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, projectName)
	require.Nil(t, err)

	t.Logf("Getting project %s with a SSH publicKey", projectName)
	resp, err := ApiGETRequest(baseSSHProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking if upstream was provisioned")
	project := models.ExpandedProject{}
	err = resp.ToJSON(&project)
	require.Nil(t, err)
	require.Contains(t, project.GitCredentials.RemoteURL, "ssh://")
	require.Equal(t, projectName, project.ProjectName)

	t.Logf("Updating project credentials")
	user := GetGiteaUser()
	privateKey, passphrase, err := GetPrivateKeyAndPassphrase()
	require.Nil(t, err)

	privateKeyPath := "private-key"
	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0777)
	require.Nil(t, err)

	defer func() {
		os.Remove(privateKeyPath)
	}()

	_, err = ExecuteCommand(fmt.Sprintf("keptn update project %s --git-remote-url=ssh://gitea-ssh:22/%s/%s.git --git-user=%s --git-private-key=%s --git-private-key-pass=%s", projectName, user, projectName, user, privateKeyPath, passphrase))
	require.Nil(t, err)

	// check if interacting with the project (e.g. adding a service) still works after updating
	t.Logf("Creating service %s in project %s", secondServiceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", secondServiceName, projectName)
	require.Nil(t, err)
}
