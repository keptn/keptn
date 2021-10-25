package go_tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	KeptnSpecVersion      = "0.2.0"
	KeptnNamespaceEnvVar  = "KEPTN_NAMESPACE"
	DefaultKeptnNamespace = "keptn"
)

type APIEventSender struct {
}

type OpenTriggeredEventsResponse struct {
	Events []*models.KeptnContextExtendedCE `json:"events"`
}

func (sender *APIEventSender) Send(ctx context.Context, event v2.Event) error {
	return sender.SendEvent(event)
}

func (sender *APIEventSender) SendEvent(event v2.Event) error {
	_, err := ApiPOSTRequest("/v1/event", event, 3)
	return err
}

func CreateProject(projectName, shipyardFilePath string, recreateIfAlreadyThere bool) error {

	retries := 3
	var err error
	var resp *req.Resp
	for i := 0; i < retries; i++ {
		if err != nil {
			<-time.After(5 * time.Second)
		}
		resp, err = ApiGETRequest("/controlPlane/v1/project/"+projectName, 3)
		if err != nil {
			continue
		}

		if resp.Response().StatusCode == http.StatusOK {
			if recreateIfAlreadyThere {
				// delete project if it exists
				_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", projectName))
				if err != nil {
					continue
				}
			} else {
				return errors.New("project already exists")
			}
		}

		_, err = ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=%s", projectName, shipyardFilePath))

		if err == nil {
			return nil
		}
	}

	return err
}

func TriggerSequence(projectName, serviceName, stageName, sequenceName string, eventData keptncommon.EventProperties) (string, error) {
	source := "golang-test"
	eventType := keptnv2.GetTriggeredEventType(stageName + "." + sequenceName)
	if eventData == nil {
		eventData = &keptnv2.EventData{}
	}
	eventData.SetProject(projectName)
	eventData.SetService(serviceName)
	eventData.SetStage(stageName)

	resp, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               eventData,
		ID:                 uuid.NewString(),
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	}, 3)
	if err != nil {
		return "", err
	}

	context := &models.EventContext{}
	err = resp.ToJSON(context)
	if err != nil {
		return "", err
	}
	return *context.KeptnContext, nil
}

func GetIntegrationWithName(name string) (models.Integration, error) {
	resp, _ := ApiGETRequest("/controlPlane/v1/uniform/registration", 3)
	integrations := []models.Integration{}
	if err := resp.ToJSON(&integrations); err != nil {
		return models.Integration{}, err
	}
	for _, r := range integrations {
		if r.Name == name {
			return r, nil
		}
	}
	return models.Integration{}, fmt.Errorf("No Keptn Integration with name %s found", name)
}

func CreateSubscription(t *testing.T, serviceName string, subscription models.EventSubscription) (string, error) {
	var fetchedIntegration models.Integration
	var err error
	require.Eventually(t, func() bool {
		fetchedIntegration, err = GetIntegrationWithName(serviceName)
		return err == nil
	}, time.Second*20, time.Second*3)

	// Integration exists - fine
	require.Nil(t, err)
	require.NotNil(t, fetchedIntegration)

	for _, s := range fetchedIntegration.Subscriptions {
		// check if the subscription for the event already exists - if yes, fine
		if s.Event == subscription.Event && reflect.DeepEqual(s.Filter, subscription.Filter) {
			return s.ID, nil
		}
	}

	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription", fetchedIntegration.ID), subscription, 3)
	require.Nil(t, err)

	subscriptionResponse := &scmodels.CreateSubscriptionResponse{}

	err = resp.ToJSON(subscriptionResponse)
	require.Nil(t, err)

	require.NotEmpty(t, subscriptionResponse.ID)

	return subscriptionResponse.ID, nil
}

func ApiDELETERequest(path string, retries int) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	return doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Delete(keptnAPIURL+path, authHeader)
	}, retries)
}

func getAuthHeader(apiToken string) req.Header {
	authHeader := req.Header{
		"Accept":  "application/json",
		"x-token": apiToken,
	}
	return authHeader
}

func ApiGETRequest(path string, retries int) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	return doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Get(keptnAPIURL+path, authHeader)
	}, retries)
}

func doHTTPRequestWithRetry(httpFunc func() (*req.Resp, error), retries int) (*req.Resp, error) {
	var reqErr error
	var r *req.Resp
	for i := 0; i <= retries; i++ {
		r, reqErr = httpFunc()
		if reqErr == nil {
			return r, nil
		}
		<-time.After(5 * time.Second)
	}
	return r, reqErr
}

func ApiPOSTRequest(path string, payload interface{}, retries int) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	return doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Post(keptnAPIURL+path, authHeader, req.BodyJSON(payload))
	}, retries)
}

func ApiPUTRequest(path string, payload interface{}, retries int) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	return doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Put(keptnAPIURL+path, authHeader, req.BodyJSON(payload))
	}, retries)
}

