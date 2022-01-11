import { ApiService } from './api-service';
import { Sequence } from '../models/sequence';
import { SequenceTypes } from '../../shared/models/sequence-types';
import { Trace } from '../models/trace';
import { Service } from '../models/service';
import { Project } from '../models/project';
import { EventState } from '../../shared/models/event-state';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Approval } from '../interfaces/approval';
import { ResultTypes } from '../../shared/models/result-types';
import { UniformRegistration } from '../models/uniform-registration';
import Yaml from 'yaml';
import { Shipyard } from '../interfaces/shipyard';
import { UniformRegistrationLocations } from '../../shared/interfaces/uniform-registration-locations';
import { WebhookConfig, WebhookConfigFilter, WebhookSecret } from '../../shared/models/webhook-config';
import { UniformRegistrationInfo } from '../../shared/interfaces/uniform-registration-info';
import { WebhookConfigYaml } from '../interfaces/webhook-config-yaml';
import { UniformSubscription, UniformSubscriptionFilter } from '../../shared/interfaces/uniform-subscription';
import axios from 'axios';
import { Resource } from '../../shared/interfaces/resource';
import { FileTree, TreeEntry } from '../../shared/interfaces/resourceFileTree';
import { EventResult } from '../interfaces/event-result';
import { Secret } from '../models/secret';
import { IRemediationAction } from '../../shared/models/remediation-action';
import { SecretScope } from '../../shared/interfaces/secret-scope';
import { KeptnService } from '../../shared/models/keptn-service';
import { SequenceState } from '../../shared/models/sequence';
import { ServiceDeploymentInformation, ServiceState } from '../../shared/models/service-state';
import { Deployment, IStageDeployment, SubSequence } from '../../shared/interfaces/deployment';
import semver from 'semver';
import { ServiceRemediationInformation } from '../../shared/interfaces/service-remediation-information';
import { Stage } from '../models/stage';
import { IServiceEvent } from '../../shared/interfaces/service';
import { Remediation } from '../models/remediation';

type TreeDirectory = ({ _: string[] } & { [key: string]: TreeDirectory }) | { _: string[] };
type FlatSecret = { path: string; name: string; key: string; parsedPath: string };
type StageRemediationInformation = {
  remediations: Remediation[];
  remediationsForStage: Sequence[];
  config?: string;
};
type StageOpenInformation = {
  openApprovals: Approval[];
  openRemediations: Remediation[];
  evaluations: Trace[];
};

export class DataService {
  private apiService: ApiService;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;
  private readonly MAX_TRACE_PAGE_SIZE = 50;
  private readonly MAX_PAGE_SIZE = 100;

  constructor(apiUrl: string, apiToken: string) {
    this.apiService = new ApiService(apiUrl, apiToken);
  }

  public async getProjects(): Promise<Project[]> {
    const response = await this.apiService.getProjects();
    return response.data.projects.map((project) => Project.fromJSON(project));
  }

  private async getPlainProject(projectName: string): Promise<Project> {
    const response = await this.apiService.getProject(projectName);
    return Project.fromJSON(response.data);
  }

  public async getProject(
    projectName: string,
    includeRemediation: boolean,
    includeApproval: boolean
  ): Promise<Project> {
    const project = await this.getPlainProject(projectName);
    let openApprovals: Approval[] = [];
    let openRemediations: Remediation[] = [];
    const allServices = Stage.getAllServices(project.stages);
    const latestDeployments = await this.getLatestDeploymentFinishedForServices(
      projectName,
      allServices,
      ResultTypes.PASSED
    );
    const latestEvaluations = await this.getLatestEvaluationResultsForServices(projectName, allServices);
    // for sequence adjustment: const latestSequences = await this.getLatestSequenceForServices(projectName, allServices);
    const cachedSequences: { [keptnContext: string]: Sequence | undefined } = {};

    if (includeRemediation) {
      openRemediations = await this.getOpenRemediations(projectName, true);
    }

    if (includeApproval) {
      openApprovals = await this.getApprovals(projectName, true);
    }

    for (const stage of project.stages) {
      const stageRemediations = openRemediations.filter((seq) => seq.stages.some((s) => s.name === stage.stageName));
      const stageApprovals = openApprovals.filter((approval) => approval.trace.data.stage === stage.stageName);
      const stageEvaluations = latestEvaluations.filter((t) => t.data.stage === stage.stageName);
      const stageDeployments = latestDeployments.filter((t) => t.data.stage === stage.stageName);
      const stageInformation: StageOpenInformation = {
        openApprovals: stageApprovals,
        openRemediations: stageRemediations,
        evaluations: stageEvaluations,
      };

      for (const service of stage.services) {
        const latestEvent = service.getLatestEvent();
        if (latestEvent) {
          const latestSequence = await this.setServiceDetails(
            service,
            stage.stageName,
            latestEvent.keptnContext,
            projectName,
            stageInformation,
            stageDeployments,
            cachedSequences[latestEvent.keptnContext]
          );
          if (latestSequence) {
            cachedSequences[latestEvent.keptnContext] = latestSequence;
          }
        }
      }
    }
    return project;
  }

  private async getLatestSequenceForServices(projectName: string, services: Service[]): Promise<Sequence[]> {
    let sequences: Sequence[] = [];

    for (let i = 0; i < services.length; i += this.MAX_PAGE_SIZE) {
      const keptnContexts: string[] = [];
      const maxLength = Math.min(i + this.MAX_PAGE_SIZE, services.length);

      for (let y = i; y < maxLength; ++y) {
        const latestServiceEvent = services[y].getLatestEvent();
        if (latestServiceEvent) {
          // string concatenation is expensive; that's why we use an array here
          keptnContexts.push(latestServiceEvent.keptnContext);
        }
      }
      if (keptnContexts.length) {
        const response = await this.apiService.getSequences(
          projectName,
          this.MAX_PAGE_SIZE,
          undefined,
          undefined,
          undefined,
          undefined,
          keptnContexts.join(',')
        );
        sequences = [...sequences, ...response.data.states];
      }
    }
    return sequences;
  }

