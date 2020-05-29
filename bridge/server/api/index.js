const express = require('express');
const axios = require('axios');
const https = require('https');

const router = express.Router();

module.exports = (params) => {
  const { apiUrl, apiToken } = params;
  const version = process.env.VERSION;

  // accepts self-signed ssl certificate
  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  router.get('/', async (req, res, next) => {
    try {
      return res.json({ version, apiUrl });
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
        data: req.params,
        headers: {
          'x-token': apiToken,
          'content-type': 'application/json'
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
