FROM --platform=$BUILDPLATFORM golang:1.18.7-alpine3.16 as builder

ARG debugBuild
ARG buildTime=unknown
ARG gitSha=unknown

WORKDIR /go/src/github.com/keptn/keptn/distributor

ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

# Target OS and architecture to build against (set automatically)
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN if [ "$GOARCH" == *'amd64'* ]; then export LINKERFLAGS="-linkmode=external"; fi

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN go build -ldflags "$LINKERFLAGS -X main.gitCommit=$gitSha -X main.buildTime=$buildTime" $BUILDFLAGS -v -o distributor ./cmd/

FROM alpine:3.16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Distributor" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/distributor/distributor /distributor

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/distributor"]
