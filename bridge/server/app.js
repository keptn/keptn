const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');

const apiRouter = require('./api');

const app = express();
const apiUrl = process.env.API_URL;
const apiToken = process.env.API_TOKEN;

// host static files (angular app)
app.use(express.static(path.join(__dirname, '../dist')));

// add some middlewares
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());

// check if we need basic authentication
if (process.env.BASIC_AUTH_USERNAME && process.env.BASIC_AUTH_PASSWORD) {
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
  console.error("Not installing authentication middleware");
}


// everything starting with /api is routed to the api implementation
app.use('/api', apiRouter({ apiUrl, apiToken }));

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
  res.status(err.status || 500);
  console.error(err);
  res.json(err);
});

module.exports = app;
