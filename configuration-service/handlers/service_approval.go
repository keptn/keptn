package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/remediation"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_approval"
)

// CreateServiceApproval creates a service approval
func CreateServiceApproval(params service_approval.CreateServiceApprovalParams) middleware.Responder {
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

func CreateRemediation(params remediation.CreateRemediationParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	err := mv.CreateRemediation(params.ProjectName, params.StageName, params.ServiceName, params.Remediation)

	if err != nil {
		return service_approval.NewCreateRemediationDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return service_approval.NewCreateRemediationOK()
}

func GetRemediations(params remediation.GetRemediationsParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return remediation.NewGetRemediationsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return remediation.NewGetRemediationsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	payload := &models.Remediations{
		PageSize:     0,
		NextPageKey:  "0",
		TotalCount:   0,
		Remediations: []*models.Remediation{},
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == params.ServiceName {
					paginationInfo := common.Paginate(len(svc.OpenRemediations), params.PageSize, params.NextPageKey)
					totalCount := len(svc.OpenRemediations)
					if paginationInfo.NextPageKey < int64(totalCount) {
						payload.Remediations = svc.OpenRemediations[paginationInfo.NextPageKey:paginationInfo.EndIndex]
					}
					payload.TotalCount = float64(totalCount)
					payload.NextPageKey = paginationInfo.NewNextPageKey
					return remediation.NewGetRemediationsOK().WithPayload(payload)
				}
			}
		}
	}
	return remediation.NewGetRemediationsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
}

func GetRemediationsForContext(params remediation.GetRemediationsForContextParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return remediation.NewGetRemediationsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return remediation.NewGetRemediationsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	payload := &models.Remediations{
		PageSize:     0,
		NextPageKey:  "0",
		TotalCount:   0,
		Remediations: []*models.Remediation{},
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == params.ServiceName {
					remediations := []*models.Remediation{}
					for _, remediation := range svc.OpenRemediations {
						if remediation.KeptnContext == params.KeptnContext {
							remediations = append(remediations, remediation)
						}
					}
					paginationInfo := common.Paginate(len(remediations), params.PageSize, params.NextPageKey)
					totalCount := len(remediations)
					if paginationInfo.NextPageKey < int64(totalCount) {
						payload.Remediations = remediations[paginationInfo.NextPageKey:paginationInfo.EndIndex]
					}
					payload.TotalCount = float64(totalCount)
					payload.NextPageKey = paginationInfo.NewNextPageKey
					return remediation.NewGetRemediationsOK().WithPayload(payload)
				}
			}
		}
	}
	return remediation.NewGetRemediationsNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
}

func CloseRemediations(params remediation.CloseRemediationsParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	err := mv.CloseOpenRemediations(params.ProjectName, params.StageName, params.ServiceName, params.KeptnContext)

	if err != nil {
		if err == common.ErrOpenRemediationNotFound {
			return remediation.NewCloseRemediationsDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String("Could not close remediation: " + err.Error())})
		}
		return remediation.NewCloseRemediationsDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not close remediation: " + err.Error())})
	}

	return remediation.NewCloseRemediationsOK()
}
