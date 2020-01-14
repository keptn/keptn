FROM node:10-alpine
ENV DATASTORE "http://mongodb-datastore.keptn-datastore.svc.cluster.local:8080/"
EXPOSE 8080

ARG version=develop

RUN npm install -g http-server
RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Prepare app directory
WORKDIR /usr/src/app

COPY package.json package-lock.json /usr/src/app/
RUN npm install

COPY . /usr/src/app

RUN npm run build
ENV NODE_ENV "production"

# Set user
USER myuser


# KEEP THE FOLLOWING LINES COMMENTED OUT!!! (they will be included within the travis-ci build)
#travis-uncomment ADD MANIFEST /
#travis-uncomment COPY entrypoint.sh /
#travis-uncomment ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/local/bin/npm", "start"]
