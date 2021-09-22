import { NextFunction, Request, Response, Router } from 'express';
import { Method } from 'axios';
import { currentPrincipal } from '../user/session';
import { axios } from '../services/axios-instance';
import { DataService } from '../services/data-service';
import { WebhookConfig } from '../../shared/interfaces/webhook-config';

const router = Router();

function apiRouter(params:
                     { apiUrl: string, apiToken: string, cliDownloadLink: string, integrationsPageLink: string, authType: string },
): Router {
  // fetch parameters for bridgeInfo endpoint
  const {apiUrl, apiToken, cliDownloadLink, integrationsPageLink, authType} = params;
  const enableVersionCheckFeature = process.env.ENABLE_VERSION_CHECK !== 'false';
  const showApiToken = process.env.SHOW_API_TOKEN !== 'false';
  const bridgeVersion = process.env.VERSION;
  const projectsPageSize = process.env.PROJECTS_PAGE_SIZE;
  const servicesPageSize = process.env.SERVICES_PAGE_SIZE;
  const keptnInstallationType = process.env.KEPTN_INSTALLATION_TYPE;
  const dataService = new DataService(apiUrl, apiToken);

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
      ...user && {user},
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
      });
      return res.send(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/swagger-ui/swagger.yaml', async (req, res, next) => {
    try {
      const result = await axios({
        method: req.method as Method,
        url: `${apiUrl}${req.url}`,
        headers: {
          'Content-Type': 'application/json',
        },
      });
      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/version.json', async (req, res, next) => {
    try {
      const result = await axios({
        method: req.method as Method,
        url: `https://get.keptn.sh/version.json`,
        headers: {
          'Content-Type': 'application/json',
          'User-Agent': `keptn/bridge:${process.env.VERSION}`,
        },
      });
      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/project/:projectName', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const projectName = req.params.projectName;
      const project = await dataService.getProject(projectName, req.query.remediation === 'true', req.query.approval === 'true');
      return res.json(project);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/tasks', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const projectName = req.params.projectName;
      const tasks = await dataService.getTasks(projectName);
      return res.json(tasks);
    } catch (error) {
      return next(error);
    }
  });

  router.post('/uniform/registration', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const uniformDates: { [key: string]: string } = req.body;
      const uniformRegistrations = await dataService.getUniformRegistrations(uniformDates);
      return res.json(uniformRegistrations);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/uniform/registration/webhook-service/config/:eventType', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const projectName = req.query.projectName?.toString();
      if (projectName) {
        const webhookConfig = await dataService.getWebhookConfig(req.params.eventType, projectName, req.query.stageName?.toString(), req.query.serviceName?.toString());
        return res.json(webhookConfig);
      } else {
        next(Error('project name not provided'));
      }
    } catch (error) {
      return next(error);
    }
  });

  router.post('/uniform/registration/webhook-service/config', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const webhookConfig: WebhookConfig = req.body.config;
      const result = await dataService.saveWebhookConfig(webhookConfig);
      return res.json(result);
    } catch (error) {
      return next(error);
    }
  });

  router.delete('/uniform/registration/:integrationId/subscription/:subscriptionId', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const integrationId = req.params.integrationId;
      const subscriptionId = req.params.subscriptionId;
      const deleteWebhook = req.query.isWebhookService === 'true';
      if (integrationId && subscriptionId) {
        await dataService.deleteSubscription(integrationId, subscriptionId, deleteWebhook);
      }
      return res.json();
    } catch (error) {
      return next(error);
    }
  });

  router.get('/uniform/registration/:integrationId/info', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const info = await dataService.getIsUniformRegistrationInfo(req.params.integrationId);
      return res.json(info);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/service/:serviceName/files', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const serviceResources = await dataService.getResourceFileTreesForService(req.params.projectName, req.params.serviceName);
      return res.json(serviceResources);
    } catch (error) {
      return next(error);
    }
  });

  router.post('/hasUnreadUniformRegistrationLogs', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const uniformDates: { [key: string]: string } = req.body;
      const status = await dataService.hasUnreadUniformRegistrationLogs(uniformDates);
      res.json(status);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/mongodb-datastore/event', async (req: Request, res: Response, next: NextFunction) => {
    try {
      if (req.params.roots === 'true') {
        const response = await dataService.getRoots(req.query.project?.toString(), req.query.pageSize?.toString(), req.query.serviceName?.toString(), req.query.fromTime?.toString(), req.query.beforeTime?.toString(), req.query.keptnContext?.toString());
        return res.json(response);
      } else {
        const response = await dataService.getTracesByContext(req.query.keptnContext?.toString(), req.query.project?.toString(), req.query.fromTime?.toString());
        return res.json(response);
      }
    } catch (error) {
      return next(error);
    }
  });

  router.get('/secrets/scope/:scope', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const response = await dataService.getSecretsForScope(req.params.scope);
      return res.json(response);
    } catch (error) {
      return next(error);
    }
  });

  router.all('*', async (req, res, next) => {
    try {
      const result = await axios({
        method: req.method as Method,
        url: `${apiUrl}${req.url}`,
        ...req.method !== 'GET' && {data: req.body},
        headers: {
          'x-token': apiToken,
          'Content-Type': 'application/json',
        },
      });

      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  return router;
}

export { apiRouter };
