# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.12 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/keptn/keptn/configuration-service
COPY . .

ARG DEP_VERSION=0.5.3
RUN curl -L -s https://github.com/golang/dep/releases/download/v$DEP_VERSION/dep-linux-amd64 -o ./dep && \
  chmod +x ./dep && \
  ./dep ensure

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