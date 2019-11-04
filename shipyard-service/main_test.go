package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/magiconair/properties/assert"
)

func TestGetEndpoint(t *testing.T) {
	endPoint, err := getServiceEndpoint("CONFIGURATION_SERVICE")

	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, endPoint.Path, "", "Endpoint has to be empty")

	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	endPoint, err = getServiceEndpoint("CONFIGURATION_SERVICE")

	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, endPoint.Scheme, "http", "Schema of configuration-service endpoint incorrect")
	assert.Equal(t, endPoint.Host, "configuration-service.keptn.svc.cluster.local", "Host of configuration-service endpoint incorrect")
}

// Helper function to build a test client with a httptest server
func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	server := httptest.NewServer(handler)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	return client, server.Close
}

func TestCreateProjectStatusNoContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Expect POST request")
		assert.Equal(t, r.URL.EscapedPath(), "/v1/project", "Expect /v1/project endpoint")
		w.WriteHeader(http.StatusNoContent) // 204 - StatusNoContent
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	client := newClient()
	client.httpClient = httpClient

	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")
	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	project := models.Project{}
	project.ProjectName = "sockshop"
	err := client.createProject(project, *logger)

	assert.Equal(t, err, nil, "Received unexpected error")
}

func TestCreateProjectBadRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Expect POST request")
		assert.Equal(t, r.URL.EscapedPath(), "/v1/project", "Expect /v1/project endpoint")
		w.WriteHeader(http.StatusBadRequest) // 400 - BadRequest
		io.WriteString(w, `{"code": 400, "message": "creating project failed due to error in configuration-service"}`)
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	client := newClient()
	client.httpClient = httpClient

	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")
	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	project := models.Project{}
	project.ProjectName = "sockshop"
	err := client.createProject(project, *logger)

	assert.Equal(t, err.Error(), "creating project failed due to error in configuration-service", "Expect an error")
}

func TestCreateStageStatusNoContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Expect POST request")
		assert.Equal(t, r.URL.EscapedPath(), "/v1/project/sockshop/stage", "Expect /v1/project/sockshop/stage endpoint")
		w.WriteHeader(http.StatusNoContent) // 204 - StatusNoContent
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	client := newClient()
	client.httpClient = httpClient

	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")
	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	project := models.Project{}
	project.ProjectName = "sockshop"
	stage := models.Stage{}
	stage.StageName = "production"
	err := client.createStage(project, stage, *logger)

	assert.Equal(t, err, nil, "Received unexpected error")
}

/* cannot mock the request
func TestStoreResource(t *testing.T) {
	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")
	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	project := models.Project{}
	project.ProjectName = "sockshop"

	var resourceURI = "shipyard.yaml"
	resourceContent, _ := json.Marshal([]string{"apple", "peach", "pear"})
	version, err := storeResourceForProject(project.ProjectName, resourceURI, string(resourceContent), *logger)

	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, version.Version, "as923nad", "Version not returned")
}
*/
