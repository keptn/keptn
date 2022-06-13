import { copyFileSync, createWriteStream, existsSync, mkdirSync, unlinkSync, WriteStream } from 'fs';
import { dirname, join } from 'path';
import { unlink } from 'fs/promises';
import * as https from 'https';
import * as http from 'http';
import { contentSecurityPolicy, frameguard, noSniff, permittedCrossDomainPolicies, xssFilter } from 'helmet';
import express, { Express, NextFunction, Request, Response } from 'express';
import { fileURLToPath, URL } from 'url';
import logger from 'morgan';
import { ComponentLogger } from './utils/logger';
import cookieParser from 'cookie-parser';
import AdmZip from 'adm-zip';
import { apiRouter } from './api';
import { execSync } from 'child_process';
import { AxiosError } from 'axios';
import { EnvironmentUtils } from './utils/environment.utils';
import { ClientFeatureFlags, ServerFeatureFlags } from './feature-flags';
import { setupOAuth } from './user/oauth';
import { SessionService } from './user/session';
import { ContentSecurityPolicyOptions } from 'helmet/dist/types/middlewares/content-security-policy';
import { printError } from './utils/print-utils';
import { AuthType } from '../shared/models/auth-type';

// eslint-disable-next-line @typescript-eslint/naming-convention
const __dirname = dirname(fileURLToPath(import.meta.url));
const apiUrl: string | undefined = process.env.API_URL;
let apiToken: string | undefined = process.env.API_TOKEN;
let cliDownloadLink: string | undefined = process.env.CLI_DOWNLOAD_LINK;
let integrationsPageLink: string | undefined = process.env.INTEGRATIONS_PAGE_LINK;
const lookAndFeelUrl: string | undefined = process.env.LOOK_AND_FEEL_URL;
const requestTimeLimit = +(process.env.REQUEST_TIME_LIMIT || 60) * 60 * 1000; // x minutes
const requestsWithinTime = +(process.env.REQUESTS_WITHIN_TIME || 10); // x requests within {requestTimeLimit}
const cleanBucketsInterval = +(process.env.CLEAN_BUCKET_INTERVAL || 60) * 60 * 1000; // clean buckets every x minutes
const throttleBucket: { [ip: string]: number[] } = {};
const rootFolder = join(__dirname, process.env.NODE_ENV === 'test' ? '../' : '../../../');
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

