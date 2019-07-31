# Clean up keptn namespace
kubectl delete services,deployments,pods,secrets --all -n keptn --ignore-not-found
kubectl delete namespace keptn --ignore-not-found

# Clean up monitoring namespace
kubectl delete daemonsets,deployments,services,configmaps --all -n knative-monitoring --ignore-not-found
kubectl delete namespace knative-monitoring --ignore-not-found

# Clean up istio namespace
kubectl delete -f ../../manifests/istio/crd-10.yaml --ignore-not-found
kubectl delete -f ../../manifests/istio/crd-11.yaml --ignore-not-found
kubectl delete -f ../../manifests/istio/crd-12.yaml --ignore-not-found
kubectl delete -f ../../manifests/istio/crd-certmanager-10.yaml --ignore-not-found
kubectl delete -f ../../manifests/istio/crd-certmanager-11.yaml --ignore-not-found
# Delete tiller
kubectl delete -f ../../manifests/tiller/tiller.yaml --ignore-not-found