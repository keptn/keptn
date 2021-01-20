package controller

import (
	"encoding/base64"
	"errors"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"testing"
)

func TestCheckAndSetServiceName(t *testing.T) {

	data := helm.CreateTestHelmChartData(t)

	testCases := []struct {
		name        string
		event       *keptnv2.ServiceCreateFinishedEventData
		error       error
		serviceName string
	}{
		{"Mismatch", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts-1"},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			errors.New("Provided Keptn service name \"carts-1\" does not match Kubernetes service name \"carts\""), "carts-1"},
		{"Mismatch2", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts-1"},
			Helm: keptnv2.Helm{Chart: "%%%"}},
			errors.New("Error when decoding the Helm Chart: illegal base64 data at input byte 0"), "carts-1"},
		{"Match", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: ""},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
		{"Set", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: ""},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := CheckAndSetServiceName(tt.event)
			if res == nil && res != tt.error {
				t.Errorf("got nil, want %s", tt.error.Error())
			} else if res != nil && tt.error != nil && res.Error() != tt.error.Error() {
				t.Errorf("got %s, want %s", res.Error(), tt.error.Error())
			} else if res != nil && tt.error == nil {
				t.Errorf("got %s, want nil", res.Error())
			}

			if tt.event.Service != tt.serviceName {
				t.Errorf("got %s, want %s", tt.event.Service, tt.serviceName)
			}
		})
	}
}