  private async getLatestDeploymentFinishedForServices(
    projectName: string,
    services: Service[],
    resultType?: ResultTypes
  ): Promise<Trace[]> {
    return this.getLatestTracesOfMultipleServices(
      projectName,
      services,
      EventTypes.DEPLOYMENT_FINISHED,
      (service: Service) => service.lastEventTypes[EventTypes.DEPLOYMENT_FINISHED],
      undefined,
      resultType
    );
  }

  private async getLatestEvaluationResultsForServices(projectName: string, services: Service[]): Promise<Trace[]> {
    return this.getLatestTracesOfMultipleServices(
      projectName,
      services,
      EventTypes.EVALUATION_FINISHED,
      (service: Service) => {
        const latestEvent = service.getLatestEvent();
        const latestEvaluation = service.lastEventTypes[EventTypes.EVALUATION_FINISHED];
        return latestEvent?.keptnContext === latestEvaluation?.keptnContext ? latestEvaluation : undefined;
      },
      KeptnService.LIGHTHOUSE_SERVICE
    );
  }

  /**
   * Fetches the latest event provided by latestEventTypes of a service.
   * If resultType is provided, the result is filtered and if it does not match, it fetches the one with the right result.
   * This will result in API-Calls:
   *  best case: O(1)
   *  worst case: O(N) (exceptional case for deployment.finished is fail)
   * @param projectName
   * @param services
   * @param eventType
   * @param getServiceEvent
   * @param source
   * @param resultType
   * @private
   */
  private async getLatestTracesOfMultipleServices(
    projectName: string,
    services: Service[],
    eventType: EventTypes,
    getServiceEvent: (service: Service) => IServiceEvent | undefined,
    source?: KeptnService,
    resultType?: ResultTypes
  ): Promise<Trace[]> {
    let traces: Trace[] = [];

    for (let i = 0; i < services.length; i += this.MAX_PAGE_SIZE) {
      const eventIds: string[] = [];
      const maxLength = Math.min(i + this.MAX_PAGE_SIZE, services.length);
      for (let y = i; y < maxLength; ++y) {
        const latestServiceEvent = getServiceEvent(services[y]);
        if (latestServiceEvent) {
          // string concatenation is expensive; that's why we use an array here
          eventIds.push(latestServiceEvent.eventId);
        }
      }
      if (eventIds.length) {
        const response = await this.apiService.getTracesOfMultipleServices(projectName, eventType, eventIds.join(','));
        if (resultType || source) {
          await this.checkAndSetEventsWithResult(response.data.events, resultType, source);
        }
        traces = [...traces, ...response.data.events];
      }
    }
    return traces;
  }

  private async checkAndSetEventsWithResult(
    traces: Trace[],
    resultType?: ResultTypes,
    source?: KeptnService
  ): Promise<void> {
    for (let i = 0; i < traces.length; ++i) {
      const trace = traces[i];
      if ((resultType && trace.data.result !== resultType) || (source && trace.source !== source)) {
        const response = await this.apiService.getTracesWithResultAndSource(
          trace.type as EventTypes,
          1,
          trace.data.project as string,
          trace.data.stage as string,
          trace.data.service as string,
          resultType,
          source
        );
        if (response.data.events.length) {
          traces[i] = response.data.events[0];
        }
      }
    }
  }

  private async setServiceDetails(
    service: Service,
    stageName: string,
    keptnContext: string,
    projectName: string,
    stageInformation: StageOpenInformation,
    stageDeployments: Trace[],
    cachedSequence: Sequence | undefined
  ): Promise<Sequence | undefined> {
    let latestSequence;
    try {
      latestSequence = cachedSequence ?? (await this.getSequence(projectName, stageName, keptnContext));
    } catch (error) {
      console.error(error);
      return;
    }
    service.latestSequence = latestSequence ? Sequence.fromJSON(latestSequence) : undefined;
    service.latestSequence?.reduceToStage(stageName);

    if (!cachedSequence && service.latestSequence) {
      const stage = service.latestSequence.stages.find((seq) => seq.name === stageName);
      if (stage) {
        stage.latestEvaluationTrace = stageInformation.evaluations.find(
          (t) => t.data.service === service.latestSequence?.service
        );
      }
    }
    const deployment = stageDeployments.find((t) => t.data.service === service.serviceName);
    if (deployment) {
      service.deploymentInformation = {
        deploymentUrl: Trace.fromJSON(deployment).getDeploymentUrl(),
        image: service.getShortImage(),
      };
    }

    const serviceRemediations = stageInformation.openRemediations.filter(
      (remediation) => remediation.service === service.serviceName
    );
    for (const remediation of serviceRemediations) {
      remediation.reduceToStage(stageName);
    }
    service.openRemediations = serviceRemediations;
    service.openApprovals = stageInformation.openApprovals.filter(
      (approval) => approval.trace.data.service === service.serviceName
    );
    return latestSequence;
  }

  public async getSequence(
    projectName: string,
    stageName?: string,
    keptnContext?: string
  ): Promise<Sequence | undefined> {
    const response = await this.apiService.getSequences(
      projectName,
      1,
      undefined,
      undefined,
      undefined,
      undefined,
      keptnContext
    );
    let sequence = response.data.states[0];
    if (sequence) {
      sequence = Sequence.fromJSON(sequence);
    }
    return sequence;
  }

  public async getSequences(
    projectName: string,
    sequenceName: string,
    sequenceState?: SequenceState
  ): Promise<Sequence[]> {
    const response = await this.apiService.getSequences(
      projectName,
      this.MAX_SEQUENCE_PAGE_SIZE,
      sequenceName,
      sequenceState
    );

    return response.data.states.map((sequence) => Sequence.fromJSON(sequence));
  }

