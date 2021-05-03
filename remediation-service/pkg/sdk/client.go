package sdk

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"log"
)

func GetHTTPClient(options ...http.Option) cloudevents.Client {
	p, err := cloudevents.NewHTTP(options...)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, _ := cloudevents.NewClient(p)
	return c
}
