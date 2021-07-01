package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

type RollbackFunc func() error

func Stringp(s string) *string {
	return &s
}

func ValidateCreateProjectParams(createProjectParams *operations.CreateProjectParams) error {

	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*createProjectParams.Name) {
		errorMsg := "Project name contains upper case letter(s) or special character(s).\n"
		errorMsg += "Keptn relies on the following conventions: "
		errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
		errorMsg += "Please update project name and try again."
		return errors.New(errorMsg)
	}
	if createProjectParams.Shipyard == nil || *createProjectParams.Shipyard == "" {
		return errors.New("shipyard must contain a valid shipyard spec encoded in base64")
	}
	shipyard := &keptnv2.Shipyard{}
	decodeString, err := base64.StdEncoding.DecodeString(*createProjectParams.Shipyard)
	if err != nil {
		return errors.New("could not decode shipyard content using base64 decoder: " + err.Error())
	}

	err = yaml.Unmarshal(decodeString, shipyard)
	if err != nil {
		return fmt.Errorf("could not unmarshal provided shipyard content: %s", err.Error())
	}

	if err := ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	return nil
}

func ValidateUpdateProjectParams(updateProjectParams *operations.UpdateProjectParams) error {

	if updateProjectParams.Name == nil || *updateProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*updateProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}

	return nil
}

func ValidateCreateServiceParams(params *operations.CreateServiceParams) error {
	if params.ServiceName == nil || *params.ServiceName == "" {
		return errors.New("Must provide a service name")
	}
	if !keptncommon.ValidateUnixDirectoryName(*params.ServiceName) {
		return errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}
	return nil
}

func Merge(in1, in2 interface{}) interface{} {
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
				in1[k] = Merge(v1, v2)
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
