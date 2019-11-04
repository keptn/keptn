# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder
ARG version=develop

# Copy local code to the container image.
WORKDIR /go/src/github.com/keptn/keptn/helm-service

# Force the go compiler to use modules 
ENV GO111MODULE=on

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy local code to the container image. 
COPY . .

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o helm-service

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine
RUN apk add --no-cache ca-certificates

ARG HELM_VERSION=2.12.3
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  cp linux-amd64/helm /bin/helm

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/helm-service/helm-service /helm-service
ADD MANIFEST /

# Run the web service on container startup.
CMD ["sh", "-c", "cat MANIFEST && /helm-service"]