package handler

import (
	"errors"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllStages_GettingProjectFromDBFails(t *testing.T) {

	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	p, err := instance.GetAllStages("my-project")
	assert.Nil(t, p)
	assert.NotNil(t, err)
}

func TestGetAllStages_ProjectNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetAllStages("my-project")
	assert.Nil(t, stage)
	assert.Equal(t, errProjectNotFound, err)

}

func TestGetAllStages(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		s1 := &models.ExpandedStage{
			StageName: "stage1",
		}
		s2 := &models.ExpandedStage{
			StageName: "stage2",
		}
		p := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{s1, s2},
		}
		return p, nil
	}

	stages, err := instance.GetAllStages("my-project")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, "my-project", stagesDbOperations.GetProjectCalls()[0].ProjectName)
}

func TestGetStage_GettingProjectFromDBFails(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.NotNil(t, err)
}

func TestGetStage_ProjectNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.Equal(t, errProjectNotFound, err)

}

func TestGetStage_StageNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	instance := NewStageManager(stagesDbOperations)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		s1 := &models.ExpandedStage{
			StageName: "stage1",
		}
		s2 := &models.ExpandedStage{
			StageName: "stage2",
		}
		p := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{s1, s2},
		}
		return p, nil
	}
	stage, err := instance.GetStage("my-project", "unknown-stage")
	assert.Nil(t, stage)
	assert.Equal(t, errStageNotFound, err)
}
