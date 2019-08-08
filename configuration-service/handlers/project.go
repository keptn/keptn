package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
)

// GetProjectHandlerFunc gets a list of projects
func GetProjectHandlerFunc(params project.GetProjectParams) middleware.Responder {
	return middleware.NotImplemented("operation project.GetProject has not yet been implemented")
}

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {
	return middleware.NotImplemented("operation project.PostProject has not yet been implemented")
}

// GetProjectProjectNameHandlerFunc gets a project by its name
func GetProjectProjectNameHandlerFunc(params project.GetProjectProjectNameParams) middleware.Responder {
	return middleware.NotImplemented("operation project.GetProjectProjectName has not yet been implemented")
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	return middleware.NotImplemented("operation project.PutProjectProjectName has not yet been implemented")
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams) middleware.Responder {
	return middleware.NotImplemented("operation project.DeleteProjectProjectName has not yet been implemented")
}
