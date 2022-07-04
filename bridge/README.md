# Keptn Bridge

Keptn bridge allows to browse the Keptn events.

Note that yarn dependencies are separated into two parts. Root level `package.json` contains dependencies for angular and other general requirements. Express server dependencies are located inside `package.json` located in server folder.

## Installation

The Keptn bridge is installed as a part of [Keptn](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the bridge in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

### Environment variables
- `ENABLE_VERSION_CHECK` - If disabled, versions.json is not loaded and the version info will not be displayed.
- `SHOW_API_TOKEN` - If disabled, the API token will not be shown in the Bridge info.
- `PROJECTS_PAGE_SIZE` - Determines, how many projects will be fetched for the Bridge. If not set, it defaults to 50.
- `KEPTN_INSTALLATION_TYPE` - Can take the values: `QUALITY_GATES`, `CONTINUOUS_OPERATIONS`, `CONTINUOUS_DELIVERY` and determines the mode in which the Bridge will be started. If only `QUALITY_GATES` is set, only functionalities and data specific for the Quality Gates Only use case will be displayed.

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

### Setting up login via OpenID

To set up a login via OpenID you have to register an application with the identity provider you want, in order to get an ID (`CLIENT_ID`) and a secret (`CLIENT_SECRET`).
After this is done, the following environment variables have to be set:

- `OAUTH_ENABLED` - Flag to enable login via OpenID. To enable it set it to `true`.
- `OAUTH_BASE_URL` - URL of the bridge (e.g. `http://localhost:3000` or `https://myBridgeInstallation.com`).
- `OAUTH_DISCOVERY` - Discovery URL of the identity provider (e.g. https://api.login.yahoo.com/.well-known/openid-configuration).
- `OAUTH_CLIENT_ID` - Client ID.
- `OAUTH_CLIENT_SECRET` (optional) - Client secret. Some identity providers require using the client secret.
- `OAUTH_ID_TOKEN_ALG` (optional) - Algorithm that is used to verify the ID token (e.g. `ES256`). Default is `RS256`.
- `OAUTH_SCOPE` (optional) - Additional scopes that should be added to the authentication flow (e.g. `profile email`), separated by space.
- `OAUTH_NAME_PROPERTY` (optional) - The property of the ID token that identifies the user. Default is `name` and fallback to `nickname`, `preferred_username` and `email`.
- `OAUTH_ALLOWED_LOGOUT_URLS` (optional) - Allowed URLs for the redirect of the end_session endpoint separated by space. Some browsers require to also add the URL the end_session endpoint is redirecting to.

#### Additional information:

- Make sure you add the redirect URI `https://${yourDomain}/${pathToBridge}/oauth/redirect` to your identity provider.
- The identity provider has to support the grant types `authorization_code` and `refresh_token` and provide the endpoints `authorization_endpoint`, `token_endpoint` and `jwks_uri`.
- The refresh of the token is done by the bridge server on demand.
- If the identity provider provides the endpoint `end_session_endpoint`, it will be used for the logout.
- The bridge server itself is a confidential client.

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

1. Run `yarn install` from bridge root level.
2. Run `yarn install` from server folder.
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
4. Run `yarn start:dev` from bridge root level to start the express server and the Angular app.
5. Access the web through the url shown on the console (e.g., http://localhost:3000/ ).

### UI testing with Cypress

[Cypress Reference](https://docs.cypress.io/api/table-of-contents)

UI tests in Keptn Bridge must not require any API call. When writing tests, please make sure to mock every outgoing request to `/api` with `cy.intercept`<br/>
[Cypress Intercept Reference](https://docs.cypress.io/api/commands/intercept)

To test your UI tests locally, use following commands:

- `yarn cypress:open` (Linux, macOS), `yarn cypress:open:win32` (Windows) - Used for the local development of tests.<br/>This opens a browser, where you can run your tests and inspect them. The tests will re-run automatically on every code change made on the `*.spec.ts` files.
- or `yarn test:ui` (Linux, macOS), `yarn test:ui:win32` (Windows) - This starts the headless browser mode that is also used in CI.<br/>This will just run the tests on a headless browser without the possibility to inspect them.

Both commands serve Angular on port 3000 with no live reload - this will also ensure that no API connection is made.

#### Known issues

- Currently, our UI tests are flaky because of some bugs in Cypress. You can find more information in [Known Issues](https://github.com/keptn/keptn/issues/7107).  
- One UI test will fail if you are on Windows and in a different time zone than Europe/Berlin due to a bug in Cypress.

### Bundle Size Report

The Keptn Bridge is bundled with the Angular CLI. To analyze the current bundle size, first build the application with `ng build --stats-json`, to generate the [`stats.json`](https://webpack.js.org/api/stats/) file. Then run the [Webpack Bundle Analyzer](https://github.com/webpack-contrib/webpack-bundle-analyzer) with `npm run bundle-report` to create an interactive treemap visualization of the contents of all your bundles.

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

#### Testing Assertions

Assertions should target a whole object, for example, given the following Object definition

```
interface OAuthConfig {
  enabled: boolean;
  discoveryURL: string;
  clientSecret?: string;
}
```

An assertion on an object with type `OAuthConfig` should look like for the whole object and not to the single properties:

```
const c = getOAuthConfig();
// GOOD
expect(c).toEqual({
  enabled: true;
  discoveryURL: "localhost";
});
// BAD
expect(c.enabled).toStrictEqual(true);
expect(c.discoveryURL).toStrictEqual("localhost");
expect(c.clientSecret).toStrictEqual(undefined);
```

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

- Optimize Imports
- Run eslint --fix

#### Visual Studio Code

The save actions are handled within the `settings.json` file. Please follow the guide for importing the file to enable the feature.

### Enable ESLint

#### Git

We use `\n` for line-endings and this is also configured in ESLint. To tell Git to use the right line-ending, execute the command `git config --global core.autocrlf input`

#### IntelliJ

`Editor -> Inspections -> JavaScript and TypeScript -> Code quality tools`
Enable `ESLint` and disable `TSLint` if enabled

The automatic ESLint configuration detects automatically the `.eslintrc.json` file in the bridge directory and applies the set rules.

#### Visual Studio Code

- Install the [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) extension
- Run `yarn add eslint`
- Close and re-open VSCode for ESLint to work properly

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

- [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)
- [Angular Language Service](https://marketplace.visualstudio.com/items?itemName=Angular.ng-template)
- [Jest](https://marketplace.visualstudio.com/items?itemName=Orta.vscode-jest)
- [Jest Runner](https://marketplace.visualstudio.com/items?itemName=firsttris.vscode-jest-runner)
- [Snyk Vulnerability Scanner](https://marketplace.visualstudio.com/items?itemName=snyk-security.snyk-vulnerability-scanner)
- [SonarLint](https://marketplace.visualstudio.com/items?itemName=SonarSource.sonarlint-vscode)
- [Conventional Commit](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits)
  <br/>Conventional Commit settings are also already defined in `settings.json`

## Production deployment

See [Dockerfile](Dockerfile) for the latest instructions. By default, the process will listen on port 3000.
