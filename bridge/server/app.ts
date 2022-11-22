import { copyFileSync, createWriteStream, existsSync, mkdirSync, unlinkSync, WriteStream } from 'fs';
import { dirname, join } from 'path';
import { unlink } from 'fs/promises';
import * as https from 'https';
import * as http from 'http';
import { contentSecurityPolicy, frameguard, noSniff, permittedCrossDomainPolicies, xssFilter } from 'helmet';
import express, { Express, NextFunction, Request, Response } from 'express';
import { fileURLToPath, URL } from 'url';
import { ComponentLogger } from './utils/logger';
import cookieParser from 'cookie-parser';
import AdmZip from 'adm-zip';
import { apiRouter } from './api';
import { AxiosError } from 'axios';
import { EnvironmentUtils } from './utils/environment.utils';
import { ClientFeatureFlags, ServerFeatureFlags } from './feature-flags';
import { setupOAuth } from './user/oauth';
import { SessionService } from './user/session';
import { ContentSecurityPolicyOptions } from 'helmet/dist/types/middlewares/content-security-policy';
import { printError } from './utils/print-utils';
import { AuthType } from '../shared/models/auth-type';
import { AuthConfig, BridgeConfiguration, EnvType } from './interfaces/configuration';

// eslint-disable-next-line @typescript-eslint/naming-convention
const __dirname = dirname(fileURLToPath(import.meta.url));
const throttleBucket: { [ip: string]: number[] } = {};
const oneWeek = 7 * 24 * 3_600_000; // 3600000msec == 1hour
const defaultContentSecurityPolicyOptions: Readonly<ContentSecurityPolicyOptions> = {
  useDefaults: true,
  directives: {
    'script-src': [
      "'self'",
      "'unsafe-eval'",
      "'sha256-9Ts7nfXdJQSKqVPxtB4Jwhf9pXSA/krLvgk8JROkI6g='", // script to set base-href inside index.html
    ],
    'style-src': [`'self'`, `'unsafe-inline'`, 'https://fonts.googleapis.com'],
    'upgrade-insecure-requests': null,
  },
};

const log = new ComponentLogger('App');

async function init(configuration: BridgeConfiguration): Promise<Express> {
  const app = express();
  const serverFeatureFlags = new ServerFeatureFlags();
  const clientFeatureFlags = new ClientFeatureFlags();
  EnvironmentUtils.setFeatureFlags(process.env, serverFeatureFlags);
  EnvironmentUtils.setFeatureFlags(process.env, clientFeatureFlags);

  const { mode, urls } = configuration;

  const rootFolder = join(__dirname, mode === EnvType.TEST ? '../' : '../../../');

  if (mode !== EnvType.TEST) {
    setupDefaultLookAndFeel(mode, rootFolder);
  }
  if (urls.lookAndFeel) {
    setupLookAndFeel(urls.lookAndFeel, rootFolder);
  }

  // UI static files - Angular application
  app.use(
    express.static(join(rootFolder, 'dist'), {
      maxAge: oneWeek, // cache files for one week
      etag: true, // Just being explicit about the default.
      lastModified: true, // Just being explicit about the default.
      setHeaders: (res: Response, path: string) => {
        // however, do not cache .html files (e.g., index.html)
        if (path.endsWith('.html')) {
          res.setHeader('Cache-Control', 'no-cache');
        }
      },
    })
  );

  // add some middlewares
  const logExpress = new ComponentLogger('Express');
  app.use((req: Request, res: Response, next: NextFunction) => {
    logExpress.info(`${req.method} ${req.url} ${res.statusCode} :: ${logExpress.prettyPrint(req.rawHeaders)}`);
    next();
  });
  app.use(express.json());
  app.use(express.urlencoded({ extended: false }));
  app.use(cookieParser());
  // OAUTH requires special security policies which are set later in ./user/oauth.ts
  if (!serverFeatureFlags.OAUTH_ENABLED) {
    app.use(contentSecurityPolicy(defaultContentSecurityPolicyOptions));
  }
  app.use(noSniff());
  app.use(permittedCrossDomainPolicies());
  app.use(frameguard());
  app.use(xssFilter());
  // Remove the X-Powered-By headers, has to be done via express and not helmet
  app.disable('x-powered-by');

  const { authType, session } = await setAuth(app, configuration);

  // everything starting with /api is routed to the api implementation
  app.use(
    '/api',
    apiRouter({
      authType,
      clientFeatureFlags,
      session,
      configuration,
    })
  );

  // fallback: go to index.html
  app.use((req, res) => {
    log.error('Not found: ' + req.url);
    res.sendFile(join(rootFolder, 'dist/index.html'), { maxAge: 0 });
  });

  // error handler
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  app.use((err: Error | AxiosError, req: Request, res: Response, _next: NextFunction) => {
    const status: number = handleError(err, req, res, authType);
    res.status(status).send(err.message);
  });

  return app;
}

async function setBasicAUTH(app: Express, auth: AuthConfig): Promise<void> {
  log.warning('Installing Basic authentication - please check environment variables!');

  setInterval(cleanIpBuckets, auth.cleanBucketIntervalMs);

  app.use((req, res, next) => {
    // parse login and password from headers
    const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
    const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');
    let userIP = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
    userIP = userIP instanceof Array ? userIP[0] : userIP;

    if (userIP && isIPThrottled(userIP, auth)) {
      log.error(`Request limit reached for IP ${userIP}`);
      res.status(429).send('Reached request limit');
      return;
    } else if (
      // if username and password are not set or wrong
      !(login && password && login === auth.basicUsername && password === auth.basicPassword)
    ) {
      updateBucket(!!(login || password), auth, userIP);

      log.error(`Access denied for IP: ${userIP}`);
      res.set('WWW-Authenticate', 'Basic realm="Keptn"');
      next({ response: { status: 401 } });
      return;
    }

    // Access granted
    return next();
  });
}

