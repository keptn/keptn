package handler

import "github.com/gin-gonic/gin"

type IRegistrationHandler interface {
	CreateRegistration(context *gin.Context)
	DeleteRegistration(context *gin.Context)
	GetRegistrations(context *gin.Context)
}

type RegistrationHandler struct{}

func (RegistrationHandler) CreateRegistration(context *gin.Context) {
	panic("implement me")
}

func (RegistrationHandler) DeleteRegistration(context *gin.Context) {
	panic("implement me")
}

func (RegistrationHandler) GetRegistrations(context *gin.Context) {
	panic("implement me")
}
