# Keptn Bridge

Keptn bridge allows to browse the Keptn events.

Note that npm dependencies are separated into two parts. Root level `package.json` contains dependencies for angular and other general requirements. Express server dependencies are located inside `package.json` located in server folder.

## Installation

The Keptn bridge is installed as a part of [Keptn](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the bridge in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

### Setting up Basic Authentication

Keptn Bridge comes with a very simple basic authentication feature, which can be controlled by setting the following two environment variables:

- `BASIC_AUTH_USERNAME` - username
- `BASIC_AUTH_PASSWORD` - password

To enable it within your Kubernetes cluster, we recommend first creating a secret which holds the two variables, and then apply this secret within the Kubernetes deployment for Keptn Bridge.

1. Create the secret using

   ```console
   kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=<USERNAME>" --from-literal="BASIC_AUTH_PASSWORD=<PASSWORD>"
   ```

   _Note: Replace `<USERNAME>` and `<PASSWORD>` with the desired credentials._

2. In case you are using Keptn 0.6.1 or older, edit the deployment using

   ```console
   kubectl -n keptn edit deployment bridge
   ```

   and add the secret to the `bridge` container, such that the container spec looks like this:

   ```yaml
       ...
       spec:
         containers:
         - name: bridge
           image: keptn/bridge2:0.6.1
           imagePullPolicy: Always
           # EDIT STARTS HERE
           envFrom:
             - secretRef:
                 name: bridge-credentials
                 optional: true
           # EDIT ENDS HERE
           ports:
           - containerPort: 3000
           ...
   ```

**Note 1**: To disable authentication, just delete the secret using `kubectl -n keptn delete secret bridge-credentials`.

**Note 2**: If you delete or edit the secret, you need to restart the respective pod by executing

```console
kubectl -n keptn scale deployment bridge --replicas=0
kubectl -n keptn scale deployment bridge --replicas=1
```

#### Throttling

Per default login attempts are throttled to 10 requests within 60 minutes. This behavior can be adjusted with the following environment variables:

- `REQUESTS_WITHIN_TIME` - how many login attempts in which timespan `REQUEST_TIME_LIMIT` (in minutes) are allowed per IP.
- `CLEAN_BUCKET_INTERVAL` - the interval (in minutes) the saved login attempts should be checked and deleted if the last request of an IP is older than `REQUEST_TIME_LIMIT` minutes. Default is 60 minutes.

### Custom Look And Feel

You can change the Look And Feel of the Keptn Bridge by creating a zip archive with your resources and make it downloadable from an URL.

When the `LOOK_AND_FEEL_URL` environment variable is set and points to a zip archive the Keptn Bridge will download that file on startup and extract its content into `/assets/branding`.

The zip archive must contain an `app-config.json` and can have optionally a logo and a stylesheet. The `app-config.json`
must define an `appTitle`, `logoUrl` and `logoInvertedUrl` and can have optionally a `stylesheetUrl`. The `logoUrl` will be used as logo in the app header, `logoInvertedUrl` will be used as app favicon and as placeholder in some empty state messages. If a `stylesheetUrl` is provided, the stylesheet
will be injected in the app header on page load.

```app-config.json
{
  "appTitle": "custom title",
  "logoUrl": "assets/branding/logo.svg",
  "logoInvertedUrl": "assets/branding/logo.svg",
  "stylesheetUrl": "assets/branding/style.css"
}
```

If no `LOOK_AND_FEEL_URL` was provided, the Bridge will use the default `logo.png`, `logo_inverted.png` and an `app-config.json`.

### Delete in your Kubernetes cluster

To delete a deployed bridge, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## Local development

1. Run `npm install` from bridge root level.
2. Run `npm install` from server folder.
3. Set `API_URL` and `API_TOKEN` environment variables, depending on your Keptn installation and operating system:
   **Linux/MacOS**
   ```console
   export API_URL=http://keptn.127.0.0.1.nip.io/api
   export API_TOKEN=1234-exam-ple
   ```
   **Windows**
   ```console
   set API_URL=http://keptn.127.0.0.1.nip.io/api
   set API_TOKEN=1234-exam-ple
   ```
4. Run `npm run start:dev` from bridge root level to start the express server and the Angular app.
5. Access the web through the url shown on the console (e.g., http://localhost:3000/ ).

## IDE Setup

### Create workspace

Before creating the workspace, make sure that you already cloned the keptn/keptn repository to your local file system. For IntelliJ we would recommend cloning the entire repository to the `Idea Projects` folder as it is easier to import modules.

#### IntelliJ

##### From the start screen, when no project is opened

1. In the menu make sure that `Projects` is selected
2. Click on the button `New Project`
3. Select `Empty Project` (the last entry in the list) and click `Next`
4. Type in the project name (e.g. keptn-bridge) and select the location where you cloned the source to. <br/>
   Select the bridge folder in the keptn source
5. Click on `Finish`
6. Close the Project Structure dialog that pops up after project creation by clicking `Cancel`
7. `File > New > Module from Existing Sources...`
8. Select the bridge folder
9. Select `Create module from existing sources`
10. On the next screen make sure that the bridge folder is checked and click `Next`
11. On the next screen click `Finish`
12. (optional) You can also remove the keptn-bridge module by right-clicking on it and selecting `Remove Module` from the context menu.

If the imported module doesn't show up immediately, close the IDE and re-open it. The bridge module should now be visible in the project.

##### From an opened workspace

1. Select `File > New > Project from Existing Sources...`
2. Follow steps 3 and onward from the description above

#### Visual Studio Code

1. `File > Open Folder...`
2. Select the bridge folder in the keptn source

### Code Style

#### IntelliJ

`File > Settings... > Editor > Code Style > TypeScript`

- Select `Project` as scheme
- Click on the cogwheel next to Scheme. A dropdown menu opens.
- In the opened menu click on `Import Scheme > IntellJ IDEA code style XML`
- Open the provided `IntelliJ.xml` file - the code styles get applied on project scope
- Apply the changes

#### Visual Studio Code

In the opened folder:

- If there is no `.vscode` directory create one
- Copy the `settings.json` to the `.vscode` directory

### Save Actions

#### IntelliJ

`File > Settings... > Tools > Actions on Save`

Make sure that following items are checked:

- Reformat Code (Whole file)
- Optimize Imports

#### Visual Studio Code

The save actions are handled within the `settings.json` file. Please follow the guide for importing the file to enable the feature.

### Enable TSLint

#### IntelliJ

`File > Settings... > Languages & Tools > TypeScript > TSLint`

Select `Automatic TSLint configuration`

The automatic TSLint configuration detects automatically the `tslint.json` file in the bridge directory and applies the set rules.

#### Visual Studio Code

- Install the [TSLint](https://marketplace.visualstudio.com/items?itemName=ms-vscode.vscode-typescript-tslint-plugin) extension
- Run `npm i tslint`
- Open the command palette (Windows: `Ctrl + Shift + P`)
- Type in `workspace execution`
- Select `TSLint: Manage workspace library execution`
- Select `Enable workspace library execution`
- Close and re-open VSCode for TSLint to work properly

### Additional configurations

#### IntelliJ

- Disable Checkbox: `File > Settings... > Editor > Inspections > JavaScript and TypeScript > General > Method can be static`

### Recommended Plugins

#### IntelliJ

- [Snyk Vulnerability Scanner](https://plugins.jetbrains.com/plugin/10972-snyk-vulnerability-scanner)
- [SonarLint](https://plugins.jetbrains.com/plugin/7973-sonarlint)
- [Conventional Commit](https://plugins.jetbrains.com/plugin/13389-conventional-commit)
  For our commit guidelines please consult our Contributing Guide for [making pull requests](https://github.com/keptn/keptn/blob/master/CONTRIBUTING.md#make-a-pull-request) and [commit types and scopes](https://github.com/keptn/keptn/blob/master/CONTRIBUTING.md#commit-types-and-scopes)

#### Visual Studio Code

- [TSLint](https://marketplace.visualstudio.com/items?itemName=ms-vscode.vscode-typescript-tslint-plugin)
- [Angular Language Service](https://marketplace.visualstudio.com/items?itemName=Angular.ng-template)
- [Jest](https://marketplace.visualstudio.com/items?itemName=Orta.vscode-jest)
- [Jest Runner](https://marketplace.visualstudio.com/items?itemName=firsttris.vscode-jest-runner)
- [Snyk Vulnerability Scanner](https://marketplace.visualstudio.com/items?itemName=snyk-security.snyk-vulnerability-scanner)
- [SonarLint](https://marketplace.visualstudio.com/items?itemName=SonarSource.sonarlint-vscode)
- [Conventional Commit](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits)
  <br/>Conventional Commit settings are also already defined in `settings.json`

## Production deployment

See [Dockerfile](Dockerfile) for the latest instructions. By default, the process will listen on port 3000.
