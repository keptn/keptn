package execute

import (
	"fmt"
	"io"
	"net/http"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/webhook.go . integrationIdRetriever:MockIntegrationIdRetriever

type integrationIdRetriever interface {
	GetIntegrationIDsByName(name string) ([]string, error)
}

func NewWebhookSubscriptionHandler(retriever integrationIdRetriever) *webhookSubscriptionHandler {
	return &webhookSubscriptionHandler{
		idRetriever: retriever,
	}
}

type webhookSubscriptionHandler struct {
	idRetriever integrationIdRetriever
}

func (w webhookSubscriptionHandler) CreateRequest(
	_ model.TaskContext, host string, body io.Reader,
) (*http.Request, error) {
	const webhookIntegrationName = "webhook-service"
	subscriptionIds, err := w.idRetriever.GetIntegrationIDsByName(webhookIntegrationName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving integration id for name %s: %w", webhookIntegrationName, err)
	}
	if len(subscriptionIds) == 0 {
		return nil, fmt.Errorf("no integration found for name %s", webhookIntegrationName)
	}

	request, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf(
			"%s/v1/uniform/registration/%s/subscription",
			host,
			subscriptionIds[0],
		), body,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	return request, nil
}
