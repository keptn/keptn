# Clean up cicd namespace
kubectl delete services,deployments,pods --all -n cicd
kubectl delete namespace cicd

# Clean up dev namespace
kubectl delete services,deployments,pods --all -n dev
kubectl delete namespace dev

# Clean up staging namespace
kubectl delete services,deployments,pods --all -n staging
kubectl delete namespace staging

# Clean up production namespace
kubectl delete services,deployments,pods --all -n production
kubectl delete namespace production

# Clean up dynatrace namespace
kubectl delete services,deployments,pods --all -n dynatrace
kubectl delete namespace dynatrace

# Clean up tower namespace
kubectl delete services,deployments,pods --all -n tower
kubectl delete namespace tower

# Clean up istio namespace
kubectl delete services,deployments,pods --all -n istio-system
kubectl delete namespace istio-system

# Clean up istio components
kubectl delete -f ../repositories/k8s-deploy-production/istio/gateway.yml
kubectl delete -f ../repositories/k8s-deploy-production/istio/destination_rule.yml
kubectl delete -f ../repositories/k8s-deploy-production/istio/virtual_service.yml

# Remove Kubernetes cluster to Dynatrace connection
./deleteKubernetesEndpoint.sh

# Verification
kubectl delete clusterrolebindings.rbac.authorization.k8s.io dynatrace-cluster-admin-binding
kubectl delete clusterrolebindings.rbac.authorization.k8s.io jenkins-rbac
kubectl delete -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/master/deploy/kubernetes.yaml
