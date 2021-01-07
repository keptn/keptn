package mesh

import (
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"os"
	"reflect"
	"testing"
)

func TestGetPublicDeploymentURI(t *testing.T) {
	type args struct {
		event v0_2_0.EventData
	}
	tests := []struct {
		name                  string
		args                  args
		want                  []string
		hostNameTemplate      string
		ingressProtocol       string
		ingressHostNameSuffix string
		ingressPort           string
	}{
		{
			name: "use default values",
			args: args{
				event: v0_2_0.EventData{
					Project: "test-project",
					Stage:   "test-stage",
					Service: "test-service",
				},
			},
			want: []string{"http://test-service.test-project-test-stage.svc.cluster.local:80"},
		},
		{
			name: "use default template with custom values",
			args: args{
				event: v0_2_0.EventData{
					Project: "test-project",
					Stage:   "test-stage",
					Service: "test-service",
				},
			},
			want:                  []string{"https://test-service.test-project-test-stage.test-hostname-suffix.dev:8090"},
			hostNameTemplate:      "",
			ingressProtocol:       "https",
			ingressHostNameSuffix: "test-hostname-suffix.dev",
			ingressPort:           "8090",
		},
		{
			name: "use custom template with custom values",
			args: args{
				event: v0_2_0.EventData{
					Project: "test-project",
					Stage:   "test-stage",
					Service: "test-service",
				},
			},
			want:                  []string{"https://test-service-test-project-test-stage.test-hostname-suffix.dev:8090"},
			hostNameTemplate:      "${INGRESS_PROTOCOL}://${service}-${project}-${stage}.${INGRESS_HOSTNAME_SUFFIX}:${INGRESS_PORT}",
			ingressProtocol:       "https",
			ingressHostNameSuffix: "test-hostname-suffix.dev",
			ingressPort:           "8090",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("INGRESS_HOSTNAME_SUFFIX", tt.ingressHostNameSuffix)
			os.Setenv("INGRESS_PORT", tt.ingressPort)
			os.Setenv("INGRESS_PROTOCOL", tt.ingressProtocol)
			os.Setenv("HOSTNAME_TEMPLATE", tt.hostNameTemplate)
			if got := GetPublicDeploymentURI(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPublicDeploymentURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
