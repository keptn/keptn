package go_tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testingShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
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

func TestBackupRestore(t *testing.T) {
	repoLocalDir := "../assets/sockshop/"
	keptnProjectName := "sockshop"
	serviceCarts := "carts"
	serviceCartsDb := "carts-db"
	resourceCarts := "carts.tgz"
	resourceCartsDb := "carts-db.tgz"

	shipyardFilePath, err := CreateTmpShipyardFile(testingShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("Creating a new project %s without a GIT Upstream", keptnProjectName)
	err = CreateProject(keptnProjectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Creating a new service %s in project %s", serviceCarts, keptnProjectName)
	_, err = ExecuteCommand("keptn create service %s --project=%s", serviceCarts, keptnProjectName)
	require.Nil(t, err)

	t.Logf("Adding resource %s to service %s in project %s", resourceCarts, serviceCarts, keptnProjectName)
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --all-stages --resource=%s --resourceUri=%s", keptnProjectName, serviceCarts, repoLocalDir+resourceCarts, "helm/"+resourceCarts)
	require.Nil(t, err)

	t.Logf("Creating a new service %s in project %s", serviceCartsDb, keptnProjectName)
	_, err = ExecuteCommand("keptn create service %s --project=%s", serviceCartsDb, keptnProjectName)
	require.Nil(t, err)

	t.Logf("Adding resource %s to service %s in project %s", resourceCartsDb, serviceCartsDb, keptnProjectName)
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --all-stages --resource=%s --resourceUri=%s", keptnProjectName, serviceCartsDb, repoLocalDir+resourceCartsDb, "helm/"+resourceCartsDb)
	require.Nil(t, err)

	t.Logf("Triggering delivery of service %s in project %s", serviceCarts, keptnProjectName)
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=docker.io/keptnexamples/carts --tag=0.12.1", keptnProjectName, serviceCarts)
	require.Nil(t, err)

	t.Logf("Triggering delivery of service %s in project %s", serviceCartsDb, keptnProjectName)
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=docker.io/mongo --tag=4.2.2", keptnProjectName, serviceCartsDb)
	require.Nil(t, err)
}
