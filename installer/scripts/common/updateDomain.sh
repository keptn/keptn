source ./utils.sh 
DOMAIN=$1

# Configure keptn virtual services
rm -f ../../manifests/keptn/gen/keptn-api-virtualservice.yaml
cat ../../manifests/keptn/keptn-api-virtualservice.yaml | \
  sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' > ../../manifests/keptn/gen/keptn-api-virtualservice.yaml

kubectl apply -f ../../manifests/keptn/gen/keptn-api-virtualservice.yaml
verify_kubectl $? "Deploying keptn api virtualservice failed."

kubectl delete secret -n istio-system istio-ingressgateway-certs

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$DOMAIN"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
#verify_kubectl $? "Creating secret for istio-ingressgateway-certs failed."

rm key.pem
rm certificate.pem

# Add config map in keptn namespace that contains the domain - this will be used by other services as well
rm -f ../../manifests/gen/keptn-domain-configmap.yaml
cat ../../manifests/keptn/keptn-domain-configmap.yaml | \
  sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' >> ../../manifests/gen/keptn-domain-configmap.yaml

kubectl apply -f ../../manifests/gen/keptn-domain-configmap.yaml
verify_kubectl $? "Creating configmap keptn-domain in keptn namespace failed."

# re-deploy github service

kubectl delete deployment github-service -n keptn
kubectl apply -f ../../manifests/keptn/uniform-services.yaml



KEPTN_ENDPOINT=https://api.keptn.$(kubectl get cm -n keptn keptn-domain -oyaml | yq - r data.app_domain)
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode)

print_info "keptn endpoint: $KEPTN_ENDPOINT"
print_info "keptn api-token: $KEPTN_API_TOKEN"

keptn auth --endpoint=$KEPTN_ENDPOINT --api-token=$KEPTN_API_TOKEN