# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.16.2-alpine as builder
ARG version=develop

WORKDIR /go/src/github.com/keptn/keptn/api

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
RUN sed -i "s/version: develop/${REPLACE}/g" /go/src/github.com/keptn/keptn/api/swagger.yaml


# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v cmd/api-server/main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/api/main /api
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger.yaml /swagger-ui/swagger.yaml

RUN sed -i "s|basePath: /v1|basePath: /api/v1 |g" /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the travis-ci build)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/api", "--host=0.0.0.0", "--port=8080"]
