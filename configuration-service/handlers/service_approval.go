package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_approval"
)

// CreateServiceApproval creates a service approval
func CreateServiceApproval(params service_approval.CreateServiceApprovalParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	mv := common.GetProjectsMaterializedView()
	err := mv.CreateOpenApproval(params.ProjectName, params.StageName, params.ServiceName, params.Approval)

	if err != nil {
		return service_approval.NewCreateServiceApprovalDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return service_approval.NewCreateServiceApprovalOK()
}

// GetServiceApprovals returns all service approvals
func GetServiceApprovals(params service_approval.GetServiceApprovalsParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return service_approval.NewGetServiceApprovalsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return service_approval.NewGetServiceApprovalsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	payload := &models.Approvals{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Approvals:   []*models.Approval{},
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == params.ServiceName {
					paginationInfo := common.Paginate(len(svc.OpenApprovals), params.PageSize, params.NextPageKey)
					totalCount := len(svc.OpenApprovals)
					if paginationInfo.NextPageKey < int64(totalCount) {
						payload.Approvals = svc.OpenApprovals[paginationInfo.NextPageKey:paginationInfo.EndIndex]
					}
					payload.TotalCount = float64(totalCount)
					payload.NextPageKey = paginationInfo.NewNextPageKey
					return service_approval.NewGetServiceApprovalsOK().WithPayload(payload)
				}
			}
		}
	}
	return service_approval.NewGetServiceApprovalsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
}

// GetServiceApproval returns a service approval by its id
func GetServiceApproval(params service_approval.GetServiceApprovalParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return service_approval.NewGetServiceApprovalsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return service_approval.NewGetServiceApprovalsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == params.ServiceName {
					for _, approval := range svc.OpenApprovals {
						if approval.EventID == params.ApprovalID {
							return service_approval.NewGetServiceApprovalOK().WithPayload(approval)
						}
					}
				}
			}
		}
	}
	return service_approval.NewGetServiceApprovalNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})

}

// CloseServiceApproval closes a service approval
func CloseServiceApproval(params service_approval.CloseServiceApprovalParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	mv := common.GetProjectsMaterializedView()

	err := mv.CloseOpenApproval(params.ProjectName, params.StageName, params.ServiceName, params.ApprovalID)

	if err != nil {
		if err == common.ErrOpenApprovalNotFound {
			return service_approval.NewCloseServiceApprovalDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("Could not close approval: " + err.Error())})
		}
		return service_approval.NewCloseServiceApprovalDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not close approval: " + err.Error())})
	}

	return service_approval.NewCloseServiceApprovalOK()
}
