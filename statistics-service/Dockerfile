# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13.7-alpine as builder

RUN apk add --no-cache gcc libc-dev git

WORKDIR /go/src/github.com/keptn/keptn/statistics-service

ARG version=develop
ENV VERSION="${version}"

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

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

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v main.go

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.11
ENV ENV=production

# Install extra packages
# See https://github.com/gliderlabs/docker-alpine/issues/136#issuecomment-272703023

RUN apk update && apk upgrade \
	&& apk add ca-certificates libc6-compat \
	&& update-ca-certificates \
	&& rm -rf /var/cache/apk/*

ARG version=develop
ENV VERSION="${version}"

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/statistics-service/main /statistics-service

COPY --from=builder /go/src/github.com/keptn/keptn/statistics-service/swagger-ui /swagger-ui

COPY --from=builder /go/src/github.com/keptn/keptn/statistics-service/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/statistics-service/docs/swagger.yaml /swagger-ui/swagger.yaml

# Replace contents for api proxy
RUN sed -i "s|basePath: /v1|basePath: /api/statistics/v1 |g" /swagger-ui/swagger.yaml
EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the build pipeline)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

# Run the web service on container startup.
ENV GIN_MODE=release
CMD ["/statistics-service"]