func GetApiCredentials() (string, string, error) {
	apiToken, err := keptnkubeutils.GetKeptnAPITokenFromSecret(false, GetKeptnNameSpaceFromEnv(), "keptn-api-token")
	if err != nil {
		return "", "", err
	}
	keptnAPIURL := os.Getenv("KEPTN_ENDPOINT")
	if keptnAPIURL == "" {
		serviceIP, err := keptnkubeutils.GetKeptnEndpointFromService(false, GetKeptnNameSpaceFromEnv(), "api-gateway-nginx")
		if err != nil {
			return "", "", err
		}
		keptnAPIURL = "http://" + serviceIP + "/api"
	}
	return apiToken, keptnAPIURL, nil
}

func ScaleDownUniform(deployments []string) error {
	for _, deployment := range deployments {
		if err := keptnkubeutils.ScaleDeployment(false, deployment, GetKeptnNameSpaceFromEnv(), 0); err != nil {
			// log the error but continue
			fmt.Println("could not scale down deployment: " + err.Error())
		}
	}
	return nil
}

func ScaleUpUniform(deployments []string) error {
	for _, deployment := range deployments {
		if err := keptnkubeutils.ScaleDeployment(false, deployment, GetKeptnNameSpaceFromEnv(), 1); err != nil {
			// log the error but continue
			fmt.Println("could not scale up deployment: " + err.Error())
		}
	}
	return nil
}

func RestartPod(deploymentName string) error {
	return keptnkubeutils.RestartPodsWithSelector(false, GetKeptnNameSpaceFromEnv(), "app.kubernetes.io/name="+deploymentName)
}

func CreateTmpShipyardFile(shipyardContent string) (string, error) {
	return CreateTmpFile("shipyard-*.yaml", shipyardContent)
}

func CreateTmpFile(fileNamePattern, fileContent string) (string, error) {
	file, err := ioutil.TempFile("", fileNamePattern)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file.Name(), []byte(fileContent), os.ModeAppend); err != nil {
		err = os.Remove(file.Name())
		if err != nil {
			return "", err
		}
		return "", err
	}
	return file.Name(), nil
}

func CreateTmpDir() (string, error) {
	return ioutil.TempDir("", "")
}

func ExecuteCommand(cmd string) (string, error) {
	split := strings.Split(cmd, " ")
	if len(split) == 0 {
		return "", errors.New("invalid command")
	}
	return keptnkubeutils.ExecuteCommand(split[0], split[1:])
}

func ExecuteCommandf(cmd string, a ...interface{}) (string, error) {
	cmdf := fmt.Sprintf(cmd, a...) //nolint:govet
	return ExecuteCommand(cmdf)
}

func GetKeptnNameSpaceFromEnv() string {
	return osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace)
}

func GetLatestEventOfType(keptnContext, projectName, stage, eventType string) (*models.KeptnContextExtendedCE, error) {
	resp, err := ApiGETRequest("/mongodb-datastore/event?project="+projectName+"&keptnContext="+keptnContext+"&stage="+stage+"&type="+eventType, 3)
	if err != nil {
		return nil, err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return nil, err
	}
	if len(events.Events) > 0 {
		return events.Events[0], nil
	}
	return nil, nil
}

func GetEventsOfType(keptnContext, projectName, stage, eventType string) ([]*models.KeptnContextExtendedCE, error) {
	resp, err := ApiGETRequest("/mongodb-datastore/event?project="+projectName+"&keptnContext="+keptnContext+"&stage="+stage+"&type="+eventType, 3)
	if err != nil {
		return nil, err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return nil, err
	}
	if len(events.Events) > 0 {
		return events.Events, nil
	}
	return nil, nil
}

func GetEventTraceForContext(keptnContext, projectName string) ([]*models.KeptnContextExtendedCE, error) {
	resp, err := ApiGETRequest("/mongodb-datastore/event?project="+projectName+"&keptnContext="+keptnContext, 3)
	if err != nil {
		return nil, err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return nil, err
	}
	if len(events.Events) > 0 {
		return events.Events, nil
	}
	return nil, nil
}

func IsEqual(t *testing.T, expected, actual interface{}, property string) bool {
	if expected != actual {
		t.Logf("%s: expected %v, got %v", property, expected, actual)
		return false
	}
	return true
}

func StringArr(el ...string) []string {
	return el
}

func VerifySequenceEndsUpInState(t *testing.T, projectName string, context *models.EventContext, timeout time.Duration, desiredStates []string) {
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, desiredStates)
	require.Eventuallyf(t, func() bool {
		states, _, err := GetState(projectName)
		if err != nil {
			return false
		}
		for _, state := range states.States {
			if doesSequenceHaveOneOfTheDesiredStates(state, context, desiredStates) {
				return true
			}
		}
		return false
	}, timeout, 10*time.Second, GetDiagnostics("shipyard-controller"))
}

