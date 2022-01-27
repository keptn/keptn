package cmd

import (
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	mocking = true
	internal.APIProvider = func(baseURL string, authToken string, httpClient ...*http.Client) (*apiutils.APISet, error) {
		return apiutils.New(baseURL, apiutils.WithAuthToken(authToken), apiutils.WithHTTPClient(&http.Client{}))
	}
	code := m.Run()
	os.Exit(code)
}
