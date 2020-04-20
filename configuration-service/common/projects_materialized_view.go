package common

import (
	"errors"
	"fmt"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/models"
	"strconv"
	"time"
)

var instance *projectsMaterializedView

type projectsMaterializedView struct {
	ProjectRepo ProjectRepo
}

func GetProjectsMaterializedView() *projectsMaterializedView {
	if instance == nil {
		instance = &projectsMaterializedView{
			ProjectRepo: &MongoDBProjectRepo{},
		}
	}
	return instance
}

func (mv *projectsMaterializedView) CreateProject(prj *models.Project) error {
	existingProject, err := mv.GetProject(prj.ProjectName)
	if existingProject != nil {
		return nil
	}
	err = mv.createProject(prj)
	if err != nil {
		return err
	}
	return nil
}

func (mv *projectsMaterializedView) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent

	return mv.updateProject(existingProject)
}

func (mv *projectsMaterializedView) GetProjects() ([]*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProjects()
}

func (mv *projectsMaterializedView) GetProject(projectName string) (*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProject(projectName)
}

func (mv *projectsMaterializedView) DeleteProject(projectName string) error {
	return mv.ProjectRepo.DeleteProject(projectName)
}

func (mv *projectsMaterializedView) CreateStage(project string, stage string) error {
	fmt.Println("Adding stage " + stage + " to project " + project)
	prj, err := mv.GetProject(project)

	if err != nil {
		fmt.Sprintf("Could not add stage %s to project %s : %s\n", stage, project, err.Error())
		return err
	}

	stageAlreadyExists := false
	for _, stg := range prj.Stages {
		if stg.StageName == stage {
			stageAlreadyExists = true
			break
		}
	}

	if stageAlreadyExists {
		fmt.Println("Stage " + stage + " already exists in project " + project)
		return nil
	}

	prj.Stages = append(prj.Stages, &models.ExpandedStage{
		Services:  []*models.ExpandedService{},
		StageName: stage,
	})

	err = mv.updateProject(prj)
	if err != nil {
		return err
	}

	fmt.Println("Added stage " + stage + " to project " + project)
	return nil
}

func (mv *projectsMaterializedView) createProject(prj *models.Project) error {
	expandedProject := &models.ExpandedProject{
		CreationDate: strconv.FormatInt(time.Now().UnixNano(), 10),
		GitRemoteURI: prj.GitRemoteURI,
		GitUser:      prj.GitUser,
		ProjectName:  prj.ProjectName,
		Shipyard:     "",
		Stages:       nil,
	}

	err := mv.ProjectRepo.CreateProject(expandedProject)
	if err != nil {
		fmt.Println("Could not create project " + prj.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}

func (mv *projectsMaterializedView) updateProject(prj *models.ExpandedProject) error {
	return mv.ProjectRepo.UpdateProject(prj)
}

func (mv *projectsMaterializedView) DeleteStage(project string, stage string) error {
	fmt.Println("Deleting stage " + stage + " from project " + project)
	prj, err := mv.GetProject(project)

	if err != nil {
		fmt.Sprintf("Could not delete stage %s from project %s : %s\n", stage, project, err.Error())
		return err
	}

	stageIndex := -1

	for idx, stg := range prj.Stages {
		if stg.StageName == stage {
			stageIndex = idx
			break
		}
	}
	if stageIndex < 0 {
		return nil
	}

	copy(prj.Stages[stageIndex:], prj.Stages[stageIndex+1:])
	prj.Stages[len(prj.Stages)-1] = nil
	prj.Stages = prj.Stages[:len(prj.Stages)-1]

	err = mv.updateProject(prj)
	return nil
}

func (mv *projectsMaterializedView) CreateService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		fmt.Println("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			for _, svc := range stg.Services {
				if svc.ServiceName == service {
					break
				}
			}
			stg.Services = append(stg.Services, &models.ExpandedService{
				CreationDate:  strconv.FormatInt(time.Now().UnixNano(), 10),
				DeployedImage: "",
				ServiceName:   service,
			})
			fmt.Println("Adding " + service + " to stage " + stage + " in project " + project + " in database")
			err := mv.updateProject(existingProject)
			if err != nil {
				fmt.Println("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not update project: " + err.Error())
			}
			break
		}
	}
	fmt.Println("Service " + service + " already exists in stage " + stage + " in project " + project)
	return nil
}

func (mv *projectsMaterializedView) DeleteService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		fmt.Println("Could not delete service " + service + " from stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			serviceIndex := -1
			for idx, svc := range stg.Services {
				if svc.ServiceName == service {
					serviceIndex = idx
				}
			}
			if serviceIndex < 0 {
				fmt.Println("Could not delete service " + service + " from stage " + stage + " in project " + project + ". Service not found in database")
				return nil
			}
			copy(stg.Services[serviceIndex:], stg.Services[serviceIndex+1:])
			stg.Services[len(stg.Services)-1] = nil
			stg.Services = stg.Services[:len(stg.Services)-1]
			break
		}
	}
	err = mv.updateProject(existingProject)
	if err != nil {
		fmt.Println("Could not delete service " + service + " from stage " + stage + " in project " + project + ": " + err.Error())
		return err
	}
	fmt.Println("Deleted service " + service + " from stage " + stage + " in project " + project)
	return nil
}

func (mv *projectsMaterializedView) UpdateEventOfService(keptnBase *keptn.KeptnBase, eventType string, keptnContext string) error {
	existingProject, err := mv.GetProject(keptnBase.Project)
	if err != nil {
		fmt.Println("Could not update service " + keptnBase.Service + " in stage " + keptnBase.Stage + " in project " + keptnBase.Project + ". Could not load project: " + err.Error())
		return err
	}

	contextInfo := &models.EventContext{
		KeptnContext: keptnContext,
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	err = updateServiceInStage(existingProject, keptnBase.Stage, keptnBase.Service, func(service *models.ExpandedService) error {
		switch eventType {
		case keptn.ConfigurationChangeEventType:
			service.LastConfigurationChangedEvent = contextInfo
		case keptn.DeploymentFinishedEventType:
			service.LastDeploymentFinishedEvent = contextInfo
			if keptnBase.Image != nil && keptnBase.Tag != nil {
				service.DeployedImage = *keptnBase.Image + ":" + *keptnBase.Tag
			}
		case keptn.TestsFinishedEventType:
			service.LastTestsFinishedEvent = contextInfo
		case keptn.EvaluationDoneEventType:
			service.LastEvaluationDoneEvent = contextInfo
		case keptn.ProblemOpenEventType:
			service.LastProblemEvent = contextInfo
		}
		return nil
	})

	if err != nil {
		fmt.Println("Could not update image of service " + keptnBase.Service + ": " + err.Error())
		return err
	}
	err = mv.updateProject(existingProject)
	if err != nil {
		fmt.Println("Could not update " + keptnBase.Project + ": " + err.Error())
		return err
	}
	return nil
}

type serviceUpdateFunc func(service *models.ExpandedService) error

func updateServiceInStage(project *models.ExpandedProject, stage string, service string, fn serviceUpdateFunc) error {
	for _, stg := range project.Stages {
		if stg.StageName == stage {
			serviceIndex := -1
			for idx, svc := range stg.Services {
				if svc.ServiceName == service {
					serviceIndex = idx
				}
			}
			if serviceIndex < 0 {
				return errors.New("service not found")
			}
			return fn(stg.Services[serviceIndex])
		}
	}
	return errors.New("service not found")
}
