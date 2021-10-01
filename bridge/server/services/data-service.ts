import { ApiService } from './api-service';
import { Sequence } from '../models/sequence';
import { SequenceTypes } from '../../shared/models/sequence-types';
import { Trace } from '../models/trace';
import { DeploymentInformation, Service } from '../models/service';
import { Project } from '../models/project';
import { EventState } from '../../shared/models/event-state';
import { Remediation } from '../models/remediation';
import { EventTypes } from '../../shared/interfaces/event-types';
import { Approval } from '../interfaces/approval';
import { ResultTypes } from '../../shared/models/result-types';
import { UniformRegistration } from '../models/uniform-registration';
import Yaml from 'yaml';
import { Shipyard } from '../interfaces/shipyard';
import { UniformRegistrationLocations } from '../../shared/interfaces/uniform-registration-locations';
import { WebhookConfig, WebhookConfigFilter } from '../../shared/models/webhook-config';
import { UniformRegistrationInfo } from '../../shared/interfaces/uniform-registration-info';
import { WebhookConfigYaml } from '../interfaces/webhook-config-yaml';
import { UniformSubscription, UniformSubscriptionFilter } from '../../shared/interfaces/uniform-subscription';
import axios from 'axios';
import { Resource } from '../../shared/interfaces/resource';
import { FileTree, TreeEntry } from '../../shared/interfaces/resourceFileTree';
import { EventResult } from '../interfaces/event-result';
import { IRemediationAction } from '../../shared/models/remediation-action';

type TreeDirectory = ({ _: string[] } & { [key: string]: TreeDirectory }) | { _: string[] };

export class DataService {
  private apiService: ApiService;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;
  private readonly MAX_TRACE_PAGE_SIZE = 50;

  constructor(apiUrl: string, apiToken: string) {
    this.apiService = new ApiService(apiUrl, apiToken);
  }

  public async getProjects(): Promise<Project[]> {
    const response = await this.apiService.getProjects();
    return response.data.projects.map(project => Project.fromJSON(project));
  }

  public async getProject(projectName: string, includeRemediation: boolean, includeApproval: boolean): Promise<Project> {
    const response = await this.apiService.getProject(projectName);
    const project = Project.fromJSON(response.data);
    let remediations: Remediation[] = [];

    if (includeRemediation) {
      remediations = await this.getRemediations(projectName);
    }
    const lastSequences: { [key: string]: Sequence } = {};
    for (const stage of project.stages) {
      for (const service of stage.services) {
        const keptnContext = service.getLatestSequence();
        if (keptnContext) {
          try {
            const latestSequence = await this.fetchServiceDetails(service, stage.stageName, keptnContext, projectName, includeApproval, remediations, lastSequences[service.serviceName]);
            if (latestSequence) {
              lastSequences[service.serviceName] = latestSequence;
            }
          } catch (error) {
            console.error(error);
          }
        }
      }
    }
    return project;
  }

  private async fetchServiceDetails(service: Service, stageName: string, keptnContext: string, projectName: string, includeApproval: boolean, remediations: Remediation[], sequenceBefore: Sequence | undefined): Promise<Sequence | undefined> {
    const latestSequence = sequenceBefore?.shkeptncontext === keptnContext ? sequenceBefore : await this.getSequence(projectName, stageName, keptnContext, true);
    service.latestSequence = latestSequence ? Sequence.fromJSON(latestSequence) : undefined;
    service.latestSequence?.reduceToStage(stageName);
    service.deploymentInformation = await this.getDeploymentInformation(service, projectName, stageName);

    const serviceRemediations = remediations.filter(remediation => remediation.service === service.serviceName && remediation.stages.some(s => s.name === stageName));
    for (const remediation of serviceRemediations) {
      remediation.reduceToStage(stageName);
    }
    service.openRemediations = serviceRemediations;

    if (includeApproval) {
      service.openApprovals = await this.getApprovals(projectName, stageName, service.serviceName);
    }
    return latestSequence;
  }

  public async getSequence(projectName: string, stageName?: string, keptnContext?: string, includeEvaluationTrace = false): Promise<Sequence | undefined> {
    const response = await this.apiService.getSequences(projectName, 1, undefined, undefined, undefined, undefined, keptnContext);
    let sequence = response.data.states[0];
    if (sequence) {
      sequence = Sequence.fromJSON(sequence);
      if (includeEvaluationTrace && stageName) { // we just need the result of a stage
        const stage = sequence.stages.find(s => s.name === stageName);
        if (stage) { // get latest evaluation
          const evaluationTraces = await this.getEvaluationResults(projectName, sequence.service, stageName, 1, sequence.shkeptncontext);
          if (evaluationTraces) {
            stage.latestEvaluationTrace = evaluationTraces.shift();
          }
        }
      }
    }
    return sequence;
  }

