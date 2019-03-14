# Clean up cicd namespace
kubectl delete services,deployments,pods --all -n cicd
kubectl delete namespace cicd

# Clean up dynatrace namespace
kubectl delete services,deployments,pods --all -n dynatrace
kubectl delete namespace dynatrace

# Clean up tower namespace
kubectl delete services,deployments,pods --all -n tower
kubectl delete namespace tower

# Clean up keptn namespace
kubectl delete services,deployments,pods --all -n keptn
kubectl delete namespace keptn

# Clean up knative components
kubectl delete --filename https://github.com/knative/serving/releases/download/v0.4.0/serving.yaml
kubectl delete --filename https://github.com/knative/build/releases/download/v0.4.0/build.yaml
kubectl delete --filename https://github.com/knative/eventing/releases/download/v0.4.0/in-memory-channel.yaml
kubectl delete --filename https://github.com/knative/eventing/releases/download/v0.4.0/release.yaml
kubectl delete --filename https://github.com/knative/eventing-sources/releases/download/v0.4.0/release.yaml
kubectl delete --filename https://github.com/knative/serving/releases/download/v0.4.0/monitoring.yaml
kubectl delete --filename https://raw.githubusercontent.com/knative/serving/v0.4.0/third_party/config/build/clusterrole.yaml

# Clean up istio namespace
kubectl delete -f ../manifests/istio/istio-knative.yaml
kubectl delete -f ../manifests/istio/istio-crds-knative.yaml
kubectl delete services,deployments,pods --all -n istio-system
kubectl delete namespace istio-system

# Verification
kubectl delete clusterrolebindings.rbac.authorization.k8s.io dynatrace-cluster-admin-binding
kubectl delete clusterrolebindings.rbac.authorization.k8s.io jenkins-rbac
kubectl delete -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/master/deploy/kubernetes.yaml
