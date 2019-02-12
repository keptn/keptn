# keptn-onboard
This tool will help you to onboard your existing apps, defined as Helm charts, into your keptn environment.

## Prerequisites
1. Prerequisites:
    1. Github Orgnaization containing source repos and config repo (Helm chart)
    2. Source repos: contain a Pipeline Definition to build the container images and push them to a registry
    3. Config repo (suffixed with ‘-config’): contains a folder called ‘helm-chart’ that contains the Helm chart for the application (composition of all Microservices)
    4. Example: https://github.com/keptn-onboard/sample-app-code and https://github.com/keptn-onboard/sample-app 
2. Allow Tiller access to namespaces (https://github.com/fnproject/fn-helm/issues/21)

```
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
```

## Helm Chart Naming Conventions:

For each micro service, there is an entry in values.yaml, written in camel case (if the name of the repo is ‘front-end’, the entry in the yams file is called ‘frontEnd’) e.g.:

```
carts:
	…
frontEnd:	
	…
```

For each of those micro services, there is a deployment and a service yaml in the helm-chart/templates directory:
```
carts-deployment.yaml

carts-service.yaml
```

- The field **metadata.name** in each of those yaml files must be set to **{{ .Chart.Name }}-\<name-of-the-microservice>**. NOTE that in this case, the name has to be in has to be hyphen case, due to k8s naming restrictions.
- Further, to support automation of Istio setup, the spec has to include the field **metadata.labels.deployment** and **spec.template.metadata.labels.deployment** with its value set to the hyphenated micro service name as well.
- The properties **metadata.labels.app** and **spec.template.metadata.labels.app** have to be set to **{{ .Chart.Name }}-selector-\<name-of-the-microservice>**.


### Deployment.yaml:
```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: "{{ .Chart.Name }}-front-end” # !!!
  labels:
   app: "{{ .Chart.Name }}-selector-front-end” # !!!
   deployment: front-end # !!!
   chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
   release: "{{ .Release.Name }}"
   heritage: "{{ .Release.Service }}"
spec:
  replicas: {{ .Values.sampleAppCode.replicaCount }}
  template:
    metadata:
      labels:
        app: "{{ .Chart.Name }}-selector-front-end” # !!!
        deployment: front-end # !!!

…
``` 

### Service.yaml

```
apiVersion: v1
kind: Service
metadata:
  name: "{{ .Chart.Name }}-front-end” # !!!
  labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
spec:
  type: {{ .Values.sampleAppCode.service.type }}
  ports:
  - port: {{ .Values.sampleAppCode.service.externalPort }}
    targetPort: {{ .Values.sampleAppCode.service.internalPort }}
    protocol: TCP
    name: {{ .Values.sampleAppCode.service.name }}
  selector:
    app: "{{ .Chart.Name }}-selector-front-end” # !!!
```
The following diagram illustrates the setup steps executed during the onboarding process:

![keptn-infra](./res/keptn-onboard-infra.png)