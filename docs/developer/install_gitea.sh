#!/bin/bash

helm repo add gitea-charts https://dl.gitea.io/charts/
helm install --values test/assets/gitea/values.yaml gitea gitea-charts/gitea -n keptn --wait --version v5.0.0

GITEA_ADMIN_USER=$(kubectl get pod -n keptn gitea-0 -ojsonpath='{@.spec.initContainers[?(@.name=="configure-gitea")].env[?(@.name=="GITEA_ADMIN_USERNAME")].value}')
GITEA_ADMIN_PASSWORD=$(kubectl get pod -n keptn gitea-0 -ojsonpath='{@.spec.initContainers[?(@.name=="configure-gitea")].env[?(@.name=="GITEA_ADMIN_PASSWORD")].value}')

curl -SL https://raw.githubusercontent.com/keptn/keptn/master/test/assets/squid/squid.conf --output squid.conf
curl -SL https://raw.githubusercontent.com/keptn/keptn/master/test/assets/squid/squid.yaml --output squid.yaml
kubectl create configmap squid.conf --from-file=squid.conf -n keptn
kubectl apply -f squid.yaml -n keptn

sleep 30

ssh-keygen -t rsa -C "gitea-http" -f "rsa_gitea" -P "myGiteaPassPhrase"
GITEA_PRIVATE_KEY=$(cat rsa_gitea)
GITEA_PUBLIC_KEY=$(cat rsa_gitea.pub)
GITEA_PRIVATE_KEY_PASSPHRASE=myGiteaPassPhrase

sleep 30

kubectl port-forward -n keptn svc/gitea-http 3000:3000 &
kubectl port-forward -n keptn svc/gitea-ssh 3001:22 &

sleep 30

curl -vkL --silent --user "${GITEA_ADMIN_USER}":"${GITEA_ADMIN_PASSWORD}" -X POST "http://localhost:3000/api/v1/users/${GITEA_ADMIN_USER}/tokens" -H "accept: application/json" -H "Content-Type: application/json; charset=utf-8" -d "{ \"name\": \"my-token\" }" -o gitea-token.txt
curl -vkL --silent --user "${GITEA_ADMIN_USER}":"${GITEA_ADMIN_PASSWORD}" -X POST "http://localhost:3000/api/v1/user/keys" -H "accept: application/json" -H "Content-Type: application/json; charset=utf-8" -d "{ \"key\": \"$GITEA_PUBLIC_KEY\",  \"title\": \"public-key-gitea\"}"

GITEA_TOKEN=$(jq -r .sha1 < gitea-token.txt)

kubectl create secret generic gitea-access -n keptn --from-literal=username="${GITEA_ADMIN_USER}" --from-literal=password="${GITEA_TOKEN}" --from-literal=private-key="${GITEA_PRIVATE_KEY}" --from-literal=private-key-pass="${GITEA_PRIVATE_KEY_PASSPHRASE}"

rm gitea-token.txt squid.conf squid.yaml
