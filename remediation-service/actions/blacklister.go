package actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/keptn/keptn/remediation-service/pkg/apis/networking/istio"
	"github.com/keptn/keptn/remediation-service/pkg/utils"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const ipHandler = "/template/blacklistip-handler.yaml"
const ipRule = "/template/checkip-rule.yaml"
const ipInstance = "/template/sourceip-instance.yaml"

type BlackLister struct {
}

// NewBlackLister creates a new BlackLister
func NewBlackLister() *BlackLister {
	return &BlackLister{}
}

func (b BlackLister) GetAction() string {
	return "blacklist"
}

func (b BlackLister) ExecuteAction(problem *keptn.ProblemEventData, shkeptncontext string,
	action *keptn.RemediationAction) error {

	ip, err := getIP(problem)
	if err != nil {
		return fmt.Errorf("failed to parse ip from ProblemDetails: %v", err)
	}

	changedFiles := make(map[string]string)

	containsMixer, err := b.containsMixers(problem.Project, problem.Stage)
	if err != nil {
		return fmt.Errorf("failed to check if mixers exist: %v", err)
	}

	if containsMixer {
		// Add IP in blacklistip-handler
		handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
		resource, err := handler.GetStageResource(problem.Project, problem.Stage, ipHandler)
		if err != nil {
			return fmt.Errorf("failed to get mixer resource: %v", err)
		}
		changedFiles[ipHandler] = resource.ResourceContent
	} else {
		changedFiles, err = b.getMixerTemplates("")
		if err != nil {
			return fmt.Errorf("failed ot read mixer templates: %v", err)
		}
	}

	newContent, err := b.addIP(changedFiles[ipHandler], ip)
	if err != nil {
		return fmt.Errorf("failed to update mixer resource: %v", err)
	}
	changedFiles[ipHandler] = newContent

	data := keptn.ConfigurationChangeEventData{
		Project:                  problem.Project,
		Service:                  problem.Service,
		Stage:                    problem.Stage,
		FileChangesUmbrellaChart: changedFiles,
	}

	err = utils.CreateAndSendConfigurationChangedEvent(problem, shkeptncontext, data)
	if err != nil {
		return fmt.Errorf("failed to send configuration change event: %v", err)
	}
	return nil
}

func (b BlackLister) ResolveAction(problem *keptn.ProblemEventData, shkeptncontext string,
	action *keptn.RemediationAction) error {
	return errors.New("no resolving action for action " + b.GetAction() + "implemented")
}

func (b BlackLister) getMixerTemplates(prefixPath string) (map[string]string, error) {

	changedFiles := make(map[string]string)

	for _, i := range []string{ipHandler, ipRule, ipInstance} {
		content, err := ioutil.ReadFile(prefixPath + i)
		if err != nil {
			return nil, err
		}
		changedFiles[i] = string(content)
	}
	return changedFiles, nil
}

func (b BlackLister) addIP(inputContent, ip string) (string, error) {

	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(inputContent))
	listChecker := istio.ListChecker{}
	err := dec.Decode(&listChecker)
	if err != nil {
		return "", fmt.Errorf("Failed to decode list checker: %v", err)
	}
	listChecker.Spec.Params.Overrides = append(listChecker.Spec.Params.Overrides, ip+"/32")

	yamlData, err := yaml.Marshal(listChecker)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// containsMixers checks whether the project and stage already contain mixers implementing
// the blacklist
func (b BlackLister) containsMixers(project, stage string) (bool, error) {

	handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
	resources, err := handler.GetAllStageResources(project, stage)
	if err != nil {
		return false, err
	}

	requiredTemplates := map[string]string{
		ipHandler:  "",
		ipRule:     "",
		ipInstance: "",
	}
	for _, resource := range resources {
		if _, ok := requiredTemplates[*resource.ResourceURI]; ok {
			delete(requiredTemplates, *resource.ResourceURI)
		}
	}
	return len(requiredTemplates) == 0, nil
}
