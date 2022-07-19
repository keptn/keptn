FROM cypress/browsers:node14.17.6-chrome100-ff98 as builder-test-ui
WORKDIR /usr/src/app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY ./ ./
CMD ./cypress/run-tests.sh

FROM node:14-alpine3.15 as builder-test-base
WORKDIR /usr/src/app
COPY package.json yarn.lock /usr/src/app/
RUN yarn install --frozen-lockfile
COPY . /usr/src/app

FROM builder-test-base as builder-code-style
CMD yarn prettier:check && yarn lint:check

FROM builder-test-base as builder-test-unit
CMD yarn test && mv ./coverage/ /shared/coverage

FROM node:14-alpine3.15 as builder-base

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app

# copy package.json only to install angular dependencies first
COPY package.json yarn.lock /usr/src/app/
RUN yarn install --frozen-lockfile
COPY . /usr/src/app

FROM builder-base as bridge-builder
RUN yarn build

FROM node:14-alpine3.15 as bridge-server-builder-base

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app/server

# copy package.json only to install dependencies first
COPY ./server/package.json ./server/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY ./server ./
COPY ./shared /usr/src/app/shared

FROM bridge-server-builder-base as bridge-server-test
CMD yarn test && \
  mv ./coverage/ /shared/coverage

FROM bridge-server-builder-base as bridge-server-builder
RUN yarn build && \
  yarn install --frozen-lockfile --production

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM gcr.io/distroless/nodejs:16 as production
ARG version=develop
LABEL org.opencontainers.image.source="https://github.com/keptn/keptn" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Bridge" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.version="${version}"

ENV VERSION="${version}"
ENV NODE_ENV "production"
ENV API_URL "http://api-gateway-nginx.keptn.svc.cluster.local"
ENV API_TOKEN ""

WORKDIR /usr/src/app

# copy angular output from angularBuilder
COPY --from=bridge-builder --chown=65532:65532 /usr/src/app/dist /usr/src/app/dist
COPY --from=bridge-server-builder --chown=65532:65532 /usr/src/app/server/dist /usr/src/app/server/dist
COPY --from=bridge-server-builder --chown=65532:65532 /usr/src/app/server/package.json /usr/src/app/server/
COPY --from=bridge-server-builder --chown=65532:65532 /usr/src/app/server/node_modules /usr/src/app/server/node_modules

# Set user
USER 65532:65532

EXPOSE 3000
CMD ["npm", "start", "--prefix", "./server"]
