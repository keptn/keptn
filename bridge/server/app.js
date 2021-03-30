const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const {execSync} = require('child_process');

const apiRouter = require('./api');
const sessionInit = require('./user/session').initialize;

const app = express();
let apiUrl = process.env.API_URL;
let apiToken = process.env.API_TOKEN;
let cliDownloadLink = process.env.CLI_DOWNLOAD_LINK;
let integrationsPageLink = process.env.INTEGRATIONS_PAGE_LINK;

if(!apiToken) {
  console.log("API_TOKEN was not provided. Fetching from kubectl.");
  apiToken = Buffer.from(execSync('kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token}').toString(), 'base64').toString();
}

if (!cliDownloadLink) {
  console.log("CLI Download Link was not provided, defaulting to github.com/keptn/keptn releases")
  cliDownloadLink = "https://github.com/keptn/keptn/releases";
}

if(!integrationsPageLink) {
  console.log("Integrations page Link was not provided, defaulting to get.keptn.sh/integrations.html")
  integrationsPageLink = "https://get.keptn.sh/integrations.html";
}

const oneWeek       = 7*24*3600000;    // 3600000msec == 1hour
// host static files (angular app)
app.use(express.static(path.join(__dirname, '../dist'), { maxAge: oneWeek }));

// add some middlewares
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());

let authType;

// todo : This needs to be set as a complete authentication mechanism
if (process.env.OAUTH_ENABLED === 'true') {
  authType = "OAUTH";
  sessionInit(app);

  // todo : handle by session
  app.use('/', (req, resp, next) => {
    req.session.authenticated = false;
    return next();
  });
}

if (process.env.BASIC_AUTH_USERNAME && process.env.BASIC_AUTH_PASSWORD) {
  authType = 'BASIC';

  console.error("Installing Basic authentication - please check environment variables!");
  app.use((req, res, next) => {
    // parse login and password from headers
    const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
    const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');

    // Verify login and password are set and correct
    if (!(login && password && login === process.env.BASIC_AUTH_USERNAME && password === process.env.BASIC_AUTH_PASSWORD)) {
      // Access denied
      console.error("Access denied");
      res.set('WWW-Authenticate', 'Basic realm="Keptn"');
      res.status(401).send('Authentication required.'); // custom message
      return;
    }

    // Access granted
    return next();
  });
} else {
  authType = 'NONE';
  console.error("Not installing authentication middleware");
}


// everything starting with /api is routed to the api implementation
app.use('/api', apiRouter({ apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType }));

// fallback: go to index.html
app.use((req, res, next) => {
  console.error("Not found: " + req.url);
  res.sendFile(path.join(`${__dirname}/../dist/index.html`));
});

// error handler
// eslint-disable-next-line no-unused-vars
app.use((err, req, res, next) => {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};
  // render the error page
  res.status(err.status || 500).send();
  console.error(err);
  // res.json(err);
});

module.exports = app;
