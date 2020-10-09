package configuration_changer

import (
	"helm.sh/helm/v3/pkg/chart"
)

// ValuesManipulator allows to manipulate the values of a Helm chart
type ValuesManipulator struct {
	values map[string]interface{}
}

// NewValuesManipulator creates a new ValuesManipulator
func NewValuesManipulator(values map[string]interface{}) *ValuesManipulator {
	return &ValuesManipulator{
		values: values,
	}
}

// Manipulate updates the values
func (v *ValuesManipulator) Manipulate(ch *chart.Chart) error {

	// Change values
	for k, v := range v.values {
		ch.Values[k] = v
	}
	return nil
}
