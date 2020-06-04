FROM golang:1.13.7-alpine as builder
ARG version=develop

WORKDIR /go/src/github.com/keptn/keptn/

# Force the go compiler to use modules
ENV GO111MODULE=on
ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev git

# Copy local code to the container image.
COPY . .

RUN cd mongodb-upgrade
# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies

# Download dependencies
RUN cd mongodb-upgrade && go mod download


# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN cd mongodb-upgrade && GOOS=linux go build -ldflags '-linkmode=external' $BUILDFLAGS -v -o upgradecollections

FROM alpine:3.11

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG HELM_VERSION=2.12.3
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64/helm /bin/helm && \
  rm -rf linux-amd64

ARG HELM3_VERSION=3.1.2
RUN wget https://get.helm.sh/helm-v$HELM3_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM3_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64/helm /bin/helm3 && \
  rm -rf linux-amd64

RUN apk update && apk upgrade && apk add --no-cache git bash curl

# Copy core and install
COPY scripts .
RUN chmod +x /*.sh

COPY --from=builder /go/src/github.com/keptn/keptn/mongodb-upgrade/upgradecollections upgradecollections

RUN ls -la
# Start the app
CMD ["bash", "upgradeKeptn.sh"]
