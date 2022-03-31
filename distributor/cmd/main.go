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
	"github.com/keptn/keptn/distributor/pkg/api"
	"github.com/keptn/keptn/distributor/pkg/clientget"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/forwarder"
	"github.com/keptn/keptn/distributor/pkg/poller"
	"github.com/keptn/keptn/distributor/pkg/receiver"
	"github.com/keptn/keptn/distributor/pkg/uniform/controlplane"
	"github.com/keptn/keptn/distributor/pkg/uniform/log"
	"github.com/keptn/keptn/distributor/pkg/uniform/watch"
	"github.com/keptn/keptn/distributor/pkg/utils"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var gitCommit string
var buildTime string

func main() {
	env := config.EnvConfig{}
	if err := envconfig.Process("", &env); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}

	bark(env)

	executionContext := createExecutionContext()
	eventSender, err := createEventSender(env)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize event sender.")
	}

	httpClient, err := clientget.CreateClientGetter(env).Get()
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize http client.")
	}

	apiset, err := createKeptnAPI(httpClient, env)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize API set.")
	}

	controlPlane := controlplane.New(apiset.UniformV1(), env.PubSubConnectionType(), env)
	uniformWatch := watch.New(controlPlane, env)
	forwarder := forwarder.New(apiset.APIV1(), httpClient, env)

	// Start event forwarder
	logger.Info("Starting Event Forwarder")
	forwarder.Start(executionContext)

	// Eventually start registration process
	if env.ValidateRegistrationConstraints() {
		id, err := uniformWatch.Start(executionContext)
		if err != nil {
			logger.Fatal(err)
		}
		uniformLogger := log.New(id, apiset.LogsV1())
		uniformLogger.Start(executionContext, forwarder.EventChannel)
	}

	//explicitly wait for service to be up before starting NATS/pollers
	for waitForService(env) {
	}

	logger.Infof("Connection type: %s", env.PubSubConnectionType())
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		err := env.ValidateKeptnAPIEndpointURL()
		if err != nil {
			logger.Fatalf("No valid URL configured for keptn api endpoint: %s", err)
		}
		logger.Info("Starting HTTP event poller")
		httpEventPoller := poller.New(env, apiset.ShipyardControlV1(), eventSender)
		uniformWatch.RegisterListener(httpEventPoller)
		if err := httpEventPoller.Start(executionContext); err != nil {
			logger.Fatalf("Could not start HTTP event poller: %v", err)
		}
	} else {
		logger.Info("Starting NATS event receiver")
		natsEventReceiver := receiver.New(env, eventSender, env.ValidateRegistrationConstraints())
		uniformWatch.RegisterListener(natsEventReceiver)
		if err := natsEventReceiver.Start(executionContext); err != nil {
			logger.Fatalf("Could not start NATS event receiver: %v", err)
		}
	}
	executionContext.Wg.Wait()
}

func bark(env config.EnvConfig) {
	padR := func(str string) string {
		width := 40
		if len(str) >= width {
			return str
		}
		buf := bytes.NewBufferString(str)
		for i := 0; i < width-len(str); i++ {
			buf.WriteByte('.')
		}
		return buf.String()
	}

	strOrUnknown := func(s string) string {
		if s == "" {
			return "unknown"
		}
		return s
	}
	printFmtStr := "%s%s\n"
	printFmtBool := "%s%t\n"
	fmt.Printf(printFmtStr, padR("Git commit"), strOrUnknown(gitCommit))
	fmt.Printf(printFmtStr, padR("Build time"), strOrUnknown(buildTime))
	fmt.Printf(printFmtStr, padR("Start time"), time.Now().UTC().String())
	fmt.Printf(printFmtBool, padR("Remote execution plane"), env.PubSubConnectionType() == config.ConnectionTypeHTTP)
	fmt.Printf(printFmtBool, padR("Oauth enabled"), env.OAuthEnabled())
	fmt.Printf(printFmtStr, padR("Keptn API endpoint"), strOrUnknown(env.KeptnAPIEndpoint))
	fmt.Printf(printFmtStr, padR("Api proxy path"), strOrUnknown(env.APIProxyPath))
	fmt.Printf(printFmtStr, padR("Api proxy port"), strOrUnknown(strconv.Itoa(env.APIProxyPort)))
	fmt.Printf(printFmtStr, padR("PubSub URL"), strOrUnknown(env.PubSubURL))
	fmt.Printf(printFmtStr, padR("PubSub group"), strOrUnknown(env.PubSubGroup))
	fmt.Printf(printFmtStr, padR("K8S node name"), strOrUnknown(env.K8sNodeName))
	fmt.Printf(printFmtStr, padR("K8S namespace"), strOrUnknown(env.K8sNamespace))
	fmt.Printf(printFmtStr, padR("K8S deployment name"), strOrUnknown(env.K8sDeploymentName))
	fmt.Printf(printFmtStr, padR("K8S pod name"), strOrUnknown(env.K8sPodName))
}

func createExecutionContext() *utils.ExecutionContext {
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
	executionContext := utils.ExecutionContext{
		Context:  ctx,
		Wg:       wg,
		CancelFn: cancel,
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

	return api.NewInternal(httpClient)
}

func createEventSender(env config.EnvConfig) (poller.EventSender, error) {
	eventSender, err := keptnv2.NewHTTPEventSender(env.PubSubRecipientURL())
	if err != nil {
		return nil, err
	}
	return eventSender, nil
}

func waitForService(env config.EnvConfig) bool {
	host := env.PubSubRecipient
	pathToHealth := path.Join(env.PubSubRecipientPath, "/health")
	port := env.PubSubRecipientPort
	serviceurl := "http://" + host + ":" + port + pathToHealth

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(serviceurl)

	if err != nil || response.StatusCode != 200 {
		logger.Infof("%s %s %s\n", serviceurl, "Service not responding", err.Error())
		return true
	}
	logger.Infof("%s %s response code %d\n", serviceurl, "Service UP and RUNNING", response.StatusCode)
	return false
}
