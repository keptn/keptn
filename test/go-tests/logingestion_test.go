package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_LogIngestion(t *testing.T) {
	myLogID := uuid.New().String()
	myErrorLogs := models.CreateLogsRequest{Logs: []models.LogEntry{
		{
			IntegrationID: myLogID,
			Message:       "an error happened",
		},
		{
			IntegrationID: myLogID,
			Message:       "another error happened",
		},
		{
			IntegrationID: myLogID,
			Message:       "yet another error happened",
		},
	}}

	// store our error logs via the API
	resp, err := ApiPOSTRequest("/controlPlane/v1/log", myErrorLogs)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// retrieve the error logs
	resp, err = ApiGETRequest(fmt.Sprintf("/controlPlane/v1/log?integrationId=%s", myLogID))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	getLogsResponse := &models.GetLogsResponse{}
	err = resp.ToJSON(getLogsResponse)

	require.Nil(t, err)
	require.Len(t, getLogsResponse.Logs, 3)
}
