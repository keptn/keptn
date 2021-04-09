package configurationchanger

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
		// Merge ch.Values[k] in v
		merge(v, ch.Values[k])
		ch.Values[k] = v
	}
	return nil
}

func merge(in1, in2 interface{}) interface{} {
	switch in1 := in1.(type) {
	case []interface{}:
		in2, ok := in2.([]interface{})
		if !ok {
			return in1
		}
		return append(in1, in2...)
	case map[string]interface{}:
		in2, ok := in2.(map[string]interface{})
		if !ok {
			return in1
		}
		for k, v2 := range in2 {
			if v1, ok := in1[k]; ok {
				in1[k] = merge(v1, v2)
			} else {
				in1[k] = v2
			}
		}
	case nil:
		in2, ok := in2.(map[string]interface{})
		if ok {
			return in2
		}
	}
	return in1
}
