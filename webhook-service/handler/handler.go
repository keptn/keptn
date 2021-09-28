package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/lib"
	logger "github.com/sirupsen/logrus"
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

type webhookExecution struct {
	webhook  lib.Webhook
	onStart  func() error
	onFinish func() error
	onError  func(err error) error
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
		return th.onWebhookRetrievalError(keptnHandler, event, nedc, err)
	}

	responses := []string{}

	if sdkErr := th.onStartedWebhookExecution(keptnHandler, event, webhook); sdkErr != nil {
		return nil, sdkErr
	}

	onError := th.getErrorCallbackForWebhookConfig(keptnHandler, event, nedc, webhook)

	secretEnvVars, err := th.gatherSecretEnvVars(*webhook)
	if err != nil {
		onError(err)
		return nil, sdkError(err.Error(), err)
	}
	nedc.Add("env", secretEnvVars)
	responses, err = th.performWebhookRequests(*webhook, nedc, responses)
	if err != nil {
		onError(err)
		return nil, sdkError(err.Error(), err)
	}

	// check if the incoming event was a task.triggered event, and if the 'sendFinished'  property of the webhook was set to true
	// only in this case, the result should be sent back to Keptn in the form of a .finished event
	if keptnv2.IsTaskEventType(*event.Type) && webhook.SendFinished {
		taskName, _, err := keptnv2.ParseTaskEventType(*event.Type)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not derive task name from event type %s", *event.Type), err)
		}
		result := map[string]interface{}{
			"project": nedc.Project(),
			"stage":   nedc.Stage(),
			"service": nedc.Service(),
			"labels":  nedc.Labels(),
			taskName: map[string]interface{}{
				"responses": responses,
			},
		}
		err = keptnHandler.SendFinishedEvent(event, result)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not send finished event: %s", err.Error()), err)
		}
		return result, nil
	}

	return nil, nil
}

func (th *TaskHandler) onWebhookRetrievalError(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, nedc *lib.EventDataModifier, err error) (interface{}, *sdk.Error) {
	// in this case, send .started and .finished event immediately
	if err := keptnHandler.SendStartedEvent(event); err != nil {
		// logthe error but continue - we need to try to send the .finished event nevertheless
		logger.WithError(err).Error("could not send .started event")
	}
	result := map[string]interface{}{
		"project": nedc.Project(),
		"stage":   nedc.Stage(),
		"service": nedc.Service(),
		"labels":  nedc.Labels(),
		"result":  keptnv2.ResultFailed,
		"status":  keptnv2.StatusErrored,
		"message": err.Error(),
	}
	th.sendFinishedEvent(keptnHandler, event, result)
	return nil, sdkError(fmt.Sprintf("could not retrieve Webhook config: %s", err.Error()), err)
}

func (th *TaskHandler) getErrorCallbackForWebhookConfig(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, nedc *lib.EventDataModifier, webhook *lib.Webhook) func(err error) {
	return func(err error) {
		logger.WithError(err).Error("error during webhook execution")
		whe, ok := err.(*lib.WebhookExecutionError)

		result := map[string]interface{}{
			"project": nedc.Project(),
			"stage":   nedc.Stage(),
			"service": nedc.Service(),
			"labels":  nedc.Labels(),
			"result":  keptnv2.ResultFailed,
			"status":  keptnv2.StatusErrored,
			"message": err.Error(),
		}

		if ok && whe.PreExecutionError {
			if webhook.SendFinished {
				// if sendFinished is set, we only need to send one started event
				// the webhook service will then send a correlating .finished event with the aggregated response payloads
				th.sendFinishedEvent(keptnHandler, event, result)
			} else {
				nrOfFinishedEvents := len(webhook.Requests) - whe.ExecutedRequests
				// if sendFinished is set to false, we need to send a .started event for each webhook request to be executed
				for i := 0; i < nrOfFinishedEvents; i++ {
					th.sendFinishedEvent(keptnHandler, event, result)
				}
			}
		}
	}
}

func (th *TaskHandler) sendFinishedEvent(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, result map[string]interface{}) {
	if err := keptnHandler.SendFinishedEvent(event, result); err != nil {
		logger.WithError(err).Error("could not send .finished event: %s", err.Error())
	}
}

func (th *TaskHandler) onStartedWebhookExecution(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, webhook *lib.Webhook) *sdk.Error {
	// check if 'sendFinished' is set to true
	if webhook.SendFinished {
		// if sendFinished is set, we only need to send one started event
		// the webhook service will then send a correlating .finished event with the aggregated response payloads
		if err := keptnHandler.SendStartedEvent(event); err != nil {
			return sdkError(fmt.Sprintf("could not send .started event: %s", err.Error()), err)
		}
	} else {
		// if sendFinished is set to false, we need to send a .started event for each webhook request to be executed
		for range webhook.Requests {
			if err := keptnHandler.SendStartedEvent(event); err != nil {
				return sdkError(fmt.Sprintf("could not send .started event: %s", err.Error()), err)
			}
		}
	}
	return nil
}

func (th *TaskHandler) performWebhookRequests(webhook lib.Webhook, nedc *lib.EventDataModifier, responses []string) ([]string, error) {
	executedRequests := 0
	for _, req := range webhook.Requests {
		// parse the data from the event, together with the secret env vars
		parsedCurlCommand, err := th.templateEngine.ParseTemplate(nedc.Get(), req)
		if err != nil {
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("could not parse request '%s'", req), lib.WithNrOfExecutedRequests(executedRequests))
		}
		// perform the request
		response, err := th.curlExecutor.Curl(parsedCurlCommand)
		if err != nil {
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("could not execute request '%s': %s", req, err.Error()), lib.WithNrOfExecutedRequests(executedRequests))
		}
		executedRequests = executedRequests + 1
		responses = append(responses, response)
	}
	return responses, nil
}

func (th *TaskHandler) gatherSecretEnvVars(webhook lib.Webhook) (map[string]string, error) {
	secretEnvVars := map[string]string{}
	for _, secretRef := range webhook.EnvFrom {
		secretValue, err := th.secretReader.ReadSecret(secretRef.SecretRef.Name, secretRef.SecretRef.Key)
		if err != nil {
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("could not read secret %s.%s", secretRef.SecretRef.Name, secretRef.SecretRef.Key))
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
