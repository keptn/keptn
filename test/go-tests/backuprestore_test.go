package go_tests

import (
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"os"
	"path"
	"testing"
	"time"

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

func Test_BackupRestore(t *testing.T) {
	repoLocalDir := "../assets/podtato-head"
	keptnProjectName := "backup-restore"
	serviceName := "helloservice"
	serviceChartLocalDir := path.Join(repoLocalDir, "helm-charts", "helloserver")
	serviceJmeterDir := path.Join(repoLocalDir, "jmeter")
	keptnNamespace := GetKeptnNameSpaceFromEnv()
	serviceHealthCheckEndpoint := "/metrics"

	t.Logf("Creating a new project %s without a GIT Upstream", keptnProjectName)
	shipyardFilePath, err := CreateTmpShipyardFile(testingShipyard)
	require.Nil(t, err)
	err = CreateProject(keptnProjectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Onboarding service %s in project %s with chart %s", serviceName, keptnProjectName, serviceChartLocalDir)
	_, err = ExecuteCommandf("keptn onboard service %s --project %s --chart=%s", serviceName, keptnProjectName, serviceChartLocalDir)
	require.Nil(t, err)

	t.Log("Adding jmeter config in prod")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "prod", serviceJmeterDir+"/jmeter.conf.yaml", "jmeter/jmeter.conf.yaml")
	require.Nil(t, err)

	t.Log("Adding load test resources for jmeter in prod")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "prod", serviceJmeterDir+"/load.jmx", "jmeter/load.jmx")
	require.Nil(t, err)

	t.Logf("Trigger delivery before backup of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Verify Direct delivery before backup of %s in stage dev", serviceName)
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err := GetPublicURLOfService(serviceName, keptnProjectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Logf("Verify Direct delivery before backup of %s in stage prod", serviceName)
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "prod", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	sequenceStates, _, err := GetState(keptnProjectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	VerifySequenceEndsUpInState(t, keptnProjectName, &keptnapimodels.EventContext{KeptnContext: &sequenceStates.States[0].Shkeptncontext}, 2*time.Minute, []string{models.SequenceFinished})

	t.Log("Verify network access to public URI of helloservice in stage prod")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "prod")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	//backup Configuration Service data

	t.Logf("Creating backup directories for configuration-service")
	err = os.Chdir(repoLocalDir)
	require.Nil(t, err)
	err = os.MkdirAll("keptn-backup", os.ModePerm)
	require.Nil(t, err)
	err = os.Chdir("keptn-backup")
	require.Nil(t, err)
	err = os.MkdirAll("config-svc-backup", os.ModePerm)
	require.Nil(t, err)

	t.Logf("Executing backup of configuration-service")
	configServicePod, err := ExecuteCommandf("kubectl get pods -n %s -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'", keptnNamespace)
	require.Nil(t, err)
	configServicePod = removeQuotes(configServicePod)
	_, err = ExecuteCommandf("kubectl cp %s/%s:/data ./config-svc-backup/ -c configuration-service", keptnNamespace, configServicePod)
	require.Nil(t, err)

	//backup MongoDB Data

	t.Logf("Creating backup directories for MongoDb data")
	err = os.MkdirAll("mongodb-backup", os.ModePerm)
	require.Nil(t, err)
	_, err = ExecuteCommandf("chmod o+w mongodb-backup")
	require.Nil(t, err)

	t.Logf("Execute MongoDb database dump")
	mongoDbRootUser, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n %s -ojsonpath={.data.mongodb-root-user}", keptnNamespace)
	require.Nil(t, err)
	mongoDbRootUser, err = decodeBase64((removeQuotes(mongoDbRootUser)))
	require.Nil(t, err)

	mongoDbRootPassword, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n %s -ojsonpath={.data.mongodb-root-password}", keptnNamespace)
	require.Nil(t, err)
	mongoDbRootPassword, err = decodeBase64((removeQuotes(mongoDbRootPassword)))
	require.Nil(t, err)

	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n %s -- mongodump --authenticationDatabase admin --username %s --password %s -d keptn -h localhost --out=/tmp/dump", keptnNamespace, mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Executing backup of MongoDB database")
	mongoDbPod, err := ExecuteCommandf("kubectl get pods -n %s -lapp.kubernetes.io/name=mongo -ojsonpath='{.items[0].metadata.name}'", keptnNamespace)
	require.Nil(t, err)
	mongoDbPod = removeQuotes(mongoDbPod)
	_, err = ExecuteCommandf("kubectl cp %s/%s:/tmp/dump ./mongodb-backup/ -c mongodb", keptnNamespace, mongoDbPod)
	require.Nil(t, err)

	//deleting testing project

	t.Logf("Deleting testing project")
	_, err = ExecuteCommandf("keptn delete project %s", keptnProjectName)
	require.Nil(t, err)

	t.Logf("Sleeping for 15s...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	//restore Configuration Service data

	t.Logf("Restoring configuration-service data")
	_, err = ExecuteCommandf("kubectl cp ./config-svc-backup/config/ %s/%s:/data -c configuration-service", keptnNamespace, configServicePod)
	require.Nil(t, err)

	//restore MongoDB data

	t.Logf("Restoring MongoDB data")
	_, err = ExecuteCommandf("kubectl cp ./mongodb-backup/keptn/ %s/%s:/tmp/dump -c mongodb", keptnNamespace, mongoDbPod)
	require.Nil(t, err)

	t.Logf("Import MongoDb database dump")
	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n %s -- mongorestore --drop --preserveUUID --authenticationDatabase admin --username %s --password %s /tmp/dump", keptnNamespace, mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Sleeping for 15s...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Trigger delivery after restore of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Logf("Verify Direct delivery after restore of %s in stage dev", serviceName)
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Logf("Verify Direct delivery after restore of %s in stage prod", serviceName)
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "prod", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "prod")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

}
