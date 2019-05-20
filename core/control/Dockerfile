FROM node:10-alpine
ENV NODE_ENV "production"
ENV PORT 8080
EXPOSE 8080
RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app

# Prepare app directory
WORKDIR /usr/src/app

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

# Start the app
CMD ["sh", "-c", "cat MANIFEST && /usr/local/bin/npm start"]