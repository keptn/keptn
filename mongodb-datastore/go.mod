module github.com/keptn/keptn/mongodb-datastore

go 1.12

require (
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/loads v0.19.2
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.4
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jeremywohl/flatten v0.0.0-20190921043622-d936035e55cf
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.0.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.1.1
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
	k8s.io/api v0.0.0-20191016225839-816a9b7df678 // indirect
	k8s.io/apimachinery v0.0.0-20191020214737-6c8691705fc5 // indirect
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
)

replace github.com/keptn/go-utils => ../../go-utils
