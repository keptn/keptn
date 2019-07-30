const express = require('express');

const router = express.Router();

module.exports = (params) => {
  const { elasticService } = params;

  router.get('/roots', async (req, res, next) => {
    try {
      const roots = await elasticService.getRoots();
      return res.json(roots);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/roots/:contextId', async (req, res, next) => {
    try {
      const roots = await elasticService.findRoots(req.params.contextId);
      return res.json(roots);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/traces/:contextId', async (req, res, next) => {
    try {
      const traces = await elasticService.getTraces(req.params.contextId);
      return res.json(traces);
    } catch (err) {
      return next(err);
    }
  });
  return router;
};
