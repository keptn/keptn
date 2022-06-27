package importer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
  "creationDate": "1655467882190850007",
  "gitRemoteURI": "http://somehost.somedomain/somenamespace/somerepo.git",
  "gitUser": "git-user",
  "projectName": "%s",
  "shipyard": "apiVersion: \"spec.keptn.sh/0.2.2\"\nkind: \"Shipyard\"\nmetadata:\n  name: \"shipyard-jes-sample\"\nspec:\n  stages:\n    - name: \"production\"\n      sequences:\n        - name: \"example-seq\"\n          tasks:\n            - name: \"remote-task\"",
  "shipyardVersion": "spec.keptn.sh/0.2.2",
  "insecureSkipTLS": false,
  "stages": [
    {
      "services": [],
      "stageName": "production"
    }
  ]
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

	checker := NewControlPlaneProjectChecker(testServer.URL)

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

	checker := NewControlPlaneProjectChecker(testServer.URL)

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

	checker := NewControlPlaneProjectChecker(testServer.URL)

	exists, err := checker.ProjectExists("foobar-exists")
	assert.Error(t, err)
	assert.False(t, exists)
}
