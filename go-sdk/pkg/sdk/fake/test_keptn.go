package fake

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
)

type FakeKeptn struct {
	TestResourceHandler sdk.ResourceHandler
	Keptn               *sdk.Keptn
}

func (f *FakeKeptn) GetResourceHandler() sdk.ResourceHandler {
	if f.TestResourceHandler == nil {
		return &TestResourceHandler{}
	}
	return f.TestResourceHandler
}

func (f *FakeKeptn) GetTaskRegistry() *sdk.TaskRegistry {
	return f.Keptn.GetTaskRegistry()
}

func (f *FakeKeptn) SetConfigurationServiceURL(configurationServiceURL string) {
	panic("implement me")
}

func (f *FakeKeptn) NewEvent(event cloudevents.Event) {
	testReceiver := f.Keptn.EventReceiver.(*TestReceiver)
	testReceiver.NewEvent(event)
}

func (f *FakeKeptn) GetEventSender() *TestSender {
	return f.Keptn.EventSender.(*TestSender)
}

func NewFakeKeptn(source string, opts ...sdk.KeptnOption) *FakeKeptn {
	eventReceiver := &TestReceiver{}
	eventSender := &TestSender{}
	resourceHandler := &TestResourceHandler{}

	var fakeKeptn = &FakeKeptn{
		Keptn: &sdk.Keptn{
			EventSender:     eventSender,
			EventReceiver:   eventReceiver,
			ResourceHandler: resourceHandler,
			Source:          source,
			TaskRegistry:    sdk.NewTasksMap(),
			SyncProcessing:  true,
		},
	}
	for _, opt := range opts {
		opt(fakeKeptn)
	}
	return fakeKeptn
}

func (f *FakeKeptn) Start() error {
	return f.Keptn.Start()
}
