package go_tests

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	archiver "github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/require"
)

const testingShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
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

const resetGitRepos = `
#!/bin/sh

cd /data/config/
for FILE in *; do
    if [ -d "$FILE" ]; then
        cd "$FILE"
        git reset --hard
        cd ..
    fi
done`

// NOTE: When changing this test (especially the reset-get-repos.sh),
// please update the Keptn documentation for Backup & Restore accordingly.

func Test_BackupRestoreConfigService(t *testing.T) {
	serviceName := "configuration-service"
	BackupRestoreTestGeneric(t, serviceName)
}

func Test_BackupRestoreResourceService(t *testing.T) {
	serviceName := "resource-service"
	BackupRestoreTestGeneric(t, serviceName)
}

func BackupRestoreTestGeneric(t *testing.T, serviceUnderTestName string) {
	repoLocalDir := "../assets/podtato-head"
	projectName := "backup-restore"
	serviceName := "helloservice"
	chartFileName := "helloservice.tgz"
	serviceChartSrcPath := path.Join(repoLocalDir, "helm-charts", "helloservice")
	serviceChartArchivePath := path.Join(repoLocalDir, "helm-charts", chartFileName)
	serviceJmeterDir := path.Join(repoLocalDir, "jmeter")
	keptnNamespace := GetKeptnNameSpaceFromEnv()
	serviceHealthCheckEndpoint := "/metrics"
	secretFileName := "-credentials.yaml"
	serviceBackupFolder := "svc-backup"
	globalBackupFolder := "keptn-backup"
	mongoDBBackupFolder := "mongodb-backup"
	resetGitReposFile := "reset-git-repos.sh"

	// Delete chart archive at the end of the test
	defer func(path string) {
		err := os.RemoveAll(path)
		require.Nil(t, err)
	}(serviceChartArchivePath)

	err := archiver.Archive([]string{serviceChartSrcPath}, serviceChartArchivePath)
	require.Nil(t, err)

	t.Logf("Creating a new project %s with a Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(testingShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	secretFileName = projectName + secretFileName

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

	t.Logf("Trigger delivery before backup of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Verify Direct delivery before backup of %s in stage dev", serviceName)
	err = VerifyDirectDeployment(serviceName, projectName, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err := GetPublicURLOfService(serviceName, projectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Logf("Verify Direct delivery before backup of %s in stage prod", serviceName)
	err = VerifyDirectDeployment(serviceName, projectName, "prod", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &sequenceStates.States[0].Shkeptncontext}, 2*time.Minute, []string{models.SequenceFinished})

	t.Log("Verify network access to public URI of helloservice in stage prod")
	cartPubURL, err = GetPublicURLOfService(serviceName, projectName, "prod")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	backupGit := serviceUnderTestName == "configuration-service"

	t.Logf("Extracting name of service %s", serviceUnderTestName)
	serviceUnderTestPod, err := ExecuteCommandf("kubectl get pods -n %s -lapp.kubernetes.io/name=%s -ojsonpath='{.items[0].metadata.name}'", keptnNamespace, serviceUnderTestName)
	require.Nil(t, err)
	serviceUnderTestPod = removeQuotes(serviceUnderTestPod)

	//backup Configuration Service data

	t.Logf("Creating backup directories for %s", serviceUnderTestName)
	err = os.Chdir(repoLocalDir)
	require.Nil(t, err)

	globalBackupFolder, err = ioutil.TempDir("./", globalBackupFolder)
	require.Nil(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		require.Nil(t, err)
	}(globalBackupFolder)

	err = os.Chdir(globalBackupFolder)
	require.Nil(t, err)

	err = os.MkdirAll(serviceBackupFolder, os.ModePerm)
	require.Nil(t, err)
	defer resetTestPath(t, "../../../go-tests")

	t.Logf("Executing backup of %s", serviceUnderTestName)
	_, err = ExecuteCommandf("kubectl cp %s/%s:/data ./%s/ -c %s", keptnNamespace, serviceUnderTestPod, serviceBackupFolder, serviceUnderTestName)
	require.Nil(t, err)

	//backup MongoDB Data

	t.Logf("Creating backup directories for MongoDb data")
	err = os.MkdirAll(mongoDBBackupFolder, os.ModePerm)
	require.Nil(t, err)
	_, err = ExecuteCommandf("chmod o+w %s", mongoDBBackupFolder)
	require.Nil(t, err)

	t.Logf("Execute MongoDb database dump")
	mongoDbRootUser, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n %s -ojsonpath={.data.mongodb-root-user}", keptnNamespace)
	require.Nil(t, err)
	mongoDbRootUser, err = decodeBase64(removeQuotes(mongoDbRootUser))
	require.Nil(t, err)

	mongoDbRootPassword, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n %s -ojsonpath={.data.mongodb-root-password}", keptnNamespace)
	require.Nil(t, err)
	mongoDbRootPassword, err = decodeBase64(removeQuotes(mongoDbRootPassword))
	require.Nil(t, err)

	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n %s -- mongodump --authenticationDatabase admin --username %s --password %s -d keptn -h localhost --out=/tmp/dump", keptnNamespace, mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Executing backup of MongoDB database")
	mongoDbPod, err := ExecuteCommandf("kubectl get pods -n %s -lapp.kubernetes.io/name=mongo -ojsonpath='{.items[0].metadata.name}'", keptnNamespace)
	require.Nil(t, err)
	mongoDbPod = removeQuotes(mongoDbPod)
	_, err = ExecuteCommandf("kubectl cp %s/%s:/tmp/dump ./%s/ -c mongodb", keptnNamespace, mongoDbPod, mongoDBBackupFolder)

	if backupGit {

		//backup git-credentials

		t.Logf("Executing backup of git-credentials")
		secret, err := ExecuteCommandf("kubectl get secret -n %s git-credentials-%s -oyaml", keptnNamespace, projectName)
		require.Nil(t, err)
		err = os.WriteFile(secretFileName, []byte(secret), 0644)
		require.Nil(t, err)
	}

	if serviceUnderTestName == "resource-service" {
		//t.Logf("Deleting resource-service pod")
		//_, err = ExecuteCommandf("kubectl delete pod %s -n %s", serviceUnderTestPod, keptnNamespace)
		err := RestartPod(serviceUnderTestName)
		require.Nil(t, err)
	} else {
		t.Logf("Deleting testing project")
		_, err = ExecuteCommandf("keptn delete project %s", projectName)
		require.Nil(t, err)
	}

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	if backupGit {
		//restore git-credentials

		t.Logf("Executing restore of git-credentials")
		_, err = ExecuteCommandf("kubectl apply -f %s -n %s", secretFileName, keptnNamespace)
		require.Nil(t, err)
	}

	//restore Configuration/Resource Service data

	t.Logf("Restoring %s data", serviceUnderTestName)
	serviceUnderTestPod, err = ExecuteCommandf("kubectl get pods -n %s -lapp.kubernetes.io/name=%s -ojsonpath='{.items[0].metadata.name}'", keptnNamespace, serviceUnderTestName)
	require.Nil(t, err)
	serviceUnderTestPod = removeQuotes(serviceUnderTestPod)
	_, err = ExecuteCommandf("kubectl cp ./%s/config/ %s/%s:/data -c %s", serviceBackupFolder, keptnNamespace, serviceUnderTestPod, serviceUnderTestName)
	require.Nil(t, err)

	// reset git repositories to current HEAD

	t.Logf("Reseting git repositories to current HEAD")
	err = os.WriteFile(resetGitReposFile, []byte(resetGitRepos), 0666)
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./%s %s/%s:/data/config -c %s", resetGitReposFile, keptnNamespace, serviceUnderTestPod, serviceUnderTestName)
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec -n %s %s -c %s -- sh ./data/config/%s", keptnNamespace, serviceUnderTestPod, serviceUnderTestName, resetGitReposFile)
	require.Nil(t, err)

	////restore MongoDB data
	//
	t.Logf("Restoring MongoDB data")
	_, err = ExecuteCommandf("kubectl cp ./%s/keptn/ %s/%s:/tmp/dump -c mongodb", mongoDBBackupFolder, keptnNamespace, mongoDbPod)
	require.Nil(t, err)

	t.Logf("Import MongoDb database dump")
	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n %s -- mongorestore --drop --preserveUUID --authenticationDatabase admin --username %s --password %s /tmp/dump", keptnNamespace, mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Sleeping for 15s...")
	time.Sleep(50 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Trigger delivery after restore of helloservice:v0.1.1")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.1", "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Verify Direct delivery after restore of %s in stage dev", serviceName)
	err = VerifyDirectDeployment(serviceName, projectName, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.1")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err = GetPublicURLOfService(serviceName, projectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Logf("Verify Direct delivery after restore of %s in stage prod", serviceName)
	err = VerifyDirectDeployment(serviceName, projectName, "prod", "ghcr.io/podtato-head/podtatoserver", "v0.1.1")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod")
	cartPubURL, err = GetPublicURLOfService(serviceName, projectName, "prod")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)
}
