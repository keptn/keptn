package go_tests

import (
	b64 "encoding/base64"
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

func Test_BackupRestore(t *testing.T) {
	repoLocalDir := "../assets/potatohead/"
	keptnProjectName := "podtato-head"
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
	err = os.Chdir(repoLocalDir)
	require.Nil(t, err)
	err = os.MkdirAll("keptn-backup", os.ModePerm)
	require.Nil(t, err)
	err = os.Chdir("keptn-backup")
	require.Nil(t, err)
	err = os.MkdirAll("config-svc-backup", os.ModePerm)
	require.Nil(t, err)

	t.Logf("Executing backup of configuration-service")
	configServicePod, err := ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp keptn/%s:/data ./config-svc-backup/ -c configuration-service", removeQuotes(configServicePod))
	require.Nil(t, err)

	//backup MongoDB Data

	t.Logf("Creating backup directories for MongoDb data")
	err = os.MkdirAll("mongodb-backup", os.ModePerm)
	require.Nil(t, err)

	t.Logf("Logging in to MongoDb database")
	mongoDbRootUser, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.user}")
	require.Nil(t, err)
	mongoDbRootUserByte, err := b64.StdEncoding.DecodeString(removeQuotes(mongoDbRootUser))
	require.Nil(t, err)
	mongoDbRootUser = string(mongoDbRootUserByte)

	mongoDbRootPassword, err := ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.admin_password}")
	require.Nil(t, err)
	mongoDbRootPasswordByte, err := b64.StdEncoding.DecodeString(removeQuotes(mongoDbRootPassword))
	require.Nil(t, err)
	mongoDbRootPassword = string(mongoDbRootPasswordByte)

	_, err = ExecuteCommandf("kubectl exec svc/mongodb -n keptn -- mongodump --authenticationDatabase admin --username admin --password %s -d keptn -h localhost --out=/tmp/dump", mongoDbRootPassword)
	require.Nil(t, err)

	t.Logf("Executing backup of MongoDB database")
	mongoDbPod, err := ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=mongodb -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp keptn/%s:/tmp/dump ./mongodb-backup/ -c mongodb", removeQuotes(mongoDbPod))
	require.Nil(t, err)

	//backup git credentials

	//t.Logf("Executing backup of git credentials")
	//_, err = ExecuteCommandf("kubectl get secret -n keptn git-credentials-%s -oyaml > %s-credentials.yaml", keptnProjectName, keptnProjectName)
	//require.Nil(t, err)

	//deleting testing project

	t.Logf("Deleting testing project")
	_, err = ExecuteCommandf("keptn delete project %s", keptnProjectName)
	require.Nil(t, err)

	//restore Configuration Service data
	t.Logf("Restoring configuration-service data")
	configServicePod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./config-svc-backup/* keptn/%s:/data -c configuration-service", removeQuotes(configServicePod))
	require.Nil(t, err)

	t.Logf("Reseting git repository HEAD")
	_, err = ExecuteCommandf(resetGitReposScript)
	require.Nil(t, err)
	configServicePod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./reset-git-repos.sh keptn/%s:/ -c configuration-service", removeQuotes(configServicePod))
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec -n keptn %s -c configuration-service -- chmod +x -R ./reset-git-repos.sh", removeQuotes(configServicePod))
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl exec -n keptn %s -c configuration-service -- ./reset-git-repos.sh", removeQuotes(configServicePod))
	require.Nil(t, err)

	//restore MongoDB data

	t.Logf("Restoring MongoDB data")
	mongoDbPod, err = ExecuteCommandf("kubectl get pods -n keptn -lapp.kubernetes.io/name=mongodb -ojsonpath='{.items[0].metadata.name}'")
	require.Nil(t, err)
	_, err = ExecuteCommandf("kubectl cp ./mongodb-backup/ keptn/%s:/opt/dump -c mongodb", removeQuotes(mongoDbPod))
	require.Nil(t, err)

	t.Logf("Logging in to MongoDb database")
	mongoDbRootUser, err = ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.user}")
	require.Nil(t, err)
	mongoDbRootUserByte, err = b64.StdEncoding.DecodeString(removeQuotes(mongoDbRootUser))
	require.Nil(t, err)
	mongoDbRootUser = string(mongoDbRootUserByte)

	mongoDbRootPassword, err = ExecuteCommandf("kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.admin_password}")
	require.Nil(t, err)
	mongoDbRootPasswordByte, err = b64.StdEncoding.DecodeString(removeQuotes(mongoDbRootPassword))
	require.Nil(t, err)
	mongoDbRootPassword = string(mongoDbRootPasswordByte)

	_, err = ExecuteCommandf("kubectl exec svc/mongodb -n keptn -- mongorestore --drop --preserveUUID --authenticationDatabase admin --username admin --password %s /opt/dump", mongoDbRootPassword)
	require.Nil(t, err)

	//restore git credentials
	//t.Logf("Restoring git credentials")
	//_, err = ExecuteCommandf("kubectl apply -f %s-credentials.yaml", keptnProjectName)
	//require.Nil(t, err)

	t.Log("Verify Direct delivery of helloservice in stage dev")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", "gabrieltanner/hello-server", "v0.1.1")
	require.Nil(t, err)

	t.Log("Verify Direct delivery of helloservice in stage staging")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "staging", "gabrieltanner/hello-server", "v0.1.1")
	require.Nil(t, err)

	t.Log("Verify Direct delivery of helloservice in stage production")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "production", "gabrieltanner/hello-server", "v0.1.1")
	require.Nil(t, err)

}
