FROM golang:1.18.4-alpine3.16 as builder-base

RUN apk add --no-cache gcc libc-dev

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

# Pre-Req for gin-gonic/go-swagger: Install swag cli to generate swagger.yaml and swagger.json
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.0

WORKDIR /go/src/github.com/keptn/keptn/secret-service

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

FROM builder-base as builder-test
ENV GOTESTSUM_FORMAT=testname

RUN go install gotest.tools/gotestsum@v1.7.0
CMD gotestsum --no-color=false -- -coverprofile=coverage.txt -covermode=atomic -v ./... && mv ./coverage.txt /shared/coverage.txt

FROM builder-base as builder
ARG version=develop
ENV CGO_ENABLED=0

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# generate swagger docs
RUN GOOS=linux swag init --parseDependency

# replace github.com/alecthomas/template; with text/template in docs/docs.go
RUN sed -i "s|github.com/alecthomas/template|text/template|g" docs/docs.go
# replace version "develop" with actual version
RUN sed -i "s/version: develop/version: ${version}/g" docs/swagger.yaml

RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -v main.go
RUN sed -i "s|basePath: /v1|basePath: /api/secrets/v1 |g" docs/swagger.yaml

FROM gcr.io/distroless/static-debian11 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Secret Service" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

ENV env=production

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/secret-service/main /secret-service
COPY --from=builder /go/src/github.com/keptn/keptn/secret-service/swagger-ui /swagger-ui

COPY --from=builder /go/src/github.com/keptn/keptn/secret-service/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/secret-service/docs/swagger.yaml /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

# Run the web service on container startup.
ENV GIN_MODE=release

ADD scopes.yaml /data/scopes.yaml

USER nonroot

CMD ["/secret-service"]