  public async getOpenRemediations(
    projectName: string,
    includeActions: boolean,
    serviceName?: string
  ): Promise<Remediation[]> {
    const sequences = await this.getSequences(projectName, SequenceTypes.REMEDIATION, SequenceState.STARTED);
    const remediations: Remediation[] = [];
    for (const sequence of sequences) {
      const stageName = sequence.stages[0]?.name;
      // there could be invalid sequences that don't have a stage because the triggered sequence was not present in the shipyard file
      if (stageName && (!serviceName || sequence.service === serviceName)) {
        const stage = { ...sequence.stages[0], actions: [] };
        const remediation: Remediation = Remediation.fromJSON({ ...sequence, stages: [stage] });

        if (includeActions) {
          await this.loadRemediationActions(
            remediation,
            projectName,
            stageName,
            sequence.service,
            sequence.shkeptncontext
          );
        }
        remediations.push(remediation);
      }
    }
    return remediations;
  }

  private async loadRemediationActions(
    remediation: Remediation,
    projectName: string,
    stageName: string,
    serviceName: string,
    keptnContext: string
  ): Promise<void> {
    const response = await this.apiService.getTraces(
      undefined,
      this.MAX_TRACE_PAGE_SIZE,
      projectName,
      stageName,
      serviceName,
      keptnContext
    );
    const traces = response.data.events;
    const actions = this.getRemediationActions(Trace.traceMapper(traces));
    remediation.stages[0].actions.push(...actions);
  }

  private getRemediationActions(traces: Trace[]): IRemediationAction[] {
    const actions: IRemediationAction[] = [];
    for (const trace of traces) {
      const actionTriggeredTrace = trace.findTrace((t) => t.type === EventTypes.ACTION_TRIGGERED && !!t.data.action);
      if (actionTriggeredTrace && actionTriggeredTrace.data.action) {
        const finishedAction = actionTriggeredTrace.traces.find((t) =>
          t.findTrace((tt) => tt.type === EventTypes.ACTION_FINISHED)
        );
        const startedAction = actionTriggeredTrace.traces.find((t) =>
          t.findTrace((tt) => tt.type === EventTypes.ACTION_STARTED)
        );
        let state: EventState;
        if (finishedAction) {
          state = EventState.FINISHED;
        } else if (startedAction) {
          state = EventState.STARTED;
        } else {
          state = EventState.TRIGGERED;
        }

        actions.push({ ...actionTriggeredTrace.data.action, state, result: finishedAction?.data.result });
      }
    }
    return actions;
  }

  private async getTrace(
    keptnContext?: string,
    projectName?: string,
    stageName?: string,
    serviceName?: string,
    eventType?: EventTypes,
    eventSource?: KeptnService
  ): Promise<Trace | undefined> {
    const response = await this.apiService.getTraces(
      eventType,
      1,
      projectName,
      stageName,
      serviceName,
      keptnContext,
      eventSource
    );
    return response.data.events.shift();
  }

  public async getApprovals(
    projectName: string,
    includeEvaluationTrace: boolean,
    stageName?: string,
    serviceName?: string
  ): Promise<Approval[]> {
    let tracesTriggered: Trace[];
    try {
      const response = await this.apiService.getOpenTriggeredEvents(
        projectName,
        EventTypes.APPROVAL_TRIGGERED,
        stageName,
        serviceName
      );
      tracesTriggered = response.data.events ?? [];
    } catch {
      // status 500 if no events are found
      tracesTriggered = [];
    }
    const approvals: Approval[] = [];
    // for each approval the latest evaluation trace (before this event) is needed
    for (const trace of tracesTriggered) {
      let evaluationTrace: Trace | undefined;
      if (includeEvaluationTrace) {
        evaluationTrace = await this.getTrace(
          trace.shkeptncontext,
          projectName,
          stageName,
          serviceName,
          EventTypes.EVALUATION_FINISHED,
          KeptnService.LIGHTHOUSE_SERVICE
        );
      }
      approvals.push({
        evaluationTrace,
        trace,
      });
    }
    return approvals;
  }

  public async hasUnreadUniformRegistrationLogs(uniformDates: { [key: string]: string }): Promise<boolean> {
    const response = await this.apiService.getUniformRegistrations();
    const registrations = await this.getValidRegistrations(response.data);
    let status = false;
    for (let i = 0; i < registrations.length && !status; ++i) {
      const registration = registrations[i];
      const logResponse = await this.apiService.getUniformRegistrationLogs(
        registration.id,
        uniformDates[registration.id],
        1
      );
      if (logResponse.data.logs.length !== 0) {
        status = true;
      }
    }
    return status;
  }

  public async getUniformRegistrations(uniformDates: { [key: string]: string }): Promise<UniformRegistration[]> {
    const response = await this.apiService.getUniformRegistrations();
    const registrations = await this.getValidRegistrations(response.data);
    for (const registration of registrations) {
      const logResponse = await this.apiService.getUniformRegistrationLogs(
        registration.id,
        uniformDates[registration.id]
      );
      registration.unreadEventsCount = logResponse.data.logs.length;
    }
    return registrations;
  }

  private async getValidRegistrations(registrations: UniformRegistration[]): Promise<UniformRegistration[]> {
    const currentDate = new Date().getTime();
    const validRegistrations: UniformRegistration[] = [];
    for (const registration of registrations) {
      const diffMins = (currentDate - new Date(registration.metadata.lastseen).getTime()) / 60_000;
      if (diffMins < 2) {
        validRegistrations.push(registration);
      }
    }
    return validRegistrations;
  }

  public async getIsUniformRegistrationInfo(integrationId: string): Promise<UniformRegistrationInfo> {
    const response = await this.apiService.getUniformRegistrations(integrationId);
    const uniformRegistration = UniformRegistration.fromJSON(response.data.shift());

    return {
      isControlPlane: uniformRegistration.metadata.location === UniformRegistrationLocations.CONTROL_PLANE,
      isWebhookService: uniformRegistration.isWebhookService,
    };
  }

