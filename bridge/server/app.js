const fs = require('fs');
const path = require('path');
const http = require('http');
const https = require('https');
const urlParser = require('url');
const logger = require('morgan');
const express = require('express');
const cookieParser = require('cookie-parser');
const {execSync} = require('child_process');
const admZip = require('adm-zip');
const helmet = require("helmet");
const apiRouter = require('./api');

const app = express();
let apiUrl = process.env.API_URL;
let apiToken = process.env.API_TOKEN;
let cliDownloadLink = process.env.CLI_DOWNLOAD_LINK;
let integrationsPageLink = process.env.INTEGRATIONS_PAGE_LINK;
let lookAndFeelUrl = process.env.LOOK_AND_FEEL_URL;
const requestTimeLimit = (+process.env.REQUEST_TIME_LIMIT || 60) * 60 * 1000; // x minutes
const requestsWithinTime = +process.env.REQUESTS_WITHIN_TIME || 10; // x requests within {requestTimeLimit}
const cleanBucketsInterval = (+process.env.CLEAN_BUCKET_INTERVAL || 60) * 60 * 1000; // clean buckets every x minutes
const throttleBucket /* {[ip: string]: number[]} */ = {};

if(!apiToken) {
  console.log("API_TOKEN was not provided. Fetching from kubectl.");
  apiToken = Buffer.from(execSync('kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token}').toString(), 'base64').toString();
}

if (!cliDownloadLink) {
  console.log("CLI Download Link was not provided, defaulting to github.com/keptn/keptn releases")
  cliDownloadLink = "https://github.com/keptn/keptn/releases";
}

if(!integrationsPageLink) {
  console.log("Integrations page Link was not provided, defaulting to get.keptn.sh/integrations.html")
  integrationsPageLink = "https://get.keptn.sh/integrations.html";
}


try {
  console.log("Installing default Look-and-Feel");

  const destDir = path.join(__dirname, '../dist/assets/branding');
  const srcDir = path.join(__dirname, '../client/assets/default-branding');
  const brandingFiles = ["app-config.json", "logo.png", "logo_inverted.png"];

  if(!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
  }

  brandingFiles.forEach((file) => {
    fs.copyFileSync(path.join(srcDir, file), path.join(destDir, file));
  });
} catch (e) {
  console.error(`Error while downloading custom Look-and-Feel file. Cause : ${e}`);
  process.exit(1);
}
if(lookAndFeelUrl) {
  let file;

  try {
    console.log("Downloading custom Look-and-Feel file from", lookAndFeelUrl);

    const destDir = path.join(__dirname, '../dist/assets/branding');
    const destFile = path.join(destDir, '/lookandfeel.zip');

    if(!fs.existsSync(destDir)) {
      fs.mkdirSync(destDir, { recursive: true });
    }

    file = fs.createWriteStream(destFile);
    let parsedUrl = urlParser.parse(lookAndFeelUrl);
    let lib = parsedUrl.protocol === "https:" ? https : http;

    lib.get(lookAndFeelUrl, (response) => {
      response.pipe(file);
      file.on('finish', () => {
        file.close(() => {
          try {
            let zip = new admZip(destFile);
            zip.extractAllToAsync(destDir, true, () => {
              fs.unlinkSync(destFile);
              console.log("Custom Look-and-Feel downloaded and extracted successfully");
            });
          } catch (err) {
            console.error(`[ERROR] Error while extracting custom Look-and-Feel file. ${err}`);
          }
        });
      });
      file.on("error", (err) => {
        file.end();
        fs.unlink(destFile, () => {
          console.error(`[ERROR] Error while saving custom Look-and-Feel file. ${err}`);
        });
      });
    }).on('error', (err) => {
      file.end();
      console.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
    });
  } catch (err) {
    file.end();
    console.error(`[ERROR] Error while downloading custom Look-and-Feel file. ${err}`);
  }
}

const oneWeek = 7 * 24 * 3600000;    // 3600000msec == 1hour

