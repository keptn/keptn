// This file is safe to edit. Once it exists it will not be overwritten

package utils

import (
	"context"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
)

func getEventBrokerURL() string {
	uri := os.Getenv("EVENTBROKER_URI")

	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		return uri
	}
	return "http://" + uri
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
	_, sent, err := c.Send(context.Background(), event)
	return sent, err
}
