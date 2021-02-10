# Keptn CLI

The `keptn` cli is a command line interface for running commands against a Keptn installation.

## Development

Using Go 1.12 (or newer), ensure that you have GO Modules enabled by executing
```console
export GO111MODULE=on
```

You can build the CLI using
```console
go build -o keptn
```

You can execute unit tests using
```console
go test ./...
```

If you want to make sure tests don't influence your local environment (or vice versa), you can run them in a Docker container:
```console
docker run --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.13 go test -race -v ./...
```

### Structure

The cli consists of 

* the entrypoint defined in [main.go](main.go), 
* the root command defined in [cmd/root.go](cmd/root.go),
* all the other commands defined in the [cmd/](cmd/) folder, and
* some utility and helper functions in the [pkg/](pkg/) folder.

## Usage

Use the following syntax to run Keptn commands from your terminal window:

```console
keptn [command] [entity] [name] [flags]
```

where **command**, **entity**, **name**, and **flags** are:

- **command**: Specifies the operation that you want to perform, for example, install, create, onboard, send.

- **entity**: Specifies the entity type. For example, the following commands run a create, onboard, and update operation on the project, service, and domain entity:

    ```console
    keptn create project 
    keptn onboard service
    keptn configure domain
    ```

- **name**: Specifies the name of the entity. Names are case-sensitive. 

- **flags**: Specifies additional parameters and flags the command requires.

If you need help, just run `keptn --help` help from the terminal window.

### Operations

The following table includes short descriptions and the general syntax for all of the `keptn` operations:

| Command  | Description  |
|:---:|---|
| `add-resource`  | Adds a resource to a service within your project in the specified stage |
| `auth`  | Authenticate the Keptn CLI against a Keptn installation  |
| `create`  | Create currently allows to create a project |
| `help`  | Help about any command |
| `install`  | Install Keptn on your Kubernetes cluster |
| `onboard`  | Onboard allows to onboard a new service |
| `send`  | Send a Keptn event in combination with the subcommand *event* |
| `status`  | Checks the status of the CLI |
| `uninstall`  | Uninstalls Keptn on your Kubernetes cluster |
| `version`  | Prints the CLI version for the current context |

## Examples: Common operations
Use the following set of examples to help you familiarize yourself with running the commonly used `keptn` operations:

- Install Keptn on a plain Kubernetes cluster
  ```console
  keptn install --platform=kubernetes
  ```

- Create a project using the definition in a shipyard.yaml
  ```console
  keptn create project my-first-project shipyard.yaml
  ```

- Onboard a (micro)service to the created project
  ```console
  keptn onboard service my-service values.yaml
  ```

- Send a new artifact event for the onboarded service
  ```console
  keptn send event new-artifact --project=my-first-project --service=my-service --image=docker.io/keptnexamples/my-service --tag=0.1.0
  ```
