# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.18.4-alpine3.16 as builder-base

WORKDIR /go/src/github.com/keptn/keptn/api

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev git

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

FROM builder-base as builder-test
ENV GOTESTSUM_FORMAT=testname

RUN go install gotest.tools/gotestsum@v1.7.0
CMD gotestsum --no-color=false -- -coverprofile=coverage.txt -covermode=atomic -v ./handlers/... ./importer/... ./utils/... && mv ./coverage.txt /shared/coverage.txt

FROM builder-base as builder
ARG version=develop
ARG debugBuild

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

RUN sed -i "s/version: develop/version: ${version}/g" /go/src/github.com/keptn/keptn/api/swagger.yaml

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v cmd/api-server/main.go


FROM alpine:3.16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn API" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}" \
    io.snyk.containers.image.dockerfile="api/Dockerfile"

RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/api/main /api
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger.yaml /swagger-ui/swagger.yaml

RUN sed -i "s|basePath: /v1|basePath: /api/v1 |g" /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/api", "--host=0.0.0.0", "--port=8080"]
