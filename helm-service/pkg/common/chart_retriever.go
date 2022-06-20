package common

import (
	"fmt"
	"net/url"

	goutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/kubeutils"
	"helm.sh/helm/v3/pkg/chart"
)

//chartRetriever is able to store a helm chart
type chartRetriever struct {
	resourceHandler *goutils.ResourceHandler
}

//RetrieveChartOptions are the parameters to obtain a chart
type RetrieveChartOptions struct {
	Project   string
	Service   string
	Stage     string
	ChartName string
	CommitID  string
}

//NewChartRetriever creates a new chartRetriever instance
func NewChartRetriever(resourceHandler *goutils.ResourceHandler) *chartRetriever {
	return &chartRetriever{
		resourceHandler: resourceHandler,
	}
}

func (cs chartRetriever) Retrieve(chartOpts RetrieveChartOptions) (*chart.Chart, string, error) {
	option := url.Values{}
	if chartOpts.CommitID != "" {
		option.Add("gitCommitID", chartOpts.CommitID)
	}
	resource, err := cs.resourceHandler.GetResource(
		*goutils.NewResourceScope().
			Project(chartOpts.Project).
			Service(chartOpts.Service).
			Resource(kubeutils.GetHelmChartURI(chartOpts.ChartName)).
			Stage(chartOpts.Stage),
		goutils.AppendQuery(option))

	if err != nil {
		return nil, "", fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartOpts.ChartName, chartOpts.Project, err.Error())
	}
	ch, err := LoadChart([]byte(resource.ResourceContent))
	if err != nil {
		return nil, "", fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartOpts.ChartName, chartOpts.Project, err.Error())
	}
	if chartOpts.CommitID == "" {
		return ch, resource.Metadata.Version, nil
	}

	return ch, chartOpts.CommitID, nil
}
