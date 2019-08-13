FROM node:10-alpine
ENV DATASTORE "http://mongodb-datastore.keptn-datastore.svc.cluster.local:8080/"
EXPOSE 8080

RUN npm install -g http-server
RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Prepare app directory
WORKDIR /usr/src/app

COPY . /usr/src/app
RUN npm install
RUN npm run build
ENV NODE_ENV "production"

USER myuser

CMD ["sh", "-c", "cat MANIFEST && /usr/local/bin/npm start"]