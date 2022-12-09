package go_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_CreateProjectWithCustomBranchName(t *testing.T) {
	projectMaster := "project-master"
	projectMain := "project-main-custom"
	keptnNamespace := GetKeptnNameSpaceFromEnv()

	defer func() {
		logs, err := PrintLogsWithDeploymentName("app.kubernetes.io/name=resource-service")
		require.Nil(t, err)
		t.Log("logs from RService: ")
		t.Log(logs)

		logsShippy, err := PrintLogsWithDeploymentName("app.kubernetes.io/name=shipyard-controller")
		require.Nil(t, err)
		t.Log("logs from Shippy: ")
		t.Log(logsShippy)
	}()

	t.Logf("Creating a new project %s with a Gitea Upstream", projectMaster)
	shipyardFilePath, err := CreateTmpShipyardFile(testingShipyard)
	require.Nil(t, err)
	projectName, err := CreateProject(projectMaster, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("Checking if main branch of project %s is master", projectMaster)
	equal, err := VerifyMainRepositoryBranchName(projectName, "master")
	require.Nil(t, err)
	require.True(t, equal)

	defer func() {
		_, err = ExecuteCommandf("kubectl set env deployment/resource-service DEFAULT_REMOTE_GIT_BRANCH=%s -n %s", "master", keptnNamespace)
		require.Nil(t, err)
		time.Sleep(75 * time.Second)
	}()

	t.Logf("Setting up DEFAULT_REMOTE_GIT_BRANCH env variable to main-custom")
	_, err = ExecuteCommandf("kubectl set env deployment/resource-service DEFAULT_REMOTE_GIT_BRANCH=%s -n %s", "main-custom", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 75s to make sure resource-service pod is restarted...")
	time.Sleep(75 * time.Second)

	t.Logf("Creating a new project %s with a Gitea Upstream", projectMain)
	projectName, err = CreateProject(projectMain, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("Checking if main branch of project %s is main-custom", projectMaster)
	equal, err = VerifyMainRepositoryBranchName(projectName, "main-custom")
	require.Nil(t, err)
	require.True(t, equal)
}
