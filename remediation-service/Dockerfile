FROM golang:1.18.4-alpine3.16 as builder-base

WORKDIR /go/src/github.com/keptn/keptn/remediation-service

ENV BUILDFLAGS=""
ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache gcc libc-dev

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
ENV CGO_ENABLED=0

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS

# Build the command inside the container.
# (You may fetch or manage dependencies here, either manually or with a tool like "godep".)
RUN GOOS=linux go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -v -o remediation-service

FROM gcr.io/distroless/static-debian11 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Remediation Service" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/keptn/keptn/remediation-service/remediation-service /remediation-service

EXPOSE 8080

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all

USER nonroot

# Run the web service on container startup.
CMD ["/remediation-service"]
