package handler

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/internal/db/mock"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllStages_GettingProjectFromDBFails(t *testing.T) {

	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	p, err := instance.GetAllStages("my-project")
	assert.Nil(t, p)
	assert.NotNil(t, err)
}

func TestGetAllStages_ProjectNotFound(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetAllStages("my-project")
	assert.Nil(t, stage)
	assert.Equal(t, common.ErrProjectNotFound, err)

}

func TestGetAllStages(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		s1 := &apimodels.ExpandedStage{
			StageName: "stage1",
		}
		s2 := &apimodels.ExpandedStage{
			StageName: "stage2",
		}
		p := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{s1, s2},
		}
		return p, nil
	}

	stages, err := instance.GetAllStages("my-project")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetStage_GettingProjectFromDBFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.NotNil(t, err)
}

func TestGetStage_ProjectNotFound(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.Equal(t, common.ErrProjectNotFound, err)

}

func TestGetStage_StageNotFound(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	instance := NewStageManager(projectMVRepo)

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		s1 := &apimodels.ExpandedStage{
			StageName: "stage1",
		}
		s2 := &apimodels.ExpandedStage{
			StageName: "stage2",
		}
		p := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{s1, s2},
		}
		return p, nil
	}
	stage, err := instance.GetStage("my-project", "unknown-stage")
	assert.Nil(t, stage)
	assert.Equal(t, common.ErrStageNotFound, err)
}
