package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/api-service/model"
	"github.com/keptn/keptn/api-service/utils"

	logger "github.com/sirupsen/logrus"
)

type IEventHandler interface {
	ForwardEvent(c *gin.Context)
	GetEvent(c *gin.Context)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

type EventHandler struct {
}

func (eh *EventHandler) ForwardEvent(c *gin.Context) {
	logger.Info("API received a keptn event")

	event := model.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	keptnContext := createOrApplyKeptnContext(event.ShkeptnContext)

	var source *url.URL
	var err error
	if len(event.Source) > 0 {
		source, err = url.Parse(string(event.Source))
		if err != nil {
			logger.Info("Unable to parse source from the received CloudEvent")
		}
	}

	if source == nil {
		// Use this URL as fallback source
		source, _ = url.Parse("https://github.com/keptn/keptn/api-service")
	}

	err = utils.SendEvent(keptnContext, event.TriggeredID, event.GitcommitID, string(event.Type), source.String(), event.Data)

	if err != nil {
		SetInternalServerErrorResponse(err, c, ErrCreation)
		return
	}

	eventContext := model.EventContext{KeptnContext: &keptnContext}

	c.JSON(http.StatusOK, eventContext)
}

func (eh *EventHandler) GetEvent(c *gin.Context) {
	logger.Info("API received a GET keptn event")

	event := model.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	eventHandler := keptnapi.NewEventHandler(utils.GetDatastoreURL())
	ef := keptnapi.EventFilter{
		EventType:    string(event.Type),
		KeptnContext: event.ShkeptnContext,
	}
	cloudEvent, errObj := eventHandler.GetEvents(&ef)
	if errObj != nil {
		if errObj.Code == 404 {
			SetNotFoundErrorResponse(fmt.Errorf("No "+string(event.Type)+" event found for Keptn context: "+event.ShkeptnContext), c, "Unable to get event")
			return
		}
		SetNotFoundErrorResponse(fmt.Errorf("%s", *errObj.Message), c, "Unable to get event")
		return
	}

	if cloudEvent == nil || len(cloudEvent) == 0 {
		SetNotFoundErrorResponse(fmt.Errorf("No "+string(event.Type)+" event found for Keptn context: "+event.ShkeptnContext), c, "Unable to get event")
		return
	}

	eventByte, err := json.Marshal(cloudEvent[0])
	if err != nil {
		SetInternalServerErrorResponse(err, c, ErrCreation)
		return
	}

	apiEvent := &model.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, apiEvent)
	if err != nil {
		SetInternalServerErrorResponse(err, c, ErrCreation)
		return
	}

	c.JSON(http.StatusOK, apiEvent)
}

func createOrApplyKeptnContext(eventKeptnContext string) string {
	uuid.SetRand(nil)
	keptnContext := uuid.New().String()
	if eventKeptnContext != "" {
		_, err := uuid.Parse(eventKeptnContext)
		if err != nil {
			if len(eventKeptnContext) < 16 {
				paddedContext := fmt.Sprintf("%-16v", eventKeptnContext)
				uuid.SetRand(strings.NewReader(paddedContext))
			} else {
				uuid.SetRand(strings.NewReader(eventKeptnContext))
			}

			keptnContext = uuid.New().String()
			uuid.SetRand(nil)
		} else {
			keptnContext = eventKeptnContext
		}
	}
	return keptnContext
}
