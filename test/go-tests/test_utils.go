package go_tests

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	v12 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/keptn/go-utils/pkg/common/strutils"

	"github.com/keptn/go-utils/pkg/common/retry"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
)

const (
	KeptnSpecVersion      = "0.2.0"
	KeptnNamespaceEnvVar  = "KEPTN_NAMESPACE"
	DefaultKeptnNamespace = "keptn"
)

type APICaller struct {
	baseURL string
	token   string
}

var errProjectAlreadyExists = errors.New("project already exists")

func NewAPICallerWithBaseURL(baseURL string) (*APICaller, error) {
	token, _, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}
	return &APICaller{
		baseURL: baseURL,
		token:   token,
	}, nil
}

func NewAPICaller() (*APICaller, error) {
	token, baseURL, err := GetApiCredentials()
	if err != nil {
		return nil, err
	}
	return &APICaller{
		baseURL: baseURL,
		token:   token,
	}, nil
}

func (a *APICaller) Get(path string, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		req.SetTimeout(15 * time.Second)
		return req.Get(a.baseURL+path, a.getAuthHeader())
	}, retries)
}

func (a *APICaller) Delete(path string, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		req.SetTimeout(15 * time.Second)
		return req.Delete(a.baseURL+path, a.getAuthHeader())
	}, retries)
}

func (a *APICaller) Put(path string, payload interface{}, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		req.SetTimeout(15 * time.Second)
		return req.Put(a.baseURL+path, a.getAuthHeader(), req.BodyJSON(payload))
	}, retries)
}

func (a *APICaller) Post(path string, payload interface{}, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		req.SetTimeout(15 * time.Second)
		return req.Post(a.baseURL+path, a.getAuthHeader(), req.BodyJSON(payload))
	}, retries)
}

func (a *APICaller) getAuthHeader() req.Header {
	authHeader := req.Header{
		"Accept":  "application/json",
		"x-token": a.token,
	}
	return authHeader
}

