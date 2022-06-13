import { Express, NextFunction, Request, Response, Router } from 'express';
import expressSession, { MemoryStore } from 'express-session';
import random from 'crypto-random-string';
import { IdTokenClaims, TokenSet, TokenSetParameters } from 'openid-client';
import MongoStore from 'connect-mongo';
import { Collection, Db, MongoClient } from 'mongodb';
import { Crypto } from './crypto';
import { getRootLocation } from './oauth-routes';
import { getOAuthMongoExternalConnectionString, getOAuthSecrets } from './secrets';
import { ComponentLogger } from '../utils/logger';

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

// before using it "init" has to be called!
export class SessionService {
  private validationCollection: Collection<ValidationType> | undefined;
  private sessionConfig: (expressSession.SessionOptions & { cookie: expressSession.CookieOptions }) | undefined;
  private readonly crypto: Crypto;
  private readonly SESSION_TIME_SECONDS: number;
  private readonly DEFAULT_TRUST_PROXY = 1;
  private readonly cookieName = 'KTSESSION';
  private readonly cookieLength = 10;
  private readonly sessionSecret: string;
  private readonly SESSION_VALIDATING_DATA_SECONDS;
  private readonly databaseSecret: string;
  private readonly validationCollectionName = 'validation';
  private readonly sessionCollectionName = 'sessions';
  private readonly log = new ComponentLogger('OAuth');

  constructor() {
    this.SESSION_TIME_SECONDS = this.getOrDefaultSessionTimeout(60); // session timeout, default to 60 minutes
    this.SESSION_VALIDATING_DATA_SECONDS = this.getOrDefaultValidatingDataTimeout(60);
    const errorSuffix =
      'must be defined when OAuth based login (OAUTH_ENABLED) is activated. Please check your bridge-oauth secret.';
    const secrets = getOAuthSecrets();
    if (!secrets.sessionSecret) {
      throw Error(`session_secret ${errorSuffix}`);
    }

    if (!secrets.databaseEncryptSecret) {
      throw Error(`database_encrypt_secret ${errorSuffix}`);
    } else if (secrets.databaseEncryptSecret.length !== 32) {
      throw Error(`The length of the secret "database_encrypt_secret" must be 32`);
    }
    this.sessionSecret = secrets.sessionSecret;
    this.databaseSecret = secrets.databaseEncryptSecret;
    this.crypto = new Crypto(this.databaseSecret);
  }

  public async init(): Promise<SessionService> {
    let store: MemoryStore | MongoStore;

    if (process.env.NODE_ENV === 'test') {
      store = new MemoryStore();
    } else {
      store = await this.setupMongoDB();
    }

    /**
     * Uses a session cookie backed by in-memory cookies store.
     *
     * Cookie store is a LRU cache, hence session removal will occur when there are stale instances.
     */
    this.sessionConfig = {
      cookie: {
        path: '/',
        httpOnly: true,
        sameSite: true, // if true or 'strict', a redirect to oauth/redirect may lead to a missing cookie
        secure: false,
      },
      name: this.cookieName,
      secret: this.sessionSecret,
      resave: false,
      saveUninitialized: false,
      genid: (): string => random({ length: this.cookieLength, type: 'url-safe' }),
      store,
    };
    return this;
  }

  public async setupMongoDB(): Promise<MongoStore> {
    const mongoCredentials = {
      user: process.env.MONGODB_USER,
      password: process.env.MONGODB_PASSWORD,
      host: process.env.MONGODB_HOST,
      database: process.env.MONGODB_DATABASE || 'openid',
    };
    const externalConnectionString = getOAuthMongoExternalConnectionString();

    if (!externalConnectionString && !mongoCredentials.user && !mongoCredentials.password && !mongoCredentials.host) {
      this.log.error(
        'could not construct mongodb connection string: env vars "MONGODB_HOST", "MONGODB_USER" and "MONGODB_PASSWORD" have to be set'
      );
      process.exit(1);
    }

    const mongoClient = new MongoClient(
      externalConnectionString ||
        `mongodb://${mongoCredentials.user}:${mongoCredentials.password}@${mongoCredentials.host}/${mongoCredentials.database}`
    );
    await mongoClient.connect();
    const db = mongoClient.db();
    await this.initCollections(db);

    this.validationCollection = db.collection(this.validationCollectionName);
    await this.initValidationTTLIndex(this.validationCollection);
    this.log.info('Successfully connected to database');

    return MongoStore.create({
      client: mongoClient,
      ttl: this.SESSION_TIME_SECONDS, // if inactive for $SESSION_TIME_SECONDS seconds, session is destroyed
      collectionName: this.sessionCollectionName,
      crypto: {
        secret: this.databaseSecret,
      },
      touchAfter: this.SESSION_TIME_SECONDS / 2, // session is only updated every {this.SESSION_TIME_SECONDS / 2} seconds
    });
  }