  public async getDeploymentInformation(service: Service, projectName: string, stageName: string): Promise<DeploymentInformation | undefined> {
    const result = await this.apiService.getTracesWithResult(EventTypes.DEPLOYMENT_FINISHED, 1, projectName, stageName, service.serviceName, ResultTypes.PASSED);
    const traceData = result.data.events[0];
    let deploymentInformation: DeploymentInformation | undefined;
    if (traceData) {
      const trace = Trace.fromJSON(traceData);
      deploymentInformation = {
        deploymentUrl: trace.getDeploymentUrl(),
        image: trace.getShortImageName(),
      };
    }
    return deploymentInformation;
  }

  private async getEvaluationResults(projectName: string, serviceName: string, stageName: string, pageSize: number, keptnContext?: string): Promise<Trace[]> {
    const response = await this.apiService.getEvaluationResults(projectName, serviceName, stageName, pageSize, keptnContext);
    return response.data.events.map(trace => Trace.fromJSON(trace));
  }

  public async getSequences(projectName: string, sequenceName: string, stageName?: string, keptnContext?: string): Promise<Sequence[]> {
    const response = await this.apiService.getSequences(projectName, this.MAX_SEQUENCE_PAGE_SIZE, sequenceName, undefined, undefined, undefined, keptnContext);
    const sequences = response.data.states;
    for (const sequence of sequences) {
      if (stageName) { // we just need the result of a stage
        if (sequence.name === SequenceTypes.REMEDIATION) { // if the sequence is a remediation also return the problemTitle
          const traceResponse = await this.apiService.getTraces(this.buildRemediationEvent(stageName), 1, projectName, stageName, sequence.service);
          const trace = traceResponse.data.events[0];
          sequence.problemTitle = trace?.data.problem?.ProblemTitle;
        }
        sequence.stages = sequence.stages.filter(stage => stage.name === stageName);
      }
    }
    return sequences.map(sequence => Sequence.fromJSON(sequence));
  }

  public async getRemediations(projectName: string): Promise<Remediation[]> {
    const sequences = await this.getSequences(projectName, SequenceTypes.REMEDIATION);
    const remediations: Remediation[] = [];
    for (const sequence of sequences) {
      const stageName = sequence.stages[0]?.name;
      // there could be invalid sequences that don't have a stage because the triggered sequence was not present in the shipyard file
      if (stageName) {
        const response = await this.apiService.getTraces(this.buildRemediationEvent(stageName), this.MAX_TRACE_PAGE_SIZE, projectName, stageName, sequence.service);
        const traces = response.data.events;
        const stage = {...sequence.stages[0], actions: []};
        const remediation: Remediation = Remediation.fromJSON({...sequence, stages: [stage]});
        const actions = this.getRemediationActions(traces);

        remediation.problemTitle = traces[0]?.data.problem?.ProblemTitle;
        remediation.stages[0].actions.push(...actions);
        remediations.push(remediation);
      }
    }
    return remediations;
  }

  private getRemediationActions(traces: Trace[]): IRemediationAction[] {
    const actions: IRemediationAction[] = [];
    for (const trace of traces) {
      if (trace.type === EventTypes.ACTION_TRIGGERED && trace.data.action) {
        const finishedAction = traces.find(t => t.triggeredid === trace.id && t.type === EventTypes.ACTION_FINISHED);
        const startedAction = traces.find(t => t.triggeredid === trace.id && t.type === EventTypes.ACTION_STARTED);
        let state: EventState;
        if (finishedAction) {
          state = EventState.FINISHED;
        } else if (startedAction) {
          state = EventState.STARTED;
        } else {
          state = EventState.TRIGGERED;
        }

        actions.push({...trace.data.action, state, result: finishedAction?.data.result});
      }
    }
    return actions;
  }

  private async getTrace(keptnContext: string, projectName: string, stageName: string, serviceName: string, eventType: EventTypes): Promise<Trace | undefined> {
    const response = await this.apiService.getTraces(eventType, 1, projectName, stageName, serviceName, keptnContext);
    return response.data.events.shift();
  }

