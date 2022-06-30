package execute

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/executor_mock.go . KeptnEndpointProvider:KeptnEndpointProviderMock

type httpdoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type endpointHandler interface {
	ExecuteAPI(doer httpdoer, ate model.APITaskExecution) (any, error)
}

var /*const*/ ErrEndpointNotDefined = errors.New("endpoint not defined")

type KeptnAPIExecutor struct {
	doer             httpdoer
	endpointMappings map[string]endpointHandler
}

type KeptnEndpointProvider interface {
	GetControlPlaneEndpoint() string
}

func NewKeptnExecutor(kep KeptnEndpointProvider) *KeptnAPIExecutor {
	return newKeptnExecutor(kep, nil)
}

func newKeptnExecutor(kep KeptnEndpointProvider, hd httpdoer) *KeptnAPIExecutor {
	if hd == nil {
		hd = new(otelWrappedHttpClient)
	}

	instance := new(KeptnAPIExecutor)
	instance.doer = hd
	instance.endpointMappings = make(map[string]endpointHandler)
	instance.registerEndpoints(kep)
	return instance
}

func (kae *KeptnAPIExecutor) registerEndpoints(kep KeptnEndpointProvider) {
	kae.endpointMappings["keptn-api-v1-create-service"] = &defaultEndpointHandler{
		requestFactory: &projectRenderRequestFactory{
			httpMethod: http.MethodPost,
			path:       `/v1/project/[[project]]/service`,
		},
		endpoint: kep.GetControlPlaneEndpoint(),
	}
}

func (kae *KeptnAPIExecutor) ExecuteAPI(ate model.APITaskExecution) (any, error) {

	endpointHandler, ok := kae.endpointMappings[ate.EndpointID]
	if !ok {
		return nil, fmt.Errorf("error executing api call for endpoint %s: %w", ate.EndpointID, ErrEndpointNotDefined)
	}

	return endpointHandler.ExecuteAPI(kae.doer, ate)
}
