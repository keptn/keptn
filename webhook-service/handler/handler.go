package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/lib"
)

const webhookConfigFileName = "webhook.yaml"

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

func (th *TaskHandler) Execute(keptnHandler sdk.IKeptn, data interface{}, eventType string) (interface{}, *sdk.Error) {

	eventData := &keptnv2.EventData{}
	if err := keptnv2.Decode(data, eventData); err != nil {
		return nil, sdkError("could not decode incoming event payload", err)
	}

	eventDataMap := map[string]interface{}{}
	// apply the EventData attributes to the result
	if err := keptnv2.Decode(data, &eventDataMap); err != nil {
		return nil, sdkError("could not apply attributes from incoming event", err)
	}

	resource, err := th.getWebHookConfigResource(keptnHandler, eventData)
	if err != nil {
		return nil, sdkError("could not find webhook config", err)
	}

	whConfig, err := keptnv2.DecodeWebHookConfigYAML([]byte(resource.ResourceContent))
	if err != nil {
		return nil, sdkError("could not decode webhook config", err)
	}

	responses := []string{}
	for _, webhook := range whConfig.Spec.Webhooks {
		if webhook.Type == eventType {
			secretEnvVars := map[string]string{}
			for _, secretRef := range webhook.EnvFrom {
				secretValue, err := th.secretReader.ReadSecret(secretRef.SecretRef.Name, secretRef.SecretRef.Key)
				if err != nil {
					return nil, sdkError(fmt.Sprintf("could not read secret %s.%s", secretRef.SecretRef.Name, secretRef.SecretRef.Key), err)
				}
				secretEnvVars[secretRef.Name] = secretValue
			}
			eventDataMap["env"] = secretEnvVars
			for _, req := range webhook.Requests {
				// parse the data from the event, together with the secret env vars
				parsedCurlCommand, err := th.templateEngine.ParseTemplate(eventDataMap, req)
				if err != nil {
					return nil, sdkError(fmt.Sprintf("could not parse request '%s'", req), err)
				}

				// perform the request
				response, err := th.curlExecutor.Curl(parsedCurlCommand)
				if err != nil {
					return nil, sdkError(fmt.Sprintf("could not execute request '%s'", req), err)
				}
				responses = append(responses, response)
			}
		}
	}

	eventDataMap[eventType] = map[string]interface{}{
		"responses": responses,
	}

	return eventDataMap, nil
}

func sdkError(msg string, err error) *sdk.Error {
	return &sdk.Error{
		StatusType: keptnv2.StatusErrored,
		ResultType: keptnv2.ResultFailed,
		Message:    msg,
		Err:        err,
	}
}

func (th *TaskHandler) performCurlRequest() (string, error) {
	return "", nil
}

func (th *TaskHandler) getWebHookConfigResource(keptnHandler sdk.IKeptn, eventData *keptnv2.EventData) (*models.Resource, error) {
	// first try to retrieve the webhook config at the service level
	resource, err := keptnHandler.GetResourceHandler().GetServiceResource(eventData.Project, eventData.Stage, eventData.Service, webhookConfigFileName)
	if err == nil && resource != nil {
		return resource, nil
	}

	// if we didn't find a config in the service directory, look at the stage
	resource, err = keptnHandler.GetResourceHandler().GetStageResource(eventData.Project, eventData.Stage, webhookConfigFileName)
	if err == nil && resource != nil {
		return resource, nil
	}

	// finally, look at project level
	resource, err = keptnHandler.GetResourceHandler().GetProjectResource(eventData.Project, webhookConfigFileName)
	if err == nil && resource != nil {
		return resource, nil
	}
	return nil, errors.New("no webhook config found")
}
