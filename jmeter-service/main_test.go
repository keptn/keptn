package main

import (
	"log"
	"net/url"
	"reflect"
	"testing"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var serviceURLTests = []struct {
	name  string
	event keptnv2.TestTriggeredEventData
	url   *url.URL
}{
	{"local", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "https://carts.sockshop-dev.svc.cluster.local/test", "carts.sockshop-dev.mydomain.com"), getURL("https://carts.sockshop-dev.svc.cluster.local/test")},
	{"public", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", "http://carts.sockshop-dev.mydomain.com:8080/myendpoint"), getURL("http://carts.sockshop-dev.mydomain.com:8080/myendpoint")},
}

func getURL(urlString string) *url.URL {
	url, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}
	return url
}

func deploymentFinishedEventInitHelper(project, service, stage, deploymentStrategy,
	deploymentURILocal, deploymentURIPublic string) keptnv2.TestTriggeredEventData {
	return keptnv2.TestTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: project,
			Service: service,
			Stage:   stage,
		},
		Deployment: keptnv2.TestTriggeredDeploymentDetails{
			DeploymentURIsLocal:  []string{deploymentURILocal},
			DeploymentURIsPublic: []string{deploymentURIPublic},
		},
	}
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
