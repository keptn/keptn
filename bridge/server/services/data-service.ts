import { ApiService } from './api-service';
import { Sequence } from '../models/sequence';
import { SequenceTypes } from '../../shared/models/sequence-types';
import { Trace } from '../models/trace';
import { Service } from '../models/service';
import { Project } from '../models/project';
import { EventState } from '../../shared/models/event-state';
import { EventTypes } from '../../shared/interfaces/event-types';
import { ResultTypes } from '../../shared/models/result-types';
import { UniformRegistration } from '../models/uniform-registration';
import { parse as parseYaml } from 'yaml';
import { IShipyardSequence, IShipyardTask, Shipyard } from '../interfaces/shipyard';
import { UniformRegistrationLocations } from '../../shared/interfaces/uniform-registration-locations';
import { IWebhookConfigFilter } from '../interfaces/webhook-config';
import { UniformRegistrationInfo } from '../../shared/interfaces/uniform-registration-info';
import { WebhookConfigYaml } from '../models/webhook-config-yaml';
import { IUniformSubscription, IUniformSubscriptionFilter } from '../../shared/interfaces/uniform-subscription';
import axios from 'axios';
import { Resource } from '../../shared/interfaces/resource';
import { FileTree, TreeEntry } from '../../shared/interfaces/resourceFileTree';
import { EventResult } from '../interfaces/event-result';
import { IRemediationAction } from '../../shared/models/remediation-action';
import { KeptnService } from '../../shared/models/keptn-service';
import { SequenceState } from '../../shared/interfaces/sequence';
import { ServiceDeploymentInformation, ServiceState } from '../../shared/models/service-state';
import { Deployment, IStageDeployment, SubSequence } from '../../shared/interfaces/deployment';
import semver from 'semver';
import { ServiceRemediationInformation } from '../../shared/interfaces/service-remediation-information';
import { Stage } from '../models/stage';
import { IServiceEvent } from '../../shared/interfaces/service';
import { Remediation } from '../models/remediation';
import { IStage } from '../../shared/interfaces/stage';
import { ISequencesFilter } from '../../shared/interfaces/sequencesFilter';
import { SecretScope, SecretScopeDefault } from '../../shared/interfaces/secret-scope';
import { ICustomSequences } from '../../shared/interfaces/custom-sequences';
import { IWebhookConfigYamlResult, IWebhookSecret } from '../interfaces/webhook-config-yaml-result';
import {
  mapBridgeSecretsToYamlSecrets,
  mapYamlSecretsToBridgeSecrets,
  migrateWebhook,
  parseToClientWebhookRequest,
} from '../models/webhook-config.utils';
import { IWebhookConfigClient } from '../../shared/interfaces/webhook-config';
import { EnvType } from '../interfaces/configuration';
import { IClientSecret } from '../../shared/interfaces/secret';
import { IServerSequenceStage } from '../interfaces/sequence-stage';

type TreeDirectory = ({ _: string[] } & { [key: string]: TreeDirectory }) | { _: string[] };
type StageRemediationInformation = {
  remediations: Remediation[];
  remediationsForStage: Sequence[];
};
type StageOpenInformation = {
  openApprovals: Trace[];
  openRemediations: Remediation[];
  evaluations: Trace[];
};

interface IEventStateDict {
  [EventState.TRIGGERED]: string[];
  [EventState.STARTED]: string[];
  [EventState.FINISHED]: string[];
}

export interface SequenceOptions {
  pageSize: string;
  name?: string;
  state?: SequenceState;
  fromTime?: string;
  beforeTime?: string;
  keptnContext?: string;
}

export interface TraceOptions {
  pageSize: string;
  project: string;
  service: string;
  stage: string;
  type?: EventTypes;
  keptnContext?: string;
  result?: ResultTypes;
  source?: KeptnService;
}

export class DataService {
  private apiService: ApiService;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;
  private readonly MAX_TRACE_PAGE_SIZE = 50;
  private readonly MAX_PAGE_SIZE = 100;

  constructor(apiUrl: string, apiToken: string | undefined, mode: EnvType) {
    this.apiService = new ApiService(apiUrl, apiToken, mode);
  }

