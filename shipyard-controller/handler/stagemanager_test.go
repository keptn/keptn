package handler

import (
	"errors"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllStages_GettingProjectFromDBFails(t *testing.T) {

	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("Whoops...")
	}

	p, err := instance.GetAllStages("my-project")
	assert.Nil(t, p)
	assert.NotNil(t, err)
}

func TestGetAllStages_ProjectNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetAllStages("my-project")
	assert.Nil(t, stage)
	assert.Equal(t, errProjectNotFound, err)

}

func TestGetAllStages(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

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
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("Whoops...")
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.NotNil(t, err)
}

func TestGetStage_ProjectNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

	stagesDbOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	stage, err := instance.GetStage("my-project", "the-stage")
	assert.Nil(t, stage)
	assert.Equal(t, errProjectNotFound, err)

}

func TestGetStage_StageNotFound(t *testing.T) {
	stagesDbOperations := &db_mock.StagesDbOperationsMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewStageManager(stagesDbOperations, logger)

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
