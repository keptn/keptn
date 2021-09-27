package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/lib"
)

const webhookConfigFileName = "webhook/webhook.yaml"

type DistributorData struct {
	SubscriptionID string `json:"subscriptionID"`
}

type TemporaryData struct {
	TemporaryData struct {
		Distributor DistributorData `json:"distributor"`
	} `json:"temporaryData"`
}

type SecretEnv struct {
	Env map[string]string
}

type TaskHandler struct {
	templateEngine lib.ITemplateEngine
	curlExecutor   lib.ICurlExecutor
	secretReader   lib.ISecretReader
}

func NewTaskHandler(templateEngine lib.ITemplateEngine, curlExecutor lib.ICurlExecutor, secretReader lib.ISecretReader) *TaskHandler {
	return &TaskHandler{
		templateEngine: templateEngine,
		curlExecutor:   curlExecutor,
		secretReader:   secretReader,
	}
}

func (th *TaskHandler) Execute(keptnHandler sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	nedc, err := lib.NewEventDataModifier(event)
	if err != nil {
		return nil, sdkError("could not parse incoming event", err)
	}
	webhook, err := th.getWebHookConfig(keptnHandler, nedc)
	if err != nil {
		return nil, sdkError(fmt.Sprintf("could not retrieve Webhook config: %s", err.Error()), err)
	}

	responses := []string{}

	secretEnvVars, sdkErr := th.gatherSecretEnvVars(*webhook)
	if sdkErr != nil {
		return nil, sdkErr
	}
	nedc.Add("env", secretEnvVars)
	responses, sdkErr = th.performWebhookRequests(*webhook, nedc, responses)
	if sdkErr != nil {
		return nil, sdkErr
	}

	// check if the incoming event was a task.triggered event
	// only in this case, the result should be sent back to Keptn in the form of a .finished event
	if keptnv2.IsTaskEventType(*event.Type) {
		taskName, _, err := keptnv2.ParseTaskEventType(*event.Type)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not derive task name from event type %s", *event.Type), err)
		}
		return map[string]interface{}{
			"project": nedc.Project(),
			"stage":   nedc.Stage(),
			"service": nedc.Service(),
			"labels":  nedc.Labels(),
			taskName: map[string]interface{}{
				"responses": responses,
			},
		}, nil
	}

	return nil, nil
}

func (th *TaskHandler) performWebhookRequests(webhook lib.Webhook, nedc *lib.EventDataModifier, responses []string) ([]string, *sdk.Error) {
	for _, req := range webhook.Requests {
		// parse the data from the event, together with the secret env vars
		parsedCurlCommand, err := th.templateEngine.ParseTemplate(nedc.Get(), req)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not parse request '%s'", req), err)
		}
		// perform the request
		response, err := th.curlExecutor.Curl(parsedCurlCommand)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not execute request '%s': %s", req, err.Error()), err)
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (th *TaskHandler) gatherSecretEnvVars(webhook lib.Webhook) (map[string]string, *sdk.Error) {
	secretEnvVars := map[string]string{}
	for _, secretRef := range webhook.EnvFrom {
		secretValue, err := th.secretReader.ReadSecret(secretRef.SecretRef.Name, secretRef.SecretRef.Key)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not read secret %s.%s", secretRef.SecretRef.Name, secretRef.SecretRef.Key), err)
		}
		secretEnvVars[secretRef.Name] = secretValue
	}
	return secretEnvVars, nil
}

func sdkError(msg string, err error) *sdk.Error {
	return &sdk.Error{
		StatusType: keptnv2.StatusErrored,
		ResultType: keptnv2.ResultFailed,
		Message:    msg,
		Err:        err,
	}
}

func (th *TaskHandler) getWebHookConfig(keptnHandler sdk.IKeptn, nedc *lib.EventDataModifier) (*lib.Webhook, error) {
	// TODO: use go-utils methods for that once available
	tmpData := &TemporaryData{}
	if err := keptnv2.Decode(nedc.Get()["data"], tmpData); err != nil {
		return nil, fmt.Errorf("could not decode temporary data of incoming event: %s", err)
	}

	if tmpData.TemporaryData.Distributor.SubscriptionID == "" {
		return nil, errors.New("incoming event does not contain a subscription ID")
	}
	subscriptionID := tmpData.TemporaryData.Distributor.SubscriptionID
	// first try to retrieve the webhook config at the service level
	resource, err := keptnHandler.GetResourceHandler().GetServiceResource(nedc.Project(), nedc.Stage(), nedc.Service(), webhookConfigFileName)
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}

	// if we didn't find a config in the service directory, look at the stage
	resource, err = keptnHandler.GetResourceHandler().GetStageResource(nedc.Project(), nedc.Stage(), webhookConfigFileName)
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}

	// finally, look at project level
	resource, err = keptnHandler.GetResourceHandler().GetProjectResource(nedc.Project(), webhookConfigFileName)
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}
	return nil, errors.New("no webhook config found")
}

func getMatchingWebhookFromResource(resource *models.Resource, subscriptionID string) *lib.Webhook {
	whConfig, err := lib.DecodeWebHookConfigYAML([]byte(resource.ResourceContent))
	if err != nil {
		return nil
	}
	for _, webhook := range whConfig.Spec.Webhooks {
		if webhook.SubscriptionID == subscriptionID {
			return &webhook
		}
	}
	return nil
}
