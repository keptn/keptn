FROM alpine:3.9

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG HELM_VERSION=2.12.3
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64/helm /bin/helm && \
  rm -rf linux-amd64

ARG HELM3_VERSION=3.1.2
RUN wget https://get.helm.sh/helm-v$HELM3_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM3_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64/helm /bin/helm3 && \
  rm -rf linux-amd64

RUN apk update && apk upgrade && apk add --no-cache git bash curl

# Copy core and install
COPY scripts .
RUN chmod +x /*.sh

# Start the app
CMD ["bash", "upgradeKeptn.sh"]
