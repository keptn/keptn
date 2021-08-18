// Copyright 2012-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"crypto/tls"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	"github.com/keptn/keptn/distributor/pkg/lib/events"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := envconfig.Process("", &config.Global); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(config.Global))
}

func _main(env config.EnvConfig) int {
	connectionType := config.GetPubSubConnectionType()
	executionContext := createExecutionContext()
	eventSender := setupEventSender()
	httpClient := setupHTTPClient()

	uniformHandler, uniformLogHandler := getUniformHandlers(connectionType)
	controlPlane := controlplane.NewControlPlane(uniformHandler, connectionType)
	uniformWatch := setupUniformWatch(controlPlane)
	forwarder := events.NewForwarder(httpClient)

	// Start event forwarder
	logger.Info("Starting Event Forwarder")
	go forwarder.Start(executionContext)

	// Eventually start registration process
	if shallRegister() {
		id := uniformWatch.Start(executionContext)
		uniformLogger := controlplane.NewEventUniformLog(id, uniformLogHandler)
		uniformLogger.Start(executionContext, forwarder.EventChannel)

		defer func() {
			err := controlPlane.Unregister()
			if err != nil {
				logger.Warnf("Unable to unregister from Keptn's control plane: %v", err)
			} else {
				logger.Infof("Unregistered Keptn Integration")
			}
		}()
	}

	logger.Infof("Connection type: %s", connectionType)
	if connectionType == config.ConnectionTypeHTTP {
		err := env.ValidateKeptnAPIEndpointURL()
		if err != nil {
			logger.Fatalf("No valid URL configured for keptn api endpoint: %s", err)
		}
		logger.Info("Starting HTTP event poller")
		httpEventPoller := events.NewPoller(env, eventSender, httpClient)
		uniformWatch.RegisterListener(httpEventPoller)
		go httpEventPoller.Start(executionContext)
	} else {
		logger.Info("Starting NATS event Receiver")
		natsEventReceiver := events.NewNATSEventReceiver(env, eventSender)
		uniformWatch.RegisterListener(natsEventReceiver)
		go natsEventReceiver.Start(executionContext)
	}
	executionContext.Wg.Wait()
	return 0
}

func createExecutionContext() *events.ExecutionContext {
	// Prepare signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	context := events.ExecutionContext{
		Context: ctx,
		Wg:      wg,
	}
	return &context
}

func shallRegister() bool {
	if config.Global.DisableRegistration {
		logger.Infof("Registration to Keptn's control plane disabled")
		return false
	}

	if config.Global.K8sNamespace == "" || config.Global.K8sDeploymentName == "" {
		logger.Warn("Skipping Registration because not all mandatory environment variables are set: K8S_NAMESPACE, K8S_DEPLOYMENT_NAME")
		return false
	}

	if isOneOfFilteredServices(config.Global.K8sDeploymentName) {
		logger.Infof("Skipping Registration because service name %s is actively filtered", config.Global.K8sDeploymentName)
		return false
	}

	return true
}

func isOneOfFilteredServices(serviceName string) bool {
	switch serviceName {
	case
		"statistics-service",
		"api-service",
		"mongodb-datastore",
		"configuration-service",
		"secret-service",
		"shipyard-controller":
		return true
	}
	return false
}

func getUniformHandlers(connectionType config.ConnectionType) (*keptnapi.UniformHandler, *keptnapi.LogHandler) {
	if connectionType == config.ConnectionTypeHTTP {
		scheme := "http" // default
		parsed, _ := url.Parse(config.Global.KeptnAPIEndpoint)
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
		uniformHandler := keptnapi.NewAuthenticatedUniformHandler(config.Global.KeptnAPIEndpoint+"/controlPlane", config.Global.KeptnAPIToken, "x-token", nil, scheme)
		uniformLogHandler := keptnapi.NewAuthenticatedLogHandler(config.Global.KeptnAPIEndpoint+"/controlPlane", config.Global.KeptnAPIToken, "x-token", nil, scheme)
		return uniformHandler, uniformLogHandler
	}
	return keptnapi.NewUniformHandler(config.DefaultShipyardControllerBaseURL), keptnapi.NewLogHandler(config.DefaultShipyardControllerBaseURL)
}

func setupUniformWatch(controlPlane *controlplane.ControlPlane) *events.UniformWatch {
	return events.NewUniformWatch(controlPlane)
}

func setupHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.Global.VerifySSL}, //nolint:gosec
	}
	client := &http.Client{Transport: tr}
	return client
}

func setupCEClient() cloudevents.Client {
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c
}

func setupEventSender() events.EventSender {
	ceClient := setupCEClient()
	return &keptnv2.HTTPEventSender{
		Client: ceClient,
	}
}