func (a *APICaller) doHTTPRequestWithRetry(httpFunc func() (*req.Resp, error), retries int) (*req.Resp, error) {
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

func ApiDELETERequest(path string, retries int) (*req.Resp, error) {
	caller, err := NewAPICaller()
	if err != nil {
		return nil, err
	}
	return caller.Delete(path, retries)
}

func ApiPOSTRequest(path string, payload interface{}, retries int) (*req.Resp, error) {
	caller, err := NewAPICaller()
	if err != nil {
		return nil, err
	}
	return caller.Post(path, payload, retries)
}

func ApiPUTRequest(path string, payload interface{}, retries int) (*req.Resp, error) {
	caller, err := NewAPICaller()
	if err != nil {
		return nil, err
	}
	return caller.Put(path, payload, retries)
}

func ApiGETRequest(path string, retries int) (*req.Resp, error) {
	caller, err := NewAPICaller()
	if err != nil {
		return nil, err
	}
	return caller.Get(path, retries)
}

func GetInternalKeptnAPI(ctx context.Context, internalService, localPort string, remotePort string) (*APICaller, error) {
	err := KubeCtlPortForwardSvc(ctx, internalService, localPort, remotePort)
	if err != nil {
		return nil, err
	}
	keptnInternalAPI, err := NewAPICallerWithBaseURL("http://127.0.0.1:" + localPort)
	if err != nil {
		return nil, err
	}
	return keptnInternalAPI, nil
}

func RecreateProjectUpstream(newProjectName string) error {
	resp, err := ApiGETRequest("/controlPlane/v1/project/"+newProjectName, 3)
	if err != nil {
		return err
	}

	if resp.Response().StatusCode == http.StatusOK {
		// delete project if it exists
		_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", newProjectName))
		if err != nil {
			return err
		}
	}

	err = RecreateGitUpstreamRepository(newProjectName)
	if err != nil {
		// retry if repo creation failed (gitea might not be available)
		return err
	}

	return nil
}

func CreateProject(projectName string, shipyardFilePath string) (string, error) {
	// The project name is prefixed with the keptn test namespace to avoid name collisions during parallel integration test runs on CI
	newProjectName := AddNamespaceToName(projectName)

	err := retry.Retry(func() error {
		if err := RecreateProjectUpstream(newProjectName); err != nil {
			return err
		}

		user := GetGiteaUser()
		token, err := GetGiteaToken()
		if err != nil {
			return err
		}
		var out string
		// apply the k8s job for creating the git upstream
		out, err = ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=%s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s", newProjectName, shipyardFilePath, user, newProjectName, user, token))

		if !strings.Contains(out, "created successfully") {
			return fmt.Errorf("unable to create project: %s", out)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return newProjectName, nil
}

func AddNamespaceToName(projectName string) string {
	return osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace) + "-" + projectName
}

func CreateProjectWithSSH(projectName string, shipyardFilePath string) (string, error) {
	// The project name is prefixed with the keptn test namespace to avoid name collisions during parallel integration test runs on CI
	newProjectName := osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace) + "-" + projectName

	err := retry.Retry(func() error {
		if err := RecreateProjectUpstream(newProjectName); err != nil {
			return err
		}

		user := GetGiteaUser()

		privateKey, passphrase, err := GetPrivateKeyAndPassphrase()
		if err != nil {
			return err
		}

		privateKeyPath := "private-key"
		err = os.WriteFile(privateKeyPath, []byte(privateKey), 0777)
		if err != nil {
			return err
		}

		defer func() {
			os.Remove(privateKeyPath)
		}()

		// apply the k8s job for creating the git upstream
		out, err := ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=%s --git-remote-url=ssh://gitea-ssh:22/%s/%s.git --git-user=%s --git-private-key=%s --git-private-key-pass=%s", newProjectName, shipyardFilePath, user, newProjectName, user, privateKeyPath, passphrase))

		if !strings.Contains(out, "created successfully") {
			return fmt.Errorf("unable to create project: %s", out)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return newProjectName, nil

}

func CreateProjectWithProxy(projectName string, shipyardFilePath string, proxyURL string) (string, error) {
	// The project name is prefixed with the keptn test namespace to avoid name collisions during parallel integration test runs on CI
	namespace := osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace)
	newProjectName := namespace + "-" + projectName

	err := retry.Retry(func() error {
		if err := RecreateProjectUpstream(newProjectName); err != nil {
			return nil
		}

		user := GetGiteaUser()
		token, err := GetGiteaToken()
		if err != nil {
			return err
		}

		// apply the k8s job for creating the git upstream
		out, err := ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=%s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s --git-proxy-url=%s --git-proxy-scheme=http --insecure-skip-tls", newProjectName, shipyardFilePath, user, newProjectName, user, token, proxyURL))

		if !strings.Contains(out, "created successfully") {
			return fmt.Errorf("unable to create project: %s", out)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return newProjectName, nil

}

func GetServiceExternalIP(namespace string, service string) (string, error) {
	ipAddr, err := ExecuteCommand(fmt.Sprintf("kubectl get svc %s -n %s -ojsonpath='{.status.loadBalancer.ingress[0].ip}'", service, namespace))
	if err != nil {
		return "", err
	}

	return removeQuotes(ipAddr), nil
}

func GetPrivateKeyAndPassphrase() (string, string, error) {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return "", "", err
	}

	giteaAccessSecret, err := clientset.CoreV1().Secrets(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), "gitea-access", v1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	privateKey := string(giteaAccessSecret.Data["private-key"])
	privateKeyPass := string(giteaAccessSecret.Data["private-key-pass"])
	if privateKey == "" || privateKeyPass == "" {
		return "", "", errors.New("no private key found")
	}

	return privateKey, privateKeyPass, nil
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

	eventContext := &models.EventContext{}
	err = resp.ToJSON(eventContext)
	if err != nil {
		return "", err
	}
	return *eventContext.KeptnContext, nil
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
	}, time.Minute, time.Second*3)

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

	subscriptionResponse := &models.CreateSubscriptionResponse{}

	err = resp.ToJSON(subscriptionResponse)
	require.Nil(t, err)

	require.NotEmpty(t, subscriptionResponse.ID)

	return subscriptionResponse.ID, nil
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

func ScaleUpUniform(deployments []string, replicas int) error {
	for _, deployment := range deployments {
		if err := keptnkubeutils.ScaleDeployment(false, deployment, GetKeptnNameSpaceFromEnv(), int32(replicas)); err != nil {
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

func storeWithCommit(t *testing.T, projectName, stage, serviceName, content, uri string) string {

	ctx, closeInternalKeptnAPI := context.WithCancel(context.Background())
	defer closeInternalKeptnAPI()
	internalKeptnAPI, err := GetInternalKeptnAPI(ctx, "service/configuration-service", "8889", "8080")
	require.Nil(t, err)
	t.Log("Storing new slo file")
	resp, err := internalKeptnAPI.Post(basePath+"/"+projectName+"/stage/"+stage+"/service/"+serviceName+"/resource", models.Resources{
		Resources: []*models.Resource{
			{
				ResourceContent: b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", content))),
				ResourceURI:     strutils.Stringp(uri),
			},
		},
	}, 3)
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)

	t.Logf("Received response %s", resp.String())
	require.Equal(t, 201, resp.Response().StatusCode)

	response := struct {
		CommitID string `json:"commitID"`
	}{}
	resp.ToJSON(&response)
	t.Log("Saved with commitID", response.CommitID)
	return response.CommitID
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
	}, timeout, 10*time.Second, GetDiagnostics("shipyard-controller", ""))
}

func doesSequenceHaveOneOfTheDesiredStates(state models.SequenceState, context *models.EventContext, desiredStates []string) bool {
	if state.Shkeptncontext == *context.KeptnContext {
		for _, desiredState := range desiredStates {
			if state.State == desiredState {
				return true
			}
		}
	}
	return false
}

func GetState(projectName string) (*models.SequenceStates, *req.Resp, error) {
	states := &models.SequenceStates{}

	resp, err := ApiGETRequest("/controlPlane/v1/sequence/"+projectName, 3)
	if err != nil {
		return nil, nil, err
	}
	err = resp.ToJSON(states)

	return states, resp, err
}

func GetStateByContext(projectName, keptnContext string) (*models.SequenceStates, *req.Resp, error) {
	states := &models.SequenceStates{}

	resp, err := ApiGETRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s?keptnContext=%s", projectName, keptnContext), 3)
	if err != nil {
		return nil, nil, err
	}
	err = resp.ToJSON(states)

	return states, resp, err
}

func GetProject(projectName string) (*models.ExpandedProject, error) {
	project := &models.ExpandedProject{}

	resp, err := ApiGETRequest("/controlPlane/v1/project/"+projectName, 3)
	if err != nil {
		return nil, err
	}

	err = resp.ToJSON(project)
	return project, err
}

func GetDiagnostics(service string, container string) string {
	if container == "" {
		container = service
	}
	outputBuilder := strings.Builder{}
	getLogsCmd := fmt.Sprintf("kubectl logs -n %s deployment/%s -c %s", GetKeptnNameSpaceFromEnv(), service, container)

	outputBuilder.WriteString(fmt.Sprintf("Logs of  of %s: \n\n", service))
	logOutput, err := ExecuteCommand(getLogsCmd)
	if err != nil {
		outputBuilder.WriteString(err.Error())
	}

	outputBuilder.WriteString(logOutput)
	outputBuilder.WriteString("\n-------------------------\n")
	getLogsCmd = fmt.Sprintf("kubectl logs -n %s deployment/%s -c %s --previous", GetKeptnNameSpaceFromEnv(), service, container)

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
	return WaitAndCheckDeployment(serviceName, projectName+"-"+stageName, time.Minute*5, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag})
}

func VerifyBlueGreenDeployment(serviceName, projectName, stageName, artifactImage, artifactTag string) error {
	if err := WaitAndCheckDeployment(serviceName, projectName+"-"+stageName, time.Minute*3, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag}); err != nil {
		return err
	}
	return WaitAndCheckDeployment(serviceName+"-primary", projectName+"-"+stageName, time.Minute*3, WaitForDeploymentOptions{WithImageName: artifactImage + ":" + artifactTag})
}

func VerifyTaskStartedEventExists(t *testing.T, keptnContext, projectName, stage string, taskName string) {
	var startedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Logf("verifying that "+taskName+".finished event for context %s does exist", keptnContext)
		taskStarted, err := GetLatestEventOfType(keptnContext, projectName, stage, keptnv2.GetStartedEventType(taskName))
		if err != nil || taskStarted == nil {
			return false
		}
		startedEvent = taskStarted
		return true
	}, 1*time.Minute, 10*time.Second)
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(startedEvent.Data, eventData)

	require.Nil(t, err)
}

func GetPublicURLOfService(serviceName, projectName, stageName string) (string, error) {
	ingressHostnameSuffix, err := GetFromConfigMap(GetKeptnNameSpaceFromEnv(), "ingress-config", func(data map[string]string) string { return data["ingress_hostname_suffix"] })
	if err != nil {
		return "", fmt.Errorf("unable to get public URL of service %s: %w", serviceName, err)
	}

	return fmt.Sprintf("http://%s.%s-%s.%s", serviceName, projectName, stageName, ingressHostnameSuffix), nil

}

// SetShipyardControllerEnvVar sets the provided value of the shipyard-controller deployment.
// This function is specific to the shipyard-controller, and eventually we should avoid setting env vars of deployments in general, as this
// leads to the respective pod being restarted, which increases the duration of the integration tests and prevents us from executing tests in parallel
func SetShipyardControllerEnvVar(t *testing.T, envVarName, envVarValue string) error {

	k8sClient, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return err
	}

	shipyardDeployment, err := k8sClient.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), "shipyard-controller", v1.GetOptions{})
	if err != nil {
		return err
	}

	if len(shipyardDeployment.Spec.Template.Spec.Containers) == 0 {
		return errors.New("shipyard deployment does not contain any container")
	}
	envVarFound := false
	for index, ev := range shipyardDeployment.Spec.Template.Spec.Containers[0].Env {
		if ev.Name == envVarName {
			envVarFound = true
			shipyardDeployment.Spec.Template.Spec.Containers[0].Env[index].Value = envVarValue
		}
	}

	if !envVarFound {
		shipyardDeployment.Spec.Template.Spec.Containers[0].Env = append(shipyardDeployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  envVarName,
			Value: envVarValue,
		})
	}

	_, err = k8sClient.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), shipyardDeployment, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	require.Eventually(t, func() bool {
		get, err := k8sClient.CoreV1().Pods(GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{LabelSelector: "app.kubernetes.io/name=shipyard-controller"})
		if err != nil {
			return false
		}

		if *shipyardDeployment.Spec.Replicas != int32(len(get.Items)) {
			// make sure only one pod is running
			return false
		}

		for _, item := range get.Items {
			if len(item.Spec.Containers) == 0 {
				continue
			}
			for _, env := range item.Spec.Containers[0].Env {
				if env.Name == envVarName && env.Value == envVarValue {
					if len(item.Status.ContainerStatuses) == 0 {
						return false
					}
					if item.Status.ContainerStatuses[0].State.Running != nil {
						return true
					}
				}
			}
		}
		return false
	}, 3*time.Minute, 10*time.Second)

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

