package handler

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// IShipyardRetriever godoc
//go:generate moq -pkg fake -skip-ensure -out ./fake/shipyardretriever_mock.go . IShipyardRetriever
type IShipyardRetriever interface {
	GetShipyard(projectName string) (*keptnv2.Shipyard, string, error)
	GetCachedShipyard(projectName string) (*keptnv2.Shipyard, error) //TODO: what to do with this?
}

type ShipyardRetriever struct {
	configurationStore common.ConfigurationStore
	projectRepo        db.ProjectMVRepo
}

func NewShipyardRetriever(configurationStore common.ConfigurationStore, projectRepo db.ProjectMVRepo) *ShipyardRetriever {
	return &ShipyardRetriever{
		configurationStore: configurationStore,
		projectRepo:        projectRepo,
	}
}

func (sr *ShipyardRetriever) GetShipyard(projectName string) (*keptnv2.Shipyard, string, error) {
	resource, err := sr.configurationStore.GetProjectResource(projectName, "shipyard.yaml")
	if err != nil {
		return nil, "", fmt.Errorf("could not retrieve shipyard.yaml for project %s: %w", projectName, err)
	}

	log.Infof("current resource: %s", resource)
	shipyard, err := common.UnmarshalShipyard(resource.ResourceContent)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal shipyard.yaml of project %s: %w", projectName, err)
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
		return nil, "", fmt.Errorf("invalid shipyard version: %w", err)
	}

	return shipyard, resource.Metadata.Version, nil
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
