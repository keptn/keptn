package controller

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/common"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	keptntypes "github.com/keptn/keptn/helm-service/pkg/types"
)

// ReleaseHandler is a handler for releasing a service
type ReleaseHandler struct {
	Handler
	mesh                  mesh.Mesh
	generatedChartHandler helm.ChartGenerator
	configurationChanger  configurationchanger.IConfigurationChanger
	chartStorer           keptntypes.IChartStorer
	chartPackager         keptntypes.IChartPackager
}

// NewReleaseHandler creates a ReleaseHandler
func NewReleaseHandler(keptnHandler Handler,
	mesh mesh.Mesh,
	configurationChanger configurationchanger.IConfigurationChanger,
	chartGenerator helm.ChartGenerator,
	chartStorer keptntypes.IChartStorer,
	chartPackager keptntypes.IChartPackager) *ReleaseHandler {
	//generatedChartHandler := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	return &ReleaseHandler{
		Handler:               keptnHandler,
		mesh:                  mesh,
		generatedChartHandler: chartGenerator,
		configurationChanger:  configurationChanger,
		chartStorer:           chartStorer,
		chartPackager:         chartPackager,
	}
}

// HandleEvent handles release.triggered events and either promotes or aborts an artifact
func (h *ReleaseHandler) HandleEvent(ce cloudevents.Event) {
	e := keptnv2.ReleaseTriggeredEventData{}
	if err := ce.DataAs(&e); err != nil {
		err = fmt.Errorf("Failed to unmarshal data: unable to convert json data from cloudEvent to release event")
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
	// retrieve commitId from sequence
	commitID := retrieveCommit(ce)

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

	if deploymentStrategy == keptnevents.Duplicate {
		// Only in case of a duplicate deployment strategy, the user-chart has to be promoted/aborted and
		// a traffic switch is necessary
		if e.Result == keptnv2.ResultPass || e.Result == keptnv2.ResultWarning {
			h.getKeptnHandler().Logger.Info(fmt.Sprintf("Promote service %s in stage %s of project %s",
				e.Service, e.Stage, e.Project))
			commitID, err = h.promoteDeployment(e.EventData, commitID)
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
	data := h.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded, e.Result, "Finished release")
	if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), data); err != nil {
		h.handleError(ce.ID(), err, keptnv2.ReleaseTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Finished release for service %s in stage %s and project %s", e.Service, e.Stage, e.Project))
}

func (h *ReleaseHandler) promoteDeployment(e keptnv2.EventData, commitID string) (string, error) {

	//configChanger := configurationchanger.NewConfigurationChanger(h.getConfigServiceURL())

	// Switch weight to 100% canary, 0% primary
	canaryWeightTo100Updater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 100)
	genChart, _, err := h.getGeneratedChart(e, commitID)
	if err != nil {
		return "", err
	}
	genChart, newCommit, err := h.configurationChanger.UpdateLoadedChart(genChart, e, true, canaryWeightTo100Updater)
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
	genChart, _, err = h.getGeneratedChart(e, newCommit)
	if err != nil {
		return "", err
	}
	genChart, commitID, err = h.configurationChanger.UpdateLoadedChart(genChart, e, true, canaryWeightTo0Updater)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChart(genChart, e, keptnevents.Duplicate); err != nil {
		return "", err
	}

	// Scale down replicas of user chart
	userChart, _, err := h.getUserChart(e, commitID)
	if err != nil {
		return "", err
	}
	if err := h.upgradeChartWithReplicas(userChart, e, keptnevents.Duplicate, 0); err != nil {
		return "", err
	}
	return commitID, nil
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
	genChartData, err := h.chartPackager.Package(newGenChart)
	if err != nil {
		return err
	}

	storeOpts := common.StoreChartOptions{
		Project:   e.Project,
		Service:   e.Service,
		Stage:     e.Stage,
		ChartName: helm.GetChartName(e.Service, true),
		HelmChart: genChartData,
	}

	if _, err := h.chartStorer.Store(storeOpts); err != nil {
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
	message string) keptnv2.ReleaseFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.ReleaseFinishedEventData{
		EventData: inEventData,
	}
}

func (h *ReleaseHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.ReleaseFinishedEventData {
	return h.getFinishedEventData(eventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
}
