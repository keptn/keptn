#!/usr/bin/env node

/**
 * Module dependencies.
 */

const app = require('./app');
const debug = require('debug')('server-stub:server');

/**
 * Get port from environment and store in Express.
 */

const PORT = normalizePort(process.env.PORT || '3000');
const HOST = process.env.HOST || '0.0.0.0';

app.set('port', PORT);

/**
 * Listen on provided port, on all network interfaces.
 */

console.log(`Running on http://${HOST}:${PORT}`);
app.listen(PORT, HOST);

/**
 * Normalize a port into a number, string, or false.
 */

function normalizePort (val) {
  const parsedPort = parseInt(val, 10);

  if (isNaN(parsedPort)) {
    // named pipe
    return val;
  }

  if (parsedPort >= 0) {
    // port number
    return parsedPort;
  }

  return false;
}
