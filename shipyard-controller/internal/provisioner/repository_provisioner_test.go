package provisioner

import (
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/keptn/go-utils/pkg/common/testutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
)

func TestProvideRepository(t *testing.T) {

	type args struct {
		project   string
		namespace string
	}

	client := &testutils.HTTPClientMock{}
	provisioner := New("som-url", client)
	tt := []struct {
		args          args
		Body          string
		StatusCode    int
		expResult     *models.ProvisioningData
		expReqPayload map[string]interface{}
		expErr        error
	}{
		{
			args: args{
				project:   "project",
				namespace: "keptn",
			},
			Body:       `{"gitRemoteURL":"http://some-url.com","gitToken":"token","gitUser":"user"}`,
			StatusCode: http.StatusCreated,
			expResult: &models.ProvisioningData{
				GitRemoteURL: "http://some-url.com",
				GitToken:     "token",
				GitUser:      "user",
			},
			expErr: nil,
			expReqPayload: map[string]interface{}{
				"project":   "project",
				"namespace": "keptn",
			},
		}, {
			args: args{
				project:   "project",
				namespace: "keptn",
			},
			Body:       "",
			StatusCode: http.StatusConflict,
			expResult:  nil,
			expErr:     fmt.Errorf(common.UnableProvisionInstance, http.StatusText(http.StatusConflict)),
		},
		{
			args: args{
				project:   "project",
				namespace: "keptn",
			},
			Body:       `invalid body`,
			StatusCode: http.StatusCreated,
			expResult:  nil,
			expErr:     fmt.Errorf(common.UnableUnMarshallProvisioningData, "invalid character 'i' looking for beginning of value"),
		},
	}

	for _, test := range tt {
		var receivedPayload map[string]interface{}
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			payloadBytes, err := ioutil.ReadAll(r.Body)
			require.Nil(t, err)

			err = json.Unmarshal(payloadBytes, &receivedPayload)
			require.Nil(t, err)

			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}
		r, err := provisioner.ProvideRepository(test.args.project, test.args.namespace)

		assert.Equal(t, test.expErr, err)
		assert.Equal(t, test.expResult, r)

		if test.expReqPayload != nil {
			require.Eventually(t, func() bool {
				return receivedPayload != nil
			}, 1*time.Second, 10*time.Millisecond)

			require.Equal(t, test.expReqPayload, receivedPayload)
		}
	}
}

func TestDeleteRepository(t *testing.T) {
	client := &testutils.HTTPClientMock{}
	provisioner := New("http://som-url.com", client)
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
			expErr:     fmt.Errorf(common.UnableProvisionDelete, http.StatusText(http.StatusNotFound)),
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
