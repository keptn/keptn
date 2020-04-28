FROM alpine:3.9

ARG KUBE_VERSION=1.14.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG OC_VERSION=3.11.0
RUN wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz -O oc.tar.gz && \
  tar xzvf oc.tar.gz && \
  mv openshift-origin-client-tools-*/oc /bin/oc && \
  rm -rf openshift-origin-client-tools-* && \
  rm oc.tar.gz

RUN apk update && apk upgrade && apk add --no-cache bash curl

# Copy core and install
COPY scripts .
RUN chmod +x /*.sh

# Start the app
CMD ["bash", "upgradeKeptn.sh"]
