FROM node:14-alpine3.14 as bridgeBuilder

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app

# copy package.json only to install angular dependencies first
COPY package.json package-lock.json /usr/src/app/
RUN npm ci
COPY . /usr/src/app
RUN npm run build

FROM node:14-alpine3.14 as bridgeServerBuilder

ARG version=develop
ENV VERSION="${version}"




# Prepare app directory
WORKDIR /usr/src/app/server

# copy package.json only to install dependencies first
COPY ./server/package.json ./server/package-lock.json ./
RUN npm i
COPY ./server ./
COPY ./shared /usr/src/app/shared
RUN npm run build && \
  npm ci --production

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM node:14-alpine3.14 as bridgeProduction

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

CMD ["/usr/local/bin/npm", "start", "--prefix", "./server"]
