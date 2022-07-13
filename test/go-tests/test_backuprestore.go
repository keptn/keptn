package go_tests

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

const testingShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-backup-restore"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "mytask"
`

func Test_BackupRestore(t *testing.T) {
	projectName := "backup-restore"
	serviceName := "helloservice"
	keptnNamespace := GetKeptnNameSpaceFromEnv()
	secretFileName := "-credentials.yaml"

	mongoDBBackupFolder, err := ioutil.TempDir("", "mongodb-backup")
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

	t.Log("Triggering sequence")
	keptnContextID, err := TriggerSequence(projectName, serviceName, "dev", "mysequence", nil)

	require.Nil(t, err)

	t.Log("Verifying if sequence reaches 'started' state")
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContextID}, 2*time.Minute, []string{models.SequenceStartedState})

	t.Log("Waiting for triggered event to be available")
	require.Eventually(t, func() bool {
		triggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
		if err != nil || triggeredEvent == nil {
			return false
		}
		return true
	}, 1*time.Minute, 5*time.Second)

	//backup MongoDB Data

	t.Logf("Creating backup directories for MongoDb data")
	err = os.MkdirAll(mongoDBBackupFolder, os.ModePerm)
	require.Nil(t, err)
	_, err = ExecuteCommandf("chmod o+w %s", mongoDBBackupFolder)
	require.Nil(t, err)

	t.Logf("Execute MongoDb database dump")
	mongoDbRootUser, mongoDbRootPassword, err := GetMongoDBCredentials()

	require.Nil(t, err)

	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n %s -- mongodump --authenticationDatabase admin --username %s --password %s -d keptn -h localhost --out=/tmp/dump", keptnNamespace, mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Executing backup of MongoDB database")

	mongoDbPods, err := GetPodNamesOfDeployment("app.kubernetes.io/name=mongo")
	require.Nil(t, err)
	require.NotEmpty(t, mongoDbPods)
	mongoDbPod := mongoDbPods[0]

	_, err = ExecuteCommandf("kubectl cp %s/%s:/tmp/dump ./%s/ -c mongodb", keptnNamespace, mongoDbPod, mongoDBBackupFolder)

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

}
