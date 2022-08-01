package execute

import (
	"errors"
	"fmt"
	"io"
	"k8s.io/utils/strings/slices"
	"net/http"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/executor_mock.go . KeptnEndpointProvider:KeptnEndpointProviderMock httpdoer:MockHTTPDoer resourcePusher:MockResourcePusher

type httpdoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type endpointHandler interface {
	ExecuteAPI(doer httpdoer, ate model.APITaskExecution) (any, error)
}

type resourcePusher interface {
	PushToStage(project string, stage string, content io.ReadCloser, resourceURI string) (any, error)
	PushToService(project string, stage string, service string, content io.ReadCloser, resourceURI string) (any, error)
}

var /*const*/ ErrEndpointNotDefined = errors.New("endpoint not defined")

type KeptnAPIExecutor struct {
	doer             httpdoer
	endpointMappings map[string]endpointHandler
	resourcePusher   resourcePusher
}

type KeptnControlPlaneEndpointProvider interface {
	GetControlPlaneEndpoint() string
}

type KeptnConfigurationServiceEndpointProvider interface {
	GetConfigurationServiceEndpoint() string
}

type KeptnSecretsEndpointProvider interface {
	GetSecretsServiceEndpoint() string
}

type KeptnEndpointProvider interface {
	KeptnControlPlaneEndpointProvider
	KeptnConfigurationServiceEndpointProvider
	KeptnSecretsEndpointProvider
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
	instance.resourcePusher = &KeptnResourcePusher{
		endpointProvider: kep,
		doer:             hd,
	}
	return instance
}

func (kae *KeptnAPIExecutor) registerEndpoints(kep KeptnEndpointProvider) {
	kae.endpointMappings[model.CreateServiceAction] = &defaultEndpointHandler{
		requestFactory: &projectRenderRequestFactory{
			httpMethod: http.MethodPost,
			path:       `/v1/project/[[project]]/service`,
		},
		endpoint: kep.GetControlPlaneEndpoint(),
	}

	kae.endpointMappings[model.CreateSecretAction] = &defaultEndpointHandler{
		requestFactory: &projectRenderRequestFactory{
			httpMethod: http.MethodPost,
			path:       "/v1/secret",
		},
		endpoint: kep.GetSecretsServiceEndpoint(),
	}
}

func (kae *KeptnAPIExecutor) ExecuteAPI(ate model.APITaskExecution) (any, error) {

	endpointHandler, ok := kae.endpointMappings[ate.EndpointID]
	if !ok {
		return nil, fmt.Errorf("error executing api call for endpoint %s: %w", ate.EndpointID, ErrEndpointNotDefined)
	}

	return endpointHandler.ExecuteAPI(kae.doer, ate)
}

func (kae *KeptnAPIExecutor) PushResource(rp model.ResourcePush) (any, error) {
	if rp.Service == "" {
		return kae.resourcePusher.PushToStage(rp.Context.Project, rp.Stage, rp.Content, rp.ResourceURI)
	}
	return kae.resourcePusher.PushToService(rp.Context.Project, rp.Stage, rp.Service, rp.Content, rp.ResourceURI)
}

func (kae *KeptnAPIExecutor) ActionSupported(actionName string) bool {
	if !slices.Contains(model.AllActions, actionName) {
		return false
	}

	return true
}
