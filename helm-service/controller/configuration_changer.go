package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"k8s.io/helm/pkg/proto/hapi/chart"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/ghodss/yaml"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
	"github.com/tidwall/sjson"
	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
)

func init() {
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
}

// ConfigurationChanger is a container of variables required for changing the configuration of a service
type ConfigurationChanger struct {
	mesh                  mesh.Mesh
	generatedChartHandler *helm.GeneratedChartHandler
	canaryLevelGen        helm.CanaryLevelGenerator
	logger                keptnutils.LoggerInterface
	keptnDomain           string
}

// NewConfigurationChanger creates a new ConfigurationChanger
func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, canaryLevelGen, keptnDomain)
	return &ConfigurationChanger{mesh: mesh, generatedChartHandler: generatedChartHandler, canaryLevelGen: canaryLevelGen,
		logger: logger, keptnDomain: keptnDomain}
}

// ChangeAndApplyConfiguration changes the configuration and applies it in the cluster
func (c *ConfigurationChanger) ChangeAndApplyConfiguration(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()

	e := &keptnevents.ConfigurationChangeEventData{}
	if err := ce.DataAs(e); err != nil {
		c.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && e.Stage == "" {
		stage, err := getFirstStage(e.Project)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error when reading shipyard: %s" + err.Error()))
			return err
		}
		e.Stage = stage
	}

	genChart, err := getGeneratedChart(e)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	deploymentStrategy, err := getDeploymentStrategyOfService(genChart)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	if len(e.ValuesCanary) > 0 {

		ch, err := c.updateChart(e, false, changeValue)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		upgradeMsg, err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, false)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}

		onboarder := NewOnboarder(c.mesh, c.canaryLevelGen, c.logger, c.keptnDomain)
		if onboarder.IsGeneratedChartEmpty(genChart) {
			genChart, err = onboarder.OnboardGeneratedService(upgradeMsg, e.Project, e.Stage, e.Service, deploymentStrategy)
			if err != nil {
				c.logger.Error(err.Error())
				return err
			}
			if deploymentStrategy == keptnevents.Direct {
				_, err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, true)
				if err != nil {
					c.logger.Error(err.Error())
					return err
				}
			}
		}
	}
	if len(e.DeploymentChanges) > 0 {
		ch, err := c.updateChart(e, true, changePrimaryDeployment)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if _, err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, true); err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}

	// Change canary
	if e.Canary != nil {
		c.logger.Debug(fmt.Sprintf("Canary action %s for service %s in stage %s of project %s was received",
			e.Canary.Action, e.Service, e.Stage, e.Project))

		if deploymentStrategy == keptnevents.Duplicate {
			c.logger.Debug(fmt.Sprintf("Apply canary action %s for service %s in stage %s of project %s", e.Canary.Action, e.Service, e.Stage, e.Project))

			if err := c.changeCanary(e); err != nil {
				c.logger.Error(err.Error())
				return err
			}
		} else {
			c.logger.Debug(fmt.Sprintf("Discard canary action %s for service %s in stage %s of project %s because service is not duplicated",
				e.Canary.Action, e.Service, e.Stage, e.Project))
			if os.Getenv("PRE_WORKFLOW_ENGINE") != "true" {
				c.logger.Error(
					fmt.Sprintf("Cannot process received canary instructions as no duplicate deployment for service %s in stage %s of project %s is available",
						e.Service, e.Stage, e.Project))
			}
		}
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		strings.HasSuffix(ce.Source(), "remediation-service") {
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, "real-user", deploymentStrategy); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
		return nil
	}

	// Send deployment finished event
	// Note that this condition also stops the keptn-flow if an artifact is discarded
	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		!(e.Canary != nil && (e.Canary.Action == keptnevents.Discard || e.Canary.Action == keptnevents.Promote)) {

		testStrategy, err := getTestStrategy(e.Project, e.Stage)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, testStrategy, deploymentStrategy); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	return nil
}

func getGeneratedChart(e *keptnevents.ConfigurationChangeEventData) (*chart.Chart, error) {
	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	helmChartName := helm.GetChartName(e.Service, true)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, url.String())
}

func getDeploymentStrategyOfService(ch *chart.Chart) (keptnevents.DeploymentStrategy, error) {
	for _, keyword := range ch.Metadata.Keywords {
		if keyword == "deployment_strategy=duplicate" {
			return keptnevents.Duplicate, nil
		}
		if keyword == "deployment_strategy=direct" {
			return keptnevents.Direct, nil
		}
	}
	return keptnevents.Duplicate, errors.New("Cannot find deployment_strategy in keywords")
}

func changePrimaryDeployment(e *keptnevents.ConfigurationChangeEventData, chart *chart.Chart) error {

	for _, template := range chart.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(template.Data))
		newContent := make([]byte, 0, 0)
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
				// It is a deployment
				newDeployment := string(doc)
				for _, change := range e.DeploymentChanges {
					newDeployment, err = sjson.Set(newDeployment, change.PropertyPath, change.Value)
					if err != nil {
						return err
					}
				}
				newContent, err = objectutils.AppendJSONStringAsYaml(newContent, newDeployment)
				if err != nil {
					return err
				}

			} else {
				newContent, err = objectutils.AppendAsYaml(newContent, document)
				if err != nil {
					return err
				}
			}
		}
		template.Data = newContent
	}

	return nil
}

