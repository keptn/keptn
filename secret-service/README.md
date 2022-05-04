# Secret Service 😎

## Overview

The **SecretService** is used to manage secrets in a Keptn Cluster.
It provides a simple API for creating, updating or deleting secrets in a specific secret backend (e.g. kubernetes, vault,...)

**NOTE:** The current implementation only supports "kubernetes" as a secret backend.

## Secret and Scopes

A secret created by the secret-service is bound to a scope.
A scope contains a set of capabilities which in turn is a set of permissions.
Currently scopes are hardcoded into a file called `scopes.yaml` which is read by the secret-service during startup.

The default scope for Keptn looks like this:
```
Scopes:
  keptn-default:
    Capabilities:
      keptn-secrets-default-read:
        Permissions:
          - get
  keptn-webhook-service:
    Capabilities:
      keptn-webhook-svc-read:
        Permissions:
          - get
  dynatrace-service:
    Capabilities:
      keptn-dynatrace-svc-read:
        Permissions:
          - get
```

In Kubernetes, *scope* maps to a K8S *ServiceAccount* and a capability maps to a K8S *Role*.

Based on the `scopes.yaml` file above, when a secret with scope `keptn-webhook-service` is created, the secret-service will:
- create a K8S secret
- create a *Role* named `keptn-webhook-svc-read` containing rules to access the created secret with permissions `get`
- create a *Rolebinding* `keptn-webhook-service-rolebinding` with *subjects* set to the *ServiceAccount* named `keptn-webhook-service`
-

Thus, every K8S Pod bound to the service account *keptn-webhook-service* is able to read the secret.

**NOTE:** The `scopes.yaml` needs to be modified manually in order to add, modify or delete any scopes. Currently,
there is no API endpoint for that.

## Generate  Swagger doc from source

1. Download and install Swag for Go by calling `go install github.com/swaggo/swag/cmd/swag` in fresh terminal.
2. `cd` to the SecretService's root folder and run `swag init`
