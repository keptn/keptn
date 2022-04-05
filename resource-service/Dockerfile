FROM --platform=$BUILDPLATFORM golang:1.18.7-alpine3.16 as builder

ARG version=develop
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

RUN apk add --no-cache gcc libc-dev

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

# Pre-Req for gin-gonic/go-swagger: Install swag cli to generate swagger.yaml and swagger.json
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.0

WORKDIR /go/src/github.com/keptn/keptn/resource-service

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM builder-base as builder
ARG version=develop

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# Target OS and architecture to build against (set automatically)
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN if [ "$GOARCH" == *'amd64'* ]; then export LINKERFLAGS="-ldflags '-linkmode=external'"; fi

# generate swagger docs
RUN swag init

# replace github.com/alecthomas/template; with text/template in docs/docs.go
RUN sed -i "s|github.com/alecthomas/template|text/template|g" docs/docs.go
# replace version "develop" with actual version
RUN sed -i "s/version: develop/version: ${version}/g" docs/swagger.yaml
RUN sed -i "s|basePath: /v1|basePath: /api/resource-service/v1 |g" docs/swagger.yaml

RUN go build $LINKERFLAGS -gcflags="${SKAFFOLD_GO_GCFLAGS}" -v main.go

FROM alpine:3.16 as production
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
RUN apk add --no-cache ca-certificates libc6-compat

# Install git
RUN apk --update add --no-cache git

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/main /resource-service
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/resource-service/docs/swagger.yaml /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
ENV GIN_MODE=release
CMD ["/resource-service"]
