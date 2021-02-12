package helm

import (
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	"github.com/stretchr/testify/assert"

	"sigs.k8s.io/yaml"
)

func TestGenerateDuplicateChart(t *testing.T) {

	logger := keptn.NewLogger("", "", "test")
	generator := NewGeneratedChartGenerator(mesh.NewIstioMesh(), logger)

	ch, err := generator.GenerateDuplicateChart(userService+renderedUserDeployment, "sockshop", "dev", "carts")

	const nsPlaceholder = "NAMESPACE_PLACEHOLDER"
	const ns = "dev"

	assert.Nil(t, err)
	var expectedValues map[string]interface{}
	assert.Equal(t, expectedValues, ch.Values)

	assert.Equal(t, "templates/carts-canary-service.yaml", ch.Templates[0].Name)
	assert.Equal(t, yamlUnmarshal([]byte(GeneratedCanaryService)), yamlUnmarshal(ch.Templates[0].Data))

	assert.Equal(t, "templates/carts-canary-istio-destinationrule.yaml", ch.Templates[1].Name)
	assert.Equal(t, yamlUnmarshal([]byte(strings.Replace(GeneratedCanaryDestinationRule, nsPlaceholder, ns, -1))), yamlUnmarshal(ch.Templates[1].Data))

	assert.Equal(t, "templates/carts-primary-service.yaml", ch.Templates[2].Name)
	assert.Equal(t, yamlUnmarshal([]byte(strings.Replace(GeneratedPrimaryService, nsPlaceholder, ns, -1))), yamlUnmarshal(ch.Templates[2].Data))

	assert.Equal(t, "templates/carts-primary-istio-destinationrule.yaml", ch.Templates[3].Name)
	assert.Equal(t, yamlUnmarshal([]byte(strings.Replace(GeneratedPrimaryDestinationRule, nsPlaceholder, ns, -1))), yamlUnmarshal(ch.Templates[3].Data))

	assert.Equal(t, "templates/carts-istio-virtualservice.yaml", ch.Templates[4].Name)
	assert.Equal(t, yamlUnmarshal([]byte(strings.Replace(GeneratedVirtualService, nsPlaceholder, ns, -1))), yamlUnmarshal(ch.Templates[4].Data))

	assert.Equal(t, "templates/carts-primary-deployment.yaml", ch.Templates[5].Name)
	assert.Equal(t, yamlUnmarshal([]byte(strings.Replace(GeneratedPrimaryDeployment, nsPlaceholder, ns, -1))), yamlUnmarshal(ch.Templates[5].Data))
}

func TestGenerateDuplicateChartWithTwoServices(t *testing.T) {

	logger := keptn.NewLogger("", "", "test")
	generator := NewGeneratedChartGenerator(mesh.NewIstioMesh(), logger)

	ch, err := generator.GenerateDuplicateChart(userService+userService+renderedUserDeployment, "sockshop", "dev", "carts")

	assert.Nil(t, ch)
	assert.Equal(t, errors.New("Chart contains multiple Kubernetes services but only 1 is allowed"), err)
}

func TestGenerateDuplicateChartWithTwoDeployments(t *testing.T) {

	logger := keptn.NewLogger("", "", "test")
	generator := NewGeneratedChartGenerator(mesh.NewIstioMesh(), logger)

	ch, err := generator.GenerateDuplicateChart(userService+renderedUserDeployment+renderedUserDeployment, "sockshop", "dev", "carts")

	assert.Nil(t, ch)
	assert.Equal(t, errors.New("Chart contains multiple Kubernetes deployments but only 1 is allowed"), err)
}

func yamlUnmarshal(data []byte) interface{} {
	var obj interface{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		log.Fatal(err)
	}
	return obj
}

func Test_getVirtualServicePublicHost(t *testing.T) {
	type args struct {
		svc       string
		project   string
		stageName string
	}
	tests := []struct {
		name             string
		hostnameTemplate string
		hostnameSuffix   string
		args             args
		want             string
		wantErr          bool
	}{
		{
			name: "get default hostname",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc.prj-stg.svc.cluster.local",
			wantErr: false,
		},
		{
			name:           "get hostname based on default template and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameSuffix: "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc.prj-stg.123.xip.io",
			wantErr: false,
		},
		{
			name:             "get hostname based on custom HOSTNAME_TEMPLATE and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameTemplate: "${service}-${stage}-${project}.${INGRESS_HOSTNAME_SUFFIX}",
			hostnameSuffix:   "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:             "get hostname based on custom HOSTNAME_TEMPLATE and custom INGRESS_HOSTNAME_SUFFIX",
			hostnameTemplate: "${INGRESS_PROTOCOL}://${service}-${stage}-${project}.${INGRESS_HOSTNAME_SUFFIX}",
			hostnameSuffix:   "123.xip.io",
			args: args{
				svc:       "svc",
				project:   "prj",
				stageName: "stg",
			},
			want:    "svc-stg-prj.123.xip.io",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOSTNAME_TEMPLATE", tt.hostnameTemplate)
			os.Setenv("INGRESS_HOSTNAME_SUFFIX", tt.hostnameSuffix)
			got, err := getVirtualServicePublicHost(tt.args.svc, tt.args.project, tt.args.stageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVirtualServicePublicHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getVirtualServicePublicHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
