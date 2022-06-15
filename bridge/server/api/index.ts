import { NextFunction, Request, Response, Router } from 'express';
import { AxiosError, Method } from 'axios';
import { axios } from '../services/axios-instance';
import { DataService } from '../services/data-service';
import { KeptnInfoResult } from '../../shared/interfaces/keptn-info-result';
import { EnvironmentUtils } from '../utils/environment.utils';
import { ClientFeatureFlags } from '../feature-flags';
import { EventTypes } from '../../shared/interfaces/event-types';
import { SessionService } from '../user/session';
import { KeptnService } from '../../shared/models/keptn-service';
import { AuthType } from '../../shared/models/auth-type';
import { KeptnVersions } from '../../shared/interfaces/keptn-versions';
import { printError } from '../utils/print-utils';
import { ComponentLogger } from '../utils/logger';

const router = Router();
const log = new ComponentLogger('APIService');

const apiRouter = (params: {
  apiUrl: string;
  apiToken: string | undefined;
  cliDownloadLink: string;
  integrationsPageLink: string;
  authType: AuthType;
  clientFeatureFlags: ClientFeatureFlags;
  session: SessionService | undefined;
}): Router => {
  // fetch parameters for bridgeInfo endpoint
  const {
    apiUrl,
    apiToken,
    cliDownloadLink,
    integrationsPageLink,
    authType,
    clientFeatureFlags: featureFlags,
    session,
  } = params;
  const enableVersionCheckFeature = process.env.ENABLE_VERSION_CHECK !== 'false';
  const showApiToken = process.env.SHOW_API_TOKEN !== 'false';
  const bridgeVersion = process.env.VERSION;
  const projectsPageSize = EnvironmentUtils.getNumber(process.env.PROJECTS_PAGE_SIZE);
  const servicesPageSize = EnvironmentUtils.getNumber(process.env.SERVICES_PAGE_SIZE);
  const keptnInstallationType = process.env.KEPTN_INSTALLATION_TYPE;
  const automaticProvisioningMsg = process.env.AUTOMATIC_PROVISIONING_MSG?.trim();
  const authMsg = process.env.AUTH_MSG;
  const dataService = new DataService(apiUrl, apiToken);

  // bridgeInfo endpoint: Provide certain metadata for Bridge
  router.get('/bridgeInfo', async (req, res, next) => {
    const user = session?.getCurrentPrincipal(req);
    const bridgeInfo: KeptnInfoResult = {
      bridgeVersion,
      featureFlags,
      keptnInstallationType,
      apiUrl,
      ...(showApiToken && { apiToken }),
      cliDownloadLink,
      enableVersionCheckFeature,
      showApiToken,
      projectsPageSize,
      servicesPageSize,
      authType,
      ...(user && { user }),
      automaticProvisioningMsg,
      authMsg,
    };

    try {
      return res.json(bridgeInfo);
    } catch (err) {
      return next(err);
    }
  });

  router.get('/integrationsPage', async (req, res, next) => {
    try {
      const result = await axios({
        method: req.method as Method,
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

  router.get('/version.json', async (req, res) => {
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
      const defaultVersions: KeptnVersions = {
        cli: { stable: [], prerelease: [] },
        bridge: { stable: [], prerelease: [] },
        keptn: { stable: [] },
      };
      printError(err as AxiosError);
      return res.json(defaultVersions);
    }
  });

  router.post('/intersectEvents', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { event, eventSuffix, projectName, stages, services } = req.body;
      if (event && eventSuffix && projectName) {
        const result = await dataService.intersectEvents(
          req.session?.tokenSet?.access_token,
          event,
          eventSuffix,
          projectName,
          stages ?? [],
          services ?? []
        );
        return res.json(result);
      } else {
        log.error(`There has been a problem with /intersetEvents. Got body: ${req.body}`);
        return res.status(400).json('incorrect data');
      }
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const projectName = req.params.projectName;
      const project = await dataService.getProject(
        req.session?.tokenSet?.access_token,
        projectName,
        req.query.remediation === 'true',
        req.query.approval === 'true'
      );
      return res.json(project);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/serviceStates', async (req, res, next) => {
    try {
      const serviceStates = await dataService.getServiceStates(
        req.session?.tokenSet?.access_token,
        req.params.projectName
      );
      return res.json(serviceStates);
    } catch (err) {
      next(err);
    }
  });

  router.get('/project/:projectName/deployment/:keptnContext', async (req, res, next) => {
    try {
      const deployment = await dataService.getServiceDeployment(
        req.session?.tokenSet?.access_token,
        req.params.projectName,
        req.params.keptnContext,
        req.query.fromTime?.toString()
      );
      return res.json(deployment);
    } catch (err) {
      next(err);
    }
  });

  router.get('/project/:projectName/tasks', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const projectName = req.params.projectName;
      const tasks = await dataService.getTasks(req.session?.tokenSet?.access_token, projectName);
      return res.json(tasks);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/services', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const services = await dataService.getServiceNames(req.session?.tokenSet?.access_token, req.params.projectName);
      return res.json(services);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/customSequences', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const customSequences = await dataService.getCustomSequenceNames(
        req.session?.tokenSet?.access_token,
        req.params.projectName
      );
      return res.json(customSequences);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/project/:projectName/sequences/filter', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const data = await dataService.getSequencesFilter(req.session?.tokenSet?.access_token, req.params.projectName);
      return res.json(data);
    } catch (error) {
      return next(error);
    }
  });

  router.post('/uniform/registration', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const uniformDates: { [key: string]: string } = req.body;
      const uniformRegistrations = await dataService.getUniformRegistrations(
        req.session?.tokenSet?.access_token,
        uniformDates
      );
      return res.json(uniformRegistrations);
    } catch (error) {
      return next(error);
    }
  });

  router.get(
    '/uniform/registration/webhook-service/config/:subscriptionId',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const projectName = req.query.projectName?.toString();
        if (projectName) {
          const webhookConfig = await dataService.getWebhookConfig(
            req.session?.tokenSet?.access_token,
            req.params.subscriptionId,
            projectName,
            req.query.stageName?.toString(),
            req.query.serviceName?.toString()
          );
          return res.json(webhookConfig);
        } else {
          next(Error('project name not provided'));
        }
      } catch (error) {
        return next(error);
      }
    }
  );

  router.delete(
    '/uniform/registration/:integrationId/subscription/:subscriptionId',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const integrationId = req.params.integrationId;
        const subscriptionId = req.params.subscriptionId;
        const deleteWebhook = req.query.isWebhookService === 'true';
        if (integrationId && subscriptionId) {
          await dataService.deleteSubscription(
            req.session?.tokenSet?.access_token,
            integrationId,
            subscriptionId,
            deleteWebhook
          );
        } else {
          log.info('No available subscription or integration ID.');
        }
        return res.json();
      } catch (error) {
        return next(error);
      }
    }
  );

  router.post(
    '/uniform/registration/:integrationId/subscription',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const integrationId = req.params.integrationId;
        const subscription = req.body.subscription;
        if (integrationId && subscription) {
          await dataService.createSubscription(
            req.session?.tokenSet?.access_token,
            integrationId,
            subscription,
            req.body.webhookConfig
          );
        } else {
          log.info('No available subscription or integration ID.');
        }
        return res.json();
      } catch (error) {
        return next(error);
      }
    }
  );

  router.put(
    '/uniform/registration/:integrationId/subscription/:subscriptionId',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const integrationId = req.params.integrationId;
        const subscriptionId = req.params.subscriptionId;
        if (integrationId && subscriptionId) {
          await dataService.updateSubscription(
            req.session?.tokenSet?.access_token,
            integrationId,
            subscriptionId,
            req.body.subscription,
            req.body.webhookConfig
          );
        } else {
          log.info('No available subscription or integration ID.');
        }
        return res.json();
      } catch (error) {
        return next(error);
      }
    }
  );

  router.get('/uniform/registration/:integrationId/info', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const info = await dataService.getIsUniformRegistrationInfo(
        req.session?.tokenSet?.access_token,
        req.params.integrationId
      );
      return res.json(info);
    } catch (error) {
      return next(error);
    }
  });

  router.get(
    '/project/:projectName/service/:serviceName/files',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const serviceResources = await dataService.getResourceFileTreesForService(
          req.session?.tokenSet?.access_token,
          req.params.projectName,
          req.params.serviceName
        );
        return res.json(serviceResources);
      } catch (error) {
        return next(error);
      }
    }
  );

  router.get(
    '/project/:projectName/service/:serviceName/openRemediations',
    async (req: Request, res: Response, next: NextFunction) => {
      try {
        const serviceRemediationInformation = await dataService.getServiceRemediationInformation(
          req.session?.tokenSet?.access_token,
          req.params.projectName,
          req.params.serviceName,
          req.query.config?.toString() === 'true'
        );

        return res.json(serviceRemediationInformation);
      } catch (error) {
        return next(error);
      }
    }
  );

  router.post('/hasUnreadUniformRegistrationLogs', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const uniformDates: { [key: string]: string } = req.body;
      const status = await dataService.hasUnreadUniformRegistrationLogs(
        req.session?.tokenSet?.access_token,
        uniformDates
      );
      res.json(status);
    } catch (error) {
      return next(error);
    }
  });

  router.get('/mongodb-datastore/event', async (req: Request, res: Response, next: NextFunction) => {
    try {
      if (req.query.root === 'true') {
        const response = await dataService.getRoots(
          req.session?.tokenSet?.access_token,
          req.query.project?.toString(),
          req.query.pageSize?.toString(),
          req.query.serviceName?.toString(),
          req.query.fromTime?.toString(),
          req.query.beforeTime?.toString(),
          req.query.keptnContext?.toString()
        );
        return res.json(response);
      } else if (req.query.keptnContext && !req.query.pageSize) {
        const response = await dataService.getTracesByContext(
          req.session?.tokenSet?.access_token,
          req.query.keptnContext.toString(),
          req.query.project?.toString(),
          req.query.fromTime?.toString(),
          req.query.type?.toString() as EventTypes | undefined,
          req.query.source?.toString() as KeptnService | undefined
        );
        return res.json(response);
      } else {
        const response = await dataService.getTraces(
          req.session?.tokenSet?.access_token,
          req.query.keptnContext?.toString(),
          req.query.project?.toString(),
          req.query.stage?.toString(),
          req.query.service?.toString(),
          req.query.type?.toString() as EventTypes | undefined,
          req.query.source?.toString() as KeptnService | undefined,
          req.query.pageSize ? parseInt(req.query.pageSize.toString(), 10) : undefined
        );
        return res.json(response);
      }
    } catch (error) {
      return next(error);
    }
  });

  router.get('/secrets/scope/:scope', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const response = await dataService.getSecretsForScope(req.session?.tokenSet?.access_token, req.params.scope);
      return res.json(response);
    } catch (error) {
      return next(error);
    }
  });

  router.all('*', async (req, res, next) => {
    try {
      const accessToken = req.session?.tokenSet?.access_token;
      const result = await axios({
        method: req.method as Method,
        url: `${apiUrl}${req.url}`,
        ...(req.method !== 'GET' && { data: req.body }),
        headers: {
          ...(apiToken && { 'x-token': apiToken }),
          ...(accessToken && { Authorization: `Bearer ${accessToken}` }),
          'Content-Type': 'application/json',
        },
      });

      return res.json(result.data);
    } catch (err) {
      return next(err);
    }
  });

  return router;
};

export { apiRouter };
