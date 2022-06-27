package sdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type FakeKeptn struct {
	TestResourceHandler ResourceHandler
	SentEvents          []models.KeptnContextExtendedCE
	Keptn               *Keptn
}

func (f *FakeKeptn) GetResourceHandler() ResourceHandler {
	if f.TestResourceHandler == nil {
		return &TestResourceHandler{}
	}
	return f.TestResourceHandler
}

func (f *FakeKeptn) NewEvent(event models.KeptnContextExtendedCE) error {
	ctx := context.WithValue(context.TODO(), types.EventSenderKey, controlplane.EventSender(f.fakeSender))
	ctx = context.WithValue(ctx, gracefulShutdownKey, &nopWG{})
	return f.Keptn.OnEvent(ctx, event)
}

func (f *FakeKeptn) AssertNumberOfEventSent(t *testing.T, numOfEvents int) {
	require.Equalf(t, numOfEvents, len(f.SentEvents), "number of events expected: %d got: %d", numOfEvents, len(f.SentEvents))
}

func (f *FakeKeptn) AssertSentEvent(t *testing.T, eventIndex int, assertFn func(ce models.KeptnContextExtendedCE) bool) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	require.True(t, assertFn(f.SentEvents[eventIndex]))
}

func (f *FakeKeptn) AssertSentEventType(t *testing.T, eventIndex int, eventType string) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	require.Equalf(t, eventType, *f.SentEvents[eventIndex].Type, "event type expected: %s got %s", eventType, *f.SentEvents[eventIndex].Type)
}

func (f *FakeKeptn) AssertSentEventStatus(t *testing.T, eventIndex int, status v0_2_0.StatusType) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.SentEvents[eventIndex], &eventData)
	require.Equal(t, status, eventData.Status)
}

func (f *FakeKeptn) AssertSentEventResult(t *testing.T, eventIndex int, result v0_2_0.ResultType) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.SentEvents[eventIndex], &eventData)
	require.Equal(t, result, eventData.Result)
}

func (f *FakeKeptn) SetAutomaticResponse(autoResponse bool) {
	f.Keptn.automaticEventResponse = autoResponse
}
func (f *FakeKeptn) SetResourceHandler(handler ResourceHandler) {
	f.TestResourceHandler = handler
	f.Keptn.resourceHandler = handler
}

func (f *FakeKeptn) SetAPI(api api.KeptnInterface) {
	f.Keptn.api = api
}

func (f *FakeKeptn) AddTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.AddTaskHandlerWithSubscriptionID(eventType, handler, "", filters...)
}

func (f *FakeKeptn) AddTaskHandlerWithSubscriptionID(eventType string, handler TaskHandler, subscriptionID string, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.Keptn.taskRegistry.Add(eventType, taskEntry{taskHandler: handler, eventFilters: filters})
}

func (f *FakeKeptn) fakeSender(ce models.KeptnContextExtendedCE) error {
	f.SentEvents = append(f.SentEvents, ce)
	return nil
}

func NewFakeKeptn(source string) *FakeKeptn {
	resourceHandler := &TestResourceHandler{}
	var fakeKeptn = &FakeKeptn{
		TestResourceHandler: resourceHandler,
		Keptn: &Keptn{
			source:                 source,
			api:                    panicKeptnInterface{},
			resourceHandler:        resourceHandler,
			taskRegistry:           newTaskMap(),
			syncProcessing:         true,
			automaticEventResponse: true,
			gracefulShutdown:       false,
			logger:                 newDefaultLogger(),
			healthEndpointRunner:   noOpHealthEndpointRunner,
		},
	}
	fakeKeptn.Keptn.eventSender = fakeKeptn.fakeSender
	return fakeKeptn
}

type TestResourceHandler struct {
	Resource models.Resource
}

func (t TestResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s%s%s%s", scope.GetProjectPath(), scope.GetStagePath(), scope.GetServicePath(), scope.GetResourcePath())), nil
}

func newResourceFromFile(filename string) *models.Resource {
	filename = filepath.Clean(filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}

type StringResourceHandler struct {
	ResourceContent string
}

func (s StringResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        &models.Version{Version: "CommitID"},
		ResourceContent: s.ResourceContent,
		ResourceURI:     nil,
	}, nil
}

type FailingResourceHandler struct {
}

func (f FailingResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}

type panicKeptnInterface struct {
}

func (p panicKeptnInterface) AuthV1() api.AuthV1Interface {

	panic("No implementation of AuthV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) EventsV1() api.EventsV1Interface {

	panic("No implementation of EventsV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) LogsV1() api.LogsV1Interface {

	panic("No implementation of LogsV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) ProjectsV1() api.ProjectsV1Interface {

	panic("No implementation of ProjectsV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) ResourcesV1() api.ResourcesV1Interface {

	panic("No implementation of  ResourcesV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) SecretsV1() api.SecretsV1Interface {

	panic("No implementation of SecretsV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) SequencesV1() api.SequencesV1Interface {

	panic("No implementation of SequencesV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) ServicesV1() api.ServicesV1Interface {

	panic("No implementation of ServicesV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) StagesV1() api.StagesV1Interface {

	panic("No implementation of StagesV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) UniformV1() api.UniformV1Interface {

	panic("No implementation of UniformV1found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) ShipyardControlV1() api.ShipyardControlV1Interface {

	panic("No implementation of ShipyardControlV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}

func (p panicKeptnInterface) APIV1() api.APIV1Interface {
	panic("No implementation of APIV1 found. Please Provide a mocked implementation of KeptnInterface for Fake Keptn")
}
