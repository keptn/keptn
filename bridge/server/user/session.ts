import { Express, Request, Router } from 'express';
import expressSession from 'express-session';
import mS from 'memorystore';
import random from 'crypto-random-string';
import { IdTokenClaims, TokenSet, TokenSetParameters } from 'openid-client';

declare module 'express-session' {
  export interface SessionData {
    authenticated?: boolean;
    tokenSet?: TokenSetParameters;
    principal?: string;
  }
}
export interface ValidSession extends expressSession.Session {
  authenticated: boolean;
  tokenSet: TokenSetParameters;
  principal?: string;
}
const memoryStore = mS(expressSession);
const router = Router();
const CHECK_PERIOD = 600_000; // check every 10 minutes
const SESSION_TIME = getOrDefaultSessionTimeout(60); // session timeout, default to 60 minutes
const COOKIE_LENGTH = 10;
const COOKIE_NAME = 'KTSESSION';
const DEFAULT_TRUST_PROXY = 1;
const SESSION_SECRET = random({ length: 200 });

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
  genid: (): string => random({ length: COOKIE_LENGTH, type: 'url-safe' }),
  store: new memoryStore({
    checkPeriod: CHECK_PERIOD,
    ttl: SESSION_TIME,
  }),
};

/**
 * Filter for for authenticated sessions. Must be enforced by endpoints that require session authentication.
 */
function isAuthenticated(
  session: expressSession.Session & Partial<expressSession.SessionData>
): session is ValidSession {
  if (session.authenticated) {
    return true;
  }

  session.authenticated = false;
  return false;
}

/**
 * Set the session authenticated state for the specific principal
 *
 * We require a mandatory principal for session authentication. Logout hint is optional and only require when there is
 * logout supported from OAuth service.
 */
function authenticateSession(req: Request, tokenSet: TokenSet): Promise<void> {
  return new Promise((resolve) => {
    // Regenerate session for the successful login
    req.session.regenerate(() => {
      const userIdentifier = process.env.OAUTH_NAME_PROPERTY || 'name';
      const claims = tokenSet.claims();
      req.session.authenticated = true;
      req.session.tokenSet = tokenSet;
      req.session.principal = (claims[userIdentifier] as string | undefined) || getUserIdentifier(claims);

      resolve();
    });
  });
}

function getUserIdentifier(claims: IdTokenClaims): string | undefined {
  return claims.name || claims.nickname || claims.preferred_username || claims.email;
}

/**
 * Returns the current principal if session is authenticated. Otherwise returns undefined
 */
function getCurrentPrincipal(req: Request): string | undefined {
  return req.session?.principal;
}

/**
 * Returns the logout hint bound to this session
 */
function getLogoutHint(req: Request): string | undefined {
  return req.session?.tokenSet?.id_token;
}

/**
 * Destroy the session comes with this request
 */
function removeSession(req: Request): void {
  req.session.destroy((error) => {
    if (error) {
      console.error(error);
    }
  });
}

function sessionRouter(app: Express): Router {
  console.log(`Enabling sessions for bridge with session timeout ${SESSION_TIME}ms.`);

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

/**
 * Function to determine session timeout. Input value is in minutes and return value is in millisecond. Value can be
 * configurable through environment variable SESSION_TIMEOUT_MIN. If the configuration is invalid, fallback to
 * provided default value.
 */
function getOrDefaultSessionTimeout(defMinutes: number): number {
  if (process.env.SESSION_TIMEOUT_MIN) {
    const sTimeout = parseInt(process.env.SESSION_TIMEOUT_MIN, 10);

    if (!isNaN(sTimeout) && sTimeout > 0) {
      return sTimeout * 60 * 1000;
    }
  }

  return defMinutes * 60 * 1000;
}

export { sessionRouter };
export { isAuthenticated };
export { authenticateSession };
export { removeSession };
export { getLogoutHint };
export { getCurrentPrincipal as currentPrincipal };
