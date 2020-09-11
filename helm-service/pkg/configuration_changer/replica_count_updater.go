package configuration_changer

import (
	"bytes"
	"encoding/json"
	"github.com/ghodss/yaml"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type ReplicaCountUpdater struct {
	replicaIncrement int
}

func NewReplicaCountUpdater(replicaIncrement int) *ReplicaCountUpdater {
	return &ReplicaCountUpdater{
		replicaIncrement: replicaIncrement,
	}
}

// Update increases the replica count in the deployments by the provided replicaIncrement
func (u *ReplicaCountUpdater) Update(ch *chart.Chart) error {

	for _, template := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(template.Data))
		newContent := make([]byte, 0, 0)
		containsDepl := false
		for {
			var document interface{}
			err := dec.Decode(&document)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			doc, err := json.Marshal(document)
			if err != nil {
				return err
			}

			var depl appsv1.Deployment
			if err := json.Unmarshal(doc, &depl); err == nil && keptnutils.IsDeployment(&depl) {
				// Deployment found
				containsDepl = true
				depl.Spec.Replicas = getPtr(*depl.Spec.Replicas + int32(u.replicaIncrement))
				newContent, err = appendAsYaml(newContent, depl)
				if err != nil {
					return err
				}
			} else {
				newContent, err = appendAsYaml(newContent, document)
				if err != nil {
					return err
				}
			}
		}
		if containsDepl {
			template.Data = newContent
		}
	}

	return nil
}

func getPtr(x int32) *int32 {
	return &x
}

func appendAsYaml(content []byte, element interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}
