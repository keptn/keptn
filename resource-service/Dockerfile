# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.19.4-alpine3.17 as builder

# install additional dependencies
RUN apk update
RUN apk upgrade
RUN apk add --no-cache gcc libc-dev git cmake make openssl
RUN apk add --no-cache libgit2-dev=1.5.0-r2

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

# Pre-Req for gin-gonic/go-swagger: Install swag cli to generate swagger.yaml and swagger.json
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.0

WORKDIR /go/src/github.com/keptn/keptn/resource-service

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

# generate swagger docs
RUN GOOS=linux swag init

# replace github.com/alecthomas/template; with text/template in docs/docs.go
RUN sed -i "s|github.com/alecthomas/template|text/template|g" docs/docs.go

ARG version=develop
# replace version "develop" with actual version
RUN sed -i "s/version: develop/version: ${version}/g" docs/swagger.yaml

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# Build the command inside the container. 
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN GOOS=linux LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/usr/local/lib/ go build -ldflags '-linkmode=external' -gcflags="${SKAFFOLD_GO_GCFLAGS}" -v main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3.17 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Resource Service" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

ENV env=production

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat cmake make openssl

# Install git
#RUN apk update
#RUN apk upgrade
RUN apk add --no-cache build-base git pkgconf
RUN apk add --no-cache libgit2-dev=1.5.0-r2

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/main /resource-service
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger.yaml
# Replace contents for api proxy
RUN sed -i "s|basePath: /v1|basePath: /api/resource-service/v1 |g" /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
ENV GIN_MODE=release
CMD ["/resource-service"]
