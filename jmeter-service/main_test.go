package main

import (
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/events"
)

var serviceURLTests = []struct {
	name  string
	event keptnevents.DeploymentFinishedEventData
	res   string
}{
	{"local", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "carts.sockshop-dev.svc.cluster.local", "carts.sockshop-dev.mydomain.com"), "carts.sockshop-dev.svc.cluster.local"},
	{"public", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", "carts.sockshop-dev.mydomain.com"), "carts.sockshop-dev.mydomain.com"},
	{"educatedGuessDirect", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", ""), "carts.sockshop-dev"},
	{"educatedGuessBG", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "blue_green_service", "", ""), "carts-canary.sockshop-dev"},
}

func deploymentFinishedEventInitHelper(project, service, stage, deploymentStrategy,
	deploymentURILocal, deploymentURIPublic string) keptnevents.DeploymentFinishedEventData {
	return keptnevents.DeploymentFinishedEventData{Project: project, Service: service, Stage: stage, DeploymentStrategy: deploymentStrategy, DeploymentURILocal: deploymentURILocal, DeploymentURIPublic: deploymentURIPublic}
}

func TestGetNewerVersion(t *testing.T) {
	for _, tt := range serviceURLTests {
		t.Run(tt.name, func(t *testing.T) {
			res := getServiceURL(tt.event)
			if res != tt.res {
				t.Errorf("got %v, want %v for %s", res, tt.res, tt.name)
			}
		})
	}
}