async function setAuth(
  app: Express,
  configuration: BridgeConfiguration
): Promise<{ authType: AuthType; session?: SessionService }> {
  let authType: AuthType;
  let session: SessionService | undefined;
  if (configuration.oauth.enabled) {
    session = await setupOAuth(app, configuration);
    authType = AuthType.OAUTH;
  } else if (configuration.auth.basicUsername && configuration.auth.basicPassword) {
    authType = AuthType.BASIC;
    await setBasicAUTH(app, configuration.auth);
  } else {
    authType = AuthType.NONE;
    log.warning('No authentication middleware is installed');
  }

  return { authType, session };
}

function setupDefaultLookAndFeel(mode: EnvType, rootFolder: string): void {
  try {
    log.debug('Installing default Look-and-Feel');

    const destDir = join(rootFolder, 'dist/assets/branding');
    const srcDir = join(rootFolder, `${mode === EnvType.DEV ? 'client' : 'dist'}/assets/default-branding`);
    const brandingFiles = ['app-config.json', 'logo.png', 'logo_inverted.png'];

    if (!existsSync(destDir)) {
      mkdirSync(destDir, { recursive: true });
    }

    brandingFiles.forEach((file) => {
      copyFileSync(join(srcDir, file), join(destDir, file));
    });
  } catch (e) {
    log.error(`Error while downloading custom Look-and-Feel file. Cause : ${e}`);
    process.exit(1);
  }
}

function setupLookAndFeel(url: string, rootFolder: string): void {
  let fl: WriteStream | undefined;

  try {
    log.debug(`Downloading custom Look-and-Feel file from ${url}`);

    const destDir = join(rootFolder, 'dist/assets/branding');
    const destFile = join(destDir, '/lookandfeel.zip');

    if (!existsSync(destDir)) {
      mkdirSync(destDir, { recursive: true });
    }

    fl = createWriteStream(destFile);
    const file: WriteStream = fl;
    const parsedUrl = new URL(url);
    const lib = parsedUrl.protocol === 'https:' ? https : http;

    lib
      .get(url, async (response) => {
        response.pipe(file);
        file.on('finish', () => {
          file.end();
          try {
            const zip = new AdmZip(destFile); // throws an error if unsupported format
            zip.extractAllToAsync(destDir, true, false, (error?: Error) => {
              unlinkSync(destFile);
              if (error) {
                log.error(`Error while extracting custom Look-and-Feel file: ${error}`);
                return;
              }
              log.info('Custom Look-and-Feel downloaded and extracted successfully');
            });
          } catch (error) {
            log.error(`Error while extracting custom Look-and-Feel file: ${error}`);
          }
        });
        file.on('error', async (err) => {
          file.end();
          try {
            await unlink(destFile);
          } catch (error) {
            log.error(`Error while saving custom Look-and-Feel file. ${error}`);
          }
          log.error(`Error while saving custom Look-and-Feel file. ${err}`);
        });
      })
      .on('error', (err) => {
        file.end();
        log.error(`Error while downloading custom Look-and-Feel file. ${err}`);
      });
  } catch (err) {
    fl?.end();
    log.error(`Error while downloading custom Look-and-Feel file. ${err}`);
  }
}

function updateBucket(loginAttempt: boolean, authConfig: AuthConfig, userIP?: string): void {
  // only fill buckets if the user tries to login
  if (userIP && loginAttempt) {
    if (!throttleBucket[userIP]) {
      throttleBucket[userIP] = [];
    }
    throttleBucket[userIP].push(new Date().getTime());

    // delete old requests. Just keep the latest {requestLimitWithinTime} requests
    if (throttleBucket[userIP].length > authConfig.nRequestWithinTime) {
      throttleBucket[userIP].shift();
    }
  }
}

/**
 *
 * @param ip. The IP of the request
 * @returns true if there are more than {requestLimitWithinTime} requests
 *          and the difference between first and last request of an IP is less than {requestTimeLimit}
 */
function isIPThrottled(ip: string, authConfig: AuthConfig): boolean {
  const ipBucket = throttleBucket[ip];
  return (
    ipBucket &&
    ipBucket.length >= authConfig.nRequestWithinTime &&
    new Date().getTime() - ipBucket[0] <= authConfig.requestTimeLimitMs
  );
}

/**
 * Delete an IP from the bucket if the last request is older than {requestTimeLimit}
 */
function cleanIpBuckets(authConfig: AuthConfig): void {
  for (const ip of Object.keys(throttleBucket)) {
    const ipBucket = throttleBucket[ip];
    if (ipBucket && new Date().getTime() - ipBucket[ipBucket.length - 1] > authConfig.requestTimeLimitMs) {
      delete throttleBucket[ip];
    }
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function handleError(err: any, req: Request, res: Response, authType: AuthType): number {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  if (err.response?.data?.message) {
    err.message = err.response?.data.message;
  }
  if (err.response?.status === 401) {
    res.setHeader('keptn-auth-type', authType);
  }

  printError(err);
  log.error(`Response status ${err.response?.status} for ${req.method} ${req.url}`);

  return err.response?.status || 500;
}

export { init, defaultContentSecurityPolicyOptions };
