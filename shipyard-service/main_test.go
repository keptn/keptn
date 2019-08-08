package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/shipyard-service/models"
	"github.com/magiconair/properties/assert"
)

func TestGetEndpoint(t *testing.T) {
	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")

	endPoint, err := getEndpoint(*logger)

	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, endPoint.Path, "", "Endpoint has to be empty")

	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	endPoint, err = getEndpoint(*logger)

	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, endPoint.Path, "http://configuration-service.keptn.svc.cluster.local", "Endpoint not set in environment variable")
}

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

func TestCreateProject(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Expecte POST request")
		fmt.Println(r.URL.EscapedPath())
		assert.Equal(t, r.URL.EscapedPath(), "/project", "Expecte POST request")
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	client := NewClient()
	client.httpClient = httpClient

	logger := keptnutils.NewLogger("4711-a83b-4bc1-9dc0-1f050c7e789b", "4711-a83b-4bc1-9dc0-1f050c7e781b", "shipyard-service")
	os.Setenv("CONFIGURATION_SERVICE", "http://configuration-service.keptn.svc.cluster.local")

	project := models.Project{}
	project.ProjectName = "sockshop"
	err := client.createProject(project, *logger)

	assert.Equal(t, err, nil, "Received unexpected error")
}

func TestCreateStage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Expecte POST request")
		fmt.Println(r.URL.EscapedPath())
		assert.Equal(t, r.URL.EscapedPath(), "/project/sockshop/stage", "Expecte POST request")
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	client := NewClient()
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

func TestStoreResource(t *testing.T) {

}
