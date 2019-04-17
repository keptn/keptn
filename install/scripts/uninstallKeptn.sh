# Clean up dynatrace namespace
echo "Deleting Dynatrace OneAgent"
kubectl delete services,deployments,pods --all -n dynatrace --ignore-not-found
kubectl delete namespace dynatrace --ignore-not-found

# Clean up tower namespace
kubectl delete services,deployments,pods --all -n tower --ignore-not-found
kubectl delete namespace tower --ignore-not-found

# Clean up keptn namespace
kubectl delete services,deployments,pods,secrets --all -n keptn --ignore-not-found
kubectl delete namespace keptn --ignore-not-found

# Clean up knative components
kubectl delete --filename https://github.com/knative/serving/releases/download/v0.5.0/serving.yaml
kubectl delete --filename https://github.com/knative/build/releases/download/v0.5.0/build.yaml
kubectl delete --filename https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml
kubectl delete --filename https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml
kubectl delete --filename https://github.com/knative/serving/releases/download/v0.5.0/monitoring.yaml
kubectl delete --filename https://raw.githubusercontent.com/knative/serving/v0.5.0/third_party/config/build/clusterrole.yaml

# Clean up istio namespace
kubectl delete -f ../manifests/istio/istio-knative.yaml --ignore-not-found
kubectl delete -f ../manifests/istio/istio-crds-knative.yaml --ignore-not-found
kubectl delete services,deployments,pods --all -n istio-system --ignore-not-found
kubectl delete namespace istio-system --ignore-not-found

# Verification
kubectl delete clusterrolebinding dynatrace-cluster-admin-binding --ignore-not-found
kubectl delete clusterrolebindings.rbac.authorization.k8s.io jenkins-rbac --ignore-not-found
kubectl delete -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/master/deploy/kubernetes.yaml --ignore-not-found
