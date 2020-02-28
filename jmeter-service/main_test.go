package main

import (
	"log"
	"net/url"
	"reflect"
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/events"
)

var serviceURLTests = []struct {
	name  string
	event keptnevents.DeploymentFinishedEventData
	url   *url.URL
}{
	{"local", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "https://carts.sockshop-dev.svc.cluster.local/test", "carts.sockshop-dev.mydomain.com"), getURL("https://carts.sockshop-dev.svc.cluster.local/test")},
	{"public", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", "http://carts.sockshop-dev.mydomain.com:8080/myendpoint"), getURL("http://carts.sockshop-dev.mydomain.com:8080/myendpoint")},
	{"educatedGuessDirect", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", ""), getURL("http://carts.sockshop-dev/health")},
	{"educatedGuessBG", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "blue_green_service", "", ""), getURL("http://carts-canary.sockshop-dev/health")},
}

func getURL(urlString string) *url.URL {
	url, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}
	return url
}

func deploymentFinishedEventInitHelper(project, service, stage, deploymentStrategy,
	deploymentURILocal, deploymentURIPublic string) keptnevents.DeploymentFinishedEventData {
	return keptnevents.DeploymentFinishedEventData{Project: project, Service: service, Stage: stage, DeploymentStrategy: deploymentStrategy, DeploymentURILocal: deploymentURILocal, DeploymentURIPublic: deploymentURIPublic}
}

func TestGetServiceURL(t *testing.T) {
	for _, tt := range serviceURLTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getServiceURL(tt.event)
			if err != nil {
				t.Errorf("unexpected error")
			}
			if !reflect.DeepEqual(*res, *tt.url) {
				t.Errorf("got %v, want %v for %s", res, tt.url, tt.name)
			}
		})
	}
}
