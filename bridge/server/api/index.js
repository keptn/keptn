const express = require('express');
const axios = require('axios');
const https = require('https');

const router = express.Router();

module.exports = (params) => {
  const { datastoreService, configurationService, apiUrl } = params;

  // accepts self-signed ssl certificate
  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  router.get('/', async (req, res, next) => {
    try {
      return res.json({
        version: process.env.VERSION
      });
    } catch (err) {
      return next(err);
    }
  });

  router.get('/events', async (req, res, next) => {
    try {
      const traces = await datastoreService.getEvents(req.query);
      return res.json(traces);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/traces/:contextId', async (req, res, next) => {
    try {
      const traces = await datastoreService.getTraces(req.params.contextId, req.query.fromTime);
      return res.json(traces);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/roots/:projectName/:serviceName', async (req, res, next) => {
    try {
      const roots = await datastoreService.getRoots(req.params.projectName, req.params.serviceName, req.query.fromTime);
      return res.json(roots);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/project', async (req, res, next) => {
    try {
      const projects = await configurationService.getProjects();
      return res.json(projects);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/project/:projectName/resource', async (req, res, next) => {
    try {
      const resources = await configurationService.getProjectResources(req.params.projectName);
      return res.json(resources);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/project/:projectName/stage', async (req, res, next) => {
    try {
      const stages = await configurationService.getStages(req.params.projectName);
      return res.json(stages);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/project/:projectName/stage/:stageName/service', async (req, res, next) => {
    try {
      const services = await configurationService.getServices(req.params.projectName, req.params.stageName);
      return res.json(services);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/swagger-ui/swagger.yaml', async (req, res, next) => {
    try {
      const result = await axios({
        method: req.method,
        url: `${apiUrl}${req.url}`,
        data: req.params,
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
        data: req.params,
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

  return router;
};
