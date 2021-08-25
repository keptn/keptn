package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"
)

type IUniformIntegrationHandler interface {
	Register(context *gin.Context)
	KeepAlive(context *gin.Context)
	Unregister(context *gin.Context)
	GetRegistrations(context *gin.Context)
	GetSubscription(context *gin.Context)
	GetSubscriptions(c *gin.Context)
	CreateSubscription(c *gin.Context)
	DeleteSubscription(c *gin.Context)
	UpdateSubscription(c *gin.Context)
}

type UniformIntegrationHandler struct {
	uniformRepo db.UniformRepo
}

func NewUniformIntegrationHandler(uniformRepo db.UniformRepo) *UniformIntegrationHandler {
	return &UniformIntegrationHandler{uniformRepo: uniformRepo}
}

// Register creates or updates a uniform integration
// @Summary BETA: Register a uniform integration
// @Description Register a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integration body models.Integration true "Integration"
// @Success 201 {object} models.RegisterResponse "ok"
// @Success 200 {object} models.RegisterResponse "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [post]
func (rh *UniformIntegrationHandler) Register(c *gin.Context) {
	integration := &models.Integration{}

	if err := c.ShouldBindJSON(integration); err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}

	integrationID := keptnmodels.IntegrationID{
		Name:      integration.Name,
		Namespace: integration.MetaData.KubernetesMetaData.Namespace,
		NodeName:  integration.MetaData.Hostname,
	}

	hash, err := integrationID.Hash()
	if err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}

	//setting IDs and last seen timestamp
	integration.ID = hash
	integration.MetaData.LastSeen = time.Now().UTC()
	for i := range integration.Subscriptions {
		s := &integration.Subscriptions[i]
		s.ID = uuid.New().String()
	}

	// for backwards compatibility, we check if there is a Subscriptions field set
	// if not, we are taking the old Subscription field and map it to the new Subscriptions field
	// Note: "old" registrations will NOT get subscription IDs
	// This code block can be deleted with later versions of Keptn
	if integration.Subscriptions == nil {
		var projectFilter []string
		var stageFilter []string
		var serviceFilter []string
		if integration.Subscription.Filter.Project != "" {
			projectFilter = strings.Split(integration.Subscription.Filter.Project, ",")
		}
		if integration.Subscription.Filter.Stage != "" {
			stageFilter = strings.Split(integration.Subscription.Filter.Stage, ",")
		}
		if integration.Subscription.Filter.Service != "" {
			serviceFilter = strings.Split(integration.Subscription.Filter.Service, ",")
		}

		for _, t := range integration.Subscription.Topics {
			ts := keptnmodels.EventSubscription{
				Event: t,
				Filter: keptnmodels.EventSubscriptionFilter{
					Projects: projectFilter,
					Stages:   stageFilter,
					Services: serviceFilter,
				},
			}
			integration.Subscriptions = append(integration.Subscriptions, ts)
		}

		raw := fmt.Sprintf("%s-%s-%s-%s-%s", integration.Name, integration.MetaData.KubernetesMetaData.Namespace, integration.Subscription.Filter.Project, integration.Subscription.Filter.Stage, integration.Subscription.Filter.Service)
		hasher := sha1.New() //nolint:gosec
		hasher.Write([]byte(raw))
		hash = hex.EncodeToString(hasher.Sum(nil))
		integration.ID = hash
	}

	err = rh.uniformRepo.CreateUniformIntegration(*integration)
	if err != nil {
		// if the integration already exists, update only the last seen field
		// and return integration ID
		if errors.Is(err, db.ErrUniformRegistrationAlreadyExists) {
			_, _ = rh.uniformRepo.UpdateLastSeen(integration.ID)
			c.JSON(http.StatusOK, &models.RegisterResponse{
				ID: integration.ID,
			})
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusCreated, &models.RegisterResponse{
		ID: integration.ID,
	})
}

// Unregister deletes a uniform integration
// @Summary BETA: Unregister a uniform integration
// @Description Unregister a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Success 200
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration/{integrationID} [delete]
func (rh *UniformIntegrationHandler) Unregister(c *gin.Context) {
	integrationID := c.Param("integrationID")

	if err := rh.uniformRepo.DeleteUniformIntegration(integrationID); err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}
	c.JSON(http.StatusOK, &models.UnregisterResponse{})
}

