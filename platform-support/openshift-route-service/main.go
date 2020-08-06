package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	keptn "github.com/keptn/go-utils/pkg/lib"

	b64 "encoding/base64"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
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

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptn.NewLogger(shkeptncontext, event.Context.GetID(), "openshift-route-service")

	logger.Debug(fmt.Sprintf("Got Event Context: %+v", event.Context))

	data := &keptn.ProjectCreateEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	if event.Type() != keptn.InternalProjectCreateEventType {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	go func() {
		err := createRoutes(data)
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	return nil
}

func createRoutes(data *keptn.ProjectCreateEventData) error {
	shipyard := keptn.Shipyard{}
	decodedStr, err := b64.StdEncoding.DecodeString(data.Shipyard)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(decodedStr, &shipyard)
	if err != nil {
		return err
	}
	for _, stage := range shipyard.Stages {
		if stage.DeploymentStrategy == "blue_green_service" {
			// add required security context constraints to the generated namespace to make istio injection work
			if err := enableMesh(data.Project, stage.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func enableMesh(project string, stage string) error {
	_, err := keptn.ExecuteCommand("oc",
		[]string{
			"adm",
			"policy",
			"add-scc-to-group",
			"privileged",
			"system:serviceaccounts",
			"-n",
			project + "-" + stage,
		})
	if err != nil {
		return errors.New("Could not add security context constraint 'privileged' for namespace " + project + "-" + stage + ": " + err.Error())
	}
	out, err := keptn.ExecuteCommand("oc",
		getEnableMeshCommandArgs(project, stage))
	if err != nil {
		return errors.New("Could not add security context constraint 'anyuid' for namespace " + project + "-" + stage + ": " + err.Error())
	}
	fmt.Println("enableMesh() output: " + out)
	return nil
}

func getEnableMeshCommandArgs(project string, stage string) []string {
	return []string{
		"adm",
		"policy",
		"add-scc-to-group",
		"anyuid",
		"system:serviceaccounts",
		"-n",
		project + "-" + stage,
	}
}