  public async getTasks(projectName: string): Promise<string[]> {
    const shipyard = await this.getShipyard(projectName);
    const tasks: string[] = ['evaluation'];
    for (const stage of shipyard.spec.stages) {
      if (stage.sequences) {
        for (const sequence of stage.sequences) {
          for (const task of sequence.tasks) {
            if (!tasks.includes(task.name)) {
              tasks.push(task.name);
            }
          }
        }
      }
    }
    return tasks;
  }

  public async getServiceNames(projectName: string): Promise<string[]> {
    const resp = await this.apiService.getStages(projectName);
    const stages = resp.data.stages;
    const services: { [serviceName: string]: boolean | undefined } = {};

    for (const stage of stages) {
      for (const service of stage.services) {
        services[service.serviceName] = true;
      }
    }

    return Object.keys(services);
  }

  public async getRoots(
    projectName: string | undefined,
    pageSize: string | undefined,
    serviceName: string | undefined,
    fromTime?: string | undefined,
    beforeTime?: string | undefined,
    keptnContext?: string | undefined
  ): Promise<EventResult> {
    const response = await this.apiService.getRoots(
      projectName,
      pageSize,
      serviceName,
      fromTime,
      beforeTime,
      keptnContext
    );
    return response.data;
  }

  public async getTracesByContext(
    keptnContext: string | undefined,
    projectName?: string | undefined,
    fromTime?: string | undefined
  ): Promise<EventResult> {
    let result: EventResult = {
      events: [],
      pageSize: 0,
      nextPageKey: 0,
      totalCount: 0,
    };
    let nextPage = 0;
    do {
      const response = await this.apiService.getTracesByContext(
        keptnContext,
        projectName,
        fromTime,
        nextPage.toString()
      );
      nextPage = response.data.nextPageKey || 0;
      result = {
        events: [...result?.events, ...response.data.events],
        pageSize: result.pageSize + response.data.pageSize,
        nextPageKey: response.data.nextPageKey,
        totalCount: response.data.totalCount,
      };
    } while (nextPage !== 0);

    return result;
  }

  public async getTraces(
    keptnContext?: string,
    projectName?: string,
    stageName?: string,
    serviceName?: string,
    eventType?: string,
    pageSize?: number
  ): Promise<EventResult> {
    const response = await this.apiService.getTraces(
      eventType,
      pageSize,
      projectName,
      stageName,
      serviceName,
      keptnContext
    );
    return response.data;
  }

  public async getResourceFileTreesForService(projectName: string, serviceName: string): Promise<FileTree[]> {
    const projectRes = await this.apiService.getProject(projectName);
    const stages = projectRes.data.stages;

    const fileTrees: FileTree[] = [];

    for (const stage of stages) {
      let nextPage: string | undefined;
      let resourceResponses: Resource[] = [];
      const fileTree: FileTree = {
        stageName: stage.stageName,
        tree: [],
      };

      do {
        const resourceRes = await this.apiService.getServiceResources(
          projectName,
          stage.stageName,
          serviceName,
          nextPage || undefined
        );
        nextPage = resourceRes.data.nextPageKey;
        resourceResponses = [...resourceResponses, ...resourceRes.data.resources];
      } while (parseInt(nextPage, 10) !== 0);

      fileTree.tree = this._getResourceFileTree(resourceResponses);
      fileTrees.push(fileTree);
    }

    return fileTrees;
  }

  private async getShipyard(projectName: string): Promise<Shipyard> {
    const response = await this.apiService.getShipyard(projectName);
    const shipyard = Buffer.from(response.data.resourceContent, 'base64').toString('utf-8');
    return Yaml.parse(shipyard);
  }

  public async createSubscription(
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: WebhookConfig
  ): Promise<void> {
    const response = await this.apiService.createSubscription(integrationId, subscription);
    if (webhookConfig) {
      await this.saveWebhookConfig(webhookConfig, response.data.id);
    }
  }

  public async updateSubscription(
    integrationId: string,
    subscriptionId: string,
    subscription: UniformSubscription,
    webhookConfig?: WebhookConfig
  ): Promise<void> {
    await this.apiService.updateSubscription(integrationId, subscriptionId, subscription);
    if (webhookConfig) {
      await this.saveWebhookConfig(webhookConfig, subscriptionId);
    }
  }

