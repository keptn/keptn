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
	"github.com/keptn/go-utils/pkg/common/retry"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	"github.com/keptn/keptn/distributor/pkg/lib/events"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var env config.EnvConfig

func main() {
	if err := envconfig.Process("", &env); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(env))
}

func _main(env config.EnvConfig) int {
	context := createExecutionContext()
	eventSender := setupEventSender()
	httpClient := setupHTTPClient()
	connectionType := config.GetPubSubConnectionType(env)
	logger.Infof("Connection type: %s", connectionType)
	if connectionType == config.ConnectionTypeHTTP {
		logger.Info("Starting HTTP event poller")
		httpEventPoller :=
			events.NewPoller(env, eventSender, httpClient)
		go httpEventPoller.Start(context)
	} else {
		logger.Info("Starting Nats event Receiver")
		natsEventReceiver := events.NewNATSEventReceiver(env, eventSender)
		go natsEventReceiver.Start(context)
	}

	logger.Info("Starting Event Forwarder")
	forwarder := events.NewForwarder(env, httpClient)
	go forwarder.Start(context)

	if shallRegister(env) {
		logger.Infof("Registering Keptn Intgration")
		uniformHandler, uniformLogHandler := getUniformHandlers(connectionType)
		controlPlane := controlplane.NewControlPlane(uniformHandler, controlplane.CreateRegistrationData(connectionType, env))
		go func() {
			retry.Retry(func() error {
				id, err := controlPlane.Register()
				if err != nil {
					logger.Warnf("Unable to register to Keptn's control plane: %s", err.Error())
					return err
				}
				logger.Infof("Registered Keptn Integration with id %s", id)

				logHandler := uniformLogHandler
				uniformLogger := controlplane.NewEventUniformLog(id, logHandler)
				uniformLogger.Start(context, forwarder.EventChannel)
				logger.Infof("Started UniformLogger for Keptn Integration")
				return nil
			})
			for {
				select {
				case <-context.Done():
					return
				case <-time.After(config.GetRegistrationInterval(env)):
					_, err := controlPlane.Register()
					if err != nil {
						logger.Warnf("Unable to (re)register to Keptn's control plane: %s", err.Error())
					}
				}
			}
		}()

		defer func() {
			err := controlPlane.Unregister()
			if err != nil {
				logger.Warnf("Unable to unregister from Keptn's control plane: %v", err)
			} else {
				logger.Infof("Unregistered Keptn Integration")
			}
		}()
	}

	context.Wg.Wait()

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

func shallRegister(env config.EnvConfig) bool {
	if env.DisableRegistration {
		logger.Infof("Registration to Keptn's control plane disabled")
		return false
	}

	if env.K8sNamespace == "" || env.K8sDeploymentName == "" {
		logger.Warn("Skipping Registration because not all mandatory environment variables are set: K8S_NAMESPACE, K8S_DEPLOYMENT_NAME")
		return false
	}
	return true
}

func getUniformHandlers(connectionType config.ConnectionType) (*keptnapi.UniformHandler, *keptnapi.LogHandler) {
	if connectionType == config.ConnectionTypeHTTP {
		uniformHandler := keptnapi.NewAuthenticatedUniformHandler(env.KeptnAPIEndpoint+"/controlPlane", env.KeptnAPIToken, "x-token", nil, "http")
		uniformLogHandler := keptnapi.NewAuthenticatedLogHandler(env.KeptnAPIEndpoint+"/controlPlane", env.KeptnAPIToken, "x-token", nil, "http")
		return uniformHandler, uniformLogHandler
	}
	return keptnapi.NewUniformHandler(config.DefaultShipyardControllerBaseURL), keptnapi.NewLogHandler(config.DefaultShipyardControllerBaseURL)
}

func setupHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !env.VerifySSL}, //nolint:gosec
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

//TODO: vvv delete vvv

func getProxyHost(path string) (string, string, string) {
	return "", "", ""
}

func getHTTPPollingEndpoint() string {
	return ""
}

func hasEventBeenSent(sentEvents []string, eventID string) bool {
	return false
}

func decodeCloudEvent(data []byte) (*cloudevents.Event, error) {

	return nil, nil
}

func matchesFilter(e cloudevents.Event) bool {

	return true
}

func pollEventsForTopic(endpoint string, token string, topic string) {

}
