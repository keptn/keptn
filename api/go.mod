module github.com/keptn/keptn/api

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/gbrlsnchs/jwt/v2 v2.0.0
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/loads v0.19.2
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.4
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.6.0
	github.com/kinbiko/jsonassert v1.0.1
	github.com/magiconair/properties v1.8.1
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
)

// replace cloudevents/sdk-go latest version with 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
