package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

func stringp(s string) *string {
	return &s
}

func validateCreateProjectParams(createProjectParams *operations.CreateProjectParams) error {

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

	if err := common.ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	return nil
}

func validateUpdateProjectParams(createProjectParams *operations.CreateProjectParams) error {

	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*createProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}

	return nil
}

func validateCreateServiceParams(params *operations.CreateServiceParams) error {
	if params.ServiceName == nil || *params.ServiceName == "" {
		return errors.New("Must provide a service name")
	}
	if !keptncommon.ValididateUnixDirectoryName(*params.ServiceName) {
		return errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}
	return nil
}
