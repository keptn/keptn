import { Router, Request, Express } from 'express';
import expressSession from 'express-session';
import mS from 'memorystore';
import random from 'crypto-random-string';

declare module 'express-session' {
  export interface SessionData {
    authenticated: boolean;
    principal: string;
    logoutHint: string;
  }
}
const memoryStore = mS(expressSession);
const router = Router();
const CHECK_PERIOD = 600_000; // check every 10 minutes
const SESSION_TIME = 1_200_000; // max age is 20 minutes
const COOKIE_LENGTH = 10;
const COOKIE_NAME = 'KTSESSION';
const DEFAULT_TRUST_PROXY = 1;
const SESSION_SECRET = random({length: 200});

/**
 * Uses a session cookie backed by in-memory cookies store.
 *
 * Cookie store is a LRU cache, hence session removal will occur when there are stale instances.
 */
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
    return random({length: COOKIE_LENGTH, type: 'url-safe'});
  },
  store: new memoryStore({
    checkPeriod: CHECK_PERIOD,
    ttl: SESSION_TIME
  }),
};

/**
 * Filter for for authenticated sessions. Must be enforced by endpoints that require session authentication.
 */
function isAuthenticated(req: Request) {
  if (req.session.authenticated) {
    return true;
  }

  req.session.authenticated = false;
  return false;
}

/**
 * Set the session authenticated state for the specific principal
 *
 * We require a mandatory principal for session authentication. Logout hint is optional and only require when there is
 * logout supported from OAuth service.
 */
function setAuthenticatedSession(req: Request, principal: string, logoutHint: string) {

  if (!principal) {
    throw Error('Invalid session initialisation. Principal is mandatory.');
  }

  req.session.authenticated = true;
  req.session.principal = principal;

  if (logoutHint) {
    req.session.logoutHint = logoutHint;
  }
}

/**
 * Returns the current principal if session is authenticated. Otherwise returns undefined
 */
function getCurrentPrincipal(req: Request): string | undefined {
  if (req.session !== undefined && req.session.authenticated) {
    return req.session.principal;
  }

  return undefined;
}

/**
 * Returns the logout hint bound to this session
 */
function getLogoutHint(req: Request): string | undefined {
  return req.session?.logoutHint;
}

/**
 * Destroy the session comes with this request
 */
function removeSession(req: Request) {
  req.session.destroy(console.error);
}

function sessionRouter(app: Express)  {
  console.log('Enabling sessions for bridge.');

  if (process.env.SECURE_COOKIE === 'true') {
    console.log('Setting secure cookies. Make sure SSL is enabled for deployment & correct trust proxy value is used.');
    sessionConfig.cookie.secure = true;

    const trustProxy = process.env.TRUST_PROXY;

    if (trustProxy) {
      console.log('Using trust proxy hops value : ' + trustProxy);
      app.set('trust proxy', +trustProxy);
    } else {
      console.log('Using default trust proxy hops value : ' + DEFAULT_TRUST_PROXY);
      app.set('trust proxy', DEFAULT_TRUST_PROXY);
    }
  }

  // Register session middleware
  router.use(expressSession(sessionConfig));

  return router;
}

export { sessionRouter };
export { isAuthenticated };
export { setAuthenticatedSession as sessionAuthentication };
export { removeSession };
export { getLogoutHint };
export { getCurrentPrincipal as currentPrincipal };