func doesSequenceHaveOneOfTheDesiredStates(state scmodels.SequenceState, context *models.EventContext, desiredStates []string) bool {
	if state.Shkeptncontext == *context.KeptnContext {
		for _, desiredState := range desiredStates {
			if state.State == desiredState {
				return true
			}
		}
	}
	return false
}

func GetState(projectName string) (*scmodels.SequenceStates, *req.Resp, error) {
	states := &scmodels.SequenceStates{}

	resp, err := ApiGETRequest("/controlPlane/v1/sequence/"+projectName, 3)
	err = resp.ToJSON(states)

	return states, resp, err
}

func GetProject(projectName string) (*scmodels.ExpandedProject, error) {
	project := &scmodels.ExpandedProject{}

	resp, err := ApiGETRequest("/controlPlane/v1/project/" + projectName)
	if err != nil {
		return nil, err
	}

	err = resp.ToJSON(project)
	return project, err
}

func GetDiagnostics(service string) string {
	outputBuilder := strings.Builder{}
	getLogsCmd := fmt.Sprintf("kubectl logs -n %s deployment/%s -c %s", GetKeptnNameSpaceFromEnv(), service, service)

	outputBuilder.WriteString(fmt.Sprintf("Logs of  of %s: \n\n", service))
	logOutput, err := ExecuteCommand(getLogsCmd)
	if err != nil {
		outputBuilder.WriteString(err.Error())
	}

	outputBuilder.WriteString(logOutput)
	outputBuilder.WriteString("\n-------------------------\n")
	getLogsCmd = fmt.Sprintf("kubectl logs -n %s deployment/%s -c %s --previous", GetKeptnNameSpaceFromEnv(), service, service)

	outputBuilder.WriteString(fmt.Sprintf("Logs of crashed instances of %s: \n\n", service))
	logOutput, err = ExecuteCommand(getLogsCmd)
	if err != nil {
		outputBuilder.WriteString(err.Error())
	}

	outputBuilder.WriteString(logOutput)
	outputBuilder.WriteString("\n-------------------------\n")

	describeDeploymentCmd := fmt.Sprintf("kubectl -n %s describe deployment %s", GetKeptnNameSpaceFromEnv(), service)
	outputBuilder.WriteString(fmt.Sprintf("Description of Deployment %s", service))
	describeDeploymentOutput, err := ExecuteCommand(describeDeploymentCmd)
	if err != nil {
		outputBuilder.WriteString(err.Error())
	}
	outputBuilder.WriteString(describeDeploymentOutput)

	return outputBuilder.String()
}

func VerifyDirectDeployment(serviceName, projectName, stageName, artifactImage, artifactTag string) error {
	return WaitAndCheckDeployment(serviceName, projectName+"-"+stageName, time.Minute*10, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag})
}

func VerifyBlueGreenDeployment(serviceName, projectName, stageName, artifactImage, artifactTag string) error {
	if err := WaitAndCheckDeployment(serviceName, projectName+"-"+stageName, time.Minute*10, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag}); err != nil {
		return err
	}
	return WaitAndCheckDeployment(serviceName+"-primary", projectName+"-"+stageName, time.Minute*10, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag})
}

func GetPublicURLOfService(serviceName, projectName, stageName string) (string, error) {
	ingressHostnameSuffix, err := GetFromConfigMap(GetKeptnNameSpaceFromEnv(), "ingress-config", func(data map[string]string) string { return data["ingress_hostname_suffix"] })
	if err != nil {
		return "", fmt.Errorf("unable to get public URL of service %s: %w", serviceName, err)
	}

	return fmt.Sprintf("http://%s.%s-%s.%s", serviceName, projectName, stageName, ingressHostnameSuffix), nil

}

func SetShipyardControllerEnvVar(t *testing.T, envVar, timeoutValue string) error {
	_, err := ExecuteCommand(fmt.Sprintf("kubectl -n %s set env deployment shipyard-controller %s=%s", GetKeptnNameSpaceFromEnv(), envVar, timeoutValue))
	if err != nil {
		return err
	}

	t.Log("restarting shipyard controller pod")
	err = RestartPod("shipyard-controller")
	if err != nil {
		return err
	}

	// wait 10s to make sure we wait for the updated pod to be ready
	<-time.After(10 * time.Second)
	t.Log("waiting for shipyard controller pod to be ready again")
	err = WaitForPodOfDeployment("shipyard-controller")
	if err != nil {
		return err
	}

	// check whether the API is reachable again
	require.Eventually(t, func() bool {
		t.Log("Verifying API availability")
		// use the shipyard-controller's project endpoint to check API availability
		resp, err := ApiGETRequest("/controlPlane/v1/project", 3)
		if err != nil {
			t.Logf("got error from API: %s", err.Error())
			return false
		}

		if resp.Response().StatusCode != http.StatusOK {
			t.Logf("API response does not have expected status code")
			return false
		}

		return true
	}, 30*time.Second, 5*time.Second)

	return nil
}
