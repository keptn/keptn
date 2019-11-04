# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.12 as builder

WORKDIR /go/src/github.com/keptn/keptn/configuration-service

# Force the go compiler to use modules 
ENV GO111MODULE=on

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download 

# Copy local code to the container image.
COPY . .

# Build the command inside the container. 
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v cmd/configuration-service-server/main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine

ENV env=production
RUN apk add --no-cache ca-certificates

# Install git
RUN apk --update add --no-cache git

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/configuration-service/main /configuration-service
COPY --from=builder /go/src/github.com/keptn/keptn/configuration-service/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/configuration-service/swagger.yaml /swagger-ui/swagger.yaml
ADD MANIFEST /

# Run the web service on container startup.
CMD ["sh", "-c", "cat MANIFEST && ./configuration-service --host=0.0.0.0 --port=8080"]