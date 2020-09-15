const express = require('express');
const axios = require('axios');
const https = require('https');

const router = express.Router();

module.exports = (params) => {
  // fetch parameters for bridgeInfo endpoint
  const { apiUrl, apiToken, cliDownloadLink } = params;
  const bridgeVersion = process.env.VERSION;

  // accepts self-signed ssl certificate
  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  // bridgeInfo endpoint: Provide certain metadata for Bridge
  router.get('/bridgeInfo', async (req, res, next) => {
    try {
      return res.json({ bridgeVersion, apiUrl, apiToken, cliDownloadLink });
    } catch (err) {
      return next(err);
    }
  });

  router.get('/swagger-ui/swagger.yaml', async (req, res, next) => {
    try {
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
      const result = await axios({
        method: req.method,
        url: `${apiUrl}${req.url}`,
        ...req.method!='GET' && { data: req.body },
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
};
