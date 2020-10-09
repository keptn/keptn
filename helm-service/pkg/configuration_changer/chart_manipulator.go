package configuration_changer

import "helm.sh/helm/v3/pkg/chart"

// ChartManipulator interface for manipulating charts
type ChartManipulator interface {
	Manipulate(ch *chart.Chart) error
}
