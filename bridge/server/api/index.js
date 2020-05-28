const express = require('express');
const axios = require('axios');
const https = require('https');

const router = express.Router();

module.exports = (params) => {
  const { apiUrl } = params;

  // accepts self-signed ssl certificate
  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  router.get('/', async (req, res, next) => {
    try {
      return res.json({
        version: process.env.VERSION,
        apiUrl: process.env.API_URL
      });
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
      let method = req.method;
      let url = `${apiUrl}${req.url.replace('/api/', '')}`;
      let data = req.params;
      const result = await axios({ method, url, data });
      return res.json(result);
    } catch (err) {
      return next(err);
    }
  });

  return router;
};
