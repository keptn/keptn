const express = require('express');

const router = express.Router();

module.exports = (params) => {
  const { datastoreService, configurationService } = params;

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

  return router;
};
