package shipyardretriever

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/configurationstore"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// IShipyardRetriever godoc
//go:generate moq -pkg fake -skip-ensure -out ./fake/shipyardretriever_mock.go . IShipyardRetriever
type IShipyardRetriever interface {
	GetShipyard(projectName string) (*keptnv2.Shipyard, error)
	GetCachedShipyard(projectName string) (*keptnv2.Shipyard, error)
	GetLatestCommitID(projectName, stageName string) (string, error)
}

type ShipyardRetriever struct {
	configurationStore configurationstore.ConfigurationStore
	projectRepo        db.ProjectMVRepo
}

func NewShipyardRetriever(configurationStore configurationstore.ConfigurationStore, projectRepo db.ProjectMVRepo) *ShipyardRetriever {
	return &ShipyardRetriever{
		configurationStore: configurationStore,
		projectRepo:        projectRepo,
	}
}

func (sr *ShipyardRetriever) GetShipyard(projectName string) (*keptnv2.Shipyard, error) {
	resource, err := sr.configurationStore.GetProjectResource(projectName, "shipyard.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve shipyard.yaml for project %s: %w", projectName, err)
	}

	shipyard, err := common.UnmarshalShipyard(resource.ResourceContent)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal shipyard.yaml of project %s: %w", projectName, err)
	}

	// update the shipyard content of the project
	shipyardContent, err := yaml.Marshal(shipyard)
	if err != nil {
		// log the error but continue
		log.Errorf("could not encode shipyard file of project %s: %v", projectName, err)
	}
	if err := sr.projectRepo.UpdateShipyard(projectName, string(shipyardContent)); err != nil {
		// log the error but continue
		log.Errorf("could not update shipyard content of project %s: %v", projectName, err)
	}

	// validate the shipyard version - only shipyard files following the current keptn spec are supported by the shipyard controller
	if err = common.ValidateShipyardVersion(shipyard); err != nil {
		// if the validation has not been successful: send a <task-sequence>.finished event with status=errored
		return nil, fmt.Errorf("invalid shipyard version: %w", err)
	}

	return shipyard, nil
}

// GetCachedShipyard returns the shipyard that is stored for the project in the materialized view, instead of pulling it from the upstream
// this is done to reduce requests to the upstream and reduce the risk of running into rate limiting problems
func (sr *ShipyardRetriever) GetCachedShipyard(projectName string) (*keptnv2.Shipyard, error) {
	project, err := sr.projectRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}

	shipyard, err := common.UnmarshalShipyard(project.Shipyard)
	if err != nil {
		return nil, err
	}
	return shipyard, nil
}

func (sr *ShipyardRetriever) GetLatestCommitID(projectName, stageName string) (string, error) {
	stageMetadata, err := sr.configurationStore.GetStageResource(projectName, stageName, "metadata.yaml")
	if err != nil {
		return "", fmt.Errorf("could not determine latest commit ID for stage %s in project %s: %w", stageName, projectName, err)
	}

	if stageMetadata == nil || stageMetadata.Metadata == nil {
		return "", fmt.Errorf("could not determine latest commit ID for stage %s in project %s: stage metadata not found", stageName, projectName)
	}

	return stageMetadata.Metadata.Version, nil
}
