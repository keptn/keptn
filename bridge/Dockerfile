FROM node:14-alpine3.14 as bridgeBuilder

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app

# copy package.json only to install angular dependencies first
COPY ./package.json ./yarn.lock ./.yarnrc /usr/src/app/
COPY ./npm-modules-cache/ npm-modules-cache/
RUN yarn install --frozen-lockfile
COPY . /usr/src/app
RUN yarn build

FROM node:14-alpine3.14 as bridgeServerBuilder

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app/server

# copy package.json only to install dependencies first
COPY ./server/package.json ./server/yarn.lock ./.yarnrc ./
COPY ./server/npm-modules-cache/ npm-modules-cache/
RUN yarn install --frozen-lockfile
COPY ./server ./
COPY ./shared /usr/src/app/shared
RUN yarn build && \
  yarn install --frozen-lockfile --production

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM node:14-alpine3.14 as bridgeProduction
LABEL org.opencontainers.image.source = "https://github.com/keptn/keptn" \
    org.opencontainers.image.url = "https://keptn.sh" \
    org.opencontainers.image.title="Keptn Bridge" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.documentation="https://keptn.sh/docs/" \
    org.opencontainers.image.licenses="Apache-2.0"


EXPOSE 3000

ENV NODE_ENV "production"
WORKDIR /usr/src/app
ARG version=develop
ENV VERSION="${version}"
ENV API_URL "http://api-gateway-nginx.keptn.svc.cluster.local"
ENV API_TOKEN ""

# copy angular output from angularBuilder
COPY --from=bridgeBuilder /usr/src/app/dist /usr/src/app/dist
COPY --from=bridgeServerBuilder /usr/src/app/server/dist /usr/src/app/server/dist
COPY --from=bridgeServerBuilder /usr/src/app/server/package.json /usr/src/app/server/
COPY --from=bridgeServerBuilder /usr/src/app/server/node_modules /usr/src/app/server/node_modules

# Copy server static files
COPY ./server/views /usr/src/app/server/views

RUN addgroup mygroup --gid 65532 && adduser -D -G mygroup myuser --uid 65532 && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Set user
USER myuser

CMD ["npm", "start", "--prefix", "./server"]
