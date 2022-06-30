package execute

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/keptn/keptn/api/importer/model"
)

var /*const*/ ErrTaskFailed = errors.New("task failed")

type otelWrappedHttpClient struct{}

func (o *otelWrappedHttpClient) Do(r *http.Request) (*http.Response, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(nil),
	}
	return client.Do(r)
}

type requestFactory interface {
	CreateRequest(tCtx model.TaskContext, host string, body io.Reader) (*http.Request, error)
}

type projectRenderRequestFactory struct {
	httpMethod string
	path       string
}

func (rf *projectRenderRequestFactory) renderUrl(
	tCtx model.TaskContext, endpoint string, path string,
) string {
	return endpoint + strings.Replace(path, "[[project]]", tCtx.Project, -1)
}
func (rf *projectRenderRequestFactory) CreateRequest(
	tCtx model.TaskContext, host string, body io.Reader,
) (*http.Request,
	error) {
	req, err := http.NewRequest(
		rf.httpMethod,
		rf.renderUrl(tCtx, host, rf.path),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("error composing request for api call %s: %w", tCtx.Task.ID, err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

type defaultEndpointHandler struct {
	requestFactory
	endpoint string
}

func (ep *defaultEndpointHandler) ExecuteAPI(doer httpdoer, ate model.APITaskExecution) (any, error) {
	request, err := ep.CreateRequest(ate.Context, ep.endpoint, ate.Payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	response, err := doer.Do(request)

	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer response.Body.Close()
	responseBody := new(any)
	if response.Header.Get("Content-Type") == "application/json" {
		bytes, err := io.ReadAll(response.Body)
		if err == nil {
			json.Unmarshal(bytes, responseBody)
		}
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}

	return responseBody, fmt.Errorf(
		"received unsuccessful http status <%d: %s>:%w", response.StatusCode,
		response.Status, ErrTaskFailed,
	)
}
