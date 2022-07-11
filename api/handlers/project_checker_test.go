package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	handlers_mock "github.com/keptn/keptn/api/handlers/fake"
)

const projectEndpoint = "/v1/project/"

const errorResponseBodyTemplate = `
{
	"code": %d,
	"message": "%s"
}
`

const projectResponseBodyTemplate = `
{
  "creationDate": "1657029454223822431",
  "projectName": "%s",
  "shipyard": "apiVersion: \"spec.keptn.sh/0.2.2\"\nkind: \"Shipyard\"\nmetadata:\n  name: \"shipyard-jes-sample\"\nspec:\n  stages:\n    - name: \"production\"\n      sequences:\n        - name: \"example-seq\"\n          tasks:\n            - name: \"remote-task\"",
  "shipyardVersion": "spec.keptn.sh/0.2.2",
  "stages": [
    {
      "services": [
        {
          "creationDate": "1657030663780124926",
          "openRemediations": null,
          "serviceName": "foobar"
        },
        {
          "creationDate": "1657206127855300115",
          "openRemediations": null,
          "serviceName": "test-service"
        }
      ],
      "stageName": "dev"
    },
    {
      "services": [
        {
          "creationDate": "1657030663780124926",
          "openRemediations": null,
          "serviceName": "foobar"
        },
        {
          "creationDate": "1657206127855300115",
          "openRemediations": null,
          "serviceName": "test-service"
        }
      ],
      "stageName": "test"
    },
    {
      "services": [
        {
          "creationDate": "1657030663780124926",
          "openRemediations": null,
          "serviceName": "foobar"
        },
        {
          "creationDate": "1657206127855300115",
          "openRemediations": null,
          "serviceName": "test-service"
        }
      ],
      "stageName": "production"
    }
  ],
  "gitCredentials": {
    "remoteURL": "http://gitea-http.gitea:3000/gitea_admin/test-import.git",
    "user": "gitea_admin",
    "https": {
      "insecureSkipTLS": false
    }
  }
}
`

type fakeProjectEndPoint struct {
	f func(w http.ResponseWriter, r *http.Request)
}

func (fpe *fakeProjectEndPoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, projectEndpoint) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fpe.f(w, r)
}

func TestProjectNotFound(t *testing.T) {
	handler := fakeProjectEndPoint{
		f: func(w http.ResponseWriter, r *http.Request) {
			projectRequested := strings.TrimPrefix(r.URL.Path, projectEndpoint)
			body := fmt.Sprintf(
				errorResponseBodyTemplate, http.StatusNotFound, fmt.Sprintf(
					"Project not found: %s",
					projectRequested,
				),
			)

			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(body))
		},
	}

	testServer := httptest.NewServer(&handler)
	defer testServer.Close()

	epm := &handlers_mock.EndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return testServer.URL
		},
	}

	checker := NewControlPlaneProjectChecker(epm)

	exists, err := checker.ProjectExists("foobar")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestProjectFound(t *testing.T) {
	handler := fakeProjectEndPoint{
		f: func(w http.ResponseWriter, r *http.Request) {
			projectRequested := strings.TrimPrefix(r.URL.Path, projectEndpoint)
			body := fmt.Sprintf(
				projectResponseBodyTemplate, projectRequested,
			)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(body))
		},
	}

	testServer := httptest.NewServer(&handler)
	defer testServer.Close()

	epm := &handlers_mock.EndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return testServer.URL
		},
	}

	checker := NewControlPlaneProjectChecker(epm)

	exists, err := checker.ProjectExists("foobar-exists")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestErrorCheckingProject(t *testing.T) {
	handler := fakeProjectEndPoint{
		f: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	}

	testServer := httptest.NewServer(&handler)
	testServer.Close()

	epm := &handlers_mock.EndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return testServer.URL
		},
	}

	checker := NewControlPlaneProjectChecker(epm)

	exists, err := checker.ProjectExists("foobar-exists")
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestGetStages(t *testing.T) {
	handler := fakeProjectEndPoint{
		f: func(w http.ResponseWriter, r *http.Request) {
			projectRequested := strings.TrimPrefix(r.URL.Path, projectEndpoint)
			body := fmt.Sprintf(
				projectResponseBodyTemplate, projectRequested,
			)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(body))
		},
	}

	testServer := httptest.NewServer(&handler)
	defer testServer.Close()

	epm := &handlers_mock.EndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return testServer.URL
		},
	}

	checker := NewControlPlaneProjectChecker(epm)

	stages, err := checker.GetStages("testproject")
	assert.NoError(t, err)
	assert.Equal(t, []string{"dev", "test", "production"}, stages)
}

func TestErrorGetStages(t *testing.T) {
	handler := fakeProjectEndPoint{
		f: func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("Handler should not have been called with a closed server")
		},
	}

	testServer := httptest.NewServer(&handler)
	testServer.Close()

	epm := &handlers_mock.EndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return testServer.URL
		},
	}

	checker := NewControlPlaneProjectChecker(epm)

	stages, err := checker.GetStages("testproject")
	assert.Error(t, err)
	assert.Nil(t, stages)
}
