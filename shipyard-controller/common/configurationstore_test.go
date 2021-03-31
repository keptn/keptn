package common

import (
	"encoding/json"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigurationStore(t *testing.T) {

	t.Run("TestCreateProject_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateProject(keptnapimodels.Project{})
		assert.Nil(t, err)
	})

	t.Run("TestCreateProject_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateProject(keptnapimodels.Project{})
		assert.NotNil(t, err)
	})

	t.Run("TestUpdateProject_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.UpdateProject(keptnapimodels.Project{})
		assert.Nil(t, err)
	})

	t.Run("TestUpdateProject_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.UpdateProject(keptnapimodels.Project{})
		assert.NotNil(t, err)
	})

	t.Run("TestDelete_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.DeleteProject("my-project")
		assert.Nil(t, err)
	})

	t.Run("TestDeleteProject_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.DeleteProject("my-project")
		assert.NotNil(t, err)
	})

	t.Run("TestCreateState_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateStage("my-project", "my-stage")
		assert.Nil(t, err)
	})

	t.Run("TestCreateState_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateStage("my-project", "my-stage")
		assert.NotNil(t, err)
	})

	t.Run("TestCreateService_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateService("my-project", "my-stage", "my-service")
		assert.Nil(t, err)
	})

	t.Run("TestCreateService_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateService("my-project", "my-stage", "my-service")
		assert.NotNil(t, err)
	})

	t.Run("TestDeleteService_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.DeleteService("my-project", "my-stage", "my-service")
		assert.Nil(t, err)
	})

	t.Run("TestDeleteService_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.DeleteService("my-project", "my-stage", "my-service")
		assert.NotNil(t, err)
	})

	t.Run("TestCreateProjectShipyard_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "{}")
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateProjectShipyard("my-project", nil)
		assert.Nil(t, err)
	})

	t.Run("TestCreateProjectShipyard_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()
		instance := NewGitConfigurationStore(ts.URL)
		err := instance.CreateProjectShipyard("my-project", nil)
		assert.NotNil(t, err)
	})

	t.Run("TestUpdateProjectResource_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "{}")
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		resourceUri := "uri"
		err := instance.UpdateProjectResource("my-project", &keptnapimodels.Resource{
			ResourceContent: "",
			ResourceURI:     &resourceUri,
		})
		assert.Nil(t, err)
	})

	t.Run("TestUpdateProjectResource_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()
		instance := NewGitConfigurationStore(ts.URL)
		resourceUri := "uri"
		err := instance.UpdateProjectResource("my-project", &keptnapimodels.Resource{
			ResourceContent: "",
			ResourceURI:     &resourceUri,
		})
		assert.NotNil(t, err)
	})

	t.Run("TestGetProjectResource_Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			j, _ := json.Marshal(keptnapimodels.Resource{})
			io.WriteString(w, string(j))
		}))
		defer ts.Close()

		instance := NewGitConfigurationStore(ts.URL)
		resource, err := instance.GetProjectResource("my-project", "uri")
		assert.Nil(t, err)
		assert.Equal(t, *resource, keptnapimodels.Resource{})
	})

	t.Run("TestGetProjectResource_APIReturnsInternalServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()
		instance := NewGitConfigurationStore(ts.URL)
		resource, err := instance.GetProjectResource("my-project", "uri")
		assert.NotNil(t, err)
		assert.Nil(t, resource)
	})

}
