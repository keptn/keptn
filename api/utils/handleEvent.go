// This file is safe to edit. Once it exists it will not be overwritten

package utils

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
)

func getEventBrokerURL() string {
	return "http://" + os.Getenv("EVENTBROKER_URI")
}

const timeout = 60

const MockSend = false

// PostToEventBroker makes a post request to the eventbroker
func PostToEventBroker(event cloudevents.Event) (*cloudevents.Event, error) {

	ec := event.Context.AsV02()
	if ec.Time == nil || ec.Time.IsZero() {
		ec.Time = &types.Timestamp{Time: time.Now()}
		event.Context = ec
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(getEventBrokerURL()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	t.Client = &http.Client{Timeout: timeout * time.Second}

	if err != nil {
		return nil, err
	}

	c, err := client.New(t)
	if err != nil {
		return nil, err
	}
	if MockSend {
		return nil, nil
	}
	return c.Send(context.Background(), event)
}
