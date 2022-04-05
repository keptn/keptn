package go_tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/mholt/archiver/v3"

	"github.com/stretchr/testify/require"
)

const testingProxyShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
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
const baseProxyProjectPath = "/controlPlane/v1/project"

func Test_ProxyAuth(t *testing.T) {
	repoLocalDir := "../assets/podtato-head"
	projectName := "proxy-auth"
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
	shipyardFilePath, err := CreateTmpShipyardFile(testingProxyShipyard)
	require.Nil(t, err)
	projectName, err = CreateProjectWithProxy(projectName, shipyardFilePath)
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

	t.Logf("Getting project %s with a proxy", projectName)
	resp, err := ApiGETRequest(baseProxyProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking if upstream was provisioned")
	project := models.ExpandedProject{}
	err = resp.ToJSON(&project)
	require.Nil(t, err)
	require.Equal(t, "squid:3128", project.GitProxyURL)
	require.Equal(t, "http", project.GitProxyScheme)
	require.Equal(t, "", project.GitProxyUser)
	require.Equal(t, true, project.GitProxyInsecure)
	require.Equal(t, projectName, project.ProjectName)

	t.Logf("Updating project credentials")
	user := GetGiteaUser()
	token, err := GetGiteaToken()
	require.Nil(t, err)

	// apply the k8s job for creating the git upstream
	_, err = ExecuteCommand(fmt.Sprintf("keptn update project %s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s --git-proxy-url=squid:3128 --git-proxy-scheme=http --git-proxy-insecure", projectName, user, projectName, user, token))
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", secondServiceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", secondServiceName, projectName)
	require.Nil(t, err)

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

	//Modify the proxy settings to be certain that no other project use the proxy
	_, err = ExecuteCommand(fmt.Sprintf("keptn update project %s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s --git-proxy-url=squid:3124 --git-proxy-scheme=http --git-proxy-insecure", projectName, user, projectName, user, token))
	require.Nil(t, err)
}
