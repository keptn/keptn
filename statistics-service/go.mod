module github.com/keptn/keptn/statistics-service

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-test/deep v1.0.8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.13.1-0.20220318125157-fe974e59cc65
	github.com/mitchellh/copystructure v1.2.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.9
	github.com/tryvium-travels/memongo v0.4.0
	go.mongodb.org/mongo-driver v1.8.4
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)
