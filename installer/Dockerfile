FROM ubuntu:18.04
ARG version=develop

ENV KEPTN_INSTALL_ENV "cluster"

# install basic tools
RUN apt-get update \
  && apt-get install -y curl wget jq git \
  && rm -rf /var/lib/apt/lists/*

## Install go
# RUN mkdir -p /goroot && \
#   curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

## Set environment variables.
# ENV GOROOT /goroot
# ENV GOPATH /gopath
# ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

# Install tools:

ARG YQ_VERSION=2.3.0
RUN wget https://github.com/mikefarah/yq/releases/download/$YQ_VERSION/yq_linux_amd64 && \
  chmod +x yq_linux_amd64 && \
  mv yq_linux_amd64 /bin/yq
RUN yq --version

ARG HELM_VERSION=3.1.2
RUN wget https://get.helm.sh/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64/helm /bin/helm && \
  rm -rf linux-amd64

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG OC_VERSION=3.11.0
RUN wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz -O oc.tar.gz && \ 
  tar xzvf oc.tar.gz && \
  mv openshift-origin-client-tools-*/oc /bin/oc && \
  rm -rf openshift-origin-client-tools-* && \
  rm oc.tar.gz

# Copy core and install
WORKDIR /usr/keptn
COPY ./manifests manifests
COPY ./scripts scripts
COPY MANIFEST scripts

RUN cd ./scripts && ls -lsa
WORKDIR /usr/keptn/scripts

# Start the app
CMD ["sh", "-c", "cat MANIFEST && ./installKeptn.sh"]