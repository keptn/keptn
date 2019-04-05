FROM node:10-alpine
ENV NODE_ENV "production"
ENV PORT 8080
EXPOSE 8080
RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Prepare app directory
WORKDIR /usr/src/app

# IMPORTANT: the root of the working directory is the root of the GitHub repo, NOT the location of the Dockerfile
COPY package.json /usr/src/app/

RUN npm install -g tsc 
RUN npm install -g concurrently 
RUN npm install -g typescript
RUN npm install -g copyfiles

COPY . /usr/src/app

RUN ls -la /usr/src/app/*

RUN npm install
RUN npm run build-ts

USER myuser
CMD ["/usr/local/bin/npm", "start"]
