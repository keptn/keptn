# Keptn Core Developer Docs


This file contains helpful resources to start developing with Keptn core services.

<details>
<summary>Table of Contents</summary>

<!-- toc -->

- [Useful tools](#useful-tools)
- [Connect to MongoDB](#connect-to-mongodb)
- [Best Practises](#best-practises)
  * [Log and Error format](#log-and-error-format)
  * [Mocking](#mocking)
- [Swagger](#swagger)
- [Adding a new Service](#adding-a-new-service)
- [Installing Gitea](#installing-gitea)

<!-- tocstop -->

</details>

## Useful tools

TBD. Compass, lens, skaffold.

## Connect to MongoDB

TBD. How to port-forwarding and create the connection string.

## Best Practises

Keptn follows also some best practises that are listed in this section.

### Environment Variables

Keptn adhere to the following conventions for setting up environment variables:

- For enabling/disabling a feature, the variable name MUST be called with the feature name and suffixed by `ENABLED`. Example: `MAX_AUTH_ENABLED`


### Log and Error format

Keptn adhere to the following conventions for logging error messages:

- Log lines MUST start with capital letter
- Log lines MUST NOT end with a dot
- Logging MUST use the default Go approach using the default [log package](https://pkg.go.dev/log)
- Returned error messages MUST start with a lower case letter and MUST NOT end with a dot
- Error messages MUST use "Could not ..." for everything that is not supported, e.g. avoid using "unable to", "not able to", "not possible".
- Errors MUST provide context information wrapping errors with `%w`
- Errors MUST be compared by types using `error.Is`
- Custom errors MUST implement the `String()` method

Example:
```go
err := fmt.Errorf("Could not access resource: %w", ErrPermission)
...
if errors.Is(err, ErrPermission) {
  log.Errorf("Failed reading resource: %v", err)
}
```

The Go website provides a good resource around this topic: https://go.dev/blog/go1.13-errors


### Mocking

TBD. Describe how to mock.

## Swagger

TBD. How to generate the swagger doc.

## Adding a new Service

TBD. What it is needed to create a new service. SDK?

## Installing Gitea

TBD. Installing giteas and how to run integration tests based on different ingress types (ClusterIP/LoadBalancer)
