import { existsSync, mkdirSync, copyFileSync, createWriteStream, unlinkSync, WriteStream } from 'fs';
import { join, dirname } from 'path';
import { unlink } from 'fs/promises';
import * as https from 'https';
import * as http from 'http';
import helmet from 'helmet';
import express, { Express, NextFunction, Request, Response } from 'express';
import { parse, fileURLToPath } from 'url';
import logger from 'morgan';
import cookieParser from 'cookie-parser';
import { execSync } from 'child_process';
import admZip from 'adm-zip';
import { apiRouter } from './api/index.js';

// tslint:disable-next-line:variable-name whitespace
const __dirname = dirname(fileURLToPath(import.meta.url));
const app = express();
const apiUrl: string | undefined = process.env.API_URL;
let apiToken: string | undefined = process.env.API_TOKEN;
let cliDownloadLink: string | undefined = process.env.CLI_DOWNLOAD_LINK;
let integrationsPageLink: string | undefined = process.env.INTEGRATIONS_PAGE_LINK;
const lookAndFeelUrl: string | undefined = process.env.LOOK_AND_FEEL_URL;

try {
  console.log('Installing default Look-and-Feel');

  const destDir = join(__dirname, '../../dist/assets/branding');
  const srcDir = join(__dirname, '../../client/assets/default-branding');
  const brandingFiles = ['app-config.json', 'logo.png', 'logo_inverted.png'];

  if (!existsSync(destDir)) {
    mkdirSync(destDir, { recursive: true });
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

    const destDir = join(__dirname, '../../dist/assets/branding');
    const destFile = join(destDir, '/lookandfeel.zip');

    if (!existsSync(destDir)) {
      mkdirSync(destDir, { recursive: true });
    }

    file = createWriteStream(destFile);
    const parsedUrl = parse(lookAndFeelUrl);
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
        }
        catch (err) {
          console.error(`[ERROR] Error while extracting custom Look-and-Feel file. ${err}`);
        }
      });
      file.on('error', async (err) => {
        file.end();
        try {
          await unlink(destFile);
        }
        catch (err){
          console.error(`[ERROR] Error while saving custom Look-and-Feel file. ${err}`);
        }
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
  app.use('/static', express.static(join(__dirname, '../views/static'), {maxAge: oneWeek}));

  // UI static files - Angular application
  app.use(express.static(join(__dirname, '../../dist'), {
      maxAge: oneWeek, // cache files for one week
      etag: true, // Just being explicit about the default.
      lastModified: true,  // Just being explicit about the default.
      setHeaders: (res: Response, path: string) => {
        // however, do not cache .html files (e.g., index.html)
        if (path.endsWith('.html')) {
          res.setHeader('Cache-Control', 'no-cache');
        }
      },
    })
  );

  // Server views based on Pug
  app.set('views', join(__dirname, '../views'));
  app.set('view engine', 'pug');

  // add some middlewares
  app.use(logger('dev'));
  app.use(express.json());
  app.use(express.urlencoded({extended: false}));
  app.use(cookieParser());
  app.use(helmet.frameguard());

  let authType: string;

  if (process.env.OAUTH_ENABLED === 'true') {
    const sessionRouter = (await import('./user/session.js')).sessionRouter(app);
    const oauthRouter = await (await import('./user/oauth.js')).oauthRouter;
    const authCheck = (await import('./user/session.js')).isAuthenticated;

    authType = 'OAUTH';

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

  } else if (process.env.BASIC_AUTH_USERNAME && process.env.BASIC_AUTH_PASSWORD) {
    authType = 'BASIC';

    console.error('Installing Basic authentication - please check environment variables!');
    app.use((req, res, next) => {
      // parse login and password from headers
      const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
      const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');

      // Verify login and password are set and correct
      if (!(login && password && login === process.env.BASIC_AUTH_USERNAME && password === process.env.BASIC_AUTH_PASSWORD)) {
        // Access denied
        console.error('Access denied');
        res.set('WWW-Authenticate', 'Basic realm="Keptn"');
        res.status(401).send('Authentication required.'); // custom message
        return;
      }

      // Access granted
      return next();
    });
  } else {
    authType = 'NONE';
    console.error('Not installing authentication middleware');
  }


// everything starting with /api is routed to the api implementation
  app.use('/api', apiRouter({apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType}));

// fallback: go to index.html
  app.use((req, res, next) => {
    console.error('Not found: ' + req.url);
    res.sendFile(join(`${__dirname}/../../dist/index.html`), {maxAge: 0});
  });

// error handler
  app.use((err: any, req: Request, res: Response, _next: NextFunction) => {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};
    // render the error page
    if (err.response?.data?.message) {
      err.message = err.response.data.message;
    }
    res.status(err.status || 500).send(err.message);
    console.error(err);
  });

  return app;
}

export { init };
