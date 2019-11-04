const createError = require('http-errors');
const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');

const DatastoreService = require('./lib/services/DatastoreService');
const configs = require('./config');

const apiRouter = require('./api');

const app = express();
const config = configs[app.get('env') || 'development'];

const datastoreService = new DatastoreService(config.datastore);

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, '../dist')));

app.get('/', (req, res) => {
  res.sendFile(path.join(`${__dirname}/../dist/index.html`));
});

app.use('/api', apiRouter({ datastoreService }));


// catch 404 and forward to error handler
app.use((req, res, next) => {
  next(createError(404));
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
