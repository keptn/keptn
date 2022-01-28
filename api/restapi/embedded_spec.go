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
    "/event": {
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
    }
  },
  "definitions": {
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
        "keptnContext"
      ],
      "properties": {
        "keptnContext": {
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
          "type": "object"
        },
        "extensions": {
          "type": "object"
        },
        "gitcommitid": {
          "type": "string"
        },
        "id": {
          "type": "string"
        },
        "shkeptncontext": {
          "type": "string"
        },
        "shkeptnversion": {
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
        "keptnlabel": {
          "type": "string"
        },
        "keptnservices": {
          "type": "object"
        },
        "keptnversion": {
          "type": "string"
        },
        "namespace": {
          "type": "string"
        },
        "shipyardversion": {
          "type": "string"
        }
      }
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
    "application/cloudevents+json",
    "application/json"
  ],
  "produces": [
    "application/json"
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
    "/event": {
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
    }
  },
  "definitions": {
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
        "keptnContext"
      ],
      "properties": {
        "keptnContext": {
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
          "type": "object"
        },
        "extensions": {
          "type": "object"
        },
        "gitcommitid": {
          "type": "string"
        },
        "id": {
          "type": "string"
        },
        "shkeptncontext": {
          "type": "string"
        },
        "shkeptnversion": {
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
        "keptnlabel": {
          "type": "string"
        },
        "keptnservices": {
          "type": "object"
        },
        "keptnversion": {
          "type": "string"
        },
        "namespace": {
          "type": "string"
        },
        "shipyardversion": {
          "type": "string"
        }
      }
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
