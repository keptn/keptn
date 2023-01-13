# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.19.5-alpine3.16 as builder

WORKDIR /go/src/github.com/keptn/keptn/cli

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev git

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image.
COPY . .

ARG debugBuild
# set buildflags for debug build
RUN if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -o keptn