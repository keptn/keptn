// This file is safe to edit. Once it exists it will not be overwritten

package utils

import (
	"context"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
)

func getEventBrokerURL() string {
	return "http://" + os.Getenv("EVENTBROKER_URI")
}

// PostToEventBroker makes a post request to the eventbroker
func PostToEventBroker(event cloudevents.Event) (*cloudevents.Event, error) {

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(getEventBrokerURL()),
		cloudevents.WithEncoding(cloudevents.HTTPStructuredV02),
	)

	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(t)
	if err != nil {
		return nil, err
	}
	return c.Send(context.Background(), event)
}
