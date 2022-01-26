package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
)

//IStageManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_manager_mock.go . IStageManager
type IStageManager interface {
	CreateStage(params models.CreateStageParams) error
	DeleteStage(params models.DeleteStageParams) error
}

// TODO: implement stage manager for directory-based stage structure

type StageManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
}

func NewStageManager(git common.IGit, credentialReader common.CredentialReader) *StageManager {
	stageManager := &StageManager{
		git:              git,
		credentialReader: credentialReader,
	}
	return stageManager
}

func (s StageManager) CreateStage(params models.CreateStageParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	credentials, err := s.credentialReader.GetCredentials(params.ProjectName)
	if err != nil {
		return fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, params.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     params.ProjectName,
		Credentials: credentials,
	}

	if !s.git.ProjectExists(gitContext) {
		return errors.ErrProjectNotFound
	}

	defaultBranch, err := s.git.GetDefaultBranch(gitContext)
	if err != nil {
		return fmt.Errorf("could not determine default branch of project %s: %w", params.ProjectName, err)
	}

	// create new branch from default branch
	if err := s.git.CreateBranch(gitContext, params.StageName, defaultBranch); err != nil {
		return fmt.Errorf("could not check out new branch %s of project %s: %w", params.StageName, params.ProjectName, err)
	}

	_, err = s.git.StageAndCommitAll(gitContext, "created stage")
	if err != nil {
		return fmt.Errorf("could not push new branch %s of project %s: %w", params.StageName, params.ProjectName, err)
	}

	return nil
}

func (s StageManager) DeleteStage(params models.DeleteStageParams) error {
	panic("implement me")
}
