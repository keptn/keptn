package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

// ConfigurationChangedEvent ...
type ConfigurationChangedEvent struct {
	Service            string `json:"service"`
	Image              string `json:"image"`
	Tag                string `json:"tag"`
	Project            string `json:"project"`
	Stage              string `json:"stage"`
	GitHubOrg          string `json:"githuborg"`
	TestStrategy       string `json:"teststrategy"`
	DeploymentStrategy string `json:"deploymentstrategy"`
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "helm-service")

	logger.Debug(fmt.Sprintf("Got Event Context: %+v", event.Context))

	data := &ConfigurationChangedEvent{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if event.Type() != "sh.keptn.events.configuration-changed" {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	go doDeployment(event, logger, *data, shkeptncontext)
	return nil
}

func doDeployment(event cloudevents.Event, logger *keptnutils.Logger, data ConfigurationChangedEvent, shkeptncontext string) {
	repo, err := keptnutils.Checkout(data.GitHubOrg, data.Project, data.Stage)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when checking out configuration from GitHub: %s", err.Error()))
	}

	logger.Info("Deploying with helm upgrade")

	switch strings.ToLower(data.DeploymentStrategy) {
	case "direct":
		if checkErr(keptnutils.DoHelmUpgrade(data.Project, data.Stage), logger) != nil {
			return
		}
		if checkErr(keptnutils.WaitForDeploymentsInNamespace(true, data.Project+"-"+data.Stage),
			logger) != nil {
			return
		}
		logger.Info("Finished deploying in stage " + data.Stage)

	case "blue_green_service":
		// Move repo head one commit back
		ref, err := keptnutils.CheckoutPrevCommit(repo)
		if checkErr(err, logger) != nil {
			return
		}

		// Do helm upgrade
		if checkErr(keptnutils.DoHelmUpgrade(data.Project, data.Stage), logger) != nil {
			return
		}

		// Wait for rollout to complete
		if checkErr(keptnutils.WaitForDeploymentToBeAvailable(true, data.Service+"-blue", data.Project+"-"+data.Stage),
			logger) != nil {
			return
		}
		if checkErr(keptnutils.WaitForDeploymentToBeAvailable(true, data.Service+"-green", data.Project+"-"+data.Stage),
			logger) != nil {
			return
		}

		// Move repo head one commit forward
		if checkErr(keptnutils.CheckoutReference(repo, ref), logger) != nil {
			return
		}

		// Do helm upgrade
		if checkErr(keptnutils.DoHelmUpgrade(data.Project, data.Stage), logger) != nil {
			return
		}
		logger.Info("Finished deploying in stage " + data.Stage)

	default:
		logger.Error("Unknown deployment strategy '" + data.DeploymentStrategy + "'")
		return
	}

	checkErr(sendDeploymentFinishedEvent(shkeptncontext, event), logger)
}

func checkErr(err error, logger *keptnutils.Logger) error {
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func sendDeploymentFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event) error {

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.deployment-finished",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: incomingEvent.Data,
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget("http://event-broker.keptn.svc.cluster.local/keptn"),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("will listen on :%d%s\n", env.Port, env.Path)
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}