module.exports = (async () => {
  // server static files - Images & CSS
  app.use('/static', express.static(path.join(__dirname, 'views/static'), {maxAge: oneWeek}));

  // UI static files - Angular application
  app.use(express.static(path.join(__dirname, '../dist'), {
    maxAge: oneWeek, // cache files for one week
    etag: true, // Just being explicit about the default.
    lastModified: true,  // Just being explicit about the default.
    setHeaders: (res, path) => {
      // however, do not cache .html files (e.g., index.html)
      if (path.endsWith('.html')) {
        res.setHeader('Cache-Control', 'no-cache');
      }
    },
  })
  );

  // Server views based on Pug
  app.set('views', path.join(__dirname, 'views'));
  app.set('view engine', 'pug');

  // add some middlewares
  app.use(logger('dev'));
  app.use(express.json());
  app.use(express.urlencoded({extended: false}));
  app.use(cookieParser());
  app.use(helmet.frameguard());

  let authType;

  if (process.env.OAUTH_ENABLED === 'true') {
    const sessionRouter = require('./user/session')(app);
    const oauthRouter = await require('./user/oauth');
    const authCheck = require('./user/session').isAuthenticated;

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

    console.error("Installing Basic authentication - please check environment variables!");

    setInterval(cleanIpBuckets, cleanBucketsInterval);

    app.use((req, res, next) => {
      // parse login and password from headers
      const b64auth = (req.headers.authorization || '').split(' ')[1] || '';
      const [login, password] = Buffer.from(b64auth, 'base64').toString().split(':');
      const userIP = (req.headers['x-forwarded-for'] || req.connection.remoteAddress)?.split(',')[0].trim();

      if (isIPThrottled(userIP)) {
        console.error("Request limit reached");
        res.status(429).send('Reached request limit');
        return;
      }
      // Verify login and password are set and correct
      else if (!(login && password && login === process.env.BASIC_AUTH_USERNAME && password === process.env.BASIC_AUTH_PASSWORD)) {
        // Access denied

        // only fill buckets if the user tries to login
        if (login || password) {
          if (!throttleBucket[userIP]) {
            throttleBucket[userIP] = [];
          }
          throttleBucket[userIP].push(new Date().getTime());

          // delete old requests. Just keep the latest {requestLimitWithinTime} requests
          if (throttleBucket[userIP].length > requestsWithinTime) {
            throttleBucket[userIP].shift();
          }
        }

        console.error("Access denied");
        res.set('WWW-Authenticate', 'Basic realm="Keptn"');
        res.status(401).send('Authentication required.'); // custom message
        return;
      }

      // Access granted
      return next();
    });
  } else {
    authType = 'NONE';
    console.error("Not installing authentication middleware");
  }


// everything starting with /api is routed to the api implementation
  app.use('/api', apiRouter({apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType}));

// fallback: go to index.html
  app.use((req, res, next) => {
    console.error("Not found: " + req.url);
    res.sendFile(path.join(`${__dirname}/../dist/index.html`), {maxAge: 0});
  });

// error handler
// eslint-disable-next-line no-unused-vars
  app.use((err, req, res, next) => {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};
    // render the error page
    if(err.response && err.response.data && err.response.data.message)
      err.message = err.response.data.message;
    res.status(err.status || 500).send(err.message);
    console.error(err);
  });

  return app;
})();

/**
 *
 * @param ip. The IP of the request
 * @returns true if there are more than {requestLimitWithinTime} requests and the difference between first and last request of an IP is less than {requestTimeLimit}
 */
function isIPThrottled(ip) {
  const ipBucket = throttleBucket[ip];
  return ipBucket && ipBucket.length >= requestsWithinTime && (new Date().getTime() - ipBucket[0]) <= requestTimeLimit;
}

/**
 * Delete an IP from the bucket if the last request is older than {requestTimeLimit}
 */
function cleanIpBuckets() {
  for (const ip in throttleBucket) {
    const ipBucket = throttleBucket[ip];
    if (ipBucket && (new Date().getTime() - ipBucket[ipBucket.length - 1]) > requestTimeLimit)
    {
      delete throttleBucket[ip];
    }
  }
}
