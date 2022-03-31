package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/keptn/go-utils/pkg/common/testutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
)

func TestProvideRepository(t *testing.T) {
	client := &testutils.HTTPClientMock{}
	provisioner := NewRepositoryProvisioner("som-url", client)
	tt := []struct {
		Body       string
		StatusCode int
		expResult  *models.ProvisioningData
		expErr     error
	}{
		{
			Body:       `{"gitRemoteURL":"http://some-url.com","gitToken":"token","gitUser":"user"}`,
			StatusCode: http.StatusCreated,
			expResult: &models.ProvisioningData{
				GitRemoteURL: "http://some-url.com",
				GitToken:     "token",
				GitUser:      "user",
			},
			expErr: nil,
		}, {
			Body:       "",
			StatusCode: http.StatusConflict,
			expResult:  nil,
			expErr:     fmt.Errorf(UnableProvisionInstance, http.StatusText(http.StatusConflict)),
		},
		{
			Body:       `invalid body`,
			StatusCode: http.StatusCreated,
			expResult:  nil,
			expErr:     fmt.Errorf(UnableUnMarshallProvisioningData, "invalid character 'i' looking for beginning of value"),
		},
	}

	for _, test := range tt {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}
		r, err := provisioner.ProvideRepository("project")

		assert.Equal(t, test.expErr, err)
		assert.Equal(t, test.expResult, r)
	}
}

func TestDeleteRepository(t *testing.T) {
	client := &testutils.HTTPClientMock{}
	provisioner := NewRepositoryProvisioner("http://som-url.com", client)
	tt := []struct {
		Body       string
		StatusCode int
		expErr     error
	}{
		{
			Body:       "",
			StatusCode: http.StatusOK,
			expErr:     nil,
		}, {
			Body:       "",
			StatusCode: http.StatusNotFound,
			expErr:     fmt.Errorf(UnableProvisionDelete, http.StatusText(http.StatusNotFound)),
		},
	}

	for _, test := range tt {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}
		err := provisioner.DeleteRepository("project", "keptn")

		assert.Equal(t, test.expErr, err)
	}
}