  public async getApprovals(projectName: string, stageName: string, serviceName: string): Promise<Approval[]> {
    let tracesTriggered: Trace[];
    try {
      const response = await this.apiService.getOpenTriggeredEvents(projectName, stageName, serviceName, EventTypes.APPROVAL_TRIGGERED);
      tracesTriggered = response.data.events ?? [];
    } catch { // status 500 if no events are found
      tracesTriggered = [];
    }
    const approvals: Approval[] = [];
    // for each approval the latest evaluation trace (before this event) is needed
    for (const trace of tracesTriggered) {
      const evaluationTrace = await this.getTrace(trace.shkeptncontext, projectName, stageName, serviceName, EventTypes.EVALUATION_FINISHED);
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
      const logResponse = await this.apiService.getUniformRegistrationLogs(registration.id, uniformDates[registration.id], 1);
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
      const logResponse = await this.apiService.getUniformRegistrationLogs(registration.id, uniformDates[registration.id]);
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
    const tasks: string[] = ['service.delete', 'service.create', 'evaluation'];
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

  public async getRoots(projectName: string | undefined, pageSize: string | undefined, serviceName: string | undefined, fromTime?: string | undefined, beforeTime?: string | undefined, keptnContext?: string | undefined): Promise<EventResult> {
    const response = await this.apiService.getRoots(projectName, pageSize, serviceName, fromTime, beforeTime, keptnContext);
    return response.data;
  }

  public async getTracesByContext(keptnContext: string | undefined, projectName?: string | undefined, fromTime?: string | undefined): Promise<EventResult> {
    let result: EventResult = {
      events: [],
      pageSize: 0,
      nextPageKey: 0,
      totalCount: 0,
    };
    let nextPage = 0;
    do {
      const response = await this.apiService.getTracesByContext(keptnContext, projectName, fromTime, nextPage.toString());
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
        const resourceRes = await this.apiService.getServiceResource(projectName, stage.stageName, serviceName, nextPage || undefined);
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

  public async createSubscription(integrationId: string, subscription: UniformSubscription, webhookConfig: WebhookConfig): Promise<void> {
    const response = await this.apiService.createSubscription(integrationId, subscription);
    const subscriptionId = response.data.id;
    await this.saveWebhookConfig(webhookConfig, subscriptionId);
  }

  public async updateSubscription(integrationId: string, subscriptionId: string, subscription: UniformSubscription, webhookConfig: WebhookConfig): Promise<void> {
    await this.apiService.updateSubscription(integrationId, subscriptionId, subscription);
    await this.saveWebhookConfig(webhookConfig, subscriptionId);
  }

  private async getWebhookConfigYaml(projectName: string, stageName?: string, serviceName?: string): Promise<WebhookConfigYaml> {
    const response = await this.apiService.getWebhookConfig(projectName, stageName, serviceName);

    try {
      const webhookConfigFile = Buffer.from(response.data.resourceContent, 'base64').toString('utf-8');
      return WebhookConfigYaml.fromJSON(Yaml.parse(webhookConfigFile));
    } catch (error) {
      throw Error('Could not parse webhook.yaml');
    }
  }

  public async getWebhookConfig(subscriptionId: string, projectName: string, stageName?: string, serviceName?: string): Promise<WebhookConfig> {
    const webhookConfigYaml: WebhookConfigYaml = await this.getWebhookConfigYaml(projectName, stageName, serviceName);

    const webhookConfig = webhookConfigYaml.parsedRequest(subscriptionId);
    if (!webhookConfig) {
      throw Error('Could not parse curl command');
    }
    return webhookConfig;
  }

  public async saveWebhookConfig(webhookConfig: WebhookConfig, subscriptionId: string): Promise<boolean> {
    const currentConfig = await this.getPreviousWebhookConfig(webhookConfig.filter);

    if (webhookConfig.prevConfiguration) {
      const previousFilter = await this.getPreviousWebhookConfig(webhookConfig.prevConfiguration.filter);
      await this.removePreviousWebhooks(previousFilter, webhookConfig.prevConfiguration.type);
    }

    const curl = this.generateWebhookConfigCurl(webhookConfig);

    for (const project of currentConfig.projects) {
      for (const stage of currentConfig.stages) {
        for (const service of currentConfig.services) {
          const previousWebhookConfig: WebhookConfigYaml = await this.getOrCreateWebhookConfigYaml(project, stage, service);
          previousWebhookConfig.addWebhook(webhookConfig.type, curl, subscriptionId);
          await this.apiService.saveWebhookConfig(previousWebhookConfig.toYAML(), project, stage, service);
        }
      }
    }

    return true;
  }

  private async getOrCreateWebhookConfigYaml(project: string, stage?: string, service?: string): Promise<WebhookConfigYaml> {
    let previousWebhookConfig: WebhookConfigYaml;
    try { // fetch existing one
      previousWebhookConfig = await this.getWebhookConfigYaml(project, stage, service);
    } catch (e: unknown) {
      if (!axios.isAxiosError(e) || e.response?.status !== 404) {
        throw e;
      } else { // if it does not exist, create one
        previousWebhookConfig = new WebhookConfigYaml();
      }
    }
    return previousWebhookConfig;
  }

  private async removePreviousWebhooks(previousConfig: WebhookConfigFilter, type: string): Promise<void> {
    for (const project of previousConfig.projects) {
      for (const stage of previousConfig.stages) {
        for (const service of previousConfig.services) {
          await this.removeWebhook(type, project, stage, service);
        }
      }
    }
  }

  private async getPreviousWebhookConfig(webhookConfig: UniformSubscriptionFilter): Promise<WebhookConfigFilter> {
    return {
      projects: webhookConfig.projects?.length ? webhookConfig.projects : (await this.getProjects()).map(project => project.projectName),
      stages: webhookConfig.stages?.length ? webhookConfig.stages : [undefined],
      services: webhookConfig.services?.length ? webhookConfig.services : [undefined],
    };
  }

  private buildRemediationEvent(stageName: string): string {
    return `sh.keptn.event.${stageName}.${SequenceTypes.REMEDIATION}.${EventState.TRIGGERED}`;
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
        stringify = stringify.replace((/\r\n|\n|\r/gm), '');
      }
      params += `--data '${stringify}' `;
    }
    return `curl ${params}${webhookConfig.url}`;
  }

  public async deleteSubscription(integrationId: string, subscriptionId: string, deleteWebhook: boolean): Promise<void> {
    if (deleteWebhook) {
      const response = await this.apiService.getUniformSubscription(integrationId, subscriptionId);
      const subscription = response.data;
      const projectName = subscription.filter.projects?.[0];
      if (projectName && subscription.filter.stages?.length) {
        await this.removeWebhooks(subscriptionId, projectName, subscription.filter.stages, subscription.filter.services?.length ? subscription.filter.services : [undefined]);
      }
    }
    await this.apiService.deleteUniformSubscription(integrationId, subscriptionId);
  }

  private async removeWebhooks(subscriptionId: string, projectName: string, stages: string[], services: string[] | [undefined]): Promise<void> {
    for (const stage of stages) {
      for (const service of services) {
        await this.removeWebhook(subscriptionId, projectName, stage, service);
      }
    }
  }

  private async removeWebhook(subscriptionId: string, projectName: string, stage?: string, service?: string): Promise<void> {
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
    const directory: TreeDirectory = {_: []};

    for (const res of resources) {
      const parts = res.resourceURI.split('/').filter(item => !!item);
      this._addToTreeDirectory(directory, parts);
    }

    return this._buildTree(directory, '').children || [];
  }

  private _addToTreeDirectory(currentDirectory: TreeDirectory, parts: string[]): void {
    let index = 0;
    for (; index < parts.length - 1; ++index) {
      const part = parts[index];
      // @ts-ignore
      currentDirectory[part] ??= {_: []};
      // @ts-ignore
      currentDirectory = currentDirectory[part];
    }
    currentDirectory._.push(parts[index]);
  }

  private _buildTree(currentDirectory: TreeDirectory, fileName: string): TreeEntry {
    const tree: TreeEntry = {fileName, children: [] as TreeEntry[]};
    const dict: { [key: string]: boolean } = {_: true};
    const folders: TreeEntry[] = [];
    const files: TreeEntry[] = [];
    for (const key in currentDirectory) {
      if (dict[key]) {
        files.push(...currentDirectory._.map(item => {
          return {
            fileName: item,
          } as TreeEntry;
        }));
      } else {
        // @ts-ignore
        folders.push(this._buildTree(currentDirectory[key] as TreeDirectory, key));
      }
    }
    folders.sort((a, b) => a.fileName.localeCompare(b.fileName));
    files.sort((a, b) => a.fileName.localeCompare(b.fileName));
    tree.children = [...folders, ...files];
    return tree;
  }
}
