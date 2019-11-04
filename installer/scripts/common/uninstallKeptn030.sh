# Clean up keptn namespace
kubectl delete services,deployments,pods,secrets --all -n keptn --ignore-not-found
kubectl delete namespace keptn --ignore-not-found

# Clean up knative components
kubectl delete -f https://github.com/knative/serving/releases/download/v0.4.0/serving.yaml --ignore-not-found
kubectl delete -f https://github.com/knative/build/releases/download/v0.4.0/build.yaml --ignore-not-found
kubectl delete -f https://github.com/knative/eventing/releases/download/v0.4.0/release.yaml --ignore-not-found
kubectl delete -f https://github.com/knative/eventing-sources/releases/download/v0.4.0/release.yaml --ignore-not-found
kubectl delete -f https://github.com/knative/serving/releases/download/v0.4.0/monitoring.yaml --ignore-not-found
kubectl delete -f https://raw.githubusercontent.com/knative/serving/v0.4.0/third_party/config/build/clusterrole.yaml --ignore-not-found

# Clean up istio namespace
kubectl delete -f ../../manifests/istio/istio-knative.yaml --ignore-not-found
kubectl delete -f ../../manifests/istio/istio-crds-knative.yaml --ignore-not-found
kubectl delete services,deployments,pods --all -n istio-system --ignore-not-found
kubectl delete namespace istio-system --ignore-not-found

# Delete tiller
kubectl delete -f ../../manifests/tiller/tiller.yaml --ignore-not-found

# Delete cluster role bindings
kubectl delete clusterrolebinding keptn-cluster-admin-binding --ignore-not-found