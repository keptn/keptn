FROM alpine:3.9

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

RUN apk update && apk upgrade && apk add --no-cache mongodb mongodb-tools bash curl

# Copy core and install
COPY scripts .
RUN chmod +x /*.sh

# Start the app
CMD ["bash", "upgradeKeptn.sh"]