  public async getProjects(accessToken: string | undefined): Promise<Project[]> {
    const response = await this.apiService.getProjects(accessToken);
    return response.data.projects.map((project) => Project.fromJSON(project));
  }

  private async getPlainProject(accessToken: string | undefined, projectName: string): Promise<Project> {
    const response = await this.apiService.getProject(accessToken, projectName);
    return Project.fromJSON(response.data);
  }

  public async getProject(
    accessToken: string | undefined,
    projectName: string,
    includeRemediation: boolean,
    includeApproval: boolean
  ): Promise<Project> {
    const project = await this.getPlainProject(accessToken, projectName);
    let openApprovals: Trace[] = [];
    let openRemediations: Remediation[] = [];
    const allServices = Stage.getAllServices(project.stages);
    const latestDeployments = await this.getLatestDeploymentFinishedForServices(
      accessToken,
      projectName,
      allServices,
      ResultTypes.PASSED
    );
    const latestEvaluations = await this.getLatestEvaluationResultsForServices(accessToken, projectName, allServices);
    const latestSequences = await this.getLatestSequenceForServices(accessToken, projectName, allServices);

    if (includeRemediation) {
      openRemediations = await this.getOpenRemediations(accessToken, projectName, true);
    }

    if (includeApproval) {
      openApprovals = await this.getApprovals(accessToken, projectName);
    }

    for (const stage of project.stages) {
      const stageRemediations = openRemediations.filter((seq) => seq.stages.some((s) => s.name === stage.stageName));
      const stageApprovals = openApprovals.filter((approval) => approval.data.stage === stage.stageName);
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
          const latestSequence = latestSequences.find((seq) => seq.shkeptncontext === latestEvent.keptnContext);
          this.setServiceDetails(service, stage.stageName, stageInformation, stageDeployments, latestSequence);
        }
      }
    }
    return project;
  }

  private async getLatestSequenceForServices(
    accessToken: string | undefined,
    projectName: string,
    services: Service[]
  ): Promise<Sequence[]> {
    const sequences: Sequence[] = [];

    let keptnContexts: string[] = [];
    for (const service of services) {
      const latestServiceEvent = service.getLatestEvent();
      if (latestServiceEvent) {
        // string concatenation is expensive; that's why we use an array here
        keptnContexts.push(latestServiceEvent.keptnContext);
      }

      if (keptnContexts.length === this.MAX_PAGE_SIZE) {
        sequences.push(...(await this.getSequencesWithContexts(accessToken, projectName, keptnContexts)));
        keptnContexts = [];
      }
    }
    // get remaining sequences (mod 100)
    if (keptnContexts.length) {
      sequences.push(...(await this.getSequencesWithContexts(accessToken, projectName, keptnContexts)));
    }

    return sequences;
  }

  private async getSequencesWithContexts(
    accessToken: string | undefined,
    projectName: string,
    keptnContexts: string[]
  ): Promise<Sequence[]> {
    const response = await this.apiService.getSequences(accessToken, projectName, {
      pageSize: this.MAX_PAGE_SIZE.toString(),
      keptnContext: keptnContexts.join(','),
    });
    return response.data.states.map((seq) => Sequence.fromJSON(seq));
  }

  private async getLatestDeploymentFinishedForServices(
    accessToken: string | undefined,
    projectName: string,
    services: Service[],
    resultType?: ResultTypes
  ): Promise<Trace[]> {
    return this.getLatestTracesOfMultipleServices(
      accessToken,
      projectName,
      services,
      EventTypes.DEPLOYMENT_FINISHED,
      (service: Service) => service.lastEventTypes[EventTypes.DEPLOYMENT_FINISHED],
      undefined,
      resultType
    );
  }

  private async getLatestEvaluationResultsForServices(
    accessToken: string | undefined,
    projectName: string,
    services: Service[]
  ): Promise<Trace[]> {
    return this.getLatestTracesOfMultipleServices(
      accessToken,
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
   * @param accessToken
   * @param projectName
   * @param services
   * @param eventType
   * @param getServiceEvent
   * @param source
   * @param resultType
   * @private
   */
  private async getLatestTracesOfMultipleServices(
    accessToken: string | undefined,
    projectName: string,
    services: Service[],
    eventType: EventTypes,
    getServiceEvent: (service: Service) => IServiceEvent | undefined,
    source?: KeptnService,
    resultType?: ResultTypes
  ): Promise<Trace[]> {
    const traces = [];
    const fetchEvents = async (ids: string[]): Promise<Trace[]> => {
      const response = await this.apiService.getTracesOfMultipleServices(
        accessToken,
        projectName,
        eventType,
        ids.join(',')
      );
      const events = response.data.events;
      if (resultType || source) {
        await this.checkAndSetEventsWithResult(accessToken, events, resultType, source);
      }
      return events;
    };

    const eventIds = services
      .map((service) => getServiceEvent(service))
      .filter((serviceEvent): serviceEvent is IServiceEvent => !!serviceEvent)
      .map((serviceEvent) => serviceEvent.eventId);

    for (let i = 0; i < eventIds.length; i += this.MAX_PAGE_SIZE) {
      const eventIdChunk = eventIds.slice(i, i + this.MAX_PAGE_SIZE);
      const events = await fetchEvents(eventIdChunk);
      traces.push(...events);
    }

    return traces;
  }

  private async checkAndSetEventsWithResult(
    accessToken: string | undefined,
    traces: Trace[],
    resultType?: ResultTypes,
    source?: KeptnService
  ): Promise<void> {
    for (let i = 0; i < traces.length; ++i) {
      const trace = traces[i];
      if ((resultType && trace.data.result !== resultType) || (source && trace.source !== source)) {
        const response = await this.apiService.getTracesWithResultAndSource(accessToken, {
          type: trace.type as EventTypes,
          pageSize: '1',
          project: trace.data.project as string,
          stage: trace.data.stage as string,
          service: trace.data.service as string,
          ...(resultType && { result: resultType }),
          ...(source && { source }),
        });
        if (response.data.events.length) {
          traces[i] = response.data.events[0];
        }
      }
    }
  }

  private setServiceDetails(
    service: Service,
    stageName: string,
    stageInformation: StageOpenInformation,
    stageDeployments: Trace[],
    latestSequence?: Sequence
  ): void {
    // remove reference. reduceToStage then does not affect other sequences
    service.latestSequence = latestSequence ? Sequence.fromJSON(latestSequence) : undefined;
    service.latestSequence?.reduceToStage(stageName);

    const stage = service.latestSequence?.stages.find((seq) => seq.name === stageName);
    if (stage && !stage.latestEvaluationTrace) {
      stage.latestEvaluationTrace = stageInformation.evaluations.find(
        (t) => t.data.service === service.latestSequence?.service
      );
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
      (approval) => approval.data.service === service.serviceName
    );
  }

  public async getSequences(
    accessToken: string | undefined,
    projectName: string,
    sequenceName: string,
    sequenceState?: SequenceState
  ): Promise<Sequence[]> {
    const response = await this.apiService.getSequences(accessToken, projectName, {
      pageSize: this.MAX_SEQUENCE_PAGE_SIZE.toString(),
      name: sequenceName,
      ...(sequenceState && { state: sequenceState }),
    });

    return response.data.states.map((sequence) => Sequence.fromJSON(sequence));
  }

  public async getOpenRemediations(
    accessToken: string | undefined,
    projectName: string,
    includeActions: boolean,
    serviceName?: string
  ): Promise<Remediation[]> {
    const sequences = await this.getSequences(
      accessToken,
      projectName,
      SequenceTypes.REMEDIATION,
      SequenceState.STARTED
    );
    const remediations: Remediation[] = [];
    for (const sequence of sequences) {
      const stageName = sequence.stages[0]?.name;
      // there could be invalid sequences that don't have a stage because the triggered sequence was not present in the shipyard file
      if (stageName && (!serviceName || sequence.service === serviceName)) {
        const stage = { ...sequence.stages[0], actions: [] };
        const remediation: Remediation = Remediation.fromJSON({ ...sequence, stages: [stage] });

        if (includeActions) {
          await this.loadRemediationActions(
            accessToken,
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
    accessToken: string | undefined,
    remediation: Remediation,
    projectName: string,
    stageName: string,
    serviceName: string,
    keptnContext: string
  ): Promise<void> {
    const response = await this.apiService.getTraces(accessToken, {
      pageSize: this.MAX_TRACE_PAGE_SIZE.toString(),
      project: projectName,
      stage: stageName,
      service: serviceName,
      keptnContext,
    });
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

  public async getApprovals(
    accessToken: string | undefined,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<Trace[]> {
    try {
      const response = await this.apiService.getOpenTriggeredEvents(
        accessToken,
        projectName,
        EventTypes.APPROVAL_TRIGGERED,
        stageName,
        serviceName
      );
      return response.data.events ?? [];
    } catch {
      // status 500 if no events are found
      return [];
    }
  }

  public async hasUnreadUniformRegistrationLogs(
    accessToken: string | undefined,
    uniformDates: { [key: string]: string }
  ): Promise<boolean> {
    const response = await this.apiService.getUniformRegistrations(accessToken);
    const registrations = this.getValidRegistrations(response.data);
    let status = false;
    for (let i = 0; i < registrations.length && !status; ++i) {
      const registration = registrations[i];
      const logResponse = await this.apiService.getUniformRegistrationLogs(
        accessToken,
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

  public async getUniformRegistrations(
    accessToken: string | undefined,
    uniformDates: { [key: string]: string }
  ): Promise<UniformRegistration[]> {
    const response = await this.apiService.getUniformRegistrations(accessToken);
    const registrations = this.getValidRegistrations(response.data);
    for (const registration of registrations) {
      const logResponse = await this.apiService.getUniformRegistrationLogs(
        accessToken,
        registration.id,
        uniformDates[registration.id]
      );
      registration.unreadEventsCount = logResponse.data.logs.length;
    }
    return registrations;
  }

  private getValidRegistrations(registrations: UniformRegistration[]): UniformRegistration[] {
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

  public async getIsUniformRegistrationInfo(
    accessToken: string | undefined,
    integrationId: string
  ): Promise<UniformRegistrationInfo> {
    const response = await this.apiService.getUniformRegistrations(accessToken, integrationId);
    const uniformRegistration = UniformRegistration.fromJSON(response.data.shift());

    return {
      isControlPlane: uniformRegistration.metadata.location === UniformRegistrationLocations.CONTROL_PLANE,
      isWebhookService: uniformRegistration.isWebhookService,
    };
  }

  public async getTasks(accessToken: string | undefined, projectName: string): Promise<string[]> {
    const shipyard = await this.getShipyard(accessToken, projectName);
    try {
      const defaultTasks = new Set<string>();
      defaultTasks.add('evaluation');

      const taskSet = shipyard.spec.stages
        .reduce((sequences, stage) => [...sequences, ...(stage.sequences ?? [])], [] as IShipyardSequence[])
        .reduce((tasks, sequence) => [...tasks, ...sequence.tasks], [] as IShipyardTask[])
        .reduce((tasks, task) => {
          tasks.add(task.name);
          return tasks;
        }, defaultTasks);

      return Array.from(taskSet);
    } catch (error) {
      throw Error('Could not parse shipyard.yaml');
    }
  }

  public async getServiceNames(accessToken: string | undefined, projectName: string): Promise<string[]> {
    const resp = await this.apiService.getStages(accessToken, projectName);
    const stages = resp.data.stages;

    return this.reduceServiceNames(stages);
  }

  public async getCustomSequenceNames(accessToken: string | undefined, projectName: string): Promise<ICustomSequences> {
    const shipyard = await this.getShipyard(accessToken, projectName);
    const ignoredSequences = ['delivery', 'evaluation'];
    const sequences: ICustomSequences = {};

    for (const stage of shipyard.spec.stages) {
      sequences[stage.name] =
        stage.sequences
          ?.filter((seq) => !ignoredSequences.includes(seq.name))
          .map((seq) => seq.name)
          .sort() ?? [];
    }
    return sequences;
  }

  private reduceServiceNames(stages: IStage[]): string[] {
    const services: { [serviceName: string]: boolean | undefined } = {};

    for (const stage of stages) {
      for (const service of stage.services) {
        services[service.serviceName] = true;
      }
    }

    return Object.keys(services);
  }

  public async getTracesByContext(
    accessToken: string | undefined,
    keptnContext: string | undefined,
    projectName?: string | undefined,
    fromTime?: string | undefined,
    type?: EventTypes,
    source?: KeptnService
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
        accessToken,
        keptnContext,
        projectName,
        fromTime,
        nextPage.toString(),
        type,
        source
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
    accessToken: string | undefined,
    filterOptions?: {
      keptnContext?: string;
      projectName?: string;
      stageName?: string;
      serviceName?: string;
      eventType?: EventTypes;
      source?: KeptnService;
      pageSize?: number;
    }
  ): Promise<EventResult> {
    const response = await this.apiService.getTraces(accessToken, {
      type: filterOptions?.eventType,
      ...(filterOptions?.pageSize !== undefined && { pageSize: filterOptions.pageSize.toString() }),
      ...(filterOptions?.projectName && { project: filterOptions.projectName }),
      ...(filterOptions?.stageName && { stage: filterOptions.stageName }),
      ...(filterOptions?.serviceName && { service: filterOptions.serviceName }),
      ...(filterOptions?.keptnContext && { keptnContext: filterOptions.keptnContext }),
      ...(filterOptions?.source && { source: filterOptions.source }),
    });
    return response.data;
  }

  public async getResourceFileTreesForService(
    accessToken: string | undefined,
    projectName: string,
    serviceName: string
  ): Promise<FileTree[]> {
    const projectRes = await this.apiService.getProject(accessToken, projectName);
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
          accessToken,
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

  private async getShipyard(accessToken: string | undefined, projectName: string): Promise<Shipyard> {
    const response = await this.apiService.getShipyard(accessToken, projectName);
    const shipyard = Buffer.from(response.data.resourceContent, 'base64').toString('utf-8');
    return parseYaml(shipyard);
  }

  public async createSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscription: IUniformSubscription,
    webhookConfig?: IWebhookConfigClient
  ): Promise<void> {
    const response = await this.apiService.createSubscription(accessToken, integrationId, subscription);
    if (webhookConfig) {
      await this.saveWebhookConfig(accessToken, webhookConfig, response.data.id, subscription.filter);
    }
  }

  public async updateSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscriptionId: string,
    subscription: IUniformSubscription,
    webhookConfig?: IWebhookConfigClient
  ): Promise<void> {
    await this.apiService.updateSubscription(accessToken, integrationId, subscriptionId, subscription);
    if (webhookConfig) {
      await this.saveWebhookConfig(accessToken, webhookConfig, subscriptionId, subscription.filter);
    }
  }

  private async getWebhookConfigYaml(
    accessToken: string | undefined,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<WebhookConfigYaml> {
    const response = await this.apiService.getWebhookConfig(accessToken, projectName, stageName, serviceName);

    try {
      const webhookConfigFile: IWebhookConfigYamlResult = parseYaml(
        Buffer.from(response.data.resourceContent, 'base64').toString('utf-8')
      );
      return WebhookConfigYaml.fromJSON(migrateWebhook(webhookConfigFile));
    } catch (error) {
      throw Error('Could not parse webhook.yaml');
    }
  }

  public async getWebhookConfig(
    accessToken: string | undefined,
    subscriptionId: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Promise<IWebhookConfigClient> {
    const webhookConfigYaml: WebhookConfigYaml = await this.getWebhookConfigYaml(
      accessToken,
      projectName,
      stageName,
      serviceName
    );

    const config = parseToClientWebhookRequest(webhookConfigYaml, subscriptionId);
    if (!config) {
      throw Error('Could not parse webhook request');
    }
    if (config.secrets) {
      mapYamlSecretsToBridgeSecrets(config.webhookConfig, config.secrets);
    }
    return config.webhookConfig;
  }

  public async saveWebhookConfig(
    accessToken: string | undefined,
    webhookConfig: IWebhookConfigClient,
    subscriptionId: string,
    filter: IUniformSubscriptionFilter
  ): Promise<boolean> {
    const currentFilters = await this.getWebhookConfigFilter(accessToken, filter);

    if (webhookConfig.prevConfiguration) {
      const previousFilter = await this.getWebhookConfigFilter(accessToken, webhookConfig.prevConfiguration.filter);
      await this.removePreviousWebhooks(accessToken, previousFilter, subscriptionId);
    }

    const secrets = await this.parseAndReplaceWebhookSecret(accessToken, webhookConfig);

    for (const project of currentFilters.projects) {
      for (const stage of currentFilters.stages) {
        for (const service of currentFilters.services) {
          const previousWebhookConfig: WebhookConfigYaml = await this.getOrCreateWebhookConfigYaml(
            accessToken,
            project,
            stage,
            service
          );
          previousWebhookConfig.addWebhook(webhookConfig, subscriptionId, secrets);
          await this.apiService.saveWebhookConfig(accessToken, previousWebhookConfig.toYAML(), project, stage, service);
        }
      }
    }

    return true;
  }

  private async parseAndReplaceWebhookSecret(
    accessToken: string | undefined,
    webhookConfig: IWebhookConfigClient
  ): Promise<IWebhookSecret[]> {
    const secrets = await this.getSecretsForScope(accessToken, SecretScopeDefault.WEBHOOK);
    return mapBridgeSecretsToYamlSecrets(webhookConfig, secrets);
  }

  private async getOrCreateWebhookConfigYaml(
    accessToken: string | undefined,
    project: string,
    stage?: string,
    service?: string
  ): Promise<WebhookConfigYaml> {
    let previousWebhookConfig: WebhookConfigYaml;
    try {
      // fetch existing one
      previousWebhookConfig = await this.getWebhookConfigYaml(accessToken, project, stage, service);
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

  private async removePreviousWebhooks(
    accessToken: string | undefined,
    previousConfig: IWebhookConfigFilter,
    subscriptionId: string
  ): Promise<void> {
    for (const project of previousConfig.projects) {
      for (const stage of previousConfig.stages) {
        for (const service of previousConfig.services) {
          await this.removeWebhook(accessToken, subscriptionId, project, stage, service);
        }
      }
    }
  }

  private async getWebhookConfigFilter(
    accessToken: string | undefined,
    webhookConfig: IUniformSubscriptionFilter
  ): Promise<IWebhookConfigFilter> {
    return {
      projects: webhookConfig.projects?.length
        ? webhookConfig.projects
        : (await this.getProjects(accessToken)).map((project) => project.projectName),
      stages: webhookConfig.stages?.length ? webhookConfig.stages : [undefined],
      services: webhookConfig.services?.length ? webhookConfig.services : [undefined],
    };
  }

  public async deleteSubscription(
    accessToken: string | undefined,
    integrationId: string,
    subscriptionId: string,
    deleteWebhook: boolean
  ): Promise<void> {
    if (deleteWebhook) {
      const response = await this.apiService.getUniformSubscription(accessToken, integrationId, subscriptionId);
      const subscription = response.data;
      const projectName = subscription.filter.projects?.[0];

      if (projectName) {
        await this.removeWebhooks(
          accessToken,
          subscriptionId,
          projectName,
          subscription.filter.stages?.length ? subscription.filter.stages : [undefined],
          subscription.filter.services?.length ? subscription.filter.services : [undefined]
        );
      }
    }
    await this.apiService.deleteUniformSubscription(accessToken, integrationId, subscriptionId);
  }

  public async getSecretsForScope(accessToken: string | undefined, scope: SecretScope): Promise<IClientSecret[]> {
    const response = await this.apiService.getSecrets(accessToken);
    return response.data.Secrets.filter((secret) => secret.scope === scope);
  }

  private async removeWebhooks(
    accessToken: string | undefined,
    subscriptionId: string,
    projectName: string,
    stages: string[] | [undefined],
    services: string[] | [undefined]
  ): Promise<void> {
    for (const stage of stages) {
      for (const service of services) {
        await this.removeWebhook(accessToken, subscriptionId, projectName, stage, service);
      }
    }
  }

  private async removeWebhook(
    accessToken: string | undefined,
    subscriptionId: string,
    projectName: string,
    stage?: string,
    service?: string
  ): Promise<void> {
    try {
      const webhookConfig: WebhookConfigYaml = await this.getWebhookConfigYaml(
        accessToken,
        projectName,
        stage,
        service
      );
      if (webhookConfig.removeWebhook(subscriptionId)) {
        if (webhookConfig.hasWebhooks()) {
          await this.apiService.saveWebhookConfig(accessToken, webhookConfig.toYAML(), projectName, stage, service);
        } else {
          await this.apiService.deleteWebhookConfig(accessToken, projectName, stage, service);
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

  public async getServiceStates(accessToken: string | undefined, projectName: string): Promise<ServiceState[]> {
    const projectResponse = await this.apiService.getProject(accessToken, projectName);
    const project = Project.fromJSON(projectResponse.data);
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
          const deploymentInformation = this.getOrCreateDeploymentInformation(
            serviceState,
            service,
            latestDeploymentEvent.keptnContext
          );
          const deploymentStage = {
            name: stage.stageName,
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
    accessToken: string | undefined,
    projectName: string,
    keptnContext: string,
    fromTimeString?: string
  ): Promise<Deployment | undefined> {
    const fromTime = fromTimeString ? new Date(fromTimeString) : undefined;
    const sequenceResponse = await this.apiService.getSequences(accessToken, projectName, {
      pageSize: '1',
      keptnContext,
    });
    let deployment: Deployment | undefined;
    const iSequence = sequenceResponse.data.states[0];
    if (iSequence) {
      const sequence = Sequence.fromJSON(iSequence);
      let openRemediations: Remediation[] | undefined;
      const project = await this.getPlainProject(accessToken, projectName);
      const traceResponse = await this.apiService.getTracesByContext(accessToken, keptnContext, projectName);
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
        const { deploymentStage, remediations, image } = await this.getDeploymentStageInformation(
          accessToken,
          traces,
          stage,
          sequence,
          project,
          {
            fromTime,
            cachedRemediations: openRemediations,
          }
        );
        openRemediations = remediations;
        deployment.image ??= image;
        deployment.stages.push(deploymentStage);
      }
    }
    return deployment;
  }

  private async getDeploymentStageInformation(
    accessToken: string | undefined,
    traces: Trace[],
    stage: IServerSequenceStage,
    sequence: Sequence,
    project: Project,
    options: { fromTime?: Date; cachedRemediations: Remediation[] | undefined }
  ): Promise<{
    deploymentStage: IStageDeployment;
    remediations: Remediation[] | undefined;
    image: string | undefined;
  }> {
    const stageTraces = traces.filter((trace) => trace.data.stage === stage.name);
    const service = project.stages
      .find((st) => st.stageName === stage.name)
      ?.services.find((sv) => sv.serviceName === sequence.service);
    const latestDeploymentContext = service?.deploymentEvent?.keptnContext;
    const evaluationTrace = stageTraces.reduce(
      (event: Trace | undefined, trace) => event || trace.getEvaluationFinishedEvent(),
      undefined
    );
    if (evaluationTrace) {
      evaluationTrace.traces = []; //clear not wanted events like from the webhook-service
    }
    const lastTimeUpdated = stageTraces[stageTraces.length - 1]?.getLastTrace()?.time;
    const approvalInformation = this.getApprovalInformation(stageTraces, service?.getShortImage());
    let deploymentURL: string | undefined;
    let stageRemediationInformation: StageRemediationInformation | undefined;

    if (latestDeploymentContext === sequence.shkeptncontext) {
      deploymentURL = stageTraces.reduce(
        (url: string | undefined, tc) =>
          url || tc.findTrace((t) => t.type === EventTypes.DEPLOYMENT_FINISHED)?.getDeploymentUrl(),
        undefined
      );
      stageRemediationInformation = await this.getStageRemediationInformation(
        accessToken,
        project.projectName,
        stage.name,
        sequence.service,
        options.cachedRemediations
      );
      options.cachedRemediations = stageRemediationInformation.remediations;
    }

    return {
      deploymentStage: {
        name: stage.name,
        state: stage.state,
        lastTimeUpdated: (lastTimeUpdated ? new Date(lastTimeUpdated) : new Date()).toISOString(),
        openRemediations: stageRemediationInformation?.remediationsForStage ?? [],
        approvalInformation,
        subSequences: this.getSubSequencesForStage(stageTraces, stage.name, options.fromTime),
        deploymentURL,
        hasEvaluation: !!evaluationTrace,
        latestEvaluation:
          evaluationTrace ??
          (
            await this.apiService.getEvaluationResults(
              accessToken,
              project.projectName,
              sequence.service,
              stage.name,
              1
            )
          ).data.events[0],
      },
      remediations: options.cachedRemediations,
      image: service?.getShortImage(),
    };
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
    accessToken: string | undefined,
    projectName: string,
    stageName: string,
    serviceName: string,
    openRemediations?: Remediation[]
  ): Promise<StageRemediationInformation> {
    if (!openRemediations) {
      openRemediations = await this.getOpenRemediations(accessToken, projectName, false, serviceName);
    }
    const openRemediationsForStage = openRemediations
      .filter((seq) => seq.stages.some((st) => st.name === stageName))
      .map((seq) => {
        const { stages, ...rest } = seq;
        return Sequence.fromJSON({
          ...rest,
          stages: stages.filter((st) => st.name === stageName),
        });
      });

    return {
      remediations: openRemediations,
      remediationsForStage: openRemediationsForStage,
    };
  }

  public async getServiceRemediationInformation(
    accessToken: string | undefined,
    projectName: string,
    serviceName: string
  ): Promise<ServiceRemediationInformation> {
    const serviceRemediationInformation: ServiceRemediationInformation = { stages: [] };
    const openRemediations = await this.getOpenRemediations(accessToken, projectName, false, serviceName);
    const stageRemediations = openRemediations.reduce((stagesAcc: { [key: string]: Sequence[] }, remediation) => {
      const stageName = remediation.stages[0].name;
      if (!stagesAcc[stageName]) {
        stagesAcc[stageName] = [];
      }
      stagesAcc[stageName].push(remediation);
      return stagesAcc;
    }, {});
    for (const stage in stageRemediations) {
      serviceRemediationInformation.stages.push({ name: stage, remediations: stageRemediations[stage] });
    }
    return serviceRemediationInformation;
  }

  public async getSequencesFilter(accessToken: string | undefined, projectName: string): Promise<ISequencesFilter> {
    const res = await this.apiService.getStages(accessToken, projectName);
    const stages = res.data.stages;
    const stageNames: string[] = [];
    const serviceSet: Set<string> = new Set();

    for (const stage of stages) {
      for (const service of stage.services) {
        serviceSet.add(service.serviceName);
      }
      stageNames.push(stage.stageName);
    }

    return {
      stages: stageNames,
      services: Array.from(serviceSet),
    };
  }

  public async intersectEvents(
    accessToken: string | undefined,
    event: string,
    eventSuffix: EventState | '>',
    projectName: string,
    stages: string[],
    services: string[]
  ): Promise<Record<string, unknown>> {
    const objects = (await this.getMultipleLatestTracesOfMultipleServices(
      accessToken,
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
    accessToken: string | undefined,
    event: string,
    eventSuffix: EventState | '>',
    projectName: string,
    stages: string[],
    services: string[]
  ): Promise<Trace[]> {
    const suffixes =
      eventSuffix === '>' ? [EventState.TRIGGERED, EventState.STARTED, EventState.FINISHED] : [eventSuffix];
    const project = await this.getPlainProject(accessToken, projectName);
    const filteredStages = stages.length ? stages : project.stages.map((s) => s.stageName);
    const filteredServices = services.length ? services : project.stages[0].services.map((s) => s.serviceName);
    const eventIds = this.getLatestServiceEventIds(project, filteredStages, filteredServices, event, suffixes);
    const traces: Trace[] = [];

    for (const suffix of suffixes) {
      const array = eventIds[suffix];
      for (let i = 0; i < array.length; i += this.MAX_PAGE_SIZE) {
        const chunkIds = array.slice(i, i + this.MAX_PAGE_SIZE);
        const response = await this.apiService.getTracesOfMultipleServices(
          accessToken,
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
  ): IEventStateDict {
    const eventIds: IEventStateDict = {
      triggered: [],
      started: [],
      finished: [],
    };

    project.stages
      .filter((st) => stageFilter.includes(st.stageName))
      .reduce((services, stage) => [...services, ...stage.services], [] as Service[])
      .filter((sv) => serviceFilter.includes(sv.serviceName))
      .forEach((service) => {
        suffixes
          .map((suffix) => ({ suffix, id: service.lastEventTypes[`${event}.${suffix}`]?.eventId }))
          .filter((data): data is { suffix: EventState; id: string } => !!data.id)
          .forEach(({ suffix, id }) => {
            eventIds[suffix].push(id);
          });
      });
    return eventIds;
  }
}
