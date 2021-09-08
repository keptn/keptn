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
import { CURLParser } from 'parse-curl-js';
import { ParsedCURL } from 'parse-curl-js/dist/interface';
import { UniformRegistrationInfo } from '../../shared/interfaces/uniform-registration-info';
import { WebhookConfigYaml } from '../interfaces/webhook-config-yaml';
import { UniformSubscriptionFilter } from '../../shared/interfaces/uniform-subscription';
import axios from 'axios';

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
        const keptnContext = service.getLatestSequence(stage.stageName);
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
    if (sequence && stageName) { // we just need the result of a stage
      sequence = Sequence.fromJSON(sequence);
      if (includeEvaluationTrace) {
        const stage = sequence.stages.find(s => s.name === stageName);
        if (stage) { // get latest evaluation
          const evaluationTraces = await this.getEvaluationResults(projectName, sequence.service, stageName, 1, sequence.shkeptncontext);
          if (evaluationTraces) {
            stage.latestEvaluationTrace = evaluationTraces.shift();
          }
        }
      }
    }
    return sequence ?? Sequence.fromJSON(sequence);
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
      const stageName = sequence.stages[0].name;
      const response = await this.apiService.getTraces(this.buildRemediationEvent(stageName), this.MAX_TRACE_PAGE_SIZE, projectName, stageName, sequence.service);
      const traces = response.data.events;
      const stage = {...sequence.stages[0], actions: []};
      const remediation: Remediation = Remediation.fromJSON({...sequence, stages: [stage]});

      remediation.problemTitle = traces[0]?.data.problem?.ProblemTitle;
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

          remediation.stages[0].actions.push({...trace.data.action, state, result: finishedAction?.data.result});
        }
      }
      remediations.push(remediation);
    }
    return remediations;
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
      for (const sequence of stage.sequences) {
        for (const task of sequence.tasks) {
          if (!tasks.includes(task.name)) {
            tasks.push(task.name);
          }
        }
      }
    }
    return tasks;
  }

  private async getShipyard(projectName: string): Promise<Shipyard> {
    const response = await this.apiService.getShipyard(projectName);
    const shipyard = Buffer.from(response.data.resourceContent, 'base64').toString('utf-8');
    return Yaml.parse(shipyard);
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

  public async getWebhookConfig(projectName: string, stageName?: string, serviceName?: string): Promise<WebhookConfig> {
    const webhookConfigYaml: WebhookConfigYaml = await this.getWebhookConfigYaml(projectName, stageName, serviceName);

    let webhookCurlCommand: ParsedCURL;
    try {
      webhookCurlCommand = new CURLParser(webhookConfigYaml?.spec?.webhooks[0]?.requests[0]).parse();
    } catch (error) {
      throw Error('Could not parse curl command');
    }

    const webhookConfig: WebhookConfig = new WebhookConfig();
    webhookConfig.method = webhookCurlCommand.method;
    webhookConfig.url = webhookCurlCommand.url;
    webhookConfig.payload = JSON.stringify(webhookCurlCommand.body.data);
    for (const [name, value] of Object.entries(webhookCurlCommand.headers)) {
      webhookConfig.header?.push({
        name,
        value,
      });
    }
    return webhookConfig;
  }

  public async saveWebhookConfig(webhookConfig: WebhookConfig): Promise<boolean> {
    const previousConfig = webhookConfig.prevFilter ? await this.getPreviousWebhookConfig(webhookConfig.prevFilter) : undefined;
    const currentConfig = await this.getPreviousWebhookConfig(webhookConfig.filter);

    if (previousConfig) {
      await this.removePreviousWebhooks(previousConfig, currentConfig, webhookConfig.type);
    }

    const curl = this.generateWebhookConfigCurl(webhookConfig);

    for (const project of currentConfig.projects) {
      for (const stage of currentConfig.stages) {
        for (const service of currentConfig.services) {
          const previousWebhookConfig: WebhookConfigYaml = await this.getOrCreateWebhookConfigYaml(project, stage, service);
          previousWebhookConfig.addWebhook(webhookConfig.type, curl);
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

  private async removePreviousWebhooks(previousConfig: WebhookConfigFilter, currentConfig: WebhookConfigFilter, type: string): Promise<void> {
    for (const project of previousConfig.projects) {
      for (const stage of previousConfig.stages) {
        for (const service of previousConfig.services) {
          if (!currentConfig.projects.includes(project) || !currentConfig.stages.some((s?: string) => s === stage) || !currentConfig.services.some((s?: string) => s === service)) {
            await this.removeWebhook(type, project, stage, service);
          }
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
      params += `--data '${webhookConfig.payload.replace((/ {2}|\r\n|\n|\r/gm), '').replace((/"/g), '\\"')}' `;
    }
    return `curl ${params}${webhookConfig.url}`;
  }

  public async deleteSubscription(integrationId: string, subscriptionId: string, deleteWebhook: boolean): Promise<void> {
    if (deleteWebhook) {
      const response = await this.apiService.getUniformSubscription(integrationId, subscriptionId);
      const subscription = response.data;
      const projectName = subscription.filter.projects?.[0];
      if (projectName) {
        await this.removeWebhooks(subscription.event, projectName, subscription.filter.stages ?? [undefined], subscription.filter.services ?? [undefined]);
      }
    }
    await this.apiService.deleteUniformSubscription(integrationId, subscriptionId);
  }

  private async removeWebhooks(eventType: string, projectName: string, stages: string[] | [undefined], services: string[] | [undefined]): Promise<void> {
    for (const stage of stages) {
      for (const service of services) {
        await this.removeWebhook(eventType, projectName, stage, service);
      }
    }
  }

  private async removeWebhook(eventType: string, projectName: string, stage?: string, service?: string): Promise<void> {
    try {
      const webhookConfig: WebhookConfigYaml = await this.getWebhookConfigYaml(projectName, stage, service);
      if (webhookConfig.removeWebhook(eventType)) {
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
}
