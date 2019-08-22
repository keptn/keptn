package validator

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/keptn/keptn/cli/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const helmChartPathSuffix = "helm"

// Values represents the Helm values and contains required fields
type Values struct {
	Image string `json:"image"`
}

// ValidateHelmChart validates keptn's requirements regarding
// the values, deployment, and service file
func ValidateHelmChart(helmChart []byte) (bool, error) {

	workingPath, err := ioutil.TempDir("", "helm")
	if err != nil {
		return false, err
	}

	if err := utils.Untar(workingPath, bytes.NewReader(helmChart)); err != nil {
		return false, err
	}

	services, err := ioutil.ReadDir(workingPath)
	if err != nil {
		return false, err
	}

	for _, f := range services {

		templateFiles, err := getFiles(filepath.Join(workingPath, f.Name(), "templates"), ".yml", ".yaml")
		if err != nil {
			return false, err
		}
		if res, err := validateTemplateRequirements(templateFiles); !res || err != nil {
			return false, err
		}

		if _, err := os.Stat(filepath.Join(workingPath, f.Name(), "values.yml")); err == nil {
			validateValues(filepath.Join(workingPath, f.Name(), "values.yaml"))
		} else if _, err := os.Stat(filepath.Join(workingPath, f.Name(), "values.yaml")); err == nil {
			validateValues(filepath.Join(workingPath, f.Name(), "values.yaml"))
		} else {
			return false, nil
		}

	}
	if err := os.RemoveAll(workingPath); err != nil {
		return false, nil
	}
	return true, nil
}

func validateValues(file string) (bool, error) {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}

	dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(dat))
	var val Values
	if err := dec.Decode(&val); err != nil {
		return false, err
	}
	return val.Image != "", nil
}

func validateTemplateRequirements(files []string) (bool, error) {

	atLeastOneDeployment := false
	atLeastOneService := false

	for _, file := range files {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			return false, err
		}
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(dat))
		for {
			var svc corev1.Service
			errSvc := dec.Decode(&svc)
			if errSvc == io.EOF {
				break
			}
			if isService(svc) && validateService(svc) {
				atLeastOneService = true
			} else if isService(svc) {
				return false, nil
			}
		}

		dec = kyaml.NewYAMLToJSONDecoder(bytes.NewReader(dat))
		for {
			var dpl appsv1.Deployment
			errDpl := dec.Decode(&dpl)
			if errDpl == io.EOF {
				break
			}

			if isDeployment(dpl) && validateDeployment(dpl) {
				atLeastOneDeployment = true
			} else if isDeployment(dpl) {
				return false, nil
			}
		}
	}
	return atLeastOneDeployment && atLeastOneService, nil
}

func isService(svc corev1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

func isDeployment(dpl appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}

func validateService(svc corev1.Service) bool {

	val, ok := svc.Spec.Selector["app"]
	return isService(svc) && ok && val != ""
}

func validateDeployment(depl appsv1.Deployment) bool {

	mLabel, ok1 := depl.Spec.Selector.MatchLabels["app"]
	podLabel, ok2 := depl.Spec.Template.ObjectMeta.Labels["app"]
	return isDeployment(depl) && ok1 && ok2 && mLabel != "" && podLabel != ""
}

func getFiles(workingPath string, extensions ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(workingPath, func(path string, info os.FileInfo, err error) error {
		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	return files, err
}
