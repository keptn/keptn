const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');

const DatastoreService = require('./lib/services/DatastoreService');
const ConfigurationService = require('./lib/services/ConfigurationService');
const configs = require('./config');

const apiRouter = require('./api');

const app = express();
const config = configs[app.get('env') || 'development'];

const datastoreService = new DatastoreService(config.datastore);
const configurationService = new ConfigurationService(config.configurationService);

// host static files (angular app)
app.use(express.static(path.join(__dirname, '../dist')));

// add some middlewares
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());

// everything starting with /api is routed to the api implementation
app.use('/api', apiRouter({ datastoreService, configurationService }));

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
