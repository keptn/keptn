package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// testShipyard writes a default shipyard file or uses the value from the shipyard parameter.
// It returns a function to delete the shipyard file.
func testShipyard(t *testing.T, shipyardFilePath string, shipyard string) func() {
	if shipyard == "" {
		shipyard = `stages:
  - name: dev
    deployment_strategy: direct
  - name: staging
    deployment_strategy: blue_green_service
  - name: production
    deployment_strategy: blue_green_service`
	}

	ioutil.WriteFile(shipyardFilePath, []byte(shipyard), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	return func() {
		os.Remove(shipyardFilePath)
	}
}

// TestCreateProjectCmd tests the default use of the create project command
func TestCreateProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --mock", shipyardFilePath)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestCreateProjectCmdWithGitMissingParam tests whether the create project command aborts
// due to a missing flag for defining a git upstream
func TestCreateProjectCmdWithGitMissingParam(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --git-user=%s --git-token=%s --mock",
		shipyardFilePath, "user", "token")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, gitErrMsg) {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestCreateProjectCmdWithGitMissingParam tests a successful create project
// command with git upstream parameters
func TestCreateProjectCmdWithGit(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --git-user=%s --git-token=%s --git-remote-url=%s --mock",
		shipyardFilePath, "user", "token", "https://")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func Test_getAndParseYaml(t *testing.T) {

	testShipyard := `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: shipyard-sockshop
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: evaluation`

	var returnedShipyardContent string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(returnedShipyardContent))
	}))
	type args struct {
		arg string
		out interface{}
	}
	tests := []struct {
		name                          string
		args                          args
		shipyardContentFromHTTPServer string
		expectedShipyard              *keptnv2.Shipyard
		wantErr                       bool
	}{
		{
			name: "",
			args: args{
				arg: server.URL,
				out: &keptnv2.Shipyard{},
			},
			shipyardContentFromHTTPServer: testShipyard,
			expectedShipyard: &keptnv2.Shipyard{
				ApiVersion: "spec.keptn.sh/0.2.0",
				Kind:       "Shipyard",
				Metadata: keptnv2.Metadata{
					Name: "shipyard-sockshop",
				},
				Spec: keptnv2.ShipyardSpec{
					Stages: []keptnv2.Stage{
						{
							Name: "dev",
							Sequences: []keptnv2.Sequence{
								{
									Name: "artifact-delivery",
									Tasks: []keptnv2.Task{
										{
											Name:       "evaluation",
											Properties: nil,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedShipyardContent = tt.shipyardContentFromHTTPServer
			shipyard, err := retrieveShipyard(tt.args.arg)
			assert.Equal(t, string(shipyard), testShipyard)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

// TestCreateProjectUnknownCommand
func TestCreateProjectUnknownCommand(t *testing.T) {
	testInvalidInputHelper("create project sockshop someUnknownCommand --shipyard=shipyard.yaml", "too many arguments set", t)
}

// TestCreateProjectUnknownParameter
func TestCreateProjectUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("create project sockshop --projectt=sockshop", "unknown flag: --projectt", t)
}

// TestCreateProjectCmdTokenAndKey
func TestCreateProjectCmdTokenAndKey(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --git-user=user--git-remote-url=https://someurl.com --git-private-key=key --git-token=token", shipyardFilePath)
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, gitErrMsg) {
		t.Errorf("missing expected error, but got %v", err)
	}
}
