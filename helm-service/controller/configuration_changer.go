package controller

import (
	"bytes"
	"encoding/json"
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
	generatedChartHandler *helm.GeneratedChartHandler
	canaryLevelGen        helm.CanaryLevelGenerator
	logger                keptnutils.LoggerInterface
	keptnDomain           string
}

// NewConfigurationChanger creates a new ConfigurationChanger
func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, canaryLevelGen, keptnDomain)
	return &ConfigurationChanger{generatedChartHandler: generatedChartHandler, canaryLevelGen: canaryLevelGen,
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

	if len(e.ValuesCanary) > 0 {
		if err := c.updateChart(e, false, changeValue); err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, false); err != nil {
			return err
		}
	}
	if len(e.DeploymentChanges) > 0 {
		if err := c.updateChart(e, true, changePrimaryDeployment); err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
	}

	// Change canary
	if e.Canary != nil {
		deplStrategies, err := GetDeploymentStrategies(e.Project)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s", err.Error()))
			return err
		}
		if deplStrategies[e.Stage] == keptnevents.Duplicate {
			if err := c.changeCanary(e); err != nil {
				c.logger.Error(err.Error())
				return err
			}
		} else {
			if os.Getenv("PRE_WORKFLOW_ENGINE") != "true" {
				c.logger.Error(fmt.Sprintf("Cannot process received canary instructions as deployment strategy for stage %s is %s",
					e.Stage, deplStrategies[e.Stage]))
			}
		}
	}

	// Send deployment finished event
	// Note that this condition also stops the keptn-flow if an artifact is discarded
	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		!(e.Canary != nil && e.Canary.Action == keptnevents.Discard) {

		testStrategy, err := getTestStrategy(e.Project, e.Stage)
		if err != nil {
			return err
		}
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, testStrategy); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		strings.HasSuffix(ce.Source(), "remediation-service") {
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, "real-user"); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	return nil
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
	editChart func(*keptnevents.ConfigurationChangeEventData, *chart.Chart) error) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	helmChartName := helm.GetChartName(e.Service, generated)
	c.logger.Info(fmt.Sprintf("Start updating chart %s", helmChartName))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, url.String())
	if err != nil {
		return err
	}

	// Edit chart
	err = editChart(e, chart)
	if err != nil {
		return err
	}

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helmChartName, chartData, url.String()); err != nil {
		return err
	}
	c.logger.Info(fmt.Sprintf("Finished updating chart %s", helmChartName))
	return nil
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

func (c *ConfigurationChanger) setCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	c.logger.Info(fmt.Sprintf("Start updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), url.String())
	if err != nil {
		return err
	}

	c.generatedChartHandler.UpdateCanaryWeight(chart, canaryWeight)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, url.String()); err != nil {
		return err
	}
	c.logger.Info(fmt.Sprintf("Finished updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	return nil
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	switch e.Canary.Action {
	case keptnevents.Discard:
		if err := c.setCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
		if err := c.deleteRelease(e, false); err != nil {
			return err
		}

	case keptnevents.Promote:
		if err := c.setCanaryWeight(e, 100); err != nil {
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}

		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), url.String())
		if err != nil {
			return err
		}
		genChartData, err := c.generatedChartHandler.GenerateDuplicateManagedChart(userChart, e.Project, e.Stage)
		if err != nil {
			return err
		}
		genChart, err := keptnutils.LoadChart(genChartData)
		if err != nil {
			return err
		}
		if err := c.generatedChartHandler.UpdateCanaryWeight(genChart, int32(100)); err != nil {
			return err
		}
		genChartData, err = keptnutils.PackageChart(genChart)
		if err != nil {
			return err
		}
		if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData, url.String()); err != nil {
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}

		if err := c.setCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
		if err := c.deleteRelease(e, false); err != nil {
			return err
		}

	case keptnevents.Set:
		if err := c.setCanaryWeight(e, e.Canary.Value); err != nil {
			return err
		}
		if err := c.ApplyConfiguration(e.Project, e.Stage, e.Service, true); err != nil {
			return err
		}
	}

	return nil
}

// deleteRelease deletes a helm release
func (c *ConfigurationChanger) deleteRelease(e *keptnevents.ConfigurationChangeEventData, generated bool) error {
	c.logger.Info(fmt.Sprintf("Start deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	if err := c.canaryLevelGen.DeleteRelease(e.Project, e.Stage, e.Service, generated); err != nil {
		return fmt.Errorf("Error when deleting release %s: %s",
			helm.GetReleaseName(e.Project, e.Stage, e.Service, generated), err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	return nil
}

// ApplyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) ApplyConfiguration(project, stage, service string, generated bool) error {

	releaseName := helm.GetReleaseName(project, stage, service, generated)
	namespace := c.canaryLevelGen.GetNamespace(project, stage, generated)
	c.logger.Info(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	ch, err := keptnutils.GetChart(project, service, stage, helm.GetChartName(service, generated), url.String())
	if err != nil {
		return fmt.Errorf("Error when reading chart %s: %s", helm.GetChartName(service, generated), err.Error())
	}

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("Error when creating temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		return fmt.Errorf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error())
	}

	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return nil
}

// ApplyDirectory applies the provided directory
func ApplyDirectory(chartPath, releaseName, namespace string) error {

	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}

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
