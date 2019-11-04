# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.12 as builder
ARG version=develop

WORKDIR /go/src/github.com/keptn/keptn/api

# Force the go compiler to use modules 
ENV GO111MODULE=on

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

ENV REPLACE="version: ${version}"
RUN sed -i "s/version: develop/${REPLACE}/g" /go/src/github.com/keptn/keptn/api/swagger.yaml


# Build the command inside the container. 
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v cmd-server/main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/api/main /api
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger.yaml /swagger-ui/swagger.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/api/project_model.yaml /swagger-ui/project_model.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/api/resource_model.yaml /swagger-ui/resource_model.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/api/response_model.yaml /swagger-ui/response_model.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/api/service_model.yaml /swagger-ui/service_model.yaml

ADD MANIFEST /

# Run the web service on container startup.
CMD ["sh", "-c", "cat MANIFEST && ./api --host=0.0.0.0 --port=8080"]