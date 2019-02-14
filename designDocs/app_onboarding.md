# Application onboarding
This tool will help you to onboard an existing app, defined as Helm charts, into your keptn environment.

##### Table of Contents
 * [Prerequisites](#step-zero)
 * [Helm chart naming conventions](#step-one)
 * [Onboarding process](#step-two)

## Prerequisites <a id="step-zero"></a>

Keptn onboarding requires following items:

1. Github Orgnaization containing a *code repo* for each microservices to onboard and one *config repo* (Helm chart)
    * Code repo: Contains a pipeline definition to build the container images and push them to a registry.
    * Config repo (suffixed with ‘-config’): Contains a folder called `helm-chart` that contains the Helm chart for the application (composition of all microservices)
    * Example: [sample-app-code repo](https://github.com/keptn/examples/onboard-sample-app/sample-app-code) and [sample-app-config repo](https://github.com/keptn/examples/onboard-sample-app/sample-app-config)

2. Allow Tiller access to namespaces, due to a known [issue](https://github.com/fnproject/fn-helm/issues/21).

    ```console
    $ kubectl create serviceaccount --namespace kube-system tiller
    $ kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    $ kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
    ```

## Helm chart naming conventions <a id="step-one"></a>

For each microservice, there is an entry in values.yaml, written in camel case, i.e., if the name of the repo is *front-end*, the entry in the yaml file is called ‘frontEnd’, e.g.:

  ```yaml
  carts:
    ...
  frontEnd:	
    ...
  ```

For each of those microservices, there is a deployment and a service yaml in the `helm-chart/templates` directory:
    ```console
    $ pwd
    ~/app-config/helm-chart/templates
    $ ls
    carts-deployment.yaml
    carts-service.yaml
    frontEnd-deployment.yaml
    frontEnd-service.yaml
    ```

* The field **metadata.name** in the `*-deplyoment.yaml` and `*-service.yaml` file must be set to **{{ .Chart.Name }}-\<name-of-the-microservice>**. **Note:** The name has to be in has to be hyphen case, due to k8s naming restrictions.
* To support automation of Istio setup, the spec section of the `*-deployment.yaml` file has to include the field **metadata.labels.deployment** and **spec.template.metadata.labels.deployment** with its value set to the hyphenated microservice name as well.
* The properties **metadata.labels.app** and **spec.template.metadata.labels.app** have to be set to **{{ .Chart.Name }}-selector-\<name-of-the-microservice>**.

### Example of deployment.yaml
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: "{{ .Chart.Name }}-front-end"
  labels:
   app: "{{ .Chart.Name }}-selector-front-end" 
   deployment: front-end
   chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
   release: "{{ .Release.Name }}"
   heritage: "{{ .Release.Service }}"
spec:
  replicas: {{ .Values.sampleAppCode.replicaCount }}
  template:
    metadata:
      labels:
        app: "{{ .Chart.Name }}-selector-front-end" 
        deployment: front-end
...
``` 

### Example of service.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: "{{ .Chart.Name }}-front-end"
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
    app: "{{ .Chart.Name }}-selector-front-end"
```

## Onboarding process <a id="step-two"></a>

The following diagram illustrates the setup steps executed during the onboarding process:

![keptn-infra](./res/keptn-onboard-infra.png)

