package sdk

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"log"
	"net/http"
)

type envConfig struct {
	Port                    int    `envconfig:"RCV_PORT" default:"8080"`
	Path                    string `envconfig:"RCV_PATH" default:"/"`
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
}

func NewHTTPClientFromEnv() cloudevents.Client {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path), cloudevents.WithGetHandlerFunc(healthEndpointHandler))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c
}

func healthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health/live" {
		w.WriteHeader(http.StatusOK)
	} else if r.URL.Path == "/health/ready" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
func NewResourceHandlerFromEnv() *api.ResourceHandler {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return api.NewResourceHandler(env.ConfigurationServiceURL)
}
