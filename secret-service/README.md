# Secret Service 😎

## Overview

The **SecretService** is used to manage secrets in a Keptn Cluster.
It provides a simple API for creating, updating or deleting secrets in a specific secret backend (e.g. kubernetes, vault,...)

**NOTE:** The current implementation only supports "kubernetes" as a secret backend.

## Secrets and Scopes

A secret created by the **SecretService** is bound to a **Scope**. 
A **Scope** contains a set of **Capabilities** which in turn is a set of permissions.
Currently scopes are hardcoded into a file called `scopes.yaml`.

Default `scopes.yaml`:
```
Scopes:
  keptn-default:
    Capabilities:
      keptn-secrets-default-read:
        Permissions:
          - get
```

**NOTE:** Thus, services making use of a secret in the `default-scope` are only allowed to read the secret.
THe `scopes.yaml` needs to be modified manually in order to add, modify or delete any scopes. Currently,
there is no API endpoint for that.

## Generate  Swagger doc from source

1. Download and install Swag for Go by calling `go get -u github.com/swaggo/swag/cmd/swag` in fresh terminal (not within a Keptn source directory!).
2. `cd` to the Statistics Service's root folder and run `swag init --parseDependency`

*Note*: This updates files in the docs/ subfolder (`docs.go`, `swagger.yaml` and `swagger.json`). It is not necessary to do it manually, as this happens within the [Dockerfile](Dockerfile) during the build process.

