// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/keptn/keptn/api/restapi/operations/auth"
	"github.com/keptn/keptn/api/restapi/operations/dynatrace"

	"github.com/keptn/keptn/api/restapi/operations/project"

	"github.com/keptn/keptn/api/restapi/operations/configure"
	"github.com/keptn/keptn/api/utils"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/google/uuid"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	models "github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/restapi/operations/openws"
	"github.com/keptn/keptn/api/ws"
)

//go:generate swagger generate server --target ../../api --name  --spec ../swagger.json --principal models.Principal

var hub *ws.Hub

func configureFlags(api *operations.API) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.API) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	api.KeyAuth = func(token string) (*models.Principal, error) {
		if token == os.Getenv("keptn-api-token") {
			prin := models.Principal(token)
			return &prin, nil
		}
		api.Logger("Access attempt with incorrect api key auth: %s", token)
		return nil, errors.New(401, "incorrect api key auth")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.AuthAuthHandler = auth.AuthHandlerFunc(func(params auth.AuthParams, principal *models.Principal) middleware.Responder {
		return auth.NewAuthOK()
	})

	api.EventSendEventHandler = event.SendEventHandlerFunc(func(params event.SendEventParams, principal *models.Principal) middleware.Responder {
		if params.Body.Shkeptncontext == nil || *params.Body.Shkeptncontext == "" {
			uuidStr := uuid.New().String()
			params.Body.Shkeptncontext = &uuidStr
		}
		keptnutils.Info(*params.Body.Shkeptncontext, "API received keptn event")

		if err := utils.PostToEventBroker(params.Body, *params.Body.Shkeptncontext); err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}
		token, err := ws.CreateChannelInfo(*params.Body.Shkeptncontext)
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}

		channelInfo := getChannelInfo(params.Body.Shkeptncontext, &token)
		data, err := channelInfo.MarshalJSON()
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}
		resp := event.SendEventCreatedBody{}
		err = resp.UnmarshalJSON(data)
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}
		return event.NewSendEventCreated().WithPayload(&resp)
	})

	api.ConfigureConfigureHandler = configure.ConfigureHandlerFunc(func(params configure.ConfigureParams, principal *models.Principal) middleware.Responder {
		if params.Body.Shkeptncontext == "" {
			params.Body.Shkeptncontext = uuid.New().String()
		}
		keptnutils.Info(params.Body.Shkeptncontext, "API received configure event")

		if err := utils.PostToEventBroker(params.Body, params.Body.Shkeptncontext); err != nil {
			return configure.NewConfigureDefault(500).WithPayload(&configure.ConfigureDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}
		token, err := ws.CreateChannelInfo(params.Body.Shkeptncontext)
		if err != nil {
			return configure.NewConfigureDefault(500).WithPayload(&configure.ConfigureDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}

		channelInfo := getChannelInfo(&params.Body.Shkeptncontext, &token)
		data, err := channelInfo.MarshalJSON()
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}
		resp := configure.ConfigureCreatedBody{}
		err = resp.UnmarshalJSON(data)
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}

		return configure.NewConfigureCreated().WithPayload(&resp)
	})

	api.ProjectProjectHandler = project.ProjectHandlerFunc(func(params project.ProjectParams, principal *models.Principal) middleware.Responder {
		if params.Body.Shkeptncontext == "" {
			params.Body.Shkeptncontext = uuid.New().String()
		}
		keptnutils.Info(params.Body.Shkeptncontext, "API received project event")

		if err := utils.PostToEventBroker(params.Body, params.Body.Shkeptncontext); err != nil {
			return project.NewProjectDefault(500).WithPayload(&project.ProjectDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}
		token, err := ws.CreateChannelInfo(params.Body.Shkeptncontext)
		if err != nil {
			return project.NewProjectDefault(500).WithPayload(&project.ProjectDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}

		channelInfo := getChannelInfo(&params.Body.Shkeptncontext, &token)
		data, err := channelInfo.MarshalJSON()
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}
		resp := project.ProjectCreatedBody{}
		err = resp.UnmarshalJSON(data)
		if err != nil {
			return event.NewSendEventDefault(500).WithPayload(&event.SendEventDefaultBody{Code: 500, Message: swag.String("Marshalling error")})
		}

		return project.NewProjectCreated().WithPayload(&resp)
	})

	api.DynatraceDynatraceHandler = dynatrace.DynatraceHandlerFunc(func(params dynatrace.DynatraceParams, principal *models.Principal) middleware.Responder {
		if params.Body.Shkeptncontext == "" {
			params.Body.Shkeptncontext = uuid.New().String()
		}
		keptnutils.Info(params.Body.Shkeptncontext, "API received Dynatrace event")

		if err := utils.PostToEventBroker(params.Body, params.Body.Shkeptncontext); err != nil {
			return dynatrace.NewDynatraceDefault(500).WithPayload(&dynatrace.DynatraceDefaultBody{Code: 500, Message: swag.String(err.Error())})
		}
		return dynatrace.NewDynatraceCreated()
	})

	api.OpenwsOpenWSHandler = openws.OpenWSHandlerFunc(func(params openws.OpenWSParams, pincipal *models.Principal) middleware.Responder {

		// Verify token
		fmt.Println("Open WS")
		err := ws.VerifyToken(params.HTTPRequest.Header)
		if err != nil {
			return openws.NewOpenWSDefault(401).WithPayload(&openws.OpenWSDefaultBody{Code: 401, Message: swag.String(err.Error())})
		}

		return middleware.ResponderFunc(func(rw http.ResponseWriter, _ runtime.Producer) {
			if val, ok := params.HTTPRequest.Header["Keptn-Ws-Channel-Id"]; ok {
				fmt.Println("Serve CLI")
				ws.ServeWsCLI(hub, rw, params.HTTPRequest, val[0])
			} else {
				fmt.Println("Serve service")
				ws.ServeWs(hub, rw, params.HTTPRequest)
			}
		})
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func getChannelInfo(channelID *string, token *string) models.ChannelInfo {

	id := uuid.New().String()
	source := "api"
	specversion := "0.2"
	time := strfmt.DateTime(time.Now())
	typeInfo := "ChannelInfo"
	channelInfo := models.ChannelInfoAO1DataChannelInfo{ChannelID: channelID, Token: token}
	data := models.ChannelInfoAO1Data{ChannelInfo: &channelInfo}

	return models.ChannelInfo{Contenttype: "application/json", ID: &id, Source: &source, Specversion: &specversion, Time: time, Type: &typeInfo, Data: &data}
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
	hub := ws.NewHub()
	go hub.Run()
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/" {
			// Verify token
			err := ws.VerifyToken(r.Header)
			if err != nil {
				w.WriteHeader(401)
				return
			}

			if val, ok := r.Header["Keptn-Ws-Channel-Id"]; ok {
				err = ws.ServeWsCLI(hub, w, r, val[0])
			} else {
				err = ws.ServeWs(hub, w, r)
			}
			if err != nil {
				w.WriteHeader(500)
			}
			return
		}

		handler.ServeHTTP(w, r)
	})
}
