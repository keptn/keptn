package v2

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// InternalAPISet is an implementation of KeptnInterface
// which can be used from within the Keptn control plane
type InternalAPISet struct {
	apimap                 InClusterAPIMappings
	httpClient             *http.Client
	apiHandler             *InternalAPIHandler
	authHandler            *v2.AuthHandler
	eventHandler           *v2.EventHandler
	logHandler             *v2.LogHandler
	projectHandler         *v2.ProjectHandler
	resourceHandler        *v2.ResourceHandler
	secretHandler          *v2.SecretHandler
	sequenceControlHandler *v2.SequenceControlHandler
	serviceHandler         *v2.ServiceHandler
	stageHandler           *v2.StageHandler
	uniformHandler         *v2.UniformHandler
	shipyardControlHandler *v2.ShipyardControllerHandler
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
		shipyardControllerApiHandler: v2.NewAPIHandlerWithHTTPClient(
			apimap[ShipyardController],
			&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))}),
	}

	as.authHandler = v2.NewAuthHandlerWithHTTPClient(
		apimap[ApiService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.logHandler = v2.NewLogHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: getClientTransport(as.httpClient.Transport)})

	as.eventHandler = v2.NewEventHandlerWithHTTPClient(
		apimap[MongoDBDatastore],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.projectHandler = v2.NewProjectHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.resourceHandler = v2.NewResourceHandlerWithHTTPClient(
		apimap[ConfigurationService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.secretHandler = v2.NewSecretHandlerWithHTTPClient(
		apimap[SecretService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.sequenceControlHandler = v2.NewSequenceControlHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.serviceHandler = v2.NewServiceHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.shipyardControlHandler = v2.NewShipyardControllerHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.stageHandler = v2.NewStageHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: otelhttp.NewTransport(as.httpClient.Transport)})

	as.uniformHandler = v2.NewUniformHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: getClientTransport(as.httpClient.Transport)})

	return as, nil
}

// API retrieves the APIHandler
func (c *InternalAPISet) API() v2.APIInterface {
	return c.apiHandler
}

// Auth retrieves the AuthHandler
func (c *InternalAPISet) Auth() v2.AuthInterface {
	return c.authHandler
}

// Events retrieves the EventHandler
func (c *InternalAPISet) Events() v2.EventsInterface {
	return c.eventHandler
}

// Logs retrieves the LogHandler
func (c *InternalAPISet) Logs() v2.LogsInterface {
	return c.logHandler
}

// Projects retrieves the ProjectHandler
func (c *InternalAPISet) Projects() v2.ProjectsInterface {
	return c.projectHandler
}

// Resources retrieves the ResourceHandler
func (c *InternalAPISet) Resources() v2.ResourcesInterface {
	return c.resourceHandler
}

// Secrets retrieves the SecretHandler
func (c *InternalAPISet) Secrets() v2.SecretsInterface {
	return c.secretHandler
}

// Sequences retrieves the SequenceControlHandler
func (c *InternalAPISet) Sequences() v2.SequencesInterface {
	return c.sequenceControlHandler
}

// Services retrieves the ServiceHandler
func (c *InternalAPISet) Services() v2.ServicesInterface {
	return c.serviceHandler
}

// Stages retrieves the StageHandler
func (c *InternalAPISet) Stages() v2.StagesInterface {
	return c.stageHandler
}

// Uniform retrieves the UniformHandler
func (c *InternalAPISet) Uniform() v2.UniformInterface {
	return c.uniformHandler
}

// ShipyardControl retrieves the ShipyardControllerHandler
func (c *InternalAPISet) ShipyardControl() v2.ShipyardControlInterface {
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
	shipyardControllerApiHandler *v2.APIHandler
}

func (i *InternalAPIHandler) SendEvent(_ context.Context, _ models.KeptnContextExtendedCE, _ v2.APISendEventOptions) (*models.EventContext, *models.Error) {
	panic("SendEvent() is not not supported for internal usage")
}

func (i *InternalAPIHandler) TriggerEvaluation(ctx context.Context, project string, stage string, service string, evaluation models.Evaluation, opts v2.APITriggerEvaluationOptions) (*models.EventContext, *models.Error) {
	return i.shipyardControllerApiHandler.TriggerEvaluation(ctx, project, stage, service, evaluation, opts)
}

func (i *InternalAPIHandler) CreateProject(ctx context.Context, project models.CreateProject, opts v2.APICreateProjectOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateProject(ctx, project, opts)
}

func (i *InternalAPIHandler) UpdateProject(ctx context.Context, project models.CreateProject, opts v2.APIUpdateProjectOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.UpdateProject(ctx, project, opts)
}

func (i *InternalAPIHandler) DeleteProject(ctx context.Context, project models.Project, opts v2.APIDeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteProject(ctx, project, opts)
}

func (i *InternalAPIHandler) CreateService(ctx context.Context, project string, service models.CreateService, opts v2.APICreateServiceOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateService(ctx, project, service, opts)
}

func (i *InternalAPIHandler) DeleteService(ctx context.Context, project string, service string, opts v2.APIDeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteService(ctx, project, service, opts)
}

func (i *InternalAPIHandler) GetMetadata(_ context.Context, _ v2.APIGetMetadataOptions) (*models.Metadata, *models.Error) {
	panic("GetMetadata() is not not supported for internal usage")
}
