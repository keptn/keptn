package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/keptn/keptn/remediation-service/actions"
	"gopkg.in/yaml.v2"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"

	configmodels "github.com/keptn/go-utils/pkg/configuration-service/models"
	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const tillernamespace = "kube-system"
const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {

	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)
	if err != nil {
		log.Fatalf("failed to create transport: %v", err)
	}

	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func isProjectAndStageAvailable(problem *keptnevents.ProblemEventData) bool {
	return problem.Project != "" && problem.Stage != ""
}

// deriveFromTags allows to derive project, stage, and service information from tags
// Input example: "Tags:":"keptn_service:carts, keptn_stage:dev, keptn_stage:sockshop"
func deriveFromTags(problem *keptnevents.ProblemEventData) {

	tags := strings.Split(problem.Tags, ", ")

	for _, tag := range tags {
		if strings.HasPrefix(tag, "keptn_service:") {
			problem.Service = tag[len("keptn_service:"):]
		} else if strings.HasPrefix(tag, "keptn_stage:") {
			problem.Stage = tag[len("keptn_stage:"):]
		} else if strings.HasPrefix(tag, "keptn_project:") {
			problem.Project = tag[len("keptn_project:"):]
		}
	}
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "remediation-service")
	logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	var problemEvent *keptnevents.ProblemEventData
	if event.Type() == keptnevents.ProblemOpenEventType {
		logger.Debug("Received problem notification")
		problemEvent = &keptnevents.ProblemEventData{}
		if err := event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	if !isProjectAndStageAvailable(problemEvent) {
		deriveFromTags(problemEvent)
	}
	if !isProjectAndStageAvailable(problemEvent) {
		return errors.New("Cannot derive project and stage from tags nor impacted entity")
	}

	logger.Debug("Received problem event with state " + problemEvent.State)

	// valide if remediation should be performed
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	autoRemediate, err := isRemediationEnabled(resourceHandler, problemEvent.Project, problemEvent.Stage)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check if remediation is enabled: %s", err.Error()))
		return err
	}

	if autoRemediate {
		logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		logger.Info(fmt.Sprintf("Remediation disabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
		return nil
	}

	// get remediation.yaml
	var resource *configmodels.Resource
	if problemEvent.Service != "" {
		resource, err = resourceHandler.GetServiceResource(problemEvent.Project, problemEvent.Stage,
			problemEvent.Service, remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(problemEvent.Project, problemEvent.Stage, remediationFileName)
	}

	if err != nil {
		logger.Error("Failed to get remediation.yaml file")
		return err
	}
	logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptnmodels.Remediations
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)
	if err != nil {
		logger.Error("Could not unmarshal remediation.yaml")
		return err
	}

	actionExecutors := []actions.ActionExecutor{actions.NewScaler(), actions.NewSlower(), actions.NewAborter(), actions.NewFeatureToggler()}

	for _, remediation := range remediationData.Remediations {
		logger.Debug("Trying to map remediation '" + remediation.Name + "' to problem '" + problemEvent.ProblemTitle + "'")
		if strings.HasPrefix(problemEvent.ProblemTitle, remediation.Name) {
			logger.Debug("Remediation for problem found")
			// currently only one remediation action is supported
			for _, a := range actionExecutors {
				if a.GetAction() == remediation.Actions[0].Action {
					if strings.ToLower(problemEvent.State) == "open" {
						if err := a.ExecuteAction(problemEvent, shkeptncontext, remediation.Actions[0]); err != nil {
							logger.Error(err.Error())
							return err
						}
						logger.Info(fmt.Sprintf("Remediation action %s successfully applied",
							remediation.Actions[0].Action))
						return nil
					} else if strings.ToLower(problemEvent.State) == "resolved" ||
						strings.ToLower(problemEvent.State) == "closed" {
						if err := a.ResolveAction(problemEvent, shkeptncontext, remediation.Actions[0]); err != nil {
							logger.Error(err.Error())
							return err
						}
						logger.Info(fmt.Sprintf("Remediation action %s resolved",
							remediation.Actions[0].Action))
						return nil
					}
				}
			}
		}
	}

	logger.Info("No remediation action found")
	return nil
}

func isRemediationEnabled(rh *configutils.ResourceHandler, project string, stage string) (bool, error) {
	keptnHandler := keptnutils.NewKeptnHandler(rh)
	shipyard, err := keptnHandler.GetShipyard(project)
	if err != nil {
		return false, err
	}
	for _, s := range shipyard.Stages {
		if s.Name == stage && s.RemediationStrategy == "automated" {
			return true, nil
		}
	}

	return false, nil
}
