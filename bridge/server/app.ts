import { copyFileSync, createWriteStream, existsSync, mkdirSync, unlinkSync, WriteStream } from 'fs';
import { dirname, join } from 'path';
import { unlink } from 'fs/promises';
import * as https from 'https';
import * as http from 'http';
import helmet from 'helmet';
import express, { Express, NextFunction, Request, Response } from 'express';
import { fileURLToPath, URL } from 'url';
import logger from 'morgan';
import cookieParser from 'cookie-parser';
import admZip from 'adm-zip';
import { apiRouter } from './api';
import { execSync } from 'child_process';

// tslint:disable-next-line:variable-name whitespace
const __dirname = dirname(fileURLToPath(import.meta.url));
const app = express();
const apiUrl: string | undefined = process.env.API_URL;
let apiToken: string | undefined = process.env.API_TOKEN;
let cliDownloadLink: string | undefined = process.env.CLI_DOWNLOAD_LINK;
let integrationsPageLink: string | undefined = process.env.INTEGRATIONS_PAGE_LINK;
const lookAndFeelUrl: string | undefined = process.env.LOOK_AND_FEEL_URL;
const requestTimeLimit = (+(process.env.REQUEST_TIME_LIMIT || 60)) * 60 * 1000; // x minutes
const requestsWithinTime = +(process.env.REQUESTS_WITHIN_TIME || 10); // x requests within {requestTimeLimit}
const cleanBucketsInterval = (+(process.env.CLEAN_BUCKET_INTERVAL || 60)) * 60 * 1000; // clean buckets every x minutes
const throttleBucket: { [ip: string]: number[] } = {};
const rootFolder = join(__dirname, '../../../');
const serverFolder = join(rootFolder, 'server');

try {
  console.log('Installing default Look-and-Feel');

  const destDir = join(rootFolder, 'dist/assets/branding');
  const srcDir = join(rootFolder, `${process.env.NODE_ENV === 'development' ? 'client' : 'dist'}/assets/default-branding`);
  const brandingFiles = ['app-config.json', 'logo.png', 'logo_inverted.png'];

  if (!existsSync(destDir)) {
    mkdirSync(destDir, {recursive: true});
  }

  brandingFiles.forEach((file) => {
    copyFileSync(join(srcDir, file), join(destDir, file));
  });
} catch (e) {
  console.error(`Error while downloading custom Look-and-Feel file. Cause : ${e}`);
  process.exit(1);
}
if (lookAndFeelUrl) {
  let file: WriteStream;

  try {
    console.log('Downloading custom Look-and-Feel file from', lookAndFeelUrl);

    const destDir = join(rootFolder, 'dist/assets/branding');
    const destFile = join(destDir, '/lookandfeel.zip');

    if (!existsSync(destDir)) {
      mkdirSync(destDir, {recursive: true});
    }

    file = createWriteStream(destFile);
    const parsedUrl = new URL(lookAndFeelUrl);
    const lib = parsedUrl.protocol === 'https:' ? https : http;

    lib.get(lookAndFeelUrl, async (response) => {
      response.pipe(file);
      file.on('finish', () => {
        file.end();
        try {
          const zip = new admZip(destFile);
          zip.extractAllToAsync(destDir, true, () => {
            unlinkSync(destFile);
            console.log('Custom Look-and-Feel downloaded and extracted successfully');
          });
        } catch (err) {
          console.error(`[ERROR] Error while extracting custom Look-and-Feel file. ${err}`);
        }
      });
      file.on('error', async (err) => {
        file.end();
        try {
          await unlink(destFile);
        } catch (error) {
          console.error(`[ERROR] Error while saving custom Look-and-Feel file. ${error}`);
        }
        console.error(`[ERROR] Error while saving custom Look-and-Feel file. ${err}`);
      });
    }).on('error', (err) => {
      file.end();
      console.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
    });
  } catch (err) {
    // @ts-ignore
    file?.end();
    console.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
  }
}

const oneWeek = 7 * 24 * 3_600_000;    // 3600000msec == 1hour

