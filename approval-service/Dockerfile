# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.18.4-alpine3.16 as builder-base

WORKDIR /go/src/github.com/keptn/keptn/approval-service

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev git

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

FROM builder-base as builder-test
ENV GOTESTSUM_FORMAT=testname

RUN go install gotest.tools/gotestsum@v1.7.0
CMD gotestsum --no-color=false -- -race -coverprofile=coverage.txt -covermode=atomic -v ./... && mv ./coverage.txt /shared/coverage.txt

FROM builder-base as builder
ARG version=develop
ARG debugBuild

# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v -o approval-service

FROM alpine:3.16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Approval Service" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}" \
    io.snyk.containers.image.dockerfile="approval-service/Dockerfile"

# we need to install ca-certificates and libc6-compat for go programs to work properly
RUN apk add --no-cache ca-certificates libc6-compat
RUN apk upgrade busybox=1.35.0-r17

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/approval-service/approval-service /approval-service

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

RUN adduser -D nonroot -u 65532
USER nonroot

# Run the web service on container startup.
CMD ["/approval-service"]