func GetPodNamesOfDeployment(labelSelector string) ([]string, error) {
	k8sClient, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return nil, err
	}

	get, err := k8sClient.CoreV1().Pods(GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, nil
	}

	podNames := []string{}

	for _, pod := range get.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}

// SetRecreateUpgradeStrategyForDeployment sets the upgrade strategy of a deployment to "Recreate".
// Needed for our minishift tests right now, as there are problems with the RollingUpdate strategy of the shipyard-controller
// Should become obsolete when we switch to testing on an OpenShift 4.x cluster instead.
func SetRecreateUpgradeStrategyForDeployment(deploymentName string) error {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return err
	}

	deployment, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		return err
	}

	deployment.Spec.Strategy.Type = v12.RecreateDeploymentStrategyType
	deployment.Spec.Strategy.RollingUpdate = nil

	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), deployment, v1.UpdateOptions{})
	if err != nil {
		return nil
	}

	<-time.After(30 * time.Second)

	if err := WaitForDeploymentInNamespace(deploymentName, GetKeptnNameSpaceFromEnv()); err != nil {
		return err
	}

	return nil
}

func decodeBase64(str string) (string, error) {
	res, err := b64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func removeQuotes(str string) string {
	if str[0] == '"' || str[0] == '\'' {
		str = str[1:]
	}
	if i := len(str) - 1; str[i] == '"' || str[i] == '\'' {
		str = str[:i]
	}

	return str
}

// GetGiteaToken checks whether the GITEA_TOKEN environment variable is set. If yes, it will return the value for that var. If not, it will try to
// fetch the token from the secret 'gitea-access' in the keptn namespace
func GetGiteaToken() (string, error) {
	// if the token is set as an env var, return that
	if tokenFromEnv := os.Getenv("GITEA_TOKEN"); tokenFromEnv != "" {
		return tokenFromEnv, nil
	}
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return "", err
	}

	giteaAccessSecret, err := clientset.CoreV1().Secrets(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), "gitea-access", v1.GetOptions{})
	if err != nil {
		return "", err
	}

	token := string(giteaAccessSecret.Data["password"])
	if token == "" {
		return "", errors.New("no gitea token found")
	}

	return token, nil
}

