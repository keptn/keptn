package go_tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/mholt/archiver/v3"
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

func Test_SSHPublicKeyAuth(t *testing.T) {
	repoLocalDir := "../assets/podtato-head"
	projectName := "public-key-auth"
	serviceName := "helloservice"
	secondServiceName := "helloservice2"
	chartFileName := "helloservice.tgz"
	serviceChartSrcPath := path.Join(repoLocalDir, "helm-charts", "helloservice")
	serviceChartArchivePath := path.Join(repoLocalDir, "helm-charts", chartFileName)
	serviceJmeterDir := path.Join(repoLocalDir, "jmeter")

	// Delete chart archive at the end of the test
	defer func(path string) {
		err := os.RemoveAll(path)
		require.Nil(t, err)
	}(serviceChartArchivePath)

	err := archiver.Archive([]string{serviceChartSrcPath}, serviceChartArchivePath)
	require.Nil(t, err)

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(testingSSHShipyard)
	require.Nil(t, err)
	projectName, err = CreateProjectWithSSH(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, projectName)
	require.Nil(t, err)

	t.Logf("Adding resource for service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommandf("keptn add-resource --project %s --service=%s --all-stages --resource=%s --resourceUri=%s", projectName, serviceName, serviceChartArchivePath, path.Join("helm", chartFileName))
	require.Nil(t, err)

	t.Log("Adding jmeter config in prod")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", projectName, serviceName, "prod", serviceJmeterDir+"/jmeter.conf.yaml", "jmeter/jmeter.conf.yaml")
	require.Nil(t, err)

	t.Log("Adding load test resources for jmeter in prod")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", projectName, serviceName, "prod", serviceJmeterDir+"/load.jmx", "jmeter/load.jmx")
	require.Nil(t, err)

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

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

	t.Logf("Creating service %s in project %s", secondServiceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", secondServiceName, projectName)
	require.Nil(t, err)

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)
}
