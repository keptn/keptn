package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"

	"github.com/keptn/keptn/helm-service/controller/mesh"
)

// Onboarder is a container of variables required for onboarding a new service
type Onboarder struct {
	mesh           mesh.Mesh
	logger         keptnutils.LoggerInterface
	canaryLevelGen helm.CanaryLevelGenerator
	keptnDomain    string
}

// NewOnboarder creates a new Onboarder
func NewOnboarder(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string) *Onboarder {
	return &Onboarder{mesh: mesh, canaryLevelGen: canaryLevelGen, logger: logger, keptnDomain: keptnDomain}
}

// DoOnboard onboards a new service
func (o *Onboarder) DoOnboard(ce cloudevents.Event) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(o.mesh)

	event := &keptnevents.ServiceCreateEventData{}
	if err := ce.DataAs(event); err != nil {
		o.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" {
		deplStrategies, err := getDeploymentStrategies(event.Project)
		if err != nil {
			o.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s" + err.Error()))
			return err
		}
		event.DeploymentStrategies = deplStrategies
	}

	o.logger.Info(fmt.Sprintf("Start creating service %s in project %s", event.Service, event.Project))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	stageHandler := keptnutils.NewStageHandler(url.String())
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		o.logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	firstService, err := o.isFirstServiceOfProject(event, stages)
	if err != nil {
		o.logger.Error("Error when checking whether any service was created before: " + err.Error())
		return err
	}
	if firstService {
		o.logger.Info("Create Helm Umbrella charts")
		if err := o.initAndApplyUmbrellaChart(event, umbrellaChartHandler, stages); err != nil {
			o.logger.Error(fmt.Sprintf("Error when initalizing and applying umbrella charts for project %s: %s", event.Project, err.Error()))
			return err
		}
	}

	serviceHandler := keptnutils.NewServiceHandler(url.String())
	helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
	if err != nil {
		o.logger.Error("Error when decoding the Helm chart")
	}

	for _, stage := range stages {
		o.logger.Debug("Creating new keptn service " + event.Service + " in stage " + stage.StageName)
		serviceHandler.CreateService(event.Project, stage.StageName, event.Service)

		o.logger.Debug("Storing the Helm chart provided by the user in stage " + stage.StageName)
		if err := keptnutils.StoreChart(event.Project, event.Service, stage.StageName, helm.GetChartName(event.Service, false),
			helmChartData, url.String()); err != nil {
			o.logger.Error("Error when storing the Helm chart: " + err.Error())
			return err
		}

		o.logger.Debug("Updating the Umbrella chart with the new Helm chart in stage " + stage.StageName)
		// if err := helm.AddChartInUmbrellaRequirements(event.Project, helm.GetChartName(event.Service, false), stage, url.String()); err != nil {
		// 	o.logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
		// 	return err
		// }
		if err := umbrellaChartHandler.AddChartInUmbrellaValues(event.Project, helm.GetChartName(event.Service, false), stage); err != nil {
			o.logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
			return err
		}

		if _, ok := event.DeploymentStrategies[stage.StageName]; !ok {
			o.logger.Error("Received event does not define deployment strategy for stage " + stage.StageName +
				". Hence, a direct strategy is used.")
		}

		if event.DeploymentStrategies[stage.StageName] == keptnevents.Duplicate {
			chartGenerator := helm.NewGeneratedChartHandler(o.mesh, o.canaryLevelGen, o.keptnDomain)

			o.logger.Debug("Generating the keptn-managed Helm chart for stage " + stage.StageName)
			ch, err := keptnutils.LoadChart(helmChartData)
			if err != nil {
				o.logger.Error("Error when loading chart: " + err.Error())
				return err
			}
			generatedChartData, err := chartGenerator.GenerateManagedChart(ch, event.Project, stage.StageName)
			if err != nil {
				o.logger.Error("Error when generating the keptn managed chart: " + err.Error())
				return err
			}

			o.logger.Debug("Storing the keptn generated Helm chart in stage " + stage.StageName)
			if err := keptnutils.StoreChart(event.Project, event.Service, stage.StageName, helm.GetChartName(event.Service, true),
				generatedChartData, url.String()); err != nil {
				o.logger.Error("Error when storing the Helm chart: " + err.Error())
				return err
			}

			o.logger.Debug("Updating the Umbrella chart with the new Helm chart in stage " + stage.StageName)
			// if err := helm.AddChartInUmbrellaRequirements(event.Project, helm.GetChartName(event.Service, true), stage, url.String()); err != nil {
			// 	o.logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
			// 	return err
			// }
			if err := umbrellaChartHandler.AddChartInUmbrellaValues(event.Project, helm.GetChartName(event.Service, true), stage); err != nil {
				o.logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
				return err
			}
		}
	}

	o.logger.Info(fmt.Sprintf("Finished creating service %s in project %s", event.Service, event.Project))
	return nil
}

func (o *Onboarder) initAndApplyUmbrellaChart(event *keptnevents.ServiceCreateEventData,
	umbrellaChartHandler *helm.UmbrellaChartHandler, stages []*models.Stage) error {

	// Initalize the umbrella chart
	if err := umbrellaChartHandler.InitUmbrellaChart(event, stages); err != nil {
		return fmt.Errorf("Error when initializing the umbrella chart: %s", err.Error())
	}

	for _, stage := range stages {
		// Apply the umbrella chart
		umbrellaChart, err := ioutil.TempDir("", "")
		if err != nil {
			return fmt.Errorf("Error when creating a temporary directory: %s", err.Error())
		}
		if err := umbrellaChartHandler.GetUmbrellaChart(umbrellaChart, event.Project, stage.StageName); err != nil {
			return fmt.Errorf("Error when getting umbrella chart: %s", err)
		}

		if err := ApplyDirectory(umbrellaChart, helm.GetUmbrellaReleaseName(event.Project, stage.StageName),
			helm.GetUmbrellaNamespace(event.Project, stage.StageName)); err != nil {
			return fmt.Errorf("Error when applying umbrella chart in stage %s: %s", stage.StageName, err.Error())
		}
		if err := os.RemoveAll(umbrellaChart); err != nil {
			return err
		}
	}
	return nil
}

func (o *Onboarder) isFirstServiceOfProject(event *keptnevents.ServiceCreateEventData, stages []*models.Stage) (bool, error) {

	if len(stages) == 0 {
		return false, errors.New("Cannot onboard service because no stage is available")
	}
	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return false, err
	}
	svcHandler := keptnutils.NewServiceHandler(url.String())
	// Use any stage for checking whether there is already a service created
	services, err := svcHandler.GetAllServices(event.Project, stages[0].StageName)
	if err != nil {
		return false, err
	}
	return len(services) == 0, nil
}