func (c *ConfigurationChanger) updateChart(e *keptnevents.ConfigurationChangeEventData, generated bool,
	editChart func(*keptnevents.ConfigurationChangeEventData, *chart.Chart) error) (*chart.Chart, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	helmChartName := helm.GetChartName(e.Service, generated)
	c.logger.Info(fmt.Sprintf("Start updating chart %s of stage %s", helmChartName, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, url.String())
	if err != nil {
		return nil, err
	}

	// Edit chart
	err = editChart(e, chart)
	if err != nil {
		return nil, err
	}

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helmChartName, chartData, url.String()); err != nil {
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("Finished updating chart %s of stage %s", helmChartName, e.Stage))
	return chart, nil
}

func changeValue(e *keptnevents.ConfigurationChangeEventData, chart *chart.Chart) error {

	values := make(map[string]interface{})
	yaml.Unmarshal([]byte(chart.Values.Raw), &values)

	// Change values
	for k, v := range e.ValuesCanary {
		values[k] = v
	}

	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	chart.Values.Raw = string(valuesData)

	return nil
}

func (c *ConfigurationChanger) setCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) (*chart.Chart, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("Start updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), url.String())
	if err != nil {
		return nil, err
	}

	c.generatedChartHandler.UpdateCanaryWeight(chart, canaryWeight)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, url.String()); err != nil {
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("Finished updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	return chart, nil
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	switch e.Canary.Action {
	case keptnevents.Discard:
		ch, err := c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if _, err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
		if err := c.deleteCanaryRelease(e); err != nil {
			return err
		}

	case keptnevents.Promote:
		ch, err := c.setCanaryWeight(e, 100)
		if err != nil {
			return err
		}
		if _, err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}

		chartGenerator := helm.NewGeneratedChartHandler(c.mesh, c.canaryLevelGen, c.keptnDomain)
		upgradeMsg, err := c.SimulateApplyChart(e.Project, e.Stage, e.Service, false)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		genChart, err := chartGenerator.GenerateDuplicateManagedChart(upgradeMsg, e.Project, e.Stage, e.Service)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.generatedChartHandler.UpdateCanaryWeight(genChart, int32(100)); err != nil {
			return err
		}
		genChartData, err := keptnutils.PackageChart(genChart)
		if err != nil {
			return err
		}
		if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData, url.String()); err != nil {
			return err
		}
		if _, err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}

		genChart, err = c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if _, err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
		if err := c.deleteCanaryRelease(e); err != nil {
			return err
		}

	case keptnevents.Set:
		ch, err := c.setCanaryWeight(e, e.Canary.Value)
		if err != nil {
			return err
		}
		if _, err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
	}

	return nil
}

// deleteCanaryRelease deletes a helm release
func (c *ConfigurationChanger) deleteCanaryRelease(e *keptnevents.ConfigurationChangeEventData) error {
	c.logger.Info(fmt.Sprintf("Start deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	if err := c.canaryLevelGen.DeleteCanaryRelease(e.Project, e.Stage, e.Service); err != nil {
		return fmt.Errorf("Error when deleting release %s: %s",
			helm.GetReleaseName(e.Project, e.Stage, e.Service, false), err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	return nil
}

// SimulateApplyChart
func (c *ConfigurationChanger) SimulateApplyChart(project, stage, service string, generated bool) (string, error) {

	releaseName := helm.GetReleaseName(project, stage, service, generated)
	namespace := c.canaryLevelGen.GetNamespace(project, stage, generated)
	c.logger.Info(fmt.Sprintf("Start dry-run of chart %s in namespace %s", releaseName, namespace))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return "", err
	}

	ch, err := keptnutils.GetChart(project, service, stage, helm.GetChartName(service, generated), url.String())
	if err != nil {
		return "", fmt.Errorf("Error when reading chart %s: %s", helm.GetChartName(service, generated), err.Error())
	}

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("Error when creating temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		return "", fmt.Errorf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error())
	}

	msg, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--reset-values", "--dry-run"})
	if err != nil {
		return "", fmt.Errorf("Error when making a dry run of chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	c.logger.Debug(msg)
	return msg, nil
}

// ApplyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) ApplyChart(ch *chart.Chart, project, stage, service string, generated bool) (string, error) {

	releaseName := helm.GetReleaseName(project, stage, service, generated)
	namespace := c.canaryLevelGen.GetNamespace(project, stage, generated)
	c.logger.Info(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace))

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("Error when creating temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		return "", fmt.Errorf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error())
	}

	msg, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--reset-values", "--wait", "--force"})
	if err != nil {
		return "", fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	c.logger.Debug(msg)

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return "", fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return msg, nil
}

// ApplyDirectory applies the provided directory
func (c *ConfigurationChanger) ApplyDirectory(chartPath, releaseName, namespace string) error {

	msg, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--reset-values", "--wait", "--force"})
	if err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	c.logger.Debug(msg)

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	return nil
}

func getInClusterConfig() bool {
	if os.Getenv("ENVIRONMENT") == "production" {
		return true
	}
	return false
}
