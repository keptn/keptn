package validation

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/config"
	"github.com/keptn/keptn/shipyard-controller/models"
	"gopkg.in/yaml.v3"
)

// ConfigurableProjectNameValidator can be used
// validate a Keptn project name
type ConfigurableProjectNameValidator struct {
	projectNameMaxSize int
}

// NewProjectValidator creates a new ConfigurableProjectNameValidator
func NewProjectValidator(env config.EnvConfig) *ConfigurableProjectNameValidator {
	return &ConfigurableProjectNameValidator{projectNameMaxSize: env.ProjectNameMaxSize}
}

// Validate performs the actual validation on field level
func (p ConfigurableProjectNameValidator) Validate(fl validator.FieldLevel) bool {
	if projectName, ok := fl.Field().Interface().(string); ok {
		return len(projectName) <= p.projectNameMaxSize
	}
	return true
}

// Tag returns the go tag the validator is bound to
func (p ConfigurableProjectNameValidator) Tag() string {
	return "projectname"
}

type ProjectValidator struct {
	ProjectNameMaxSize int
}

func (p ProjectValidator) Validate(params interface{}) error {
	switch t := params.(type) {
	case *models.CreateProjectParams:
		return p.validateCreateProjectParams(t)
	case *models.UpdateProjectParams:
		return p.validateUpdateProjectParams(t)
	default:
		return nil
	}
}
func (p ProjectValidator) validateCreateProjectParams(createProjectParams *models.CreateProjectParams) error {
	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if len(*createProjectParams.Name) > p.ProjectNameMaxSize {
		return fmt.Errorf("project name exceeds maximum lenght of %d characters", p.ProjectNameMaxSize)
	}
	if !keptncommon.ValidateKeptnEntityName(*createProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}
	if createProjectParams.Shipyard == nil || *createProjectParams.Shipyard == "" {
		return errors.New("shipyard must contain a valid shipyard spec encoded in base64")
	}
	shipyard := &keptnv2.Shipyard{}
	decodeString, err := base64.StdEncoding.DecodeString(*createProjectParams.Shipyard)
	if err != nil {
		return errors.New("could not decode shipyard content")
	}

	err = yaml.Unmarshal(decodeString, shipyard)
	if err != nil {
		return fmt.Errorf("could not unmarshal provided shipyard content")
	}

	if err := common.ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateGitRemoteURL(createProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitToken != "" {
		return fmt.Errorf("privateKey and token cannot be used together")
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitProxyURL != "" {
		return fmt.Errorf("privateKey and proxy cannot be used together")
	}

	if createProjectParams.GitPrivateKey != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitPrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitPemCertificate != "" {
		return fmt.Errorf("SSH authorization and PEM Certificate be used together")
	}

	if createProjectParams.GitPemCertificate != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitPemCertificate)
		if err != nil {
			return errors.New("could not decode PEM Certificate content")
		}
	}

	return nil
}

func (p ProjectValidator) validateUpdateProjectParams(updateProjectParams *models.UpdateProjectParams) error {
	if updateProjectParams.Name == nil || *updateProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*updateProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}

	if updateProjectParams.Shipyard != nil && *updateProjectParams.Shipyard != "" {
		shipyard := &keptnv2.Shipyard{}
		decodeString, err := base64.StdEncoding.DecodeString(*updateProjectParams.Shipyard)
		if err != nil {
			return errors.New("could not decode shipyard content")
		}

		err = yaml.Unmarshal(decodeString, shipyard)
		if err != nil {
			return fmt.Errorf("could not unmarshal provided shipyard content")
		}

		if err := common.ValidateShipyardVersion(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}

		if err := common.ValidateShipyardStages(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}
	}

	if err := common.ValidateGitRemoteURL(updateProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitToken != "" {
		return fmt.Errorf("privateKey and token cannot be used together")
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitProxyURL != "" {
		return fmt.Errorf("privateKey and proxy cannot be used together")
	}

	if updateProjectParams.GitPrivateKey != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitPrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitPemCertificate != "" {
		return fmt.Errorf("SSH authorization and PEM Certificate be used together")
	}

	if updateProjectParams.GitPemCertificate != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitPemCertificate)
		if err != nil {
			return errors.New("could not decode PEM Certificate content")
		}
	}

	return nil
}
