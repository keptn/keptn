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
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	"github.com/keptn/keptn/distributor/pkg/lib/events"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	if err := envconfig.Process("", &config.Global); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}
	env := config.Global

	executionContext := createExecutionContext()
	eventSender, err := setupEventSender(env)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize event sender.")
	}

	httpClient := createHTTPClient(env)
	apiset, err := createKeptnAPI(httpClient)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize API set.")
	}

	controlPlane := controlplane.NewControlPlane(apiset.UniformV1(), config.PubSubConnectionType())
	uniformWatch := setupUniformWatch(controlPlane)
	forwarder := events.NewForwarder(apiset.APIV1(), apiset.ProxyV1())

	// Start event forwarder
	logger.Info("Starting Event Forwarder")
	forwarder.Start(executionContext)

	// Eventually start registration process
	register := shallRegister()
	if register {
		id := uniformWatch.Start(executionContext)
		if id == "" {
			logger.Fatal("Could not register Uniform")
		}
		uniformLogger := controlplane.NewEventUniformLog(id, apiset.LogsV1())
		uniformLogger.Start(executionContext, forwarder.EventChannel)
	}

	logger.Infof("Connection type: %s", config.PubSubConnectionType())
	if config.PubSubConnectionType() == config.ConnectionTypeHTTP {
		err := env.ValidateKeptnAPIEndpointURL()
		if err != nil {
			logger.Fatalf("No valid URL configured for keptn api endpoint: %s", err)
		}
		logger.Info("Starting HTTP event poller")
		httpEventPoller := events.NewPoller(env, apiset.ShipyardControlV1(), eventSender)
		uniformWatch.RegisterListener(httpEventPoller)
		if err := httpEventPoller.Start(executionContext); err != nil {
			logger.Fatalf("Could not start HTTP event poller: %v", err)
		}
	} else {
		logger.Info("Starting NATS event Receiver")
		natsEventReceiver := events.NewNATSEventReceiver(env, eventSender, register)
		uniformWatch.RegisterListener(natsEventReceiver)
		if err := natsEventReceiver.Start(executionContext); err != nil {
			logger.Fatalf("Could not start NATS event receiver: %v", err)
		}
	}
	executionContext.Wg.Wait()
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
	executionContext := events.ExecutionContext{
		Context: ctx,
		Wg:      wg,
	}
	return &executionContext
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

func createKeptnAPI(httpClient *http.Client) (keptnapi.KeptnInterface, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if config.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme := "http"
		parsed, _ := url.Parse(config.Global.KeptnAPIEndpoint)
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
		return keptnapi.New(config.Global.KeptnAPIEndpoint, keptnapi.WithScheme(scheme), keptnapi.WithHTTPClient(httpClient), keptnapi.WithAuthToken(config.Global.KeptnAPIToken))
	}
	return keptnapi.New(config.DefaultShipyardControllerBaseURL, keptnapi.WithHTTPClient(httpClient))
}

func setupUniformWatch(controlPlane controlplane.IControlPlane) *events.UniformWatch {
	return events.NewUniformWatch(controlPlane)
}

func createHTTPClient(envConfig config.EnvConfig) *http.Client {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.Global.VerifySSL}, //nolint:gosec
		},
		Timeout: 5 * time.Second,
	}

	if envConfig.UseSSO() {
		conf := clientcredentials.Config{
			ClientID:     envConfig.SSOClientID,
			ClientSecret: envConfig.SSOClientSecret,
			Scopes:       envConfig.SSOScopes,
			TokenURL:     envConfig.SSOTokenURL,
		}

		client := conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c))
		return client
	}
	return c
}

func setupEventSender(env config.EnvConfig) (events.EventSender, error) {
	eventSender, err := keptnv2.NewHTTPEventSender(env.PubSubRecipientURL())
	if err != nil {
		return nil, err
	}
	return eventSender, nil
}
