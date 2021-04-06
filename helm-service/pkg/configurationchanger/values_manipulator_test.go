package configurationchanger

import (
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"log"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestSimpleMerge(t *testing.T) {

	inputValuesContent := `
image: gw
tag: 0.0.1 
`

	inputValues := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(inputValuesContent), &inputValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	newValuesContent := `
tag: 0.0.2
`
	newValues := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(newValuesContent), &newValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	inputChart := chart.Chart{Values: inputValues}
	NewValuesManipulator(newValues).Manipulate(&inputChart)

	expectedValuesContent := `
image: gw
tag: 0.0.2 `
	expectedValues := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(expectedValuesContent), &expectedValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	assert.Equal(t, expectedValues, inputChart.Values)
}

func TestComplexMerge(t *testing.T) {

	inputValuesContent := `
gw:
  image:
    name: gw
    tag: 0.0.1 
  deployment:
    hostName: test
  port: 8080`

	inputValues := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(inputValuesContent), &inputValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	newValuesContent := `
gw:
  image:
    tag: 0.0.2
`
	newValues := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(newValuesContent), &newValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	inputChart := chart.Chart{Values: inputValues}
	NewValuesManipulator(newValues).Manipulate(&inputChart)

	expectedValuesContent := `
gw:
  image:
    name: gw
    tag: 0.0.2
  deployment:
    hostName: test
  port: 8080`
	expectedValues := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(expectedValuesContent), &expectedValues); err != nil {
		log.Fatalf("Unmarshalling error")
	}

	assert.Equal(t, expectedValues, inputChart.Values)
}
