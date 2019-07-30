package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/utils"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type stage struct {
	Name string `json:"name"`
}
type projectData struct {
	Project  string      `json:"project"`
	Registry interface{} `json:"registry"`
	Stages   []stage     `json:"stages"`
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

	utils.ServiceName = "helm-service"

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

	utils.Debug(shkeptncontext, fmt.Sprintf("Got Event Context: %+v", event.Context))

	data := &projectData{}
	if err := event.DataAs(data); err != nil {
		utils.Error(shkeptncontext, fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	if event.Type() != "create.project" {
		const errorMsg = "Received unexpected keptn event"
		utils.Error(shkeptncontext, errorMsg)
		return errors.New(errorMsg)
	}

	go createRoutes(shkeptncontext, *data)
	return nil
}

func createRoutes(shkeptncontext string, data projectData) {

	for _, stage := range data.Stages {
		exposeRoute(data.Project, stage.Name, shkeptncontext)
	}
}

func exposeRoute(project string, stage string, shkeptncontext string) {
	appDomain := os.Getenv("APP_DOMAIN")
	if appDomain == "" {
		utils.Error(shkeptncontext, "No app domain defined. Cannot create route.")
		return
	}
	// oc create route edge istio-wildcard-ingress-secure-keptn --service=istio-ingressgateway --hostname="www.keptn.ingress-gateway.$BASE_URL" --port=http2 --wildcard-policy=Subdomain --insecure-policy='Allow'
	utils.Info(shkeptncontext, "Trying to create route www."+stage+"."+project+"."+appDomain)
	out, err := utils.ExecuteCommand("oc",
		[]string{
			"create",
			"route",
			"edge",
			project + "-" + stage,
			"--service=istio-ingressgateway",
			"--hostname=www." + project + "-" + stage + "." + appDomain,
			"--port=http2",
			"--wildcard-policy=Subdomain",
			"--insecure-policy=Allow",
			"-n",
			"istio-system",
		})
	if err != nil {
		utils.Error(shkeptncontext, "Could not create route for: "+err.Error())
	}
	utils.Info(shkeptncontext, "oc create route output: "+out)
}

func checkErr(err error, shkeptncontext string) error {
	if err != nil {
		utils.Error(shkeptncontext, err.Error())
		return err
	}
	return nil
}