// GetMongoDBCredentials retrieves the credentials of the mongodb user from the mongodb credentials secret
func GetMongoDBCredentials() (string, string, error) {
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return "", "", err
	}

	mongoDBSecret, err := clientset.CoreV1().Secrets(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), "mongodb-credentials", v1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	user := string(mongoDBSecret.Data["mongodb-root-user"])
	if user == "" {
		return "", "", errors.New("no mongodb user found")
	}

	password := string(mongoDBSecret.Data["mongodb-root-password"])
	if password == "" {
		return "", "", errors.New("no mongodb password found")
	}

	return user, password, nil
}

func GetGiteaUser() string {
	if os.Getenv("GITEA_ADMIN_USER") != "" {
		return os.Getenv("GITEA_ADMIN_USER")
	}
	return "gitea_admin"
}

// RecreateGitUpstreamRepository creates a kubernetes job that (re)creates the upstream repo for a project on the internal gitea instance
func RecreateGitUpstreamRepository(project string) error {
	ctx, closeInternalKeptnAPI := context.WithCancel(context.Background())
	defer closeInternalKeptnAPI()
	internalKeptnAPI, err := GetInternalKeptnAPI(ctx, "service/gitea-http", "3002", "3000")

	if err != nil {
		return err
	}

	token, err := GetGiteaToken()
	if err != nil {
		return err
	}
	user := GetGiteaUser()

	// first, check if the repository already exists
	repoAPIPath := fmt.Sprintf("/api/v1/repos/%s/%s?access_token=%s", user, project, token)
	resp, err := internalKeptnAPI.Get(repoAPIPath, 5)
	if err != nil {
		return err
	}

	if resp.Response().StatusCode >= 200 && resp.Response().StatusCode < 300 {
		// if yes, delete the existing project
		resp, err = internalKeptnAPI.Delete(repoAPIPath, 5)
		if err != nil {
			return fmt.Errorf("could not delete existing repository for project %s: %w", project, err)
		}
	}

	projectPayload := map[string]interface{}{
		"name":           project,
		"description":    "Upstream repo for test project",
		"default_branch": "main",
	}
	// now create the repo
	post, err := internalKeptnAPI.Post(fmt.Sprintf("/api/v1/user/repos?access_token=%s", token), projectPayload, 5)

	if err != nil {
		return err
	}

	if post.Response().StatusCode < 200 || post.Response().StatusCode >= 300 {
		var responsePayload []byte
		if _, err := post.Response().Body.Read(responsePayload); err != nil {
			return fmt.Errorf("could not create upstream repo for project %s", project)
		}
		return fmt.Errorf("could not create upstream repo for project %s: %s", project, string(responsePayload))
	}

	return nil
}

func WaitForDeploymentToBeScaledDown(deploymentName string) error {
	// if the token is set as an env var, return that
	clientset, err := keptnkubeutils.GetClientset(false)
	if err != nil {
		return err
	}

	err = retry.Retry(func() error {
		pods, err := clientset.CoreV1().Pods(GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return err
		}

		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, deploymentName) {
				return fmt.Errorf("Pod for deployment %s still present", deploymentName)
			}
		}
		return nil
	}, retry.NumberOfRetries(40))

	return err
}

func checkResourceInResponse(resources models.Resources, resourceName string) error {
	for _, resource := range resources.Resources {
		if *resource.ResourceURI == resourceName {
			return nil
		}
	}

	return fmt.Errorf("Resource %s not found in received response.", resourceName)
}

func resetTestPath(t *testing.T, path string) {
	err := os.Chdir(path)
	require.Nil(t, err)
}
