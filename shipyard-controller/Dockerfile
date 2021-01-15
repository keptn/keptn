# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.13.7-alpine as builder
ARG version=develop

WORKDIR /go/src/github.com/keptn/keptn/shipyard-controller

# Force the go compiler to use modules 
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev git

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

ARG debugBuild

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# Copy local code to the container image.
COPY . .

ENV REPLACE="version: ${version}"

# Build the command inside the container. 
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3.11

ENV env=production

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat

# Install git
RUN apk --update add --no-cache git=2.24.3-r0

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/shipyard-controller/main /shipyard-controller
COPY --from=builder /go/src/github.com/keptn/keptn/shipyard-controller/swagger-ui /swagger-ui

COPY --from=builder /go/src/github.com/keptn/keptn/shipyard-controller/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/shipyard-controller/docs/swagger.yaml /swagger-ui/swagger.yaml
RUN sed -i "s|basePath: /v1|basePath: /api/shipyard-controller/v1 |g" /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the travis-ci build)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

# Run the web service on container startup.
ENV GIN_MODE=release
CMD ["/shipyard-controller"]
#CMD ["/shipyard-controller", "--host=0.0.0.0", "--port=8080"]
