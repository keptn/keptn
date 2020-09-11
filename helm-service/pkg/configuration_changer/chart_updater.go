package configuration_changer

import "helm.sh/helm/v3/pkg/chart"

type ChartUpdater interface {
	Update(ch *chart.Chart) error
}
