package configuration_changer

import (
	"helm.sh/helm/v3/pkg/chart"
)

type ValuesUpdater struct {
	values map[string]interface{}
}

func NewValuesUpdater(values map[string]interface{}) *ValuesUpdater {
	return &ValuesUpdater{
		values: values,
	}
}

func (v *ValuesUpdater) Update(ch *chart.Chart) error {

	// Change values
	for k, v := range v.values {
		ch.Values[k] = v
	}
	return nil
}
