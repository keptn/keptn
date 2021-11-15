module github.com/keptn/keptn/statistics-service

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-test/deep v1.0.8
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.10.1-0.20211108082931-55a5cc361a0a
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.4
	github.com/ugorji/go v1.1.8 // indirect
	go.mongodb.org/mongo-driver v1.7.4
	google.golang.org/protobuf v1.25.0 // indirect
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)
