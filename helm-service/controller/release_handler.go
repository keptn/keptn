package controller

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

// ReleaseHandler is a handler for releasing a service
type ReleaseHandler struct {
	Handler
	mesh                  mesh.Mesh
	generatedChartHandler helm.ChartGenerator
	configurationChanger  configurationchanger.IConfigurationChanger
	chartStorer           keptnutils.ChartStorer
	chartPackager         keptnutils.ChartPackager
}

// NewReleaseHandler creates a ReleaseHandler
func NewReleaseHandler(keptnHandler *keptnv2.Keptn,
	mesh mesh.Mesh,
	configurationChanger configurationchanger.IConfigurationChanger,
	chartGenerator helm.ChartGenerator,
	chartStorer keptnutils.ChartStorer,
	chartPackager keptnutils.ChartPackager,
	configServiceURL string) *ReleaseHandler {
	//generatedChartHandler := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	return &ReleaseHandler{
		Handler:               NewHandlerBase(keptnHandler, configServiceURL),
		mesh:                  mesh,
		generatedChartHandler: chartGenerator,
		configurationChanger:  configurationChanger,
		chartStorer:           chartStorer,
		chartPackager:         chartPackager,
	}
}

// HandleEvent handles release.triggered events and either promotes or aborts an artifact
func (h *ReleaseHandler) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) {

	defer closeLogger(h.getKeptnHandler())

	e := keptnv2.ReleaseTriggeredEventData{}
	if err := ce.DataAs(&e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Send release started event
	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Starting release for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
	if err := h.sendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.ReleaseTaskName), h.getStartedEventData(e.EventData)); err != nil {
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	deploymentStrategy, err := keptnevents.GetDeploymentStrategy(e.Deployment.DeploymentStrategy)
	if err != nil {
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	var gitVersion string
	if deploymentStrategy == keptnevents.Duplicate {
		// Only in case of a duplicate deployment strategy, the user-chart has to be promoted/aborted and
		// a traffic switch is necessary
		if e.Result == keptnv2.ResultPass || e.Result == keptnv2.ResultWarning {
			h.getKeptnHandler().Logger.Info(fmt.Sprintf("Promote service %s in stage %s of project %s",
				e.Service, e.Stage, e.Project))
			gitVersion, err = h.promoteDeployment(e.EventData)
			if err != nil {
				h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
				return
			}
		} else {
			h.getKeptnHandler().Logger.Info(fmt.Sprintf("Rollback service %s in stage %s of project %s",
				e.Service, e.Stage, e.Project))
			gitVersion, err = h.rollbackDeployment(e.EventData)
			if err != nil {
				h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
				return
			}
		}
	} else {
		h.getKeptnHandler().Logger.Info(fmt.Sprintf(
			"No release action required, as the service %s in stage %s of project %s has a direct deployment strategy",
			e.Service, e.Stage, e.Project))
	}

	// Send finished event
	data := h.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded, e.Result, "Finished release", gitVersion)
	if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), data); err != nil {
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Finished release for service %s in stage %s and project %s", e.Service, e.Stage, e.Project))
}

func (h *ReleaseHandler) rollbackDeployment(e keptnv2.EventData) (string, error) {

	canaryWeightTo0Updater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 0)
	genChart, gitVersion, err := configurationchanger.NewConfigurationChanger(h.getConfigServiceURL()).UpdateChart(e,
		true, canaryWeightTo0Updater)
	if err != nil {
		return "", err
	}

	// Upgrade generated chart
	if err := h.upgradeChart(genChart, e, keptnevents.Duplicate); err != nil {
		return "", err
	}

	userChart, err := h.getUserChart(e)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChartWithReplicas(userChart, e, keptnevents.Duplicate, 0); err != nil {
		return "", err
	}
	return gitVersion, nil
}

func (h *ReleaseHandler) promoteDeployment(e keptnv2.EventData) (string, error) {

	//configChanger := configurationchanger.NewConfigurationChanger(h.getConfigServiceURL())

	// Switch weight to 100% canary, 0% primary
	canaryWeightTo100Updater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 100)
	genChart, _, err := h.configurationChanger.UpdateChart(e, true, canaryWeightTo100Updater)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChart(genChart, e, keptnevents.Duplicate); err != nil {
		return "", err
	}

	// Update and apply new generated chart
	if err := h.updateGeneratedChart(e); err != nil {
		return "", err
	}

	// Switch weight to 0% canary, 100% primary
	canaryWeightTo0Updater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 0)
	genChart, gitVersion, err := h.configurationChanger.UpdateChart(e, true, canaryWeightTo0Updater)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChart(genChart, e, keptnevents.Duplicate); err != nil {
		return "", err
	}

	// Scale down replicas of user chart
	userChart, err := h.getUserChart(e)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChartWithReplicas(userChart, e, keptnevents.Duplicate, 0); err != nil {
		return "", err
	}
	return gitVersion, nil
}

func (h *ReleaseHandler) updateGeneratedChart(e keptnv2.EventData) error {

	canaryWeightTo100Updater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 100)
	//chartGenerator := helm.NewGeneratedChartGenerator(h.mesh, h.getKeptnHandler().Logger)
	userChartManifest, err := h.getHelmExecutor().GetManifest(helm.GetReleaseName(e.Project, e.Stage, e.Service, false),
		e.Project+"-"+e.Stage)
	if err != nil {
		return err
	}
	newGenChart, err := h.generatedChartHandler.GenerateDuplicateChart(userChartManifest, e.Project, e.Stage, e.Service)
	if err != nil {
		return err
	}
	if err := canaryWeightTo100Updater.Manipulate(newGenChart); err != nil {
		return err
	}
	genChartData, err := h.chartPackager.PackageChart(newGenChart)
	if err != nil {
		return err
	}
	if _, err := h.chartStorer.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData); err != nil {
		return err
	}
	return h.upgradeChart(newGenChart, e, keptnevents.Duplicate)
}

func (h *ReleaseHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.ReleaseStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.ReleaseStartedEventData{EventData: inEventData}
}

func (h *ReleaseHandler) getFinishedEventData(inEventData keptnv2.EventData, status keptnv2.StatusType, result keptnv2.ResultType,
	message string, gitCommit string) keptnv2.ReleaseFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.ReleaseFinishedEventData{
		EventData: inEventData,
		Release:   keptnv2.ReleaseData{GitCommit: gitCommit},
	}
}

func (h *ReleaseHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.ReleaseFinishedEventData {
	return h.getFinishedEventData(eventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error(), "")
}
