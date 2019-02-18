# Install keptn
These introductions explain how to install keptn. 

##### Table of Contents
 * [Prerequisites](#step-zero)
 * [Install](#step-install)
 * [Future work: Install](#step-install)

## Install keptn core components <a id="step-install"></a>

1. Insert information in *./scripts/creds.json* by executing `defineCredentials.sh` in the `scripts` directory. This script will prompt you for all information needed to complete the setup and populate the file *scripts/creds.json* with them. 

    **Attention:** This file will hold your personal access-token and credentials needed for the automatic setup of keptn. Take care of not leaking this file! (As a first precaution we have added this file to the `.gitignore` file to prevent committing it to GitHub.)

    ```console
    $ ./defineCredentials.sh
    ```
    
1. Execute `setupInfrastructure.sh` in the `scripts` directory. This script deploys a container registry and sets up a multi-stage environment with a *dev*, *staging*, and *production* namespaces. 

    **Attention:** The script will create several new resources for you and will also update the files shipped with keptn. Take care of not leaking any files that will hold personal information. Including:
    - `manifests/dynatrace/cr.yml`
    - `manifests/istio/service_entries.yml`

    **Note:** The script will run for some time (~10-15 min), since it will wait for Jenkins to boot before setting credentials via the Jenkins REST API.

    ```console
    $ ./setupInfrastructure.sh
    ```

1. To verify the installation, run the following `kubectl` commands that show configured parts of the K8s cluster: 

    ```console
    $ kubectl get namespaces
    NAME           STATUS    AGE
    cicd           Active    10m
    default        Active    10m
    dev            Active    10m
    dynatrace      Active    10m
    istio-system   Active    10m
    kube-public    Active    10m
    kube-system    Active    10m
    production     Active    10m
    staging        Active    10m
    ```

    ```console
    $ kubectl get pods -n istio-system
    ```

    ```console
    $ kubectl get pods -n cicd
    NAME                                 READY     STATUS    RESTARTS   AGE
    docker-deployment-5594b8c597-lgr4d   1/1       Running   0          10m
    ```

<!-- 
## Future work: Install keptn like istio <a id="step-install"></a>

    ```console
    $ kubectl apply keptn.yml
    ```
-->


