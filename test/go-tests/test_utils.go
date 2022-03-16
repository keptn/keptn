package go_tests

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/common/retry"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
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
		return req.Get(a.baseURL+path, a.getAuthHeader())
	}, retries)
}

func (a *APICaller) Delete(path string, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Delete(a.baseURL+path, a.getAuthHeader())
	}, retries)
}

func (a *APICaller) Put(path string, payload interface{}, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
		return req.Put(a.baseURL+path, a.getAuthHeader(), req.BodyJSON(payload))
	}, retries)
}

func (a *APICaller) Post(path string, payload interface{}, retries int) (*req.Resp, error) {
	return a.doHTTPRequestWithRetry(func() (*req.Resp, error) {
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

func CreateProject(projectName string, shipyardFilePath string, recreateIfAlreadyThere bool) (string, error) {

	retries := 5
	var err error
	var resp *req.Resp

	// The project name is prefixed with the keptn test namespace to avoid name collisions during parallel integration test runs on CI
	newProjectName := osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace) + "-" + projectName

	for i := 0; i < retries; i++ {
		if err != nil {
			<-time.After(10 * time.Second)
		}
		resp, err = ApiGETRequest("/controlPlane/v1/project/"+newProjectName, 3)
		if err != nil {
			continue
		}

		if resp.Response().StatusCode == http.StatusOK {
			if recreateIfAlreadyThere {
				// delete project if it exists
				_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", newProjectName))
				if err != nil {
					continue
				}
			} else {
				return "", errors.New("project already exists")
			}
		}

		err = RecreateGitUpstreamRepository(newProjectName)
		if err != nil {
			// retry if repo creation failed (gitea might not be available)
			continue
		}

		user := GetGiteaUser()
		token, err := GetGiteaToken()
		if err != nil {
			return "", err
		}

		// apply the k8s job for creating the git upstream
		_, err = ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=%s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s", newProjectName, shipyardFilePath, user, newProjectName, user, token))

		if err == nil {
			return newProjectName, nil
		}
	}

	return "", err
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
			m, _ := json.MarshalIndent(state, "", "  ")
			t.Logf("%s", m)
			if doesSequenceHaveOneOfTheDesiredStates(state, context, desiredStates) {
				return true
			}
		}
		return false
	}, timeout, 10*time.Second, GetDiagnostics("shipyard-controller", ""))
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

func GetStateByContext(projectName, keptnContext string) (*scmodels.SequenceStates, *req.Resp, error) {
	states := &scmodels.SequenceStates{}

	resp, err := ApiGETRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s?keptnContext=%s", projectName, keptnContext), 3)
	err = resp.ToJSON(states)

	return states, resp, err
}

func GetProject(projectName string) (*scmodels.ExpandedProject, error) {
	project := &scmodels.ExpandedProject{}

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

func GetGiteaUser() string {
	if os.Getenv("GITEA_ADMIN_USER") != "" {
		return os.Getenv("GITEA_ADMIN_USER")
	}
	return "gitea_admin"
}

// recreateUpstreamRepository creates a kubernetes job that (re)creates the upstream repo for a project on the internal gitea instance
func RecreateGitUpstreamRepository(project string) error {
	jobName := "recreate-upstream-repo"
	clientset, err := keptnkubeutils.GetClientset(false)

	// check if the job for recreating the project already exists (e.g. due to a previous run
	get, err := clientset.BatchV1().Jobs(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), jobName, v1.GetOptions{})
	if err == nil && get != nil {
		err := clientset.BatchV1().Jobs(GetKeptnNameSpaceFromEnv()).Delete(context.TODO(), jobName, v1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("could not delete previous instance of job %s: %w", project, err)
		}
	}

	token, err := GetGiteaToken()
	if err != nil {
		return err
	}
	user := GetGiteaUser()

	deleteCmd := fmt.Sprintf(`curl -X DELETE "http://gitea-http:3000/api/v1/repos/%s/%s?access_token=%s"`, user, project, token)
	createCmd := fmt.Sprintf(`curl -X POST "http://gitea-http:3000/api/v1/user/repos?access_token=%s" -H "accept: application/json" -H "content-type: application/json" -d "{\"name\":\"%s\", \"description\": \"Sample description\", \"default_branch\": \"main\"}"`, token, project)

	parallelPods := int32(1)

	recreateJob := &batchv1.Job{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      jobName,
			Namespace: GetKeptnNameSpaceFromEnv(),
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "go-tests",
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelPods,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "delete-previous",
							Image: "curlimages/curl:7.80.0",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								deleteCmd,
							},
						},
						{
							Name:  "create",
							Image: "curlimages/curl:7.80.0",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								createCmd,
							},
						},
					},
				},
			},
		},
	}

	_, err = clientset.BatchV1().Jobs(GetKeptnNameSpaceFromEnv()).Create(context.TODO(), recreateJob, v1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("could not create job to create upstream repo: %w", err)
	}

	defer func() {
		_ = clientset.BatchV1().Jobs(GetKeptnNameSpaceFromEnv()).Delete(context.TODO(), jobName, v1.DeleteOptions{})
	}()

	var jobStatus *batchv1.Job
	// wait for the job to be completed
	err = retry.Retry(func() error {
		get, err := clientset.BatchV1().Jobs(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), jobName, v1.GetOptions{})
		if err != nil {
			return err
		}

		if get.Status.CompletionTime == nil {
			// job not done yet
			return errors.New("job not completed")
		}

		jobStatus = get
		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(5*time.Second))
	if err != nil {
		return err
	}

	if jobStatus.Status.Failed > 1 {
		return errors.New("could not (re)create upstream repository")
	}

	return nil
}

func checkResourceInResponse(resources models.Resources, resourceName string) error {
	for _, resource := range resources.Resources {
		if *resource.ResourceURI == resourceName {
			return nil
		}
	}

	return fmt.Errorf("Resource %s not found in received response.", resourceName)
}
