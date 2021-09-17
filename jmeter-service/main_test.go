package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var serviceURLTests = []struct {
	name  string
	event keptnv2.TestTriggeredEventData
	url   *url.URL
}{
	{"local", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "https://carts.sockshop-dev.svc.cluster.local/test", "carts.sockshop-dev.mydomain.com"), getURL("https://carts.sockshop-dev.svc.cluster.local/test")},
	{"public", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", "http://carts.sockshop-dev.mydomain.com:8080/myendpoint"), getURL("http://carts.sockshop-dev.mydomain.com:8080/myendpoint")},
	{"no path but valid uri", deploymentFinishedEventInitHelper("sockshop", "carts", "dev", "direct", "", "http://carts.sockshop-dev.mydomain.com:8080"), getURL("http://carts.sockshop-dev.mydomain.com:8080/")},
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

func Test_checkEndpointAvailable(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	reachableURL, _ := url.Parse(ts.URL)
	nonReachableURL, _ := url.Parse("http://1.2.3.4:1234")
	nonReachableURLwithoutPort, _ := url.Parse("http://1.2.3.4")

	type args struct {
		timeout    time.Duration
		serviceURL *url.URL
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Url reachable", args{1 * time.Second, reachableURL}, false},
		{"Url not reachable", args{1 * time.Nanosecond, nonReachableURL}, true},
		{"nil", args{1 * time.Nanosecond, nil}, true},
		{"Url without port not reachable", args{1 * time.Nanosecond, nonReachableURLwithoutPort}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkEndpointAvailable(tt.args.timeout, tt.args.serviceURL); (err != nil) != tt.wantErr {
				t.Errorf("checkEndpointAvailable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
