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
	resp, err := ApiPOSTRequest("/controlPlane/v1/log", myErrorLogs, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// retrieve the error logs
	resp, err = ApiGETRequest(fmt.Sprintf("/controlPlane/v1/log?integrationId=%s", myLogID), 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	getLogsResponse := &models.GetLogsResponse{}
	err = resp.ToJSON(getLogsResponse)

	require.Nil(t, err)
	require.Len(t, getLogsResponse.Logs, 3)
	require.Equal(t, int64(3), getLogsResponse.TotalCount)

	// retrieve the error logs - using pagination
	resp, err = ApiGETRequest(fmt.Sprintf("/controlPlane/v1/log?integrationId=%s&pageSize=1", myLogID), 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	getLogsResponse = &models.GetLogsResponse{}
	err = resp.ToJSON(getLogsResponse)

	require.Nil(t, err)
	require.Len(t, getLogsResponse.Logs, 1)
	require.Equal(t, int64(3), getLogsResponse.TotalCount)

	// delete the logs
	resp, err = ApiDELETERequest(fmt.Sprintf("/controlPlane/v1/log?integrationId=%s", myLogID), 3)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// retrieve the error logs again -should not be there anymore
	resp, err = ApiGETRequest(fmt.Sprintf("/controlPlane/v1/log?integrationId=%s", myLogID), 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	getLogsResponse = &models.GetLogsResponse{}
	err = resp.ToJSON(getLogsResponse)

	require.Nil(t, err)
	require.Len(t, getLogsResponse.Logs, 0)
	require.Equal(t, int64(0), getLogsResponse.TotalCount)
}
