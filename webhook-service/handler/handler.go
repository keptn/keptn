package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	keptn "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/lib"
	logger "github.com/sirupsen/logrus"
)

const webhookConfigFileName = "webhook/webhook.yaml"

type SecretEnv struct {
	Env map[string]string
}

type TaskHandler struct {
	templateEngine   lib.ITemplateEngine
	curlExecutor     lib.ICurlExecutor
	requestValidator lib.RequestValidator
	secretReader     lib.ISecretReader
}

func NewTaskHandler(templateEngine lib.ITemplateEngine, curlExecutor lib.ICurlExecutor, requestValidator lib.RequestValidator, secretReader lib.ISecretReader) *TaskHandler {
	return &TaskHandler{
		templateEngine:   templateEngine,
		curlExecutor:     curlExecutor,
		requestValidator: requestValidator,
		secretReader:     secretReader,
	}
}

func (th *TaskHandler) Execute(keptnHandler sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	logger.Infof("Received event %s", event)
	eventAdapter, err := lib.NewEventDataAdapter(event)
	if err != nil {
		return nil, sdkError("could not parse incoming event", err)
	}
	subscriptionID, err := eventAdapter.SubscriptionID()
	if err != nil {
		logger.Infof("will not handle event: %s", err.Error())
		return nil, nil
	}
	webhook, err := th.getWebHookConfig(keptnHandler, eventAdapter, subscriptionID, event.GitCommitID)
	if err != nil {
		return th.onPreExecutionError(keptnHandler, event, eventAdapter, fmt.Errorf("could not retrieve Webhook config: %s", err.Error()))
	}

	responses := []string{}

	if sdkErr := th.onStartedWebhookExecution(keptnHandler, event, webhook); sdkErr != nil {
		return nil, sdkErr
	}

	onError := th.getErrorCallbackForWebhookConfig(keptnHandler, event, eventAdapter, webhook)

	secretEnvVars, err := th.gatherSecretEnvVars(*webhook)
	if err != nil {
		onError(err, secretEnvVars)
		return nil, sdkError(removeSecretsFromMessage(err.Error(), secretEnvVars), err)
	}
	eventAdapter.Add("env", secretEnvVars)
	responses, err = th.performWebhookRequests(*webhook, eventAdapter, responses)
	if err != nil {
		onError(err, secretEnvVars)
		return nil, sdkError(removeSecretsFromMessage(err.Error(), secretEnvVars), err)
	}

	// check if the incoming event was a task.triggered event, and if the 'sendFinished'  property of the webhook was set to true
	// only in this case, the result should be sent back to Keptn in the form of a .finished event
	if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && webhook.ShouldSendFinishedEvent() {
		taskName, _, err := keptnv2.ParseTaskEventType(*event.Type)
		if err != nil {
			return nil, sdkError(fmt.Sprintf("could not derive task name from event type %s", *event.Type), err)
		}
		result := map[string]interface{}{
			"project": eventAdapter.Project(),
			"stage":   eventAdapter.Stage(),
			"service": eventAdapter.Service(),
			"labels":  eventAdapter.Labels(),
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

func (th *TaskHandler) onPreExecutionError(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, eventAdapter *lib.EventDataAdapter, err error) (interface{}, *sdk.Error) {
	// in this case, send .started and .finished event immediately
	if err := keptnHandler.SendStartedEvent(event); err != nil {
		// log the error but continue - we need to try to send the .finished event nevertheless
		logger.WithError(err).Error("could not send .started event")
	}
	result := map[string]interface{}{
		"project": eventAdapter.Project(),
		"stage":   eventAdapter.Stage(),
		"service": eventAdapter.Service(),
		"labels":  eventAdapter.Labels(),
		"result":  keptnv2.ResultFailed,
		"status":  keptnv2.StatusErrored,
		"message": err.Error(),
	}
	th.sendFinishedEvent(keptnHandler, event, result)
	return nil, sdkError(err.Error(), err)
}

func (th *TaskHandler) getErrorCallbackForWebhookConfig(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, eventAdapter *lib.EventDataAdapter, webhook *lib.Webhook) func(err error, secrets map[string]string) {
	return func(err error, secrets map[string]string) {
		logger.WithError(err).Error("error during webhook execution")
		whe, ok := err.(*lib.WebhookExecutionError)

		result := map[string]interface{}{
			"project": eventAdapter.Project(),
			"stage":   eventAdapter.Stage(),
			"service": eventAdapter.Service(),
			"labels":  eventAdapter.Labels(),
			"result":  keptnv2.ResultFailed,
			"status":  keptnv2.StatusErrored,
			"message": removeSecretsFromMessage(err.Error(), secrets),
		}

		if ok && whe.PreExecutionError {
			if webhook.ShouldSendFinishedEvent() {
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

func removeSecretsFromMessage(errMsg string, secrets map[string]string) string {
	result := errMsg
	for _, val := range secrets {
		result = strings.ReplaceAll(result, val, "***")
	}
	return result
}

func (th *TaskHandler) sendFinishedEvent(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, result map[string]interface{}) {
	// ensure that finished events are only sent as responses for .triggered events
	if !keptnv2.IsTaskEventType(*event.Type) || !keptnv2.IsTriggeredEventType(*event.Type) {
		return
	}
	if err := keptnHandler.SendFinishedEvent(event, result); err != nil {
		logger.WithError(err).Error("could not send .finished event")
	}
}

func (th *TaskHandler) onStartedWebhookExecution(keptnHandler sdk.IKeptn, event sdk.KeptnEvent, webhook *lib.Webhook) *sdk.Error {
	// only send .started events for <task>.triggered events
	if !keptnv2.IsTaskEventType(*event.Type) || !keptnv2.IsTriggeredEventType(*event.Type) {
		return nil
	}

	// check if a started event should be sent for the webhook
	if !webhook.ShouldSendStartedEvent() {
		return nil
	}
	// check if 'sendFinished' is set to true
	if webhook.ShouldSendFinishedEvent() {
		// if sendFinished is set, we only need to send one started event
		// the webhook service will then send a correlating .finished event with the aggregated response payloads
		if err := keptnHandler.SendStartedEvent(event); err != nil {
			return sdkError(fmt.Sprintf("could not send .started event: %s", err.Error()), err)
		}
	} else {
		// if sendFinished is set to false, we need to send a .started event for each webhook request to be executed
		for _, req := range webhook.Requests {
			if err := keptnHandler.SendStartedEvent(event); err != nil {
				return sdkError(fmt.Sprintf("could not send .started event: request '%s': %s", req, err.Error()), err)
			}
		}
	}
	return nil
}

func (th *TaskHandler) performWebhookRequests(webhook lib.Webhook, eventAdapter *lib.EventDataAdapter, responses []string) ([]string, error) {
	executedRequests := 0
	logger.Infof("executing webhooks for subscriptionID %s", webhook.SubscriptionID)
	for _, req := range webhook.Requests {
		request, err := th.CreateRequest(req)
		if err != nil {
			logger.Infof("creating CURL request failed: %s", err.Error())
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("creating CURL request failed: %s", err.Error()), lib.WithNrOfExecutedRequests(executedRequests))
		}
		// parse the data from the event, together with the secret env vars
		parsedCurlCommand, err := th.templateEngine.ParseTemplate(eventAdapter.Get(), request)
		if err != nil {
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("could not parse request '%s' : %s", request, err.Error()), lib.WithNrOfExecutedRequests(executedRequests))
		}
		// perform the request
		response, err := th.curlExecutor.Curl(parsedCurlCommand)
		if err != nil {
			return nil, lib.NewWebhookExecutionError(true, fmt.Errorf("could not execute request '%s': %s", request, err.Error()), lib.WithNrOfExecutedRequests(executedRequests))
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

func (th *TaskHandler) CreateRequest(request interface{}) (string, error) {
	switch req := request.(type) {
	// v1alpha1 version
	case string:
		logger.Debug("creating CURL request from type string")
		if err := th.validateAlphaCurlRequest(req); err != nil {
			return "", err
		}
		return req, nil
	// v1beta1 version
	default:
		logger.Debug("creating CURL request from type Request")
		convertedRequest := lib.ConvertToRequest(request)
		if err := th.requestValidator.Validate(convertedRequest); err != nil {
			return "", err
		}
		betaRequest := buildBetaCurlRequest(convertedRequest)
		if betaRequest != "" {
			return betaRequest, nil
		}
	}

	return "", fmt.Errorf("could not create request: invalid request type")
}

func (th *TaskHandler) validateAlphaCurlRequest(curlCmd string) error {
	sanitizedCurlCmd := strings.ReplaceAll(curlCmd, "\\", "")
	denyList := lib.CreateListOfDeniedURLs(lib.GetEnv())
	for _, durl := range denyList {
		if strings.Contains(sanitizedCurlCmd, durl) {
			return fmt.Errorf("curl command contains denied URL '%s'", durl)
		}
	}
	return nil
}

func buildBetaCurlRequest(req lib.Request) string {
	tmpReq := ""
	if req.Method != "" {
		tmpReq = fmt.Sprintf("curl --request %s", req.Method)
	}
	if len(req.Headers) > 0 {
		for _, header := range req.Headers {
			tmpReq = fmt.Sprintf(tmpReq+" --header '%s: %s'", header.Key, header.Value)
		}
	}
	if req.Payload != "" {
		tmpReq = fmt.Sprintf(tmpReq+" --data '%s'", req.Payload)
	}
	if req.Options != "" {
		tmpReq = fmt.Sprintf(tmpReq+" %s", req.Options)
	}
	if req.URL != "" {
		tmpReq = fmt.Sprintf(tmpReq+" %s", req.URL)
	}
	return tmpReq
}

func sdkError(msg string, err error) *sdk.Error {
	return &sdk.Error{
		StatusType: keptnv2.StatusErrored,
		ResultType: keptnv2.ResultFailed,
		Message:    msg,
		Err:        err,
	}
}

func (th *TaskHandler) getWebHookConfig(keptnHandler sdk.IKeptn, eventAdapter *lib.EventDataAdapter, subscriptionID string, commitID string) (*lib.Webhook, error) {
	// first try to retrieve the webhook config at the service level
	commitOption := url.Values{}
	if commitID != "" {
		commitOption.Add("commitID", commitID)
	}
	resourceScope := *keptn.NewResourceScope().Project(eventAdapter.Project()).Stage(eventAdapter.Stage()).Service(eventAdapter.Service()).Resource(webhookConfigFileName)
	resource, err := keptnHandler.GetResourceHandler().GetResource(resourceScope, keptn.AppendQuery(commitOption))
	logger.Debug("searching for webhook config at service level...")
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}

	// if we didn't find a config in the service directory, look at the stage
	resourceScope = *keptn.NewResourceScope().Project(eventAdapter.Project()).Stage(eventAdapter.Stage()).Resource(webhookConfigFileName)
	resource, err = keptnHandler.GetResourceHandler().GetResource(resourceScope, keptn.AppendQuery(commitOption))
	logger.Debug("searching for webhook config at stage level...")
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}

	// finally, look at project level
	resourceScope = *keptn.NewResourceScope().Project(eventAdapter.Project()).Resource(webhookConfigFileName)
	resource, err = keptnHandler.GetResourceHandler().GetResource(resourceScope, keptn.AppendQuery(commitOption))
	logger.Debug("searching for webhook config at project level...")
	if err == nil && resource != nil {
		if matchingWebhook := getMatchingWebhookFromResource(resource, subscriptionID); matchingWebhook != nil {
			return matchingWebhook, nil
		}
	}
	if err != nil {
		logger.Debugf("no webhook config found, err: %s", err.Error())
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
