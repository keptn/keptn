package helm

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

func check(e error, t *testing.T) {
	if e != nil {
		t.Error(e)
	}
}

type GeneratedResource struct {
	URI         string
	FileContent []string
}

func Equals(actual *chart.Chart, valuesExpected GeneratedResource, templatesExpected []GeneratedResource, t *testing.T) {

	// Compare values
	jsonData, err := json.Marshal(actual.Values)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(jsonData), valuesExpected.FileContent[0])

	for _, resource := range templatesExpected {

		reader := ioutil.NopCloser(bytes.NewReader(GetTemplateByName(actual, resource.URI).Data))
		decoder := kyaml.NewDocumentDecoder(reader)

		for i := 0; ; i++ {
			b1 := make([]byte, 4096)
			n1, err := decoder.Read(b1)
			if err == io.EOF {
				break
			}
			assert.Nil(t, err, "")

			jsonData, err := objectutils.ToJSON(b1[:n1])
			if err != nil {
				t.Error(err)
			}

			ja := jsonassert.New(t)
			ja.Assertf(string(jsonData), resource.FileContent[i])
		}
	}
}

func GetTemplateByName(chart *chart.Chart, templateName string) *chart.File {

	for _, template := range chart.Templates {
		if template.Name == templateName {
			return template
		}
	}
	return nil
}

// GetTestGeneratedChart returns a sample chart representing a "generated-chart"
func GetTestGeneratedChart() chart.Chart {
	return chart.Chart{
		Raw: nil,
		Metadata: &chart.Metadata{
			Name:       "carts-generated",
			Version:    "0.1.0",
			APIVersion: "v2",
		},
		Lock: nil,
		Templates: []*chart.File{
			{
				Name: "carts-canary-istio-destinationrule.yaml",
				Data: []byte(GeneratedCanaryDestinationRule),
			},
			{
				Name: "carts-canary-service.yaml",
				Data: []byte(GeneratedCanaryService),
			},
			{
				Name: "carts-istio-virtualservice.yaml",
				Data: []byte(GeneratedVirtualService),
			},
			{
				Name: "carts-primary-deployment.yaml",
				Data: []byte(GeneratedPrimaryDeployment),
			},
			{
				Name: "carts-primary-istio-destinationrule.yaml",
				Data: []byte(GeneratedPrimaryDestinationRule),
			},
			{
				Name: "carts-primary-service.yaml",
				Data: []byte(GeneratedPrimaryService),
			},
		},
	}
}

// GetTestUserChart returns a sample chart representing a chart provided by the user
func GetTestUserChart() chart.Chart {
	return chart.Chart{
		Raw: nil,
		Metadata: &chart.Metadata{
			Name:       "carts",
			Version:    "0.1.0",
			APIVersion: "v2",
		},
		Lock: nil,
		Templates: []*chart.File{
			{
				Name: "carts-service.yaml",
				Data: []byte(userService),
			},
			{
				Name: "carts-deployment.yaml",
				Data: []byte(userDeployment),
			},
		},
	}
}
