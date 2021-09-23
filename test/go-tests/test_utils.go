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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
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

func (sender *APIEventSender) Send(ctx context.Context, event v2.Event) error {
	return sender.SendEvent(event)
}

func (sender *APIEventSender) SendEvent(event v2.Event) error {
	_, err := ApiPOSTRequest("/v1/event", event)
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
		resp, err = ApiGETRequest("/controlPlane/v1/project/" + projectName)
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

		_, err = ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=./%s", projectName, shipyardFilePath))

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
	})
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
	resp, _ := ApiGETRequest("/controlPlane/v1/uniform/registration")
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

func CreateSubscription(t *testing.T, serviceName string, subscription models.EventSubscription) error {
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
		if s.Event == subscription.Event {
			return nil
		}
	}

	_, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription", fetchedIntegration.ID), subscription)
	require.Nil(t, err)

	return nil
}

func ApiDELETERequest(path string) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	r, err := req.Delete(keptnAPIURL+path, authHeader)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func getAuthHeader(apiToken string) req.Header {
	authHeader := req.Header{
		"Accept":  "application/json",
		"x-token": apiToken,
	}
	return authHeader
}

func ApiGETRequest(path string) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	r, err := req.Get(keptnAPIURL+path, authHeader)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ApiPOSTRequest(path string, payload interface{}) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	r, err := req.Post(keptnAPIURL+path, authHeader, req.BodyJSON(payload))
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ApiPUTRequest(path string, payload interface{}) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := getAuthHeader(apiToken)

	r, err := req.Put(keptnAPIURL+path, authHeader, req.BodyJSON(payload))
	if err != nil {
		return nil, err
	}

	return r, nil
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

func WaitForPodOfDeployment(deploymentName string) error {
	return keptnkubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

func CreateTmpShipyardFile(shipyardContent string) (string, error) {
	return CreateTmpFile("shipyard-*.yaml", shipyardContent)
}

func CreateTmpFile(fileNamePattern, fileContent string) (string, error) {
	file, err := ioutil.TempFile(".", fileNamePattern)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file.Name(), []byte(fileContent), os.ModeAppend); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
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
	resp, err := ApiGETRequest("/mongodb-datastore/event?project=" + projectName + "&keptnContext=" + keptnContext + "&stage=" + stage + "&type=" + eventType)
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

func GetEventTraceForContext(keptnContext, projectName string) ([]*models.KeptnContextExtendedCE, error) {
	resp, err := ApiGETRequest("/mongodb-datastore/event?project=" + projectName + "&keptnContext=" + keptnContext)
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

	resp, err := ApiGETRequest("/controlPlane/v1/sequence/" + projectName)
	err = resp.ToJSON(states)

	return states, resp, err
}

func KubeClient(t *testing.T) *kubernetes.Clientset {
	clientset, err := keptnkubeutils.GetClientset(false)
	require.Nil(t, err)
	return clientset
}

func SetEnvVarsOfDeployment(deploymentName string, containerName string, envVars []v1.EnvVar) error {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return err
	}
	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for index, container := range depl.Spec.Template.Spec.Containers {
		if "distributor" == container.Name {
			for _, e := range envVars {
				replaced := false
				for ii, existingEnvVar := range depl.Spec.Template.Spec.Containers[index].Env {
					// if we find an already existing env war with the same name, we need to replace it
					if existingEnvVar.Name == e.Name {
						depl.Spec.Template.Spec.Containers[index].Env[ii] = e
						replaced = true
						break
					}
				}
				// if we did not replace an env var, we need to append it
				if !replaced {
					depl.Spec.Template.Spec.Containers[index].Env = append(depl.Spec.Template.Spec.Containers[index].Env, e)
					replaced = false
				}
			}
		}
	}

	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return keptnkubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

func GetImageOfDeploymentContainer(deploymentName, containerName string) (string, error) {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return "", err
	}
	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, container := range depl.Spec.Template.Spec.Containers {
		if containerName == container.Name {
			return container.Image, nil
		}
	}
	return "", fmt.Errorf("container %s not found in deployment %s", containerName, deploymentName)
}

func SetImageOfDeploymentContainer(deploymentName, containerName, image string) error {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return err
	}

	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for index, container := range depl.Spec.Template.Spec.Containers {
		if containerName == container.Name {
			depl.Spec.Template.Spec.Containers[index].Image = image
			depl.Spec.Template.Spec.Containers[index].ImagePullPolicy = "Always"
		}
	}
	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return keptnkubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
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
