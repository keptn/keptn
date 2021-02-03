package cmd

import (
	"bytes"
	"encoding/json"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mattn/go-shellwords"
)

const unexpectedErrMsg = "unexpected error, got '%v'"

func executeActionCommandC(cmd string) (string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()
	return buf.String(), err
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
		wantOutput       string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := newRedirector()
			r.redirectStdErr()

			returnedMetadata = tt.metadataResponse
			returnedMetadataStatus = tt.metadataStatus
			Version = tt.cliVersion

			runVersionCheck()

			out := r.revertStdErr()
			if out != tt.wantOutput {
				t.Errorf("unexpected output: '%s', expected '%s'", out, tt.wantOutput)
			}
		})
	}
}
