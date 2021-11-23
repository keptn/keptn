package go_tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testingShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-potatohead"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "staging"
      sequences:
        - name: "delivery-direct"
          triggeredOn:
            - event: "dev.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "production"
      sequences:
        - name: "delivery-direct"
          triggeredOn:
            - event: "staging.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"
`

const resetGitReposScript = `cat <<EOT >> reset-git-repos.sh
#!/bin/sh

cd /data/config/
for FILE in *; do
    if [ -d "\$FILE" ]; then
        cd "\$FILE"
        git reset --hard
        cd ..
    fi
done
EOT
`

func TestBackupRestore(t *testing.T) {
	repoLocalDir := "../assets/potatohead/"
	keptnProjectName := "potatohead"
	serviceName := "helloserver"
	resourceName := "helloserver.tgz"

	shipyardFilePath, err := CreateTmpShipyardFile(testingShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("Creating a new project %s without a GIT Upstream", keptnProjectName)
	err = CreateProject(keptnProjectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Creating a new service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn create service %s --project=%s", serviceName, keptnProjectName)
	require.Nil(t, err)

	t.Logf("Adding resource %s to service %s in project %s", resourceName, serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --all-stages --resource=%s --resourceUri=%s", keptnProjectName, serviceName, repoLocalDir+resourceName, "helm/"+resourceName)
	require.Nil(t, err)

	t.Logf("Triggering delivery of service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image='gabrieltanner/hello-server' --tag=v0.1.1", keptnProjectName, serviceName)
	require.Nil(t, err)

	//backup Configuration Service data

	t.Logf("Creating backup directories for configuration-service")
	err = os.MkdirAll("keptn-backup", os.ModePerm)
	require.Nil(t, err)
	err = os.Chdir("keptn-backup")
	require.Nil(t, err)
	err = os.MkdirAll("config-svc-backup", os.ModePerm)
	require.Nil(t, err)

	t.Logf("Executing backup of configuration-service")
	configServicePod, err := ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp keptn/%s:/data ./config-svc-backup/ -c configuration-service", configServicePod)
	require.Nil(t, err)

	//backup MongoDB Data

	t.Logf("Creating backup directories for MongoDb data")
	err = os.MkdirAll("mongodb-backup", os.ModePerm)
	require.Nil(t, err)

	t.Logf("Logging in to MongoDb database")
	mongoDbRootUser, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.mongodb-root-user} | base64 -d")
	require.Nil(t, err)
	mongoDbRootPassword, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.mongodb-root-password} | base64 --decode")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n keptn -- mongodump --authenticationDatabase admin --username %s--password %s -d keptn -h localhost --out=/tmp/dump", mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Executing backup of MongoDB database")
	mongoDbPod, err := ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=mongo -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp keptn/%s:/tmp/dump ./mongodb-backup/ -c mongodb", mongoDbPod)
	require.Nil(t, err)

	//backup git credentials

	t.Logf("Executing backup of git credentials")
	_, err = ExecuteCommandf("kubectl get secret -n keptn git-credentials-%s -oyaml > %s-credentials.yaml", keptnProjectName, keptnProjectName)
	require.Nil(t, err)

	//deleting testing project

	t.Logf("Deleting testing project")
	_, err = ExecuteCommandf("keptn delete project %s", keptnProjectName)
	require.Nil(t, err)

	//restore Configuration Service data
	t.Logf("Restoring configuration-service data")
	configServicePod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./config-svc-backup/* keptn/%s:/data -c configuration-service", configServicePod)
	require.Nil(t, err)

	t.Logf("Reseting git repository HEAD")
	_, err = ExecuteCommandf(resetGitReposScript)
	require.Nil(t, err)
	configServicePod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./reset-git-repos.sh keptn/%s:/ -c configuration-service", configServicePod)
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec -n keptn %s -c configuration-service -- chmod +x -R ./reset-git-repos.sh", configServicePod)
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec -n keptn %s -c configuration-service -- ./reset-git-repos.sh", configServicePod)
	require.Nil(t, err)

	//restore MongoDB data

	t.Logf("Restoring MongoDB data")
	mongoDbPod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=mongo -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./mongodb-backup/ keptn/%s:/opt/dump -c mongodb", mongoDbPod)
	require.Nil(t, err)

	t.Logf("Logging in to MongoDb database")
	mongoDbRootUser, err = ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.mongodb-root-user} | base64 -d")
	require.Nil(t, err)
	mongoDbRootPassword, err = ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.mongodb-root-password} | base64 --decode")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec svc/keptn-mongo -n keptn -- mongorestore --drop --preserveUUID --authenticationDatabase admin --username %s --password %s /opt/dump", mongoDbRootUser, mongoDbRootPassword)
	require.Nil(t, err)

	//restore git credentials
	t.Logf("Restoring git credentials")
	_, err = ExecuteCommandf("kubectl apply -f %s-credentials.yaml", keptnProjectName)
	require.Nil(t, err)

}
