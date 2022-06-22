package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

//chartPackager is able to package a helm chart
type chartPackager struct {
}

//NewChartPackager creates a new chartPackager instance
func NewChartPackager() *chartPackager {
	return &chartPackager{}
}

//packages a helm chart into its byte representation
func (pc chartPackager) Package(ch *chart.Chart) ([]byte, error) {
	helmPackage, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf(ErrPackageChartMsg, err.Error())
	}
	defer os.RemoveAll(helmPackage)

	// Marshal values into values.yaml
	// This step is necessary as chartutil.Save uses the Raw content
	for _, f := range ch.Raw {
		if f.Name == chartutil.ValuesfileName {
			f.Data, err = yaml.Marshal(ch.Values)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	name, err := chartutil.Save(ch, helmPackage)
	if err != nil {
		return nil, fmt.Errorf(ErrPackageChartMsg, err.Error())
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf(ErrPackageChartMsg, err.Error())
	}
	return data, nil
}
