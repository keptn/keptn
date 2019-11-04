# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/keptn/keptn/platform-support/openshift-route-service
COPY . .

ARG DEP_VERSION=0.5.3
RUN curl -L -s https://github.com/golang/dep/releases/download/v$DEP_VERSION/dep-linux-amd64 -o ./dep && \
  chmod +x ./dep && \
  ./dep ensure

# Build the command inside the container.
# (You may fetch or manage dependencies here, 
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o openshift-route-service

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM frolvlad/alpine-glibc:latest
RUN apk add --no-cache ca-certificates

ARG OC_VERSION=3.11.0
RUN wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz && \ 
  tar xzvf openshift*tar.gz && \
  cp openshift-origin-client-tools-*/oc /bin/oc && \
  cp openshift-origin-client-tools-*/oc /usr/local/bin && \
  rm -rf openshift-origin-client-tools-* && \
  rm -rf openshift*tar.gz

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/platform-support/openshift-route-service/openshift-route-service /openshift-route-service
ADD MANIFEST /

# Run the web service on container startup.
CMD ["sh", "-c", "cat MANIFEST && /openshift-route-service"]