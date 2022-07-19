# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.18.4-alpine3.16 as builder-base

ARG debugBuild

WORKDIR /go/src/github.com/keptn/keptn/distributor

# Force the go compiler to use modules
ENV GO111MODULE=on
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
CMD gotestsum --no-color=false -- -coverprofile=coverage.txt -covermode=atomic -v ./... && mv ./coverage.txt /shared/coverage.txt

FROM builder-base as builder

ARG debugBuild
ARG buildTime=unknown
ARG gitSha=unknown

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-linkmode=external -X main.gitCommit=$gitSha -X main.buildTime=$buildTime" $BUILDFLAGS -v -o distributor ./cmd/

FROM alpine:3.16 as production
ARG version=develop

LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Distributor" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}" \
    io.snyk.containers.image.dockerfile="distributor/Dockerfile"

RUN apk add --no-cache ca-certificates libc6-compat
RUN apk del busybox

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/distributor/distributor /distributor

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/distributor"]
