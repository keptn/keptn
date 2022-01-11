package cmd

import (
	"fmt"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("-----------------------------")
	mocking = true
	internal.APIProvider = func(baseURL string, authToken string, authHeader string, scheme string) (*apiutils.ApiSet, error) {
		return apiutils.NewApiSet(baseURL, authToken, authHeader, &http.Client{}, scheme)
	}
	code := m.Run()
	os.Exit(code)
}
