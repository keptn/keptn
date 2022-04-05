FROM --platform=$BUILDPLATFORM golang:1.18.7-alpine3.16 as builder

ARG version=develop
ARG debugBuild

WORKDIR /go/src/github.com/keptn/keptn/api

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

# Target OS and architecture to build against (set automatically)
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN if [ "$GOARCH" == *'amd64'* ]; then export LINKERFLAGS="-ldflags '-linkmode=external'"; fi

# set buildflags for debug build
# set the gcflags for all packages, disable optimizations and inlining
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

RUN sed -i "s/version: develop/version: ${version}/g" swagger.yaml
RUN sed -i "s|basePath: /v1|basePath: /api/v1 |g" swagger.yaml

RUN go build $LINKERFLAGS $BUILDFLAGS -v cmd/api-server/main.go


FROM alpine:3.16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn API" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/api/main /api
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/api/swagger.yaml /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/api", "--host=0.0.0.0", "--port=8080"]
