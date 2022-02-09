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
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/lib/client"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	"github.com/keptn/keptn/distributor/pkg/lib/events"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	env := config.EnvConfig{}
	if err := envconfig.Process("", &env); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}

	executionContext := createExecutionContext()
	eventSender, err := createEventSender(env)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize event sender.")
	}

	httpClient, err := client.CreateClientGetter(env).Get()
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize http client.")
	}

	apiset, err := createKeptnAPI(httpClient, env)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize API set.")
	}

	controlPlane := controlplane.NewControlPlane(apiset.UniformV1(), env.PubSubConnectionType(), env)
	uniformWatch := setupUniformWatch(controlPlane)
	forwarder := events.NewForwarder(apiset.APIV1(), httpClient, env)

	// Start event forwarder
	logger.Info("Starting Event Forwarder")
	forwarder.Start(executionContext)

	// Eventually start registration process
	if env.ValidateRegistrationConstraints() {
		id := uniformWatch.Start(executionContext)
		if id == "" {
			logger.Fatal("Could not register Uniform")
		}
		uniformLogger := controlplane.NewEventUniformLog(id, apiset.LogsV1())
		uniformLogger.Start(executionContext, forwarder.EventChannel)
	}

	logger.Infof("Connection type: %s", env.PubSubConnectionType())
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
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
		natsEventReceiver := events.NewNATSEventReceiver(env, eventSender, env.ValidateRegistrationConstraints())
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

func createKeptnAPI(httpClient *http.Client, env config.EnvConfig) (keptnapi.KeptnInterface, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme := "http"
		parsed, _ := url.Parse(env.KeptnAPIEndpoint)
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
		return keptnapi.New(env.KeptnAPIEndpoint, keptnapi.WithScheme(scheme), keptnapi.WithHTTPClient(httpClient), keptnapi.WithAuthToken(env.KeptnAPIToken))
	}

	return keptnapi.NewInternal(httpClient)
}

func setupUniformWatch(controlPlane controlplane.IControlPlane) *events.UniformWatch {
	return events.NewUniformWatch(controlPlane)
}

func createEventSender(env config.EnvConfig) (events.EventSender, error) {
	eventSender, err := keptnv2.NewHTTPEventSender(env.PubSubRecipientURL())
	if err != nil {
		return nil, err
	}
	return eventSender, nil
}