  private async getWebhookConfigYaml(
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<WebhookConfigYaml> {
    const response = await this.apiService.getWebhookConfig(projectName, stageName, serviceName);

    try {
      const webhookConfigFile = Buffer.from(response.data.resourceContent, 'base64').toString('utf-8');
      return WebhookConfigYaml.fromJSON(Yaml.parse(webhookConfigFile));
    } catch (error) {
      throw Error('Could not parse webhook.yaml');
    }
  }

  private replaceWithBridgeSecrets(webhookConfig: WebhookConfig): void {
    if (webhookConfig.secrets) {
      for (const webhookSecret of webhookConfig.secrets) {
        const bridgeSecret = '.secret.' + webhookSecret.secretRef.name + '.' + webhookSecret.secretRef.key;

        const regex = new RegExp('.env.' + webhookSecret.name, 'g');
        webhookConfig.url = webhookConfig.url.replace(regex, bridgeSecret);
        webhookConfig.payload = webhookConfig.payload.replace(regex, bridgeSecret);

        for (const header of webhookConfig.header) {
          header.value = header.value.replace(regex, bridgeSecret);
        }
      }
    }
  }

  public async getWebhookConfig(
    subscriptionId: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<WebhookConfig> {
    const webhookConfigYaml: WebhookConfigYaml = await this.getWebhookConfigYaml(projectName, stageName, serviceName);

    const webhookConfig = webhookConfigYaml.parsedRequest(subscriptionId);
    if (!webhookConfig) {
      throw Error('Could not parse curl command');
    }
    this.replaceWithBridgeSecrets(webhookConfig);
    return webhookConfig;
  }

  public async saveWebhookConfig(webhookConfig: WebhookConfig, subscriptionId: string): Promise<boolean> {
    const currentFilters = await this.getWebhookConfigFilter(webhookConfig.filter);

    if (webhookConfig.prevConfiguration) {
      const previousFilter = await this.getWebhookConfigFilter(webhookConfig.prevConfiguration.filter);
      await this.removePreviousWebhooks(previousFilter, subscriptionId);
    }

    const secrets = await this.parseAndReplaceWebhookSecret(webhookConfig);
    const curl = this.generateWebhookConfigCurl(webhookConfig);

    for (const project of currentFilters.projects) {
      for (const stage of currentFilters.stages) {
        for (const service of currentFilters.services) {
          const previousWebhookConfig: WebhookConfigYaml = await this.getOrCreateWebhookConfigYaml(
            project,
            stage,
            service
          );
          previousWebhookConfig.addWebhook(
            webhookConfig.type,
            curl,
            subscriptionId,
            secrets,
            webhookConfig.sendFinished
          );
          await this.apiService.saveWebhookConfig(previousWebhookConfig.toYAML(), project, stage, service);
        }
      }
    }

    return true;
  }

  private async parseAndReplaceWebhookSecret(webhookConfig: WebhookConfig): Promise<WebhookSecret[]> {
    const webhookScopeSecrets = await this.getSecretsForScope(SecretScope.WEBHOOK);
    const flatSecret = this.getSecretPathFlat(webhookScopeSecrets);

    const secrets: WebhookSecret[] = [];
    webhookConfig.url = this.addWebhookSecretsFromString(webhookConfig.url, flatSecret, secrets);
    webhookConfig.payload = this.addWebhookSecretsFromString(webhookConfig.payload, flatSecret, secrets);

    for (const head of webhookConfig.header) {
      head.value = this.addWebhookSecretsFromString(head.value, flatSecret, secrets);
    }

    return secrets;
  }

  private addWebhookSecretsFromString(
    parseString: string,
    allSecretPaths: FlatSecret[],
    existingSecrets: WebhookSecret[]
  ): string {
    const foundSecrets = allSecretPaths.filter((scrt) => parseString.includes(scrt.path));
    let replacedString = parseString;
    for (const found of foundSecrets) {
      const idx = existingSecrets.findIndex((secret) => secret.name === found.parsedPath);
      if (idx === -1) {
        const secret: WebhookSecret = {
          name: found.parsedPath,
          secretRef: {
            name: found.name,
            key: found.key,
          },
        };
        existingSecrets.push(secret);
      }

      replacedString = replacedString.replace(new RegExp('secret.' + found.path, 'g'), 'env.' + found.parsedPath);
    }

    return replacedString;
  }

  private getSecretPathFlat(secrets: Secret[]): FlatSecret[] {
    const flatScopeSecrets: FlatSecret[] = [];
    for (const secret of secrets) {
      if (secret.keys) {
        for (const key of secret.keys) {
          const sanitizedName = secret.name.replace(/[^a-zA-Z0-9]/g, '');
          const sanitizedKey = key.replace(/[^a-zA-Z0-9]/g, '');
          const flat: FlatSecret = {
            path: secret.name + '.' + key,
            name: secret.name,
            key,
            parsedPath: 'secret_' + sanitizedName + '_' + sanitizedKey,
          };
          flatScopeSecrets.push(flat);
        }
      }
    }
    return flatScopeSecrets;
  }
  private async getOrCreateWebhookConfigYaml(
    project: string,
    stage?: string,
    service?: string
  ): Promise<WebhookConfigYaml> {
    let previousWebhookConfig: WebhookConfigYaml;
    try {
      // fetch existing one
      previousWebhookConfig = await this.getWebhookConfigYaml(project, stage, service);
    } catch (e: unknown) {
      if (!axios.isAxiosError(e) || e.response?.status !== 404) {
        throw e;
      } else {
        // if it does not exist, create one
        previousWebhookConfig = new WebhookConfigYaml();
      }
    }
    return previousWebhookConfig;
  }

  private async removePreviousWebhooks(previousConfig: WebhookConfigFilter, subscriptionId: string): Promise<void> {
    for (const project of previousConfig.projects) {
      for (const stage of previousConfig.stages) {
        for (const service of previousConfig.services) {
          await this.removeWebhook(subscriptionId, project, stage, service);
        }
      }
    }
  }

  private async getWebhookConfigFilter(webhookConfig: UniformSubscriptionFilter): Promise<WebhookConfigFilter> {
    return {
      projects: webhookConfig.projects?.length
        ? webhookConfig.projects
        : (await this.getProjects()).map((project) => project.projectName),
      stages: webhookConfig.stages?.length ? webhookConfig.stages : [undefined],
      services: webhookConfig.services?.length ? webhookConfig.services : [undefined],
    };
  }

  private generateWebhookConfigCurl(webhookConfig: WebhookConfig): string {
    let params = '';
    for (const header of webhookConfig?.header || []) {
      params += `--header '${header.name}: ${header.value}' `;
    }
    params += `--request ${webhookConfig.method} `;
    if (webhookConfig.proxy) {
      params += `--proxy ${webhookConfig.proxy} `;
    }
    if (webhookConfig.payload) {
      let stringify = webhookConfig.payload;
      try {
        stringify = JSON.stringify(JSON.parse(webhookConfig.payload));
      } catch {
        stringify = stringify.replace(/\r\n|\n|\r/gm, '');
      }
      params += `--data '${stringify}' `;
    }
    return `curl ${params}${webhookConfig.url}`;
  }

  public async deleteSubscription(
    integrationId: string,
    subscriptionId: string,
    deleteWebhook: boolean
  ): Promise<void> {
    if (deleteWebhook) {
      const response = await this.apiService.getUniformSubscription(integrationId, subscriptionId);
      const subscription = response.data;
      const projectName = subscription.filter.projects?.[0];

      if (projectName) {
        await this.removeWebhooks(
          subscriptionId,
          projectName,
          subscription.filter.stages?.length ? subscription.filter.stages : [undefined],
          subscription.filter.services?.length ? subscription.filter.services : [undefined]
        );
      }
    }
    await this.apiService.deleteUniformSubscription(integrationId, subscriptionId);
  }

  public async getSecretsForScope(scope: string): Promise<Secret[]> {
    const response = await this.apiService.getSecrets();
    const secrets = response.data.Secrets.map((secret) => Secret.fromJSON(secret));
    return secrets.filter((secret) => secret.scope === scope);
  }

  private async removeWebhooks(
    subscriptionId: string,
    projectName: string,
    stages: string[] | [undefined],
    services: string[] | [undefined]
  ): Promise<void> {
    for (const stage of stages) {
      for (const service of services) {
        await this.removeWebhook(subscriptionId, projectName, stage, service);
      }
    }
  }

  private async removeWebhook(
    subscriptionId: string,
    projectName: string,
    stage?: string,
    service?: string
  ): Promise<void> {
    try {
      const webhookConfig: WebhookConfigYaml = await this.getWebhookConfigYaml(projectName, stage, service);
      if (webhookConfig.removeWebhook(subscriptionId)) {
        if (webhookConfig.hasWebhooks()) {
          await this.apiService.saveWebhookConfig(webhookConfig.toYAML(), projectName, stage, service);
        } else {
          await this.apiService.deleteWebhookConfig(projectName, stage, service);
        }
      }
    } catch (e: unknown) {
      // ignore if yaml was not found. e.g. on create
      if (!axios.isAxiosError(e) || e.response?.status !== 404) {
        throw e;
      }
    }
  }

  private _getResourceFileTree(resources: Resource[]): TreeEntry[] {
    const directory: TreeDirectory = { _: [] };

    for (const res of resources) {
      const parts = res.resourceURI.split('/').filter((item) => !!item);
      this._addToTreeDirectory(directory, parts);
    }

    return this._buildTree(directory, '').children || [];
  }

  private _addToTreeDirectory(currentDirectory: TreeDirectory, parts: string[]): void {
    let index = 0;
    for (; index < parts.length - 1; ++index) {
      const part = parts[index];
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      currentDirectory[part] ??= { _: [] };
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      currentDirectory = currentDirectory[part];
    }
    currentDirectory._.push(parts[index]);
  }

  private _buildTree(currentDirectory: TreeDirectory, fileName: string): TreeEntry {
    const tree: TreeEntry = { fileName, children: [] as TreeEntry[] };
    const dict: { [key: string]: boolean } = { _: true };
    const folders: TreeEntry[] = [];
    const files: TreeEntry[] = [];
    for (const key in currentDirectory) {
      if (dict[key]) {
        files.push(
          ...currentDirectory._.map(
            (item) =>
              ({
                fileName: item,
              } as TreeEntry)
          )
        );
      } else {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        folders.push(this._buildTree(currentDirectory[key] as TreeDirectory, key));
      }
    }
    folders.sort((a, b) => a.fileName.localeCompare(b.fileName));
    files.sort((a, b) => a.fileName.localeCompare(b.fileName));
    tree.children = [...folders, ...files];
    return tree;
  }

  public async getServiceStates(projectName: string): Promise<ServiceState[]> {
    const projectResponse = await this.apiService.getProject(projectName);
    const project = Project.fromJSON(projectResponse.data);
    const openRemediations = await this.getOpenRemediations(projectName, false);
    const serviceStates: ServiceState[] = [];
    for (const stage of project.stages) {
      for (const service of stage.services) {
        const latestDeploymentEvent = service.latestDeploymentEvent;
        let serviceState = serviceStates.find((s) => s.name === service.serviceName);
        if (!serviceState) {
          serviceState = new ServiceState(service.serviceName);
          serviceStates.push(serviceState);
        }

        if (latestDeploymentEvent) {
          const serviceRemediations = openRemediations.filter(
            (remediation) =>
              remediation.service === service.serviceName &&
              remediation.stages.some((remediationStage) => remediationStage.name === stage.stageName)
          );
          const deploymentInformation = this.getOrCreateDeploymentInformation(
            serviceState,
            service,
            latestDeploymentEvent.keptnContext
          );
          const deploymentStage = {
            name: stage.stageName,
            hasOpenRemediations: serviceRemediations.length !== 0,
            time: new Date(+latestDeploymentEvent.time / 1_000_000).toISOString(),
          };
          deploymentInformation.stages.push(deploymentStage);
        }
      }
    }
    this.sortServiceStates(serviceStates);
    return serviceStates;
  }

  private sortServiceStates(serviceStates: ServiceState[]): void {
    serviceStates.sort((a, b) => a.name.localeCompare(b.name));
    for (const serviceState of serviceStates) {
      serviceState.deploymentInformation.sort((a, b) =>
        a.version &&
        b.version &&
        semver.valid(a.version) != null &&
        semver.valid(b.version) != null &&
        semver.gt(a.version, b.version, true)
          ? -1
          : 1
      );
    }
  }

  private getOrCreateDeploymentInformation(
    serviceState: ServiceState,
    service: Service,
    latestDeploymentContext: string
  ): ServiceDeploymentInformation {
    let deploymentInformation = serviceState.deploymentInformation.find(
      (deployment) => deployment.keptnContext === latestDeploymentContext
    );
    if (!deploymentInformation) {
      deploymentInformation = {
        stages: [],
        name: service.serviceName,
        image: service.getShortImageName(),
        version: service.getImageVersion(),
        keptnContext: latestDeploymentContext,
      };
      serviceState.deploymentInformation.push(deploymentInformation);
    }
    return deploymentInformation;
  }

  public async getServiceDeployment(
    projectName: string,
    keptnContext: string,
    fromTimeString?: string
  ): Promise<Deployment | undefined> {
    const fromTime = fromTimeString ? new Date(fromTimeString) : undefined;
    const sequenceResponse = await this.apiService.getSequences(
      projectName,
      1,
      undefined,
      undefined,
      undefined,
      undefined,
      keptnContext
    );
    let deployment: Deployment | undefined;
    const iSequence = sequenceResponse.data.states[0];
    if (iSequence) {
      const sequence = Sequence.fromJSON(iSequence);
      let openRemediations: Remediation[] | undefined;
      const project = await this.getPlainProject(projectName);
      const traceResponse = await this.apiService.getTracesByContext(keptnContext, projectName);
      const traces = Trace.traceMapper(traceResponse.data.events);
      deployment = {
        state: sequence.state,
        keptnContext: sequence.shkeptncontext,
        service: sequence.service,
        stages: [],
        labels:
          traces[traces.length - 1]?.getFinishedEvent()?.data.labels ??
          traces[0]?.getFinishedEvent()?.data.labels ??
          {},
      };
      for (const stage of sequence.stages) {
        const stageTraces = traces.filter((trace) => trace.data.stage === stage.name);
        const service = project.stages
          .find((st) => st.stageName === stage.name)
          ?.services.find((sv) => sv.serviceName === sequence.service);
        const latestDeploymentContext = service?.deploymentEvent?.keptnContext;
        const evaluationTrace = stageTraces.reduce(
          (event: Trace | undefined, trace) => event || trace.getEvaluationFinishedEvent(),
          undefined
        );
        const lastTimeUpdated = stageTraces[stageTraces.length - 1]?.getLastTrace()?.time;
        const approvalInformation = this.getApprovalInformation(stageTraces, service?.getShortImage());
        let deploymentURL: string | undefined;
        let stageRemediationInformation: StageRemediationInformation | undefined;
        deployment.image ??= service?.getShortImage();

        if (latestDeploymentContext === sequence.shkeptncontext) {
          deploymentURL = stageTraces.reduce(
            (url: string | undefined, tc) =>
              url || tc.findTrace((t) => t.type === EventTypes.DEPLOYMENT_FINISHED)?.getDeploymentUrl(),
            undefined
          );
          stageRemediationInformation = await this.getStageRemediationInformation(
            projectName,
            stage.name,
            sequence.service,
            openRemediations
          );
          openRemediations = stageRemediationInformation.remediations;
        }

        deployment.stages.push({
          name: stage.name,
          state: stage.state,
          lastTimeUpdated: (lastTimeUpdated ? new Date(lastTimeUpdated) : new Date()).toISOString(),
          openRemediations: stageRemediationInformation?.remediationsForStage ?? [],
          remediationConfig: stageRemediationInformation?.config,
          approvalInformation,
          subSequences: this.getSubSequencesForStage(stageTraces, stage.name, fromTime),
          deploymentURL,
          hasEvaluation: !!evaluationTrace,
          latestEvaluation:
            evaluationTrace ??
            (await this.apiService.getEvaluationResults(projectName, sequence.service, stage.name, 1)).data.events[0],
        });
      }
    }
    return deployment;
  }

  private getApprovalInformation(traces: Trace[], deployedImage?: string): IStageDeployment['approvalInformation'] {
    let approvalInformation: IStageDeployment['approvalInformation'];
    const approvalTrace = traces.reduce(
      (approval: Trace | undefined, trace) => approval || trace.findTrace((t) => !!t.isApproval()),
      undefined
    );
    if (approvalTrace?.isApprovalPending()) {
      approvalTrace.traces = []; // remove child traces
      approvalInformation = {
        trace: approvalTrace,
        deployedImage,
      };
    }
    return approvalInformation;
  }

  private getSubSequencesForStage(stageTraces: Trace[], stageName: string, fromTime?: Date): SubSequence[] {
    const subSequences = stageTraces.filter((trace) => {
      let status = trace.data.stage === stageName;
      if (status && fromTime) {
        const lastTraceTime = trace.getLastTrace().time;
        status = !!lastTraceTime && new Date(lastTraceTime) > fromTime;
      }
      return status;
    });
    return subSequences
      .map((seq) => {
        return {
          name: seq.getLabel(),
          type: seq.type,
          result: seq.getStatus(),
          time: (seq.time ? new Date(seq.time) : new Date()).toISOString(),
          state: seq.isFinished() ? SequenceState.FINISHED : SequenceState.STARTED,
          id: seq.id,
          message: seq.getMessage(),
          hasPendingApproval: !!seq.findTrace((t) => !!t.isApproval())?.isApprovalPending(),
        };
      })
      .reverse();
  }

  private async getStageRemediationInformation(
    projectName: string,
    stageName: string,
    serviceName: string,
    openRemediations?: Remediation[]
  ): Promise<StageRemediationInformation> {
    if (!openRemediations) {
      openRemediations = await this.getOpenRemediations(projectName, false, serviceName);
    }
    let remediationConfig: string | undefined;
    const openRemediationsForStage = openRemediations
      .filter((seq) => seq.stages.some((st) => st.name === stageName))
      .map((seq) => {
        const { stages, ...rest } = seq;
        return Sequence.fromJSON({
          ...rest,
          stages: stages.filter((st) => st.name === stageName),
        });
      });
    if (openRemediationsForStage.length) {
      const resourceResponse = await this.apiService.getServiceResource(
        projectName,
        stageName,
        serviceName,
        'remediation.yaml'
      );
      remediationConfig = resourceResponse.data.resourceContent;
    }
    return {
      remediations: openRemediations,
      remediationsForStage: openRemediationsForStage,
      config: remediationConfig,
    };
  }

  public async getServiceRemediationInformation(
    projectName: string,
    serviceName: string,
    includeConfig: boolean
  ): Promise<ServiceRemediationInformation> {
    const serviceRemediationInformation: ServiceRemediationInformation = { stages: [] };
    const openRemediations = await this.getOpenRemediations(projectName, false, serviceName);
    const stageRemediations = openRemediations.reduce((stagesAcc: { [key: string]: Sequence[] }, remediation) => {
      const stageName = remediation.stages[0].name;
      if (!stagesAcc[stageName]) {
        stagesAcc[stageName] = [];
      }
      stagesAcc[stageName].push(remediation);
      return stagesAcc;
    }, {});
    for (const stage in stageRemediations) {
      let config: string | undefined;
      if (includeConfig) {
        const configResponse = await this.apiService.getServiceResource(
          projectName,
          stage,
          serviceName,
          'remediation.yaml'
        );
        config = configResponse.data.resourceContent;
      }
      serviceRemediationInformation.stages.push({ name: stage, remediations: stageRemediations[stage], config });
    }
    return serviceRemediationInformation;
  }

  public async intersectEvents(
    event: string,
    eventSuffix: EventState | '>',
    projectName: string,
    stages: string[],
    services: string[]
  ): Promise<Record<string, unknown>> {
    const objects = (await this.getMultipleLatestTracesOfMultipleServices(
      event,
      eventSuffix,
      projectName,
      stages,
      services
    )) as unknown as Record<string, unknown>[];

    let result = objects[0] ?? {};
    for (let i = 1; i < objects.length; ++i) {
      result = this.intersectObjects(result, objects[i]);
    }

    return result;
  }

  public intersectObjects(object1: Record<string, unknown>, object2: Record<string, unknown>): Record<string, unknown> {
    return Object.assign(
      {},
      ...Object.keys(object1).map((k) => {
        if (!(k in object2)) {
          return {};
        }
        const child1 = object1[k];
        const child2 = object2[k];
        if (child1 instanceof Array && child2 instanceof Array) {
          const temp = this.intersectArrays(child1, child2);
          return temp.length ? { [k]: temp } : {};
        } else {
          return this.validateAndIntersectObjects(child1, child2, k);
        }
      })
    );
  }

  private validateAndIntersectObjects(child1: unknown, child2: unknown, k: string): Record<string, unknown> {
    const isOneArray = child1 instanceof Array || child2 instanceof Array; // without this condition, it could be the case that child1 is an Array and child2 is object and typeof array is object
    if (isOneArray) {
      // type mismatch
      return {};
    } else if (this.isObject(child1) && this.isObject(child2)) {
      const temp = this.intersectObjects(child1, child2);
      return Object.keys(temp).length ? { [k]: temp } : {};
    } else if (!this.isObject(child1) && !this.isObject(child2)) {
      // string, number, null, ... is accepted
      return { [k]: child1 };
    }
    return {};
  }

  private isObject(element: unknown): element is Record<string, unknown> {
    return element !== null && typeof element === 'object';
  }

  private intersectArrays(array1: Array<unknown>, array2: Array<unknown>): Array<unknown> {
    const maxLength = Math.min(array1.length, array2.length);
    const result = [];
    for (let i = 0; i < maxLength; ++i) {
      const child1 = array1[i];
      const child2 = array2[i];
      if (child1 instanceof Array && child2 instanceof Array) {
        const array = this.intersectArrays(child1, child2);
        if (array.length) {
          result.push(array);
        }
      } else if (this.isObject(child1) && this.isObject(child2)) {
        const obj = this.intersectObjects(child1, child2);
        if (Object.keys(obj).length) {
          result.push(obj);
        }
      } else if (!this.isObject(child1) && !this.isObject(child2)) {
        result.push(child1);
      }
    }
    return result;
  }

  private async getMultipleLatestTracesOfMultipleServices(
    event: string,
    eventSuffix: EventState | '>',
    projectName: string,
    stages: string[],
    services: string[]
  ): Promise<Trace[]> {
    const suffixes =
      eventSuffix === '>' ? [EventState.TRIGGERED, EventState.STARTED, EventState.FINISHED] : [eventSuffix];
    const project = await this.getPlainProject(projectName);
    const filteredStages = stages.length ? stages : project.stages.map((s) => s.stageName);
    const filteredServices = services.length ? services : project.stages[0].services.map((s) => s.serviceName);
    const eventIds = this.getLatestServiceEventIds(project, filteredStages, filteredServices, event, suffixes);
    const traces: Trace[] = [];

    for (const suffix of suffixes) {
      const array = eventIds[suffix];
      for (let i = 0; i < array.length; i += this.MAX_PAGE_SIZE) {
        const chunkIds = array.slice(i, i + this.MAX_PAGE_SIZE);
        const response = await this.apiService.getTracesOfMultipleServices(
          projectName,
          `${event}.${suffix}` as EventTypes,
          chunkIds.join(',')
        );
        traces.push(...response.data.events);
      }
    }
    return traces;
  }

  private getLatestServiceEventIds(
    project: Project,
    stageFilter: string[],
    serviceFilter: string[],
    event: string,
    suffixes: EventState[]
  ): {
    [EventState.TRIGGERED]: string[];
    [EventState.STARTED]: string[];
    [EventState.FINISHED]: string[];
  } {
    const eventIds: {
      [EventState.TRIGGERED]: string[];
      [EventState.STARTED]: string[];
      [EventState.FINISHED]: string[];
    } = {
      triggered: [],
      started: [],
      finished: [],
    };

    for (const stage of project.stages) {
      if (stageFilter.includes(stage.stageName)) {
        for (const service of stage.services) {
          if (serviceFilter.includes(service.serviceName)) {
            for (const suffix of suffixes) {
              const id = service.lastEventTypes[`${event}.${suffix}`]?.eventId;
              if (id) {
                eventIds[suffix].push(id);
              }
            }
          }
        }
      }
    }
    return eventIds;
  }
}
