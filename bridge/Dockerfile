FROM node:10-alpine as builder

ARG version=develop
ENV VERSION="${version}"

# Prepare app directory
WORKDIR /usr/src/app

# copy package.json only to install dependencies first
COPY package.json package-lock.json /usr/src/app/
RUN npm install

COPY . /usr/src/app

RUN npm run build

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM node:10-alpine

EXPOSE 3000

ENV NODE_ENV "production"
ARG version=develop
ENV VERSION="${version}"
ENV API_URL "http://api-gateway-nginx.keptn.svc.cluster.local"
ENV API_TOKEN ""

RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Prepare app directory
WORKDIR /usr/src/app

# install additional dependencies for running the app
RUN npm install express@4.17.1 depd@1.1.2 cookie-parser@1.4.4 morgan@1.9.1 axios@0.19.2

# copy source
COPY . /usr/src/app
# copy dist from builder
COPY --from=builder /usr/src/app/dist /usr/src/app/dist

# Set user
USER myuser

# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the travis-ci build)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/local/bin/npm", "start"]
