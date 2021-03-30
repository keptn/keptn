const express = require('express');
const router = express.Router();
const expressSession = require('express-session');
const memoryStore = require('memorystore')(expressSession);
const random = require('crypto-random-string');

const CHECK_PERIOD = 1800000; // check every 30 minutes
const SESSION_TIME = 3600000; // max age is 60 minutes
const COOKIE_LENGTH = 20;
const COOKIE_NAME = 'KTSESSION';
const DEFAULT_TRUST_PROXY = 1;

const SESSION_SECRET = random({length: 50});

const sessionConfig = {
  cookie: {
    path: '/',
    httpOnly: true,
    sameSite: true,
    secure: false,
  },
  name: COOKIE_NAME,
  secret: SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  genid: () => {
    return random({length: COOKIE_LENGTH, type: 'url-safe'})
  },
  store: new memoryStore({
    checkPeriod: CHECK_PERIOD,
    ttl: SESSION_TIME
  }),
}

function initialize(app) {
  console.log('Enabling sessions for bridge.');

  if (process.env.SECURE_COOKIE === 'true') {
    console.log('Setting secure cookies. Make sure SSL is enabled for deployment & correct trust proxy value is used.');
    sessionConfig.cookie.secure = true;

    const trustProxy = process.env.TRUST_PROXY;

    if (trustProxy) {
      console.log('Using trust proxy hops value : ' + trustProxy);
      app.set('trust proxy', parseInt(trustProxy));
    } else {
      console.log('Using default trust proxy hops value : ' + DEFAULT_TRUST_PROXY);
      app.set('trust proxy', DEFAULT_TRUST_PROXY);
    }
  }

  // Register session middleware
  router.use(expressSession(sessionConfig));
  app.use(router);
}

exports.initialize = initialize;
