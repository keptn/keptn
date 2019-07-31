FROM ubuntu:18.04

ENV KEPTN_INSTALL_ENV "cluster"

RUN apt-get update \
  && apt-get install -y curl \
  && apt-get install -y wget

## Install go
# RUN mkdir -p /goroot && \
#   curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

## Set environment variables.
# ENV GOROOT /goroot
# ENV GOPATH /gopath
# ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

# Install tools:
RUN apt-get install -y jq
RUN jq --version

RUN apt-get install -y git
RUN git --version

ARG YQ_VERSION=2.3.0
RUN wget https://github.com/mikefarah/yq/releases/download/$YQ_VERSION/yq_linux_amd64 && \
  chmod +x yq_linux_amd64 && \
  cp yq_linux_amd64 /bin/yq
RUN yq --version

ARG HELM_VERSION=2.12.3
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
  cp linux-amd64/helm /bin/helm

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG OC_VERSION=3.11.0
RUN wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz -o oc.tar.gz && \ 
  tar xzvf openshift*tar.gz && \
  cp openshift-origin-client-tools-*/oc /bin/oc

# Copy core and install
WORKDIR /usr/keptn
COPY ./manifests manifests
COPY ./scripts scripts
COPY MANIFEST scripts

RUN cd ./scripts && ls -lsa

WORKDIR /usr/keptn/scripts

# Start the app
CMD ["sh", "-c", "cat MANIFEST && ./installKeptn.sh"]
