FROM node:14-alpine3.13 as angularBuilder

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app

# copy package.json only to install angular dependencies first
COPY package.json package-lock.json /usr/src/app/
RUN npm install

COPY . /usr/src/app

RUN npm run build

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM node:14-alpine3.13

EXPOSE 3000

ENV NODE_ENV "production"
ARG version=develop
ENV VERSION="${version}"
ENV API_URL "http://api-gateway-nginx.keptn.svc.cluster.local"
ENV API_TOKEN ""

# copy angular output from angularBuilder
COPY --from=angularBuilder /usr/src/app/dist /usr/src/app/dist

# copy source
COPY . /usr/src/app

# install additional npm dependencies for express server
WORKDIR /usr/src/app/server
RUN npm install
WORKDIR /usr/src/app

RUN addgroup mygroup --gid 65532 && adduser -D -G mygroup myuser --uid 65532 && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Set user
USER myuser

# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the travis-ci build)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/local/bin/npm", "start"]
