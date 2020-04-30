package handlers

import (
	"fmt"
	"github.com/keptn/keptn/configuration-service/restapi/operations/services"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	utils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service"
	"gopkg.in/yaml.v2"
)

type serviceMetadata struct {
	ServiceName       string
	CreationTimestamp string
}

// GetProjectProjectNameStageStageNameServiceHandlerFunc get list of services
func GetProjectProjectNameStageStageNameServiceHandlerFunc(params service.GetProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return stage.NewGetProjectProjectNameStageDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return stage.NewGetProjectProjectNameStageNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	var payload = &models.Services{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Services:    []*models.ExpandedService{},
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			paginationInfo := common.Paginate(len(stg.Services), params.PageSize, params.NextPageKey)
			totalCount := len(stg.Services)
			if paginationInfo.NextPageKey < int64(totalCount) {
				for _, svc := range stg.Services[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
					payload.Services = append(payload.Services, svc)
				}
			}
			payload.TotalCount = float64(totalCount)
			payload.NextPageKey = paginationInfo.NewNextPageKey
			return service.NewGetProjectProjectNameStageStageNameServiceOK().WithPayload(payload)
		}
	}
	return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
}

// GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc get the specified service
func GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.GetProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return service.NewGetProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	if prj == nil {
		return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == params.ServiceName {
					return service.NewGetProjectProjectNameStageStageNameServiceServiceNameOK().WithPayload(svc)
				}
			}
			return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
		}
	}
	return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
}

// PostProjectProjectNameStageStageNameServiceHandlerFunc creates a new service
func PostProjectProjectNameStageStageNameServiceHandlerFunc(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	servicePath := projectConfigPath + "/" + params.Service.ServiceName

	if !common.StageExists(params.ProjectName, params.StageName, false) {
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}

	if common.ServiceExists(params.ProjectName, params.StageName, params.Service.ServiceName, false) {
		return service.NewPostProjectProjectNameStageStageNameServiceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service already exists")})
	}

	logger.Debug("Creating new resource(s) in: " + projectConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	err = os.MkdirAll(servicePath, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create service directory")})
	}

	newServiceMetadata := &serviceMetadata{
		ServiceName:       params.Service.ServiceName,
		CreationTimestamp: time.Now().String(),
	}

	metadataString, err := yaml.Marshal(newServiceMetadata)
	err = common.WriteFile(servicePath+"/metadata.yaml", metadataString)

	common.StageAndCommitAll(params.ProjectName, "Added service: "+params.Service.ServiceName, true)

	mv := common.GetProjectsMaterializedView()
	err = mv.CreateService(params.ProjectName, params.StageName, params.Service.ServiceName)
	if err != nil {
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	return service.NewPostProjectProjectNameStageStageNameServiceNoContent()
}

// PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc updates a service
func PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.PutProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.PutProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc deletes a service
func DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.DeleteProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.DeleteProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}

func GetServices(params services.GetServicesParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return services.NewGetServicesDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return services.NewGetServicesNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	allServices := []*models.ExpandedServiceWithStageInfo{}

	for _, stg := range prj.Stages {
		var serviceToUpdate *models.ExpandedServiceWithStageInfo
		for _, svc := range stg.Services {
			appendService := true
			for _, sv := range allServices {
				if sv.ServiceName == svc.ServiceName {
					serviceToUpdate = sv
					appendService = false
				}
			}
			if appendService {
				serviceToUpdate = &models.ExpandedServiceWithStageInfo{
					ServiceName:  svc.ServiceName,
					CreationDate: svc.CreationDate,
					StageInfo:    []*models.InverseServiceStageInfo{},
				}
				allServices = append(allServices, serviceToUpdate)
			}
			stageInfo := &models.InverseServiceStageInfo{
				DeployedImage:  svc.DeployedImage,
				LastEventTypes: svc.LastEventTypes,
				StageName:      stg.StageName,
			}
			serviceToUpdate.StageInfo = append(serviceToUpdate.StageInfo, stageInfo)
		}
	}

	var payload = &models.ServicesWithStageInfo{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Services:    []*models.ExpandedServiceWithStageInfo{},
	}

	paginationInfo := common.Paginate(len(allServices), params.PageSize, params.NextPageKey)

	totalCount := len(allServices)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, svc := range allServices[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Services = append(payload.Services, svc)
		}
	}
	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return services.NewGetServicesOK().WithPayload(payload)
}

func GetService(params services.GetServiceParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return services.NewGetServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return services.NewGetServiceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	var serviceResult *models.ExpandedServiceWithStageInfo
	for _, stg := range prj.Stages {
		for _, svc := range stg.Services {
			if svc.ServiceName == params.ServiceName {
				if serviceResult == nil {
					serviceResult = &models.ExpandedServiceWithStageInfo{
						CreationDate: svc.CreationDate,
						ServiceName:  svc.ServiceName,
						StageInfo:    []*models.InverseServiceStageInfo{},
					}
				}
			}
			stageInfo := &models.InverseServiceStageInfo{
				DeployedImage:  svc.DeployedImage,
				LastEventTypes: svc.LastEventTypes,
				StageName:      stg.StageName,
			}
			serviceResult.StageInfo = append(serviceResult.StageInfo, stageInfo)
		}
	}

	if serviceResult != nil {
		return services.NewGetServiceOK().WithPayload(serviceResult)
	}

	return services.NewGetServiceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
}
