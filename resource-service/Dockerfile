FROM golang:1.18.4-alpine3.16 as builder-base

# install additional dependencies
RUN apk add --no-cache gcc libc-dev

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN go install github.com/swaggo/swag/cmd/swag@v1.7.0

WORKDIR /go/src/github.com/keptn/keptn/resource-service

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

FROM builder-base as builder-test
ENV GOTESTSUM_FORMAT=testname
RUN apk add --no-cache git && \
    go install gotest.tools/gotestsum@v1.7.0
CMD gotestsum --no-color=false -- -coverprofile=coverage.txt -covermode=atomic -v ./... && mv ./coverage.txt /shared/coverage.txt

FROM builder-base as builder
ARG version=develop
ENV CGO_ENABLED=0

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# generate swagger docs
RUN swag init

# replace github.com/alecthomas/template; with text/template in docs/docs.go
RUN sed -i "s|github.com/alecthomas/template|text/template|g" docs/docs.go
# replace version "develop" with actual version
RUN sed -i "s/version: develop/version: ${version}/g" docs/swagger.yaml

# Build the command inside the container. 
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -v main.go

# Replace contents for api proxy
RUN sed -i "s|basePath: /v1|basePath: /api/resource-service/v1 |g" docs/swagger.yaml

# Use a Docker multi-stage build to create a lean production image.
FROM gcr.io/distroless/static-debian11 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Resource Service" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

ENV env=production

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/main /resource-service
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

USER nonroot

# Run the web service on container startup.
ENV GIN_MODE=release
CMD ["/resource-service"]
