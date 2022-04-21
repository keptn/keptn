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

- **command**: Specifies the operation that you want to perform, for example, install, create, send.

- **entity**: Specifies the entity type. For example, the following commands run a create, get, and update operation on the project, service, and domain entity:

    ```console
    keptn create project 
    keptn get service
    keptn configure domain
    ```

- **name**: Specifies the name of the entity. Names are case-sensitive. 

- **flags**: Specifies additional parameters and flags the command requires.

If you need help, just run `keptn --help` help from the terminal window.

### Operations

The following table includes short descriptions and the general syntax for all of the `keptn` operations:

| Command  | Description  |
|:---:|---|
| `abort`  | Aborts the execution of a sequence |
| `add-resource`  | Adds a local resource to a service within your project in the specified stage |
| `auth`  | Authenticates the Keptn CLI against a Keptn installation  |
| `completion`  | Generate completion script  |
| `configure`  | Configures one of the specified parts of Keptn  |
| `create`  | Creates a new project, service or secret |
| `delete`  | Deletes a project |
| `generate`  | Generates the markdown CLI documentation or a support archive |
| `get`  | Displays an event or Keptn entities such as project, stage, or service |
| `help`  | Help about any command |
| `install`  | Installs Keptn on a Kubernetes cluster |
| `pause`  | Pauses the execution of a sequence |
| `resume`  | Resumes the execution of a sequence |
| `send`  | Sends an event to Keptn |
| `set`  | Sets flags of the CLI configuration |
| `status`  | Checks the status of the CLI |
| `trigger`  | Triggers the execution of an action in keptn |
| `uninstall`  | Uninstalls Keptn from a Kubernetes cluster |
| `update`  | Updates an existing Keptn project |
| `upgrade`  | Upgrades Keptn on a Kubernetes cluster |
| `version`  | Shows the version of Keptn and Keptn CLI |

## Examples: Common operations
Use the following set of examples to help you familiarize yourself with running the commonly used `keptn` operations:

- Install Keptn on a plain Kubernetes cluster
  ```console
  keptn install --platform=kubernetes
  ```

- Create a project using the definition in a shipyard.yaml
  ```console
  keptn create project my-first-project --shipyard=shipyard.yaml
  ```

- Create a service for the new project
  ```console
  keptn create service my-first-service --project=my-first-project
  ```

- Trigger the delivery of a new artifact for the project's new service
  ```console
  keptn trigger delivery --project=my-first-project --service=my-first-service --image=docker.io/keptnexamples/my-service:0.1.0
  ```