async function init(): Promise<Express> {
  const app = express();
  const serverFeatureFlags = new ServerFeatureFlags();
  const clientFeatureFlags = new ClientFeatureFlags();
  EnvironmentUtils.setFeatureFlags(process.env, serverFeatureFlags);
  EnvironmentUtils.setFeatureFlags(process.env, clientFeatureFlags);

  if (process.env.NODE_ENV !== 'test') {
    setupDefaultLookAndFeel();
  }
  if (lookAndFeelUrl) {
    setupLookAndFeel(lookAndFeelUrl);
  }
  if (!apiUrl) {
    throw Error('API_URL is not provided');
  }
  if (!apiToken) {
    log.warning('API_TOKEN was not provided. Fetching from kubectl.');
    apiToken =
      Buffer.from(
        execSync('kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token}').toString(),
        'base64'
      ).toString() || undefined;
  }

  if (!cliDownloadLink) {
    log.warning('CLI Download Link was not provided, defaulting to github.com/keptn/keptn releases');
    cliDownloadLink = 'https://github.com/keptn/keptn/releases';
  }

  if (!integrationsPageLink) {
    log.warning('Integrations page Link was not provided, defaulting to get.keptn.sh/integrations.html');
    integrationsPageLink = 'https://get.keptn.sh/integrations.html';
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
  app.use(logger('dev'));
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

  const { authType, session } = await setAuth(app, serverFeatureFlags.OAUTH_ENABLED);

  // everything starting with /api is routed to the api implementation
  app.use(
    '/api',
    apiRouter({
      apiUrl,
      apiToken,
      cliDownloadLink,
      integrationsPageLink,
      authType,
      clientFeatureFlags,
      session,
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

async function setOAUTH(app: Express): Promise<SessionService> {
  const errorSuffix =
    'must be defined when OAuth based login (OAUTH_ENABLED) is activated.' +
    ' Please check your environment variables.';

  if (!process.env.OAUTH_DISCOVERY) {
    throw Error(`OAUTH_DISCOVERY ${errorSuffix}`);
  }
  if (!process.env.OAUTH_CLIENT_ID) {
    throw Error(`OAUTH_CLIENT_ID ${errorSuffix}`);
  }
  if (!process.env.OAUTH_BASE_URL) {
    throw Error(`OAUTH_BASE_URL ${errorSuffix}`);
  }

  return setupOAuth(app, process.env.OAUTH_DISCOVERY, process.env.OAUTH_CLIENT_ID, process.env.OAUTH_BASE_URL);
}

async function setBasicAUTH(app: Express): Promise<void> {
  log.error('Installing Basic authentication - please check environment variables!');

  setInterval(cleanIpBuckets, cleanBucketsInterval);

  app.use((req, res, next) => {
    // parse login and password from headers
    const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
    const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');
    let userIP = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
    userIP = userIP instanceof Array ? userIP[0] : userIP;

    if (userIP && isIPThrottled(userIP)) {
      log.error('Request limit reached');
      res.status(429).send('Reached request limit');
      return;
    } else if (
      // if username and password are not set or wrong
      !(login && password && login === process.env.BASIC_AUTH_USERNAME && password === process.env.BASIC_AUTH_PASSWORD)
    ) {
      updateBucket(!!(login || password), userIP);

      log.error('Access denied');
      res.set('WWW-Authenticate', 'Basic realm="Keptn"');
      next({ response: { status: 401 } });
      return;
    }

    // Access granted
    return next();
  });
}

async function setAuth(app: Express, oAuthEnabled: boolean): Promise<{ authType: AuthType; session?: SessionService }> {
  let authType: AuthType;
  let session: SessionService | undefined;
  if (oAuthEnabled) {
    session = await setOAUTH(app);
    authType = AuthType.OAUTH;
  } else if (process.env.BASIC_AUTH_USERNAME && process.env.BASIC_AUTH_PASSWORD) {
    authType = AuthType.BASIC;
    await setBasicAUTH(app);
  } else {
    authType = AuthType.NONE;
    log.info('Not installing authentication middleware');
  }

  return { authType, session };
}

function setupDefaultLookAndFeel(): void {
  try {
    log.info('Installing default Look-and-Feel');

    const destDir = join(rootFolder, 'dist/assets/branding');
    const srcDir = join(
      rootFolder,
      `${process.env.NODE_ENV === 'development' ? 'client' : 'dist'}/assets/default-branding`
    );
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

function setupLookAndFeel(url: string): void {
  let fl: WriteStream | undefined;

  try {
    log.info(`Downloading custom Look-and-Feel file from ${lookAndFeelUrl}`);

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
                log.error(`[ERROR] Error while extracting custom Look-and-Feel file. ${error}`);
                return;
              }
              log.info('Custom Look-and-Feel downloaded and extracted successfully');
            });
          } catch (error) {
            log.error(`[ERROR] Error while extracting custom Look-and-Feel file. ${error}`);
          }
        });
        file.on('error', async (err) => {
          file.end();
          try {
            await unlink(destFile);
          } catch (error) {
            log.error(`[ERROR] Error while saving custom Look-and-Feel file. ${error}`);
          }
          log.error(`[ERROR] Error while saving custom Look-and-Feel file. ${err}`);
        });
      })
      .on('error', (err) => {
        file.end();
        log.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
      });
  } catch (err) {
    fl?.end();
    log.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
  }
}

function updateBucket(loginAttempt: boolean, userIP?: string): void {
  // only fill buckets if the user tries to login
  if (userIP && loginAttempt) {
    if (!throttleBucket[userIP]) {
      throttleBucket[userIP] = [];
    }
    throttleBucket[userIP].push(new Date().getTime());

    // delete old requests. Just keep the latest {requestLimitWithinTime} requests
    if (throttleBucket[userIP].length > requestsWithinTime) {
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
function isIPThrottled(ip: string): boolean {
  const ipBucket = throttleBucket[ip];
  return ipBucket && ipBucket.length >= requestsWithinTime && new Date().getTime() - ipBucket[0] <= requestTimeLimit;
}

/**
 * Delete an IP from the bucket if the last request is older than {requestTimeLimit}
 */
function cleanIpBuckets(): void {
  for (const ip of Object.keys(throttleBucket)) {
    const ipBucket = throttleBucket[ip];
    if (ipBucket && new Date().getTime() - ipBucket[ipBucket.length - 1] > requestTimeLimit) {
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

  return err.response?.status || 500;
}

export { init, defaultContentSecurityPolicyOptions };
