package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"

	"helm.sh/helm/v3/pkg/chart"

	"github.com/ghodss/yaml"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller/helm"
)

// ActionTriggeredHandler handles sh.keptn.events.action.triggered events for scaling
type ActionTriggeredHandler struct {
	keptnHandler     *keptn.Keptn
	helmExecutor     helm.HelmExecutor
	configServiceURL string
}

// ActionScaling is the identifier for the scaling action
const ActionScaling = "scaling"

// NewActionTriggeredHandler creates a new ActionTriggeredHandler
func NewActionTriggeredHandler(keptnHandler *keptn.Keptn,
	configServiceURL string) *ActionTriggeredHandler {
	helmExecutor := helm.NewHelmV3Executor(keptnHandler.Logger)
	return &ActionTriggeredHandler{keptnHandler: keptnHandler, helmExecutor: helmExecutor,
		configServiceURL: configServiceURL}
}

// HandleEvent takes the sh.keptn.events.action.triggered event and performs the scaling action
func (a *ActionTriggeredHandler) HandleEvent(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()
	actionTriggeredEvent := keptn.ActionTriggeredEventData{}

	err := ce.DataAs(&actionTriggeredEvent)
	if err != nil {
		errMsg := "action.triggered event not well-formed: " + err.Error()
		a.keptnHandler.Logger.Error(errMsg)
		return errors.New(errMsg)
	}

	if actionTriggeredEvent.Action.Action == ActionScaling {
		// Send action.started event
		if sendErr := a.sendEvent(ce, keptn.ActionStartedEventType, a.getActionStartedEvent(actionTriggeredEvent)); sendErr != nil {
			a.keptnHandler.Logger.Error(sendErr.Error())
			return errors.New(sendErr.Error())
		}

		resp := a.handleScaling(actionTriggeredEvent)
		if resp.Action.Status == keptn.ActionStatusErrored {
			a.keptnHandler.Logger.Error(fmt.Sprintf("action %s failed with result %s", actionTriggeredEvent.Action.Action, resp.Action.Result))
		} else {
			a.keptnHandler.Logger.Info(fmt.Sprintf("Finished action with status %s and result %s", resp.Action.Status, resp.Action.Result))
		}

		// Send action.finished event
		if sendErr := a.sendEvent(ce, keptn.ActionFinishedEventType, resp); sendErr != nil {
			a.keptnHandler.Logger.Error(sendErr.Error())
			return errors.New(sendErr.Error())
		}
	} else {
		a.keptnHandler.Logger.Info("Received unhandled action: " + actionTriggeredEvent.Action.Action + ". Exiting")
		return nil
	}

	return nil
}

func (a *ActionTriggeredHandler) getActionFinishedEvent(result keptn.ActionResultType, status keptn.ActionStatusType,
	actionTriggeredEvent keptn.ActionTriggeredEventData) keptn.ActionFinishedEventData {

	return keptn.ActionFinishedEventData{
		Project: actionTriggeredEvent.Project,
		Service: actionTriggeredEvent.Service,
		Stage:   actionTriggeredEvent.Stage,
		Action: keptn.ActionResult{
			Result: result,
			Status: status,
		},
		Labels: actionTriggeredEvent.Labels,
	}
}

func (a *ActionTriggeredHandler) getActionStartedEvent(actionTriggeredEvent keptn.ActionTriggeredEventData) keptn.ActionStartedEventData {

	return keptn.ActionStartedEventData{
		Project: actionTriggeredEvent.Project,
		Service: actionTriggeredEvent.Service,
		Stage:   actionTriggeredEvent.Stage,
		Labels:  actionTriggeredEvent.Labels,
	}
}

