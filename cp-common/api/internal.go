package api

import (
	"crypto/tls"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"time"
)

// InternalAPISet is an implementation of APISet
// which can be used from within the Keptn control plane
type InternalAPISet struct {
	apimap                 InClusterAPIMappings
	httpClient             *http.Client
	apiHandler             *InternalAPIHandler
	authHandler            *api.AuthHandler
	eventHandler           *api.EventHandler
	logHandler             *api.LogHandler
	projectHandler         *api.ProjectHandler
	resourceHandler        *api.ResourceHandler
	secretHandler          *api.SecretHandler
	sequenceControlHandler *api.SequenceControlHandler
	serviceHandler         *api.ServiceHandler
	stageHandler           *api.StageHandler
	uniformHandler         *api.UniformHandler
	shipyardControlHandler *api.ShipyardControllerHandler
}

// InternalService is used to enumerate internal Keptn services
type InternalService int

const (
	ConfigurationService InternalService = iota
	ShipyardController
	ApiService
	SecretService
	MongoDBDatastore
)

// InClusterAPIMappings maps a keptn service name to its reachable domain name
type InClusterAPIMappings map[InternalService]string

// DefaultInClusterAPIMappings gives you the default InClusterAPIMappings
var DefaultInClusterAPIMappings = InClusterAPIMappings{
	ConfigurationService: "configuration-service:8080",
	ShipyardController:   "shipyard-controller:8080",
	ApiService:           "api-service:8080",
	SecretService:        "secret-service:8080",
	MongoDBDatastore:     "mongodb-datastore:8080",
}

// NewInternal creates a new InternalAPISet usable for calling keptn services from within the control plane
func NewInternal(client *http.Client, apiMappings ...InClusterAPIMappings) (*InternalAPISet, error) {
	var apimap InClusterAPIMappings
	if len(apiMappings) > 0 {
		apimap = apiMappings[0]
	} else {
		apimap = DefaultInClusterAPIMappings
	}

	if client == nil {
		client = &http.Client{}
	}

	as := &InternalAPISet{}
	as.httpClient = client

	as.apiHandler = &InternalAPIHandler{
		shipyardControllerApiHandler: &api.APIHandler{
			BaseURL:    apimap[ShipyardController],
			HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
			Scheme:     "http",
		},
	}

	as.authHandler = &api.AuthHandler{
		BaseURL:    apimap[ApiService],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.logHandler = &api.LogHandler{
		BaseURL:      apimap[ShipyardController],
		HTTPClient:   &http.Client{Transport: getClientTransport(as.httpClient.Transport)},
		Scheme:       "http",
		LogCache:     []models.LogEntry{},
		TheClock:     clock.New(),
		SyncInterval: 1 * time.Minute,
	}

	as.eventHandler = &api.EventHandler{
		BaseURL:    apimap[MongoDBDatastore],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}

	as.projectHandler = &api.ProjectHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}

	as.resourceHandler = &api.ResourceHandler{
		BaseURL:    apimap[ConfigurationService],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.secretHandler = &api.SecretHandler{
		BaseURL:    apimap[SecretService],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.sequenceControlHandler = &api.SequenceControlHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.serviceHandler = &api.ServiceHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.shipyardControlHandler = &api.ShipyardControllerHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))},
		Scheme:     "http",
	}
	as.stageHandler = &api.StageHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: otelhttp.NewTransport(as.httpClient.Transport)},
		Scheme:     "http",
	}
	as.uniformHandler = &api.UniformHandler{
		BaseURL:    apimap[ShipyardController],
		HTTPClient: &http.Client{Transport: getClientTransport(as.httpClient.Transport)},
		Scheme:     "http",
	}
	return as, nil
}

// APIV1 retrieves the APIHandler
func (c *InternalAPISet) APIV1() api.APIV1Interface {
	return c.apiHandler
}

// AuthV1 retrieves the AuthHandler
func (c *InternalAPISet) AuthV1() api.AuthV1Interface {
	return c.authHandler
}

// EventsV1 retrieves the EventHandler
func (c *InternalAPISet) EventsV1() api.EventsV1Interface {
	return c.eventHandler
}

// LogsV1 retrieves the LogHandler
func (c *InternalAPISet) LogsV1() api.LogsV1Interface {
	return c.logHandler
}

// ProjectsV1 retrieves the ProjectHandler
func (c *InternalAPISet) ProjectsV1() api.ProjectsV1Interface {
	return c.projectHandler
}

// ResourcesV1 retrieves the ResourceHandler
func (c *InternalAPISet) ResourcesV1() api.ResourcesV1Interface {
	return c.resourceHandler
}

// SecretsV1 retrieves the SecretHandler
func (c *InternalAPISet) SecretsV1() api.SecretsV1Interface {
	return c.secretHandler
}

// SequencesV1 retrieves the SequenceControlHandler
func (c *InternalAPISet) SequencesV1() api.SequencesV1Interface {
	return c.sequenceControlHandler
}

// ServicesV1 retrieves the ServiceHandler
func (c *InternalAPISet) ServicesV1() api.ServicesV1Interface {
	return c.serviceHandler
}

// StagesV1 retrieves the StageHandler
func (c *InternalAPISet) StagesV1() api.StagesV1Interface {
	return c.stageHandler
}

// UniformV1 retrieves the UniformHandler
func (c *InternalAPISet) UniformV1() api.UniformV1Interface {
	return c.uniformHandler
}

// ShipyardControlV1 retrieves the ShipyardControllerHandler
func (c *InternalAPISet) ShipyardControlV1() api.ShipyardControlV1Interface {
	return c.shipyardControlHandler
}

func wrapOtelTransport(base http.RoundTripper) *otelhttp.Transport {
	return otelhttp.NewTransport(base)
}

func getClientTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return tr
	}
	if tr, isDefaultTransport := rt.(*http.Transport); isDefaultTransport {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.Proxy = http.ProxyFromEnvironment
		return tr
	}
	return rt

}

// InternalAPIHandler is used instead of APIHandler from go-utils because we cannot support
// (unauthenticated) internal calls to the api-service at the moment. So this implementation
// will panic as soon as a client wants to call these methods
type InternalAPIHandler struct {
	shipyardControllerApiHandler *api.APIHandler
}

func (i *InternalAPIHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	panic("SendEvent() is not not supported for internal usage")
}

func (i *InternalAPIHandler) TriggerEvaluation(project string, stage string, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error) {
	return i.shipyardControllerApiHandler.TriggerEvaluation(project, stage, service, evaluation)
}

func (i *InternalAPIHandler) CreateProject(project models.CreateProject) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateProject(project)
}

func (i *InternalAPIHandler) UpdateProject(project models.CreateProject) (string, *models.Error) {
	return i.shipyardControllerApiHandler.UpdateProject(project)
}

func (i *InternalAPIHandler) DeleteProject(project models.Project) (*models.DeleteProjectResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteProject(project)
}

func (i *InternalAPIHandler) CreateService(project string, service models.CreateService) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateService(project, service)
}

func (i *InternalAPIHandler) DeleteService(project string, service string) (*models.DeleteServiceResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteService(project, service)
}

func (i *InternalAPIHandler) GetMetadata() (*models.Metadata, *models.Error) {
	panic("GetMetadata() is not not supported for internal usage")
}
