package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cli/pkg/version"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

const unexpectedErrMsg = "unexpected error, got '%v'"

const keptnVersionResponse = `{
    "cli": {
        "stable": [ "0.7.0", "0.7.1", "0.7.2", "0.7.3", "0.8.0", "0.8.1"],
        "prerelease": [ ]
    }, 
    "bridge": {
        "stable": [ "0.7.0", "0.7.1", "0.7.2", "0.7.3", "0.8.0", "0.8.1"],
        "prerelease": [ ]
    },
    "keptn": {
        "stable": [
            {
              "version": "0.8.1",
              "upgradableVersions": [ "0.8.0" ]
            },
            {
              "version": "0.8.0",
              "upgradableVersions": [ "0.7.1", "0.7.2", "0.7.3" ]
            },
            {
              "version": "0.7.3",
              "upgradableVersions": [ "0.7.0", "0.7.1", "0.7.2" ]
            },
            {
              "version": "0.7.2",
              "upgradableVersions": [ "0.7.0", "0.7.1" ]
            },
            {
              "version": "0.7.1",
              "upgradableVersions": [ "0.7.0" ]
            }
        ]
    }
}`

func executeActionCommandC(cmd string) (string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)

	ts := getMockVersionHTTPServer()

	defer ts.Close()

	vChecker := &version.VersionChecker{
		VersionFetcherClient: &version.VersionFetcherClient{
			HTTPClient: http.DefaultClient,
			VersionURL: ts.URL,
		},
	}
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		runVersionCheck(vChecker, os.Args[1:])
	}

	rootCmd.SetOut(buf)
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	return buf.String(), err
}

func getMockVersionHTTPServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(200)
		w.Write([]byte(keptnVersionResponse))
	}))
	return ts
}

type redirector struct {
	originalStdOut *os.File
	originalStdErr *os.File
	r              *os.File
	w              *os.File
}

func newRedirector() *redirector {
	return &redirector{}
}

func (r *redirector) redirectStdOut() {
	r.originalStdOut = os.Stdout
	r.r, r.w, _ = os.Pipe()
	os.Stdout = r.w
}

func (r *redirector) redirectStdErr() {
	r.originalStdErr = os.Stderr
	r.r, r.w, _ = os.Pipe()
	os.Stderr = r.w
}

func (r *redirector) revertStdOut() string {

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r.r)
		outC <- buf.String()
	}()

	// back to normal state
	r.w.Close()
	os.Stdout = r.originalStdOut
	return <-outC
}

func (r *redirector) revertStdErr() string {

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r.r)
		outC <- buf.String()
	}()

	// back to normal state
	r.w.Close()
	os.Stderr = r.originalStdErr
	return <-outC
}

func Test_runVersionCheck(t *testing.T) {
	mocking = true
	var returnedMetadataStatus int
	var returnedMetadata keptnapimodels.Metadata

	// var returnedVersionStatus int

	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// mock metadata endpoint
		if strings.Contains(request.URL.String(), "/metadata") {
			writer.WriteHeader(returnedMetadataStatus)
			marshal, _ := json.Marshal(&returnedMetadata)
			writer.Write(marshal)
			return
		}
		// mock version endpoint
		//if strings.Contains(request.URL.String(), "/version") {
		//	writer.WriteHeader(returnedVersionStatus)
		//	marshal, _ := json.Marshal(&returnedMetadata)
		//	writer.Write(marshal)
		//	return
		//}
	}))

	defer ts.Close()
	os.Setenv("MOCK_SERVER", ts.URL)

	tests := []struct {
		name             string
		metadataStatus   int
		metadataResponse keptnapimodels.Metadata
		cliVersion       string
		// args excludes the main `keptn` command
		// e.g., keptn -q install, args would be ['-q', 'install']
		args            []string
		wantOutput      string
		doNotWantOutput string
	}{
		{
			name:           "get version",
			cliVersion:     "0.8.0",
			metadataStatus: http.StatusOK,
			metadataResponse: keptnapimodels.Metadata{
				Keptnversion: "0.8.1-dev",
			},
			wantOutput: "* Warning: Your Keptn CLI version (0.8.0) and Keptn cluster version (0.8.1-dev) don't match. This can lead to problems. Please make sure to use the same versions.\n",
		},
		{
			name:           "received error from Keptn API",
			cliVersion:     "0.8.0",
			metadataStatus: http.StatusInternalServerError,
			wantOutput:     "* Warning: could not check Keptn server version: received invalid response from Keptn API\n",
		},
		{
			name:            "skip version check for 'keptn install'",
			cliVersion:      "0.8.0",
			metadataStatus:  http.StatusServiceUnavailable,
			args:            []string{"install"},
			doNotWantOutput: "* Warning: could not check Keptn server version: Error connecting to server:",
		},
		{
			name:            "skip version check for 'keptn --any-flag install'",
			cliVersion:      "0.8.0",
			metadataStatus:  http.StatusServiceUnavailable,
			args:            []string{"--any-flag", "install"},
			doNotWantOutput: "* Warning: could not check Keptn server version: Error connecting to server:",
		},
		{
			name:           "show version check for 'keptn command-other-than-install'",
			cliVersion:     "0.8.0",
			metadataStatus: http.StatusOK,
			metadataResponse: keptnapimodels.Metadata{
				Keptnversion: "0.8.1-dev",
			},
			args:       []string{"command-other-than-install"},
			wantOutput: "* Warning: Your Keptn CLI version (0.8.0) and Keptn cluster version (0.8.1-dev) don't match.",
		},
		{
			name:           "show version check for 'keptn --any-flag command-other-than-install'",
			cliVersion:     "0.8.0",
			metadataStatus: http.StatusOK,
			metadataResponse: keptnapimodels.Metadata{
				Keptnversion: "0.8.1-dev",
			},
			args:       []string{"--any-flag", "command-other-than-install"},
			wantOutput: "* Warning: Your Keptn CLI version (0.8.0) and Keptn cluster version (0.8.1-dev) don't match.",
		},
		{
			name:           "don't show warning if the versions match",
			cliVersion:     "0.8.0",
			metadataStatus: http.StatusOK,
			metadataResponse: keptnapimodels.Metadata{
				Keptnversion: "0.8.0",
			},
			doNotWantOutput: "* Warning: Your Keptn CLI version (0.8.0) and Keptn cluster version (0.8.0) don't match.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := newRedirector()
			r.redirectStdErr()

			returnedMetadata = tt.metadataResponse
			returnedMetadataStatus = tt.metadataStatus
			Version = tt.cliVersion

			ts := getMockVersionHTTPServer()
			defer ts.Close()

			vChecker := &version.VersionChecker{
				VersionFetcherClient: &version.VersionFetcherClient{
					HTTPClient: http.DefaultClient,
					VersionURL: ts.URL,
				},
			}
			runVersionCheck(vChecker, tt.args)

			// reset version
			Version = ""

			out := r.revertStdErr()
			if tt.wantOutput != "" && !strings.Contains(out, tt.wantOutput) {
				t.Errorf("unexpected output: '%s', expected '%s'", out, tt.wantOutput)
			}

			if tt.doNotWantOutput != "" && strings.Contains(out, tt.doNotWantOutput) {
				t.Errorf("unexpected output: '%s', output should not contain '%s'", out, tt.doNotWantOutput)
			}
		})
	}
}