func (a *ActionTriggeredHandler) handleScaling(actionTriggeredEvent keptn.ActionTriggeredEventData) keptn.ActionFinishedEventData {

	value, ok := actionTriggeredEvent.Action.Value.(string)
	if !ok {
		return a.getActionFinishedEvent("could not parse action.value to string value",
			keptn.ActionStatusErrored, actionTriggeredEvent)
	}
	replicaIncrement, err := strconv.Atoi(value)
	if err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()),
			keptn.ActionStatusErrored, actionTriggeredEvent)
	}

	// Get generated chart
	helmChartName := helm.GetChartName(actionTriggeredEvent.Service, true)
	a.keptnHandler.Logger.Info(fmt.Sprintf("Retrieve chart %s of stage %s", helmChartName, actionTriggeredEvent.Stage))

	ch, err := keptnutils.GetChart(actionTriggeredEvent.Project, actionTriggeredEvent.Service, actionTriggeredEvent.Stage, helmChartName, a.configServiceURL)
	if err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()), keptn.ActionStatusErrored, actionTriggeredEvent)
	}
	deploymentStrategy, err := getDeploymentStrategyOfService(ch)
	if err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()), keptn.ActionStatusErrored, actionTriggeredEvent)
	}

	// Edit chart
	a.keptnHandler.Logger.Info(fmt.Sprintf("Edit chart %s of stage %s", helmChartName, actionTriggeredEvent.Stage))
	if err := a.increaseReplicaCount(ch, replicaIncrement); err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType("failed when editing deployment: "+err.Error()),
			keptn.ActionStatusErrored, actionTriggeredEvent)
	}

	// Upgrade chart
	a.keptnHandler.Logger.Info(fmt.Sprintf("Start upgrading chart %s of stage %s", helmChartName, actionTriggeredEvent.Stage))
	if err := a.upgradeChart(ch, actionTriggeredEvent, deploymentStrategy); err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()), keptn.ActionStatusErrored, actionTriggeredEvent)
	}
	a.keptnHandler.Logger.Info(fmt.Sprintf("Finished upgrading chart %s of stage %s", helmChartName, actionTriggeredEvent.Stage))

	// Store chart
	a.keptnHandler.Logger.Info(fmt.Sprintf("Store chart %s of stage %s", helmChartName, actionTriggeredEvent.Stage))
	chartData, err := keptnutils.PackageChart(ch)
	if err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()), keptn.ActionStatusErrored, actionTriggeredEvent)
	}
	if err := keptnutils.StoreChart(actionTriggeredEvent.Project, actionTriggeredEvent.Service, actionTriggeredEvent.Stage,
		helmChartName, chartData, a.configServiceURL); err != nil {
		return a.getActionFinishedEvent(keptn.ActionResultType(err.Error()), keptn.ActionStatusErrored, actionTriggeredEvent)
	}

	return a.getActionFinishedEvent(keptn.ActionResultPass, keptn.ActionStatusSucceeded, actionTriggeredEvent)
}

func (a *ActionTriggeredHandler) upgradeChart(ch *chart.Chart, action keptn.ActionTriggeredEventData, strategy keptn.DeploymentStrategy) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	return a.helmExecutor.UpgradeChart(ch,
		helm.GetReleaseName(action.Project, action.Stage, action.Service, generated),
		action.Project+"-"+action.Stage,
		getKeptnValues(action.Project, action.Stage, action.Service,
			getDeploymentName(strategy, generated)))
}

// increaseReplicaCount increases the replica count in the deployments by the provided replicaIncrement
func (a *ActionTriggeredHandler) increaseReplicaCount(ch *chart.Chart, replicaIncrement int) error {

	for _, template := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(template.Data))
		newContent := make([]byte, 0, 0)
		containsDepl := false
		for {
			var document interface{}
			err := dec.Decode(&document)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			doc, err := json.Marshal(document)
			if err != nil {
				return err
			}

			var depl appsv1.Deployment
			if err := json.Unmarshal(doc, &depl); err == nil && keptnutils.IsDeployment(&depl) {
				// Deployment found
				containsDepl = true
				depl.Spec.Replicas = getPtr(*depl.Spec.Replicas + int32(replicaIncrement))
				newContent, err = appendAsYaml(newContent, depl)
				if err != nil {
					return err
				}
			} else {
				newContent, err = appendAsYaml(newContent, document)
				if err != nil {
					return err
				}
			}
		}
		if containsDepl {
			template.Data = newContent
		}
	}

	return nil
}

func getPtr(x int32) *int32 {
	return &x
}

func appendAsYaml(content []byte, element interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}

func (a ActionTriggeredHandler) sendEvent(ce cloudevents.Event, eventType string, data interface{}) error {
	keptnHandler, err := keptn.NewKeptn(&ce, keptn.KeptnOpts{
		EventBrokerURL: os.Getenv("EVENTBROKER"),
	})
	if err != nil {
		a.keptnHandler.Logger.Error("Could not initialize Keptn handler: " + err.Error())
		return err
	}

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        eventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext, "triggerid": ce.ID()},
		}.AsV02(),
		Data: data,
	}

	err = keptnHandler.SendCloudEvent(event)
	if err != nil {
		a.keptnHandler.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}