// GetRegistrations Retrieve uniform integrations matching the provided filter
// @Summary BETA: Retrieve uniform integrations matching the provided filter
// @Description Retrieve uniform integrations matching the provided filter
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string false "id"
// @Param name query string false "name"
// @Param project query string false "project"
// @Param stage query string false "stage"
// @Param service query string false "service"
// @Success 200 {object} []models.Integration "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [get]
func (rh *UniformIntegrationHandler) GetRegistrations(c *gin.Context) {
	params := &models.GetUniformIntegrationsParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	uniformIntegrations, err := rh.uniformRepo.GetUniformIntegrations(*params)
	if err != nil {
		SetInternalServerErrorResponse(err, c, "Unable to query uniform integrations repository")
		return
	}

	c.JSON(http.StatusOK, uniformIntegrations)
}

// KeepAlive returns current registration data of an integration
// @Summary BETA: Endpoint for sending heartbeat messages sent from Keptn integrations to the control plane
// @Description Endpoint for sending heartbeat messages sent from Keptn integrations to the control plane
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Success 200 {object} models.Integration "ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration/{integrationID}/ping [PUT]
func (rh *UniformIntegrationHandler) KeepAlive(c *gin.Context) {
	integrationID := c.Param("integrationID")

	registration, err := rh.uniformRepo.UpdateLastSeen(integrationID)
	if err != nil {
		if errors.Is(err, db.ErrUniformRegistrationNotFound) {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, registration)

}

// CreateSubscription creates a new subscription
// @Summary BETA: Create a new subscription
// @Description  Create a new subscription
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Param subscription body models.Subscription true "Subscription"
// @Success 201
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Failure 404 {object} models.Error "Not found"
// @Router /uniform/registration/{integrationID}/subscription [post]
func (rh *UniformIntegrationHandler) CreateSubscription(c *gin.Context) {

	integrationID := c.Param("integrationID")
	subscription := &models.Subscription{}

	if err := c.ShouldBindJSON(subscription); err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}
	subscription.ID = uuid.New().String()

	err := rh.uniformRepo.CreateOrUpdateSubscription(integrationID, *subscription)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusCreated, &models.CreateSubscriptionResponse{
		ID: subscription.ID,
	})
}

// UpdateSubscription updates or creates a subscription
// @Summary BETA: Update or create a subscription
// @Description Update or create a subscription
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Param subscriptionID path string true "subscriptionID"
// @Param subscription body models.Subscription true "Subscription"
// @Success 201
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Failure 404 {object} models.Error "Not found"
// @Router /uniform/registration/{integrationID}/subscription/{subscriptionID} [put]
func (rh *UniformIntegrationHandler) UpdateSubscription(c *gin.Context) {

	integrationID := c.Param("integrationID")
	subscriptionID := c.Param("subscriptionID")

	subscription := &models.Subscription{}

	if err := c.ShouldBindJSON(subscription); err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}
	subscription.ID = subscriptionID

	err := rh.uniformRepo.CreateOrUpdateSubscription(integrationID, *subscription)
	if err != nil {
		//TODO: set appropriate http codes
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusCreated, &models.CreateSubscriptionResponse{
		ID: subscription.ID,
	})
}

// DeleteSubscription deletes a new subscription
// @Summary BETA: Delete a subscription
// @Description  Delete a subscription
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Param subscriptionID path string true "subscriptionID"
// @Success 200
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Failure 404 {object} models.Error "Not found"
// @Router /uniform/registration/{integrationID}/subscription/{subscriptionID} [delete]
func (rh *UniformIntegrationHandler) DeleteSubscription(c *gin.Context) {
	integrationID := c.Param("integrationID")
	subscriptionID := c.Param("subscriptionID")

	err := rh.uniformRepo.DeleteSubscription(integrationID, subscriptionID)
	if err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, models.DeleteSubscriptionResponse{})
}

// GetSubscription retrieves an already existing subscription
// @Summary BETA: Retrieve an already existing subscription
// @Description  Retrieve an already existing subscription
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Param subscriptionID path string true "subscriptionID"
// @Success 200 {object} models.Subscription "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Failure 404 {object} models.Error "Not found"
// @Router /uniform/registration/{integrationID}/subscription/{subscriptionID} [get]
func (rh *UniformIntegrationHandler) GetSubscription(c *gin.Context) {
	integrationID := c.Param("integrationID")
	subscriptionID := c.Param("subscriptionID")

	subscription, err := rh.uniformRepo.GetSubscription(integrationID, subscriptionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetSubscriptions retrieves all subscriptions of a uniform integration
// @Summary BETA: Retrieve all subscriptions of a uniform integration
// @Description  Retrieve all subscriptions of a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationID path string true "integrationID"
// @Success 200 {object} []models.Subscription "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Failure 404 {object} models.Error "Not found"
// @Router /uniform/registration/{integrationID}/subscription [get]
func (rh *UniformIntegrationHandler) GetSubscriptions(c *gin.Context) {
	integrationID := c.Param("integrationID")

	subscriptions, err := rh.uniformRepo.GetSubscriptions(integrationID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}
