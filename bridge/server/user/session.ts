import { Express, Request, Router } from 'express';
import expressSession, { MemoryStore } from 'express-session';
import random from 'crypto-random-string';
import { IdTokenClaims, TokenSet, TokenSetParameters } from 'openid-client';
import MongoStore from 'connect-mongo';
import { Collection, MongoClient } from 'mongodb';
import { Crypto } from './crypto';

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

interface ValidationType {
  _id: string;
  codeVerifier: string;
  nonce: string;
  creationDate: Date;
}

const router = Router();
const SESSION_TIME_SECONDS = getOrDefaultSessionTimeout(60); // session timeout, default to 60 minutes
const SESSION_VALIDATING_DATA_SECONDS = getOrDefaultValidatingDataTimeout(60);
const COOKIE_LENGTH = 10;
const COOKIE_NAME = 'KTSESSION';
const DEFAULT_TRUST_PROXY = 1;
const SESSION_SECRET = process.env.OAUTH_SESSION_SECRET || random({ length: 200 });
const DATABASE_SECRET = process.env.OAUTH_DATABASE_ENCRYPT_SECRET || random({ length: 32 });
const crypto = new Crypto(DATABASE_SECRET);
let store: MemoryStore | MongoStore;
let validationCollection: Collection<ValidationType> | undefined;

if (DATABASE_SECRET.length !== 32) {
  console.error('The length of the env variable "OAUTH_DATABASE_ENCRYPT_SECRET" must be 32');
  process.exit(1);
}

if (process.env.NODE_ENV === 'test') {
  store = new MemoryStore();
} else {
  store = await setupMongoDB();
}

/**
 * Uses a session cookie backed by in-memory cookies store.
 *
 * Cookie store is a LRU cache, hence session removal will occur when there are stale instances.
 */
const sessionConfig = {
  cookie: {
    path: '/',
    httpOnly: true,
    sameSite: true, // if true or 'strict', a redirect to oauth/redirect may lead to a missing cookie
    secure: false,
  },
  name: COOKIE_NAME,
  secret: SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  genid: (): string => random({ length: COOKIE_LENGTH, type: 'url-safe' }),
  store,
} as expressSession.SessionOptions & { cookie: expressSession.CookieOptions };

async function setupMongoDB(): Promise<MongoStore> {
  const mongoCredentials = {
    user: process.env.MONGODB_USER,
    password: process.env.MONGODB_PASSWORD,
    host: process.env.MONGODB_HOST,
    database: process.env.MONGODB_DATABASE || 'openid',
  };

  if (!mongoCredentials.user && !mongoCredentials.password && !mongoCredentials.host) {
    console.error(
      'could not construct mongodb connection string: env vars "MONGODB_HOST", "MONGODB_USER" and "MONGODB_PASSWORD" have to be set'
    );
    process.exit(1);
  }

  const mongoClient = new MongoClient(
    `mongodb://${mongoCredentials.user}:${mongoCredentials.password}@${mongoCredentials.host}`
  );
  await mongoClient.connect();

  validationCollection = mongoClient.db(mongoCredentials.database).collection('validation');

  const indexName = 'validation_index';
  const indexes = await validationCollection.indexes();
  let validationIndex = indexes.find((index) => index.name === indexName);

  if (validationIndex && validationIndex.expireAfterSeconds !== SESSION_VALIDATING_DATA_SECONDS) {
    await validationCollection.dropIndex(indexName);
    validationIndex = undefined;
  }
  if (!validationIndex) {
    await validationCollection.createIndex(
      {
        creationDate: 1,
      },
      {
        name: indexName,
        expireAfterSeconds: SESSION_VALIDATING_DATA_SECONDS,
      }
    );
  }
  console.log('Successfully connected to database');

  return MongoStore.create({
    client: mongoClient,
    ttl: SESSION_TIME_SECONDS, // if inactive for $SESSION_TIME_SECONDS seconds, session is destroyed
    dbName: mongoCredentials.database,
    collectionName: 'sessions',
    crypto: {
      secret: DATABASE_SECRET,
    },
  });
}

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
 * Saves state, code verifier and nonce for the login flow
 */
async function saveValidationData(state: string, codeVerifier: string, nonce: string): Promise<void> {
  await validationCollection?.insertOne({
    _id: state,
    codeVerifier: crypto.encrypt(codeVerifier),
    nonce: crypto.encrypt(nonce),
    creationDate: new Date(),
  });
}

/**
 * returns state, code verifier and nonce and removes it afterwards
 */
async function getAndRemoveValidationData(state: string): Promise<ValidationType | undefined> {
  const data = await validationCollection?.findOneAndDelete({ _id: state });
  try {
    return data?.value
      ? {
          _id: data.value._id,
          codeVerifier: crypto.decrypt(data.value.codeVerifier),
          nonce: crypto.decrypt(data.value.nonce),
          creationDate: data.value.creationDate,
        }
      : undefined;
  } catch (e) {
    console.error('Error wile decrypting validation data. Cause:', e);
  }
}

/**
 * Set the session authenticated state for the specific principal
 *
 * We require a mandatory principal for session authentication. Logout hint is optional and only require when there is
 * logout supported from OAuth service.
 */
function authenticateSession(req: Request, tokenSet: TokenSet): Promise<void> {
  // Regenerate session for the successful login
  return new Promise((resolve) => {
    req.session.regenerate(() => {
      const userIdentifier = process.env.OAUTH_NAME_PROPERTY;
      const claims = tokenSet.claims();
      req.session.authenticated = true;
      req.session.tokenSet = tokenSet;
      req.session.principal =
        (userIdentifier ? (claims[userIdentifier] as string | undefined) : undefined) || getUserIdentifier(claims);
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
  console.log(`Enabling sessions for bridge with session timeout ${SESSION_TIME_SECONDS} seconds.`);

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
  return getDefaultTime(process.env.SESSION_TIMEOUT_MIN, defMinutes);
}

/**
 * Function to determine validating timeout, the timeout for the temporarily saved validation data for the login.
 * Input value is in minutes and return value is in millisecond. Value can be
 * configurable through environment variable SESSION_VALIDATING_TIMEOUT_MIN. If the configuration is invalid, fallback to
 * provided default value.
 */
function getOrDefaultValidatingDataTimeout(defMinutes: number): number {
  return getDefaultTime(process.env.SESSION_VALIDATING_TIMEOUT_MIN, defMinutes);
}

function getDefaultTime(timeMin: string | undefined, defMinutes: number): number {
  if (timeMin) {
    const sTimeout = parseInt(timeMin, 10);

    if (!isNaN(sTimeout) && sTimeout > 0) {
      return sTimeout * 60;
    }
  }

  return defMinutes * 60;
}

export {
  sessionRouter,
  isAuthenticated,
  authenticateSession,
  saveValidationData,
  getAndRemoveValidationData,
  removeSession,
  getLogoutHint,
  getCurrentPrincipal as currentPrincipal,
};
