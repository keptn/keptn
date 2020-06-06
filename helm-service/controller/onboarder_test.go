package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/helmtest"

	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
)

const configBaseURL = "localhost:6060"
const projectName = "sockshop"
const serviceName = "carts"
const stage1 = "dev"
const stage2 = "staging"
const stage3 = "production"

const shipyard = `project: sockshop
stages:
    - {deployment_strategy: direct, name: dev, test_strategy: functional}
    - {deployment_strategy: blue_green_service, name: staging, test_strategy: performance}
    - {deployment_strategy: blue_green_service, name: production}`

func createTestProject(t *testing.T) {

	prjHandler := configutils.NewProjectHandler(configBaseURL)
	prj := configmodels.Project{ProjectName: projectName}
	respErr, err := prjHandler.CreateProject(prj)
	check(err, t)
	assert.Nil(t, respErr, "Creating a project failed")

	// Send shipyard
	rHandler := configutils.NewResourceHandler(configBaseURL)
	shipyardURI := "shipyard.yaml"
	shipyardResource := configmodels.Resource{ResourceURI: &shipyardURI, ResourceContent: shipyard}
	resources := []*configmodels.Resource{&shipyardResource}
	_, err2 := rHandler.CreateProjectResources(projectName, resources)
	if err2 != nil {
		t.Error(err)
	}

	// Create stages
	stageHandler := configutils.NewStageHandler(configBaseURL)
	for _, stage := range []string{stage1, stage2, stage3} {

		respErr, err := stageHandler.CreateStage(projectName, stage)
		check(err, t)
		assert.Nil(t, respErr, "Creating a stage failed")
	}
}

// TestDoOnboard tests the onboarding of a new chart. Therefore, this test requires the configuration-service
// on localhost:8080
func TestDoOnboard(t *testing.T) {

	_, err := http.Get("http://" + configBaseURL)
	if err != nil {
		t.Skip("Skipping test; no local configuration-service available")
	}
	os.Setenv("CONFIGURATION_SERVICE", configBaseURL)

	createTestProject(t)

	data := helmtest.CreateHelmChartData(t)
	encodedChart := base64.StdEncoding.EncodeToString(data)
	fmt.Println(encodedChart)
	ce := cloudevents.New("0.2")
	dataBytes, err := json.Marshal(keptnevents.ServiceCreateEventData{Project: projectName, Service: serviceName, HelmChart: encodedChart})
	if err != nil {
		t.Error(err)
	}
	ce.Data = dataBytes

	id := uuid.New().String()
	onboarder := NewOnboarder(mesh.NewIstioMesh(),
		keptnutils.NewLogger(id, "service.create", "helm-service"), "test.keptn.sh", "")
	loggingDone := make(chan bool)
	err = onboarder.DoOnboard(ce, loggingDone)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckAndSetServiceName(t *testing.T) {

	createErrorMsg := "Service name contains special character(s). " +
		"The service name has to be a valid Unix directory name. For details see " +
		"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/"

	o := NewOnboarder(nil, nil, "test.keptn.sh", "")
	data := helmtest.CreateHelmChartData(t)

	testCases := []struct {
		name        string
		event       *keptnevents.ServiceCreateEventData
		error       error
		serviceName string
	}{
		{"Mismatch", &keptnevents.ServiceCreateEventData{Service: "carts-1", HelmChart: base64.StdEncoding.EncodeToString(data)},
			errors.New("Provided Keptn service name \"carts-1\" does not match Kubernetes service name \"carts\""), "carts-1"},
		{"Match", &keptnevents.ServiceCreateEventData{Service: "carts", HelmChart: base64.StdEncoding.EncodeToString(data)},
			nil, "carts"},
		{"Set", &keptnevents.ServiceCreateEventData{Service: "", HelmChart: base64.StdEncoding.EncodeToString(data)},
			nil, "carts"},
		{"EmptyName", &keptnevents.ServiceCreateEventData{Service: ""},
			errors.New(createErrorMsg), ""},
		{"ValidName", &keptnevents.ServiceCreateEventData{Service: "a"},
			nil, "a"},
		{"ValidName", &keptnevents.ServiceCreateEventData{Service: "aa"},
			nil, "aa"},
		{"ValidName", &keptnevents.ServiceCreateEventData{Service: "aa7"},
			nil, "aa7"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := o.checkAndSetServiceName(tt.event)
			if res == nil && res != tt.error {
				t.Errorf("got nil, want %s", tt.error.Error())
			} else if res != nil && tt.error != nil && res.Error() != tt.error.Error() {
				t.Errorf("got %s, want %s", res.Error(), tt.error.Error())
			} else if res != nil && tt.error == nil {
				t.Errorf("got %s, want nil", res.Error())
			}

			if tt.event.Service != tt.serviceName {
				t.Errorf("got %s, want %s", tt.event.Service, tt.serviceName)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}

func check(e *configmodels.Error, t *testing.T) {
	if e != nil {
		t.Error(e.Message)
	}
}
