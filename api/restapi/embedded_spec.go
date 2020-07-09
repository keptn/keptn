// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json",
    "application/cloudevents+json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "keptn api",
    "version": "develop"
  },
  "basePath": "/v1",
  "paths": {
    "/auth": {
      "post": {
        "tags": [
          "Auth"
        ],
        "summary": "Checks the provided token",
        "operationId": "auth",
        "responses": {
          "200": {
            "description": "Authenticated"
          }
        }
      }
    },
    "/configure/bridge/expose": {
      "post": {
        "tags": [
          "configure"
        ],
        "summary": "Exposes the bridge",
        "parameters": [
          {
            "$ref": "#/parameters/configureBridge"
          }
        ],
        "responses": {
          "200": {
            "description": "Bridge was successfully exposed/disposed",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bridge could not be exposed/disposed",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      }
    },
    "/event": {
      "get": {
        "tags": [
          "Event"
        ],
        "summary": "Deprecated endpoint - please use /mongodb-datastore/v1/event",
        "deprecated": true,
        "parameters": [
          {
            "type": "string",
            "description": "KeptnContext of the events to get",
            "name": "keptnContext",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "description": "Type of the Keptn cloud event",
            "name": "type",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "response_model.yaml#/definitions/keptnContextExtendedCE"
            }
          },
          "404": {
            "description": "Failed. Event could not be found.",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Event"
        ],
        "summary": "Forwards the received event",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "response_model.yaml#/definitions/keptnContextExtendedCE"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Forwarded",
            "schema": {
              "$ref": "response_model.yaml#/definitions/eventContext"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      }
    },
    "/metadata": {
      "get": {
        "tags": [
          "Metadata"
        ],
        "summary": "Get keptn installation metadata",
        "operationId": "metadata",
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "response_model.yaml#/definitions/metadata"
            }
          }
        }
      }
    },
    "/project": {
      "post": {
        "tags": [
          "Project"
        ],
        "summary": "Creates a new project",
        "parameters": [
          {
            "$ref": "#/parameters/project"
          }
        ],
        "responses": {
          "200": {
            "description": "Creating of project triggered",
            "schema": {
              "$ref": "response_model.yaml#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be created",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      }
    },
    "/project/{projectName}": {
      "delete": {
        "tags": [
          "Project"
        ],
        "summary": "Deletes the specified project",
        "responses": {
          "200": {
            "description": "Deleting of project triggered",
            "schema": {
              "$ref": "response_model.yaml#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be deleted",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      },
      "parameters": [
        {
          "$ref": "#/parameters/projectName"
        }
      ]
    },
    "/project/{projectName}/service": {
      "post": {
        "tags": [
          "Service"
        ],
        "summary": "Creates a new service",
        "parameters": [
          {
            "$ref": "#/parameters/service"
          }
        ],
        "responses": {
          "200": {
            "description": "Creating of service triggered",
            "schema": {
              "$ref": "response_model.yaml#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be created",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "response_model.yaml#/definitions/error"
            }
          }
        }
      },
      "parameters": [
        {
          "$ref": "#/parameters/projectName"
        }
      ]
    }
  },
  "parameters": {
    "configureBridge": {
      "description": "Parameters for configuring the bridge access",
      "name": "configureBridge",
      "in": "body",
      "schema": {
        "$ref": "configure_model.yaml#/definitions/configureBridge"
      }
    },
    "project": {
      "description": "Project entity",
      "name": "project",
      "in": "body",
      "schema": {
        "$ref": "project_model.yaml#/definitions/project"
      }
    },
    "projectName": {
      "type": "string",
      "description": "Name of the project",
      "name": "projectName",
      "in": "path",
      "required": true
    },
    "service": {
      "description": "Service entity",
      "name": "service",
      "in": "body",
      "schema": {
        "$ref": "service_model.yaml#/definitions/service"
      }
    },
    "stageName": {
      "type": "string",
      "description": "Name of the stage",
      "name": "stageName",
      "in": "path",
      "required": true
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "security": [
    {
      "key": []
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json",
    "application/cloudevents+json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "keptn api",
    "version": "develop"
  },
  "basePath": "/v1",
  "paths": {
    "/auth": {
      "post": {
        "tags": [
          "Auth"
        ],
        "summary": "Checks the provided token",
        "operationId": "auth",
        "responses": {
          "200": {
            "description": "Authenticated"
          }
        }
      }
    },
    "/configure/bridge/expose": {
      "post": {
        "tags": [
          "configure"
        ],
        "summary": "Exposes the bridge",
        "parameters": [
          {
            "description": "Parameters for configuring the bridge access",
            "name": "configureBridge",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/configureBridge"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Bridge was successfully exposed/disposed",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bridge could not be exposed/disposed",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/event": {
      "get": {
        "tags": [
          "Event"
        ],
        "summary": "Deprecated endpoint - please use /mongodb-datastore/v1/event",
        "deprecated": true,
        "parameters": [
          {
            "type": "string",
            "description": "KeptnContext of the events to get",
            "name": "keptnContext",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "description": "Type of the Keptn cloud event",
            "name": "type",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/keptnContextExtendedCE"
            }
          },
          "404": {
            "description": "Failed. Event could not be found.",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Event"
        ],
        "summary": "Forwards the received event",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/keptnContextExtendedCE"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Forwarded",
            "schema": {
              "$ref": "#/definitions/eventContext"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/metadata": {
      "get": {
        "tags": [
          "Metadata"
        ],
        "summary": "Get keptn installation metadata",
        "operationId": "metadata",
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "$ref": "#/definitions/metadata"
            }
          }
        }
      }
    },
    "/project": {
      "post": {
        "tags": [
          "Project"
        ],
        "summary": "Creates a new project",
        "parameters": [
          {
            "description": "Project entity",
            "name": "project",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/project"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Creating of project triggered",
            "schema": {
              "$ref": "#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be created",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/project/{projectName}": {
      "delete": {
        "tags": [
          "Project"
        ],
        "summary": "Deletes the specified project",
        "responses": {
          "200": {
            "description": "Deleting of project triggered",
            "schema": {
              "$ref": "#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be deleted",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "parameters": [
        {
          "type": "string",
          "description": "Name of the project",
          "name": "projectName",
          "in": "path",
          "required": true
        }
      ]
    },
    "/project/{projectName}/service": {
      "post": {
        "tags": [
          "Service"
        ],
        "summary": "Creates a new service",
        "parameters": [
          {
            "description": "Service entity",
            "name": "service",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/service"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Creating of service triggered",
            "schema": {
              "$ref": "#/definitions/eventContext"
            }
          },
          "400": {
            "description": "Failed. Project could not be created",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "parameters": [
        {
          "type": "string",
          "description": "Name of the project",
          "name": "projectName",
          "in": "path",
          "required": true
        }
      ]
    }
  },
  "definitions": {
    "configureBridge": {
      "type": "object",
      "required": [
        "expose"
      ],
      "properties": {
        "expose": {
          "type": "boolean"
        },
        "password": {
          "type": "string"
        },
        "user": {
          "type": "string"
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "fields": {
          "type": "string"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "eventContext": {
      "type": "object",
      "required": [
        "token",
        "keptnContext"
      ],
      "properties": {
        "keptnContext": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    },
    "keptnContextExtendedCE": {
      "type": "object",
      "required": [
        "data",
        "source",
        "type"
      ],
      "properties": {
        "contenttype": {
          "type": "string"
        },
        "data": {
          "type": [
            "object",
            "string"
          ]
        },
        "extensions": {
          "type": "object"
        },
        "id": {
          "type": "string"
        },
        "shkeptncontext": {
          "type": "string"
        },
        "source": {
          "type": "string",
          "format": "uri-reference"
        },
        "specversion": {
          "type": "string"
        },
        "time": {
          "type": "string",
          "format": "date-time"
        },
        "triggeredid": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "metadata": {
      "type": "object",
      "properties": {
        "bridgeversion": {
          "type": "string"
        },
        "keptnversion": {
          "type": "string"
        },
        "namespace": {
          "type": "string"
        }
      }
    },
    "project": {
      "type": "object",
      "required": [
        "name",
        "shipyard"
      ],
      "properties": {
        "gitRemoteURL": {
          "type": "string"
        },
        "gitToken": {
          "type": "string"
        },
        "gitUser": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "shipyard": {
          "type": "string"
        }
      }
    },
    "service": {
      "type": "object",
      "required": [
        "serviceName"
      ],
      "properties": {
        "deploymentStrategies": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "helmChart": {
          "type": "string"
        },
        "serviceName": {
          "type": "string"
        }
      }
    }
  },
  "parameters": {
    "configureBridge": {
      "description": "Parameters for configuring the bridge access",
      "name": "configureBridge",
      "in": "body",
      "schema": {
        "$ref": "#/definitions/configureBridge"
      }
    },
    "project": {
      "description": "Project entity",
      "name": "project",
      "in": "body",
      "schema": {
        "$ref": "#/definitions/project"
      }
    },
    "projectName": {
      "type": "string",
      "description": "Name of the project",
      "name": "projectName",
      "in": "path",
      "required": true
    },
    "service": {
      "description": "Service entity",
      "name": "service",
      "in": "body",
      "schema": {
        "$ref": "#/definitions/service"
      }
    },
    "stageName": {
      "type": "string",
      "description": "Name of the stage",
      "name": "stageName",
      "in": "path",
      "required": true
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "security": [
    {
      "key": []
    }
  ]
}`))
}