async function init(): Promise<Express> {
  if (!apiUrl) {
    throw Error('API_URL is not provided');
  }
  if (!apiToken) {
    console.log('API_TOKEN was not provided. Fetching from kubectl.');
    apiToken = Buffer.from(execSync('kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token}')
      .toString(), 'base64').toString();
  }

  if (!cliDownloadLink) {
    console.log('CLI Download Link was not provided, defaulting to github.com/keptn/keptn releases');
    cliDownloadLink = 'https://github.com/keptn/keptn/releases';
  }

  if (!integrationsPageLink) {
    console.log('Integrations page Link was not provided, defaulting to get.keptn.sh/integrations.html');
    integrationsPageLink = 'https://get.keptn.sh/integrations.html';
  }

  // server static files - Images & CSS
  app.use('/static', express.static(join(serverFolder, 'views/static'), {maxAge: oneWeek}));

  // UI static files - Angular application
  app.use(express.static(join(rootFolder, 'dist'), {
      maxAge: oneWeek, // cache files for one week
      etag: true, // Just being explicit about the default.
      lastModified: true,  // Just being explicit about the default.
      setHeaders: (res: Response, path: string) => {
        // however, do not cache .html files (e.g., index.html)
        if (path.endsWith('.html')) {
          res.setHeader('Cache-Control', 'no-cache');
        }
      },
    }),
  );

  // Server views based on Pug
  app.set('views', join(serverFolder, 'views'));
  app.set('view engine', 'pug');

  // add some middlewares
  app.use(logger('dev'));
  app.use(express.json());
  app.use(express.urlencoded({extended: false}));
  app.use(cookieParser());
  app.use(helmet.frameguard());

  const authType: string = await setAuth();

// everything starting with /api is routed to the api implementation
  app.use('/api', apiRouter({apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType}));

// fallback: go to index.html
  app.use((req, res) => {
    console.error('Not found: ' + req.url);
    res.sendFile(join(rootFolder, 'dist/index.html'), {maxAge: 0});
  });

// error handler
  // tslint:disable-next-line:no-any
  app.use((err: any, req: Request, res: Response, _next: NextFunction) => {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};
    // render the error page
    if (err.response?.data?.message) {
      err.message = err.response.data.message;
    }
    if (err.response?.status === 401) {
      res.setHeader('keptn-auth-type', authType);
    }

    res.status(err.response?.status || 500).send(err.message);
    console.error(err);
  });

  return app;
}

async function setOAUTH(): Promise<void> {
  const sessionRouter = (await import('./user/session.js')).sessionRouter(app);
  const oauthRouter = await (await import('./user/oauth.js')).oauthRouter;
  const authCheck = (await import('./user/session.js')).isAuthenticated;

  // Initialise session middleware
  app.use(sessionRouter);
  // Initializing OAuth middleware.
  app.use(oauthRouter);

  // Authentication filter for API requests
  app.use('/api', (req, resp, next) => {
    if (!authCheck(req)) {
      resp.status(401).send('Unauthorized');
      return;
    }
    return next();
  });
}

async function setBasisAUTH(): Promise<void> {
  console.error('Installing Basic authentication - please check environment variables!');

  setInterval(cleanIpBuckets, cleanBucketsInterval);

  app.use((req, res, next) => {
    // parse login and password from headers
    const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
    const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');
    let userIP = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
    userIP = userIP instanceof Array ? userIP[0] : userIP;

    if (userIP && isIPThrottled(userIP)) {
      console.error('Request limit reached');
      res.status(429).send('Reached request limit');
      return;
    } else if (!(login && password && login === process.env.BASIC_AUTH_USERNAME && password === process.env.BASIC_AUTH_PASSWORD)) {
      updateBucket(!!(login || password), userIP);

      console.error('Access denied');
      res.set('WWW-Authenticate', 'Basic realm="Keptn"');
      res.status(401).send('Authentication required.'); // custom message
      return;
    }

    // Access granted
    return next();
  });
}

async function setAuth(): Promise<string> {
  let authType;
  if (process.env.OAUTH_ENABLED === 'true') {
    await setOAUTH();
    authType = 'OAUTH';

  } else if (process.env.BASIC_AUTH_USERNAME && process.env.BASIC_AUTH_PASSWORD) {
    authType = 'BASIC';
    await setBasisAUTH();
  } else {
    authType = 'NONE';
    console.error('Not installing authentication middleware');
  }

  return authType;
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
  return ipBucket && ipBucket.length >= requestsWithinTime && (new Date().getTime() - ipBucket[0]) <= requestTimeLimit;
}

/**
 * Delete an IP from the bucket if the last request is older than {requestTimeLimit}
 */
function cleanIpBuckets(): void {
  for (const ip of Object.keys(throttleBucket)) {
    const ipBucket = throttleBucket[ip];
    if (ipBucket && (new Date().getTime() - ipBucket[ipBucket.length - 1]) > requestTimeLimit) {
      delete throttleBucket[ip];
    }
  }
}

export { init };
