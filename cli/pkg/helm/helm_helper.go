// +build !nokubectl

// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	kubeclient "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"k8s.io/client-go/tools/clientcmd"
)

// Helper provides helper functions for common Helm operations
type Helper struct {
}

// NewHelper creates a Helper
func NewHelper() Helper {
	return Helper{}
}

// DownloadChart downloads a Helm chart using the provided repo URL
func (c Helper) DownloadChart(chartRepoURL string) (*chart.Chart, error) {

	resp, err := http.Get(chartRepoURL)
	if err != nil {
		return nil, errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}

	ch, err := keptnutils.LoadChart(bytes)
	if err != nil {
		return nil, errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}
	return ch, err
}

func newActionConfig(config *rest.Config, namespace string) (*action.Configuration, error) {

	logFunc := func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}

	restClientGetter := newConfigFlags(config, namespace)
	kubeClient := &kubeclient.Client{
		Factory: cmdutil.NewFactory(restClientGetter),
		Log:     logFunc,
	}
	client, err := kubeClient.Factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	s := driver.NewSecrets(client.CoreV1().Secrets(namespace))
	s.Log = logFunc

	return &action.Configuration{
		RESTClientGetter: restClientGetter,
		Releases:         storage.Init(s),
		KubeClient:       kubeClient,
		Log:              logFunc,
	}, nil
}

func newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

// GetHistory returns the history for a Helm release
func (c Helper) GetHistory(releaseName, namespace string) ([]*release.Release, error) {

	logging.PrintLog(fmt.Sprintf("Check availability of Helm release %s in namespace %s", releaseName, namespace), logging.VerboseLevel)

	config, err := clientcmd.BuildConfigFromFlags("", getKubeConfig())
	if err != nil {
		return nil, err
	}

	cfg, err := newActionConfig(config, namespace)
	if err != nil {
		return nil, err
	}

	histClient := action.NewHistory(cfg)

	return histClient.Run(releaseName)
}

// UpgradeChart upgrades/installs the provided chart
func (c Helper) UpgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error {

	if len(ch.Templates) > 0 {
		logging.PrintLog(fmt.Sprintf("Start upgrading Helm Chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)

		config, err := clientcmd.BuildConfigFromFlags("", getKubeConfig())
		if err != nil {
			return err
		}

		cfg, err := newActionConfig(config, namespace)
		if err != nil {
			return err
		}

		histClient := action.NewHistory(cfg)
		var release *release.Release

		if _, err = histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
			iCli := action.NewInstall(cfg)
			iCli.Namespace = namespace
			iCli.ReleaseName = releaseName
			iCli.Wait = true
			iCli.Timeout = 10 * time.Minute
			release, err = iCli.Run(ch, vals)
		} else {
			iCli := action.NewUpgrade(cfg)
			iCli.Namespace = namespace
			iCli.Wait = true
			iCli.Timeout = 10 * time.Minute
			iCli.ReuseValues = true
			release, err = iCli.Run(releaseName, ch, vals)
		}
		if err != nil {
			return fmt.Errorf("Error when installing/upgrading Helm Chart %s in namespace %s: %s",
				releaseName, namespace, err.Error())
		}
		if release != nil {
			logging.PrintLog(release.Manifest, logging.VerboseLevel)
			if err := waitForDeploymentsOfHelmRelease(release.Manifest); err != nil {
				return err
			}
		} else {
			logging.PrintLog("Release is nil", logging.InfoLevel)
		}
		logging.PrintLog(fmt.Sprintf("Finished upgrading Helm Chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)
	} else {
		logging.PrintLog("Upgrade not done since this is an empty Helm Chart", logging.InfoLevel)
	}
	return nil
}

func getKubeConfig() string {
	if os.Getenv("KUBECONFIG") != "" {
		return keptnutils.ExpandTilde(os.Getenv("KUBECONFIG"))
	}
	return filepath.Join(
		keptnutils.UserHomeDir(), ".kube", "config",
	)

}

// UninstallRelease uninstalls the provided release
func (c Helper) UninstallRelease(releaseName, namespace string) error {
	logging.PrintLog(fmt.Sprintf("Start uninstalling Helm release %s in namespace %s", releaseName, namespace), logging.InfoLevel)
	config, err := clientcmd.BuildConfigFromFlags("", getKubeConfig())
	if err != nil {
		return err
	}

	cfg, err := newActionConfig(config, namespace)
	if err != nil {
		return err
	}

	iCli := action.NewUninstall(cfg)
	_, err = iCli.Run(releaseName)

	if err != nil {
		return fmt.Errorf("Error when uninstalling Helm release %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	return nil
}

func getDeployments(helmManifest string) []*appsv1.Deployment {

	deployments := []*appsv1.Deployment{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var dpl appsv1.Deployment
		err := dec.Decode(&dpl)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsDeployment(&dpl) {
			deployments = append(deployments, &dpl)
		}
	}
	return deployments
}

func waitForDeploymentsOfHelmRelease(helmManifest string) error {
	depls := getDeployments(helmManifest)
	for _, depl := range depls {
		if err := keptnutils.WaitForDeploymentToBeRolledOut(false, depl.Name, depl.Namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", depl.Name, depl.Namespace, err.Error())
		}
	}
	return nil
}
