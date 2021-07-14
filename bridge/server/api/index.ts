import { Router } from 'express';
import axios from 'axios';
import * as https from 'https';
import { currentPrincipal } from '../user/session.js';

const router = Router();

function apiRouter(params:
                     {apiUrl: string | undefined, apiToken: string, cliDownloadLink: string, integrationsPageLink: string, authType: string}
                   ): Router {
  // fetch parameters for bridgeInfo endpoint
  const { apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType } = params;
  const enableVersionCheckFeature = process.env.ENABLE_VERSION_CHECK !== 'false';
  const showApiToken = process.env.SHOW_API_TOKEN !== 'false';
  const bridgeVersion = process.env.VERSION;
  const projectsPageSize = process.env.PROJECTS_PAGE_SIZE;
  const servicesPageSize = process.env.SERVICES_PAGE_SIZE;
  const keptnInstallationType = process.env.KEPTN_INSTALLATION_TYPE;

  // accepts self-signed ssl certificate
  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  // bridgeInfo endpoint: Provide certain metadata for Bridge
  router.get('/bridgeInfo', async (req, res, next) => {
    const user = currentPrincipal(req);
    const bridgeInfo = {
      bridgeVersion,
      keptnInstallationType,
      apiUrl, ...showApiToken && {apiToken},
      cliDownloadLink,
      enableVersionCheckFeature,
      showApiToken,
      projectsPageSize,
      servicesPageSize,
      authType,
      ... user && {user}
    };

    try {
      return res.json(bridgeInfo);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/integrationsPage', async (req, res, next) => {
    try {
      // @ts-ignore
      const result = await axios({
        method: req.method,
        url: `${integrationsPageLink}`,
        httpsAgent: agent
      });
      return res.send(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/swagger-ui/swagger.yaml', async (req, res, next) => {
    try {
      // @ts-ignore
      const result = await axios({
        method: req.method,
        url: `${apiUrl}${req.url}`,
        headers: {
          'Content-Type': 'application/json'
        },
        httpsAgent: agent
      });
      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/version.json', async (req, res, next) => {
    try {
      // @ts-ignore
      const result = await axios({
        method: req.method,
        url: `https://get.keptn.sh/version.json`,
        headers: {
          'Content-Type': 'application/json',
          'User-Agent': `keptn/bridge:${process.env.VERSION}`
        },
        httpsAgent: agent
      });
      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.all('*', async (req, res, next) => {
    try {
      // @ts-ignore
      const result = await axios({
        method: req.method,
        url: `${apiUrl}${req.url}`,
        ...req.method !== 'GET' && { data: req.body },
        headers: {
          'x-token': apiToken,
          'Content-Type': 'application/json'
        },
        httpsAgent: agent
      });

      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  return router;
}

export { apiRouter };
