FROM --platform=$BUILDPLATFORM golang:1.18.7 as builder

ARG version=develop
ARG debugBuild

WORKDIR /go/src/github.com/keptn/keptn/mongodb-datastore

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apt-get install -y gcc libc-dev

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

# Target OS and architecture to build against (set automatically)
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN if [ "$GOARCH" == *'amd64'* ]; then export LINKERFLAGS="-ldflags '-linkmode=external'"; fi

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

RUN sed -i "s/version: develop/version: ${version}/g" ./swagger.yaml
RUN sed -i "s|basePath: /|basePath: /api/mongodb-datastore |g" ./swagger.yaml
RUN sed -i '/paths:/i securityDefinitions:\n  key:\n    type: apiKey\n    in: header\n    name: x-token\nsecurity:\n  - key: []\n' ./swagger.yaml

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN cd cmd/mongodb-datastore-server && go build $LINKERFLAGS $BUILDFLAGS -v -o mongodb-datastore

FROM alpine:3.16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn MongoDB Datastore" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/mongodb-datastore/cmd/mongodb-datastore-server/mongodb-datastore /mongodb-datastore
COPY --from=builder /go/src/github.com/keptn/keptn/mongodb-datastore/swagger-ui /swagger-ui
COPY --from=builder /go/src/github.com/keptn/keptn/mongodb-datastore/swagger.yaml /swagger-ui/swagger-original.yaml
COPY --from=builder /go/src/github.com/keptn/keptn/mongodb-datastore/swagger.yaml /swagger-ui/swagger.yaml

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/mongodb-datastore", "--port=8080", "--host=0.0.0.0"]