  private async initValidationTTLIndex(collection: Collection<ValidationType>): Promise<void> {
    const indexName = 'validation_index';
    const indexes = await collection.indexes();
    let validationIndex = indexes.find((index) => index.name === indexName);

    if (validationIndex && validationIndex.expireAfterSeconds !== this.SESSION_VALIDATING_DATA_SECONDS) {
      await collection.dropIndex(indexName);
      validationIndex = undefined;
    }
    if (!validationIndex) {
      await collection.createIndex(
        {
          creationDate: 1,
        },
        {
          name: indexName,
          expireAfterSeconds: this.SESSION_VALIDATING_DATA_SECONDS,
        }
      );
    }
  }

  private async initCollections(db: Db): Promise<void> {
    const createCollections = [this.validationCollectionName, this.sessionCollectionName];
    const collections = await db
      .listCollections(
        { name: { $in: createCollections } },
        {
          nameOnly: true,
          authorizedCollections: true,
        }
      )
      .toArray();
    if (collections.length === createCollections.length) {
      return;
    }
    const newCollections = createCollections.filter(
      (col) => !collections.some((collection) => collection.name === col)
    );
    for (const createCollection of newCollections) {
      await db.createCollection(createCollection);
    }
  }

  /**
   * Filter for authenticated sessions. Must be enforced by endpoints that require session authentication.
   */
  public isAuthenticated(
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
  public async saveValidationData(state: string, codeVerifier: string, nonce: string): Promise<void> {
    await this.validationCollection?.insertOne({
      _id: state,
      codeVerifier: this.crypto.encrypt(codeVerifier),
      nonce: this.crypto.encrypt(nonce),
      creationDate: new Date(),
    });
  }

  /**
   * returns state, code verifier and nonce and removes it afterwards
   */
  public async getAndRemoveValidationData(state: string): Promise<ValidationType | undefined> {
    const data = await this.validationCollection?.findOneAndDelete({ _id: state });
    try {
      return data?.value
        ? {
            _id: data.value._id,
            codeVerifier: this.crypto.decrypt(data.value.codeVerifier),
            nonce: this.crypto.decrypt(data.value.nonce),
            creationDate: data.value.creationDate,
          }
        : undefined;
    } catch (err) {
      const msg = err instanceof Error ? `${err.name}: ${err.message}` : `${err}`;
      this.log.error(`Error wile decrypting validation data. Cause: ${msg}`);
    }
  }

  /**
   * Set the session authenticated state for the specific principal
   *
   * We require a mandatory principal for session authentication. Logout hint is optional and only require when there is
   * logout supported from OAuth service.
   */
  public authenticateSession(req: Request, tokenSet: TokenSet): Promise<void> {
    // Regenerate session for the successful login
    return new Promise((resolve) => {
      req.session.regenerate(() => {
        const userIdentifier = process.env.OAUTH_NAME_PROPERTY;
        const claims = tokenSet.claims();
        req.session.authenticated = true;
        req.session.tokenSet = tokenSet;
        req.session.principal =
          (userIdentifier ? (claims[userIdentifier] as string | undefined) : undefined) ||
          this.getUserIdentifier(claims);
        resolve();
      });
    });
  }

  private getUserIdentifier(claims: IdTokenClaims): string | undefined {
    return claims.name || claims.nickname || claims.preferred_username || claims.email;
  }

  /**
   * Returns the current principal if session is authenticated. Otherwise returns undefined
   */
  public getCurrentPrincipal(req: Request): string | undefined {
    return req.session?.principal;
  }

  /**
   * Returns the logout hint bound to this session
   */
  public getLogoutHint(req: Request): string | undefined {
    return req.session?.tokenSet?.id_token;
  }

  /**
   * Destroy the session comes with this request
   */
  public removeSession(req: Request): void {
    req.session.destroy((error) => {
      if (error) {
        this.log.error(error);
      }
    });
  }

  public sessionRouter(app: Express): Router {
    if (!this.sessionConfig) {
      throw Error('Session store is not initialized! Did you forget to call init()?');
    }
    const router = Router();
    this.log.info(`Enabling sessions for bridge with session timeout ${this.SESSION_TIME_SECONDS} seconds.`);

    if (process.env.SECURE_COOKIE === 'true') {
      this.log.info(
        'Setting secure cookies. Make sure SSL is enabled for deployment & correct trust proxy value is used.'
      );
      this.sessionConfig.cookie.secure = true;

      const trustProxy = process.env.TRUST_PROXY;

      if (trustProxy) {
        this.log.info('Using trust proxy hops value : ' + trustProxy);
        app.set('trust proxy', +trustProxy);
      } else {
        this.log.info('Using default trust proxy hops value : ' + this.DEFAULT_TRUST_PROXY);
        app.set('trust proxy', this.DEFAULT_TRUST_PROXY);
      }
    }

    // Register session middleware
    router.use(async (req, res, next: NextFunction) => {
      const status = await this.setSessionAndGetNextResponse(req, res);
      if (
        status instanceof Error &&
        (status.message === 'Encrypted session was tampered with!' || status.message.startsWith('Unexpected token'))
      ) {
        // Database encryption changed. Delete session cookie of client

        // Redirect to login does not work because the redirect to the SSO provider fails then
        // Probably related to some invalid headers

        res.cookie(this.cookieName, '', { expires: new Date() }); // on redirect res.clearCookie does not work
        if (req.path.startsWith('/api/')) {
          // Redirect to root does not work
          // Client does not reload the page and another call is not triggered immediately
          next({ response: { status: 401 } });
        } else {
          // On initial load, if pages are directly accessed like /dashboard, response must not be 401,
          // else this would not allow fetching/forwarding/accessing the page/website
          res.redirect(getRootLocation());
        }
      } else {
        return next(status);
      }
    });

    return router;
  }

  private setSessionAndGetNextResponse(req: Request, res: Response): Promise<Error | string | undefined> {
    return new Promise((resolve) => {
      expressSession(this.sessionConfig)(req, res, (status: Error | string | undefined) => {
        resolve(status);
      });
    });
  }

  /**
   * public to determine session timeout. Input value is in minutes and return value is in millisecond. Value can be
   * configurable through environment variable SESSION_TIMEOUT_MIN. If the configuration is invalid, fallback to
   * provided default value.
   */
  private getOrDefaultSessionTimeout(defMinutes: number): number {
    return this.getDefaultTime(process.env.SESSION_TIMEOUT_MIN, defMinutes);
  }

  /**
   * public to determine validating timeout, the timeout for the temporarily saved validation data for the login.
   * Input value is in minutes and return value is in millisecond. Value can be
   * configurable through environment variable SESSION_VALIDATING_TIMEOUT_MIN. If the configuration is invalid, fallback to
   * provided default value.
   */
  private getOrDefaultValidatingDataTimeout(defMinutes: number): number {
    return this.getDefaultTime(process.env.SESSION_VALIDATING_TIMEOUT_MIN, defMinutes);
  }

  private getDefaultTime(timeMin: string | undefined, defMinutes: number): number {
    if (timeMin) {
      const sTimeout = parseInt(timeMin, 10);

      if (!isNaN(sTimeout) && sTimeout > 0) {
        return sTimeout * 60;
      }
    }

    return defMinutes * 60;
  }
}
