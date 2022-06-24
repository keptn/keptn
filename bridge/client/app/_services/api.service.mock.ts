/* eslint-disable @typescript-eslint/no-unused-vars */
import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable, of } from 'rxjs';
import { KeptnInfoResult } from '../../../shared/interfaces/keptn-info-result';
import moment from 'moment';
import { KeptnVersions } from '../../../shared/interfaces/keptn-versions';
import { Project } from '../_models/project';
import { UniformRegistrationResult } from '../../../shared/interfaces/uniform-registration-result';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { UniformSubscription } from '../_models/uniform-subscription';
import { WebhookConfig, WebhookConfigMethod } from '../../../shared/interfaces/webhook-config';
import { UniformRegistrationLogResponse } from '../../../shared/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { SecretScope } from '../../../shared/interfaces/secret-scope';
import { IMetadata } from '../_interfaces/metadata';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../_interfaces/event-result';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { Trace } from '../_models/trace';
import { ServiceState } from '../../../shared/models/service-state';
import { Deployment } from '../../../shared/interfaces/deployment';
import { IServiceRemediationInformation } from '../_interfaces/service-remediation-information';
import { VersionResponseMock } from './_mockData/api-responses/version-response.mock';
import { ProjectsMock } from './_mockData/projects.mock';
import { UniformRegistrationsMock } from './_mockData/uniform-registrations.mock';
import { UniformRegistrationLogsMock } from './_mockData/uniform-registrations-logs.mock';
import { SecretsResponseMock } from './_mockData/api-responses/secrets-response.mock';
import { BridgeInfoResponseMock } from './_mockData/api-responses/bridgeInfo-response.mock';
import { MetadataResponseMock } from './_mockData/api-responses/metadata-response.mock';
import { FileTreeMock } from './_mockData/fileTree.mock';
import { SequencesMock } from './_mockData/sequences.mock';
import { TracesResponseMock } from './_mockData/api-responses/traces-response.mock';
import { EvaluationResultsResponseDataMock } from './_mockData/api-responses/evaluation-results-response.mock';
import { EventResultResponseMock } from './_mockData/api-responses/event-result-response.mock';
import { ServiceStatesResultResponseMock } from './_mockData/api-responses/service-states-results-response.mock';
import { DeploymentResponseMock } from './_mockData/api-responses/deployment-response.mock';
import { ISequencesFilter } from '../../../shared/interfaces/sequencesFilter';
import { SequenceFilterMock } from './_mockData/sequence-filter.mock';
import { TriggerResponse, TriggerSequenceData } from '../_models/trigger-sequence';
import { IService } from '../../../shared/interfaces/service';
import { IGitDataExtended } from '../../../shared/models/IProject';
import { IProjectResult } from '../../../shared/interfaces/project-result';

@Injectable({
  providedIn: null,
})
export class ApiServiceMock extends ApiService {
  private localStoreMock: Map<string, string> = new Map();

  public get environmentFilter(): { [projectName: string]: { services: string[] } } {
    const item = this.localStoreMock.get(this.ENVIRONMENT_FILTER_COOKIE);
    return item ? JSON.parse(item) : {};
  }

  public set environmentFilter(filter: { [projectName: string]: { services: string[] } }) {
    this.localStoreMock.set(this.ENVIRONMENT_FILTER_COOKIE, JSON.stringify(filter));
  }

  public get uniformLogDates(): { [key: string]: string } {
    const data = this.localStoreMock.get(this.INTEGRATION_DATES);
    return data ? JSON.parse(data) : {};
  }

  public set uniformLogDates(dates: { [key: string]: string }) {
    this.localStoreMock.set(this.INTEGRATION_DATES, JSON.stringify(dates));
  }

  public getKeptnInfo(): Observable<KeptnInfoResult> {
    return of(BridgeInfoResponseMock);
  }

  public getIntegrationsPage(): Observable<string> {
    const page =
      '<h2>Custom Integrations on Keptn Sandbox</h2><ul><li><p>GitLab: <a href="https://github.com/keptn-sandbox/gitlab-service">GitHub</a></p></ul>';
    return of(page);
  }

  public isVersionCheckEnabled(): boolean | undefined {
    const item = this.localStoreMock.get(this.VERSION_CHECK_COOKIE);
    const versionInfo = item ? JSON.parse(item) : undefined;
    let enabled = versionInfo?.enabled;
    if (!enabled && (!versionInfo?.time || moment().subtract(5, 'days').isAfter(versionInfo.time))) {
      enabled = undefined;
    }
    return enabled;
  }

  public setVersionCheck(enabled: boolean): void {
    this.localStoreMock.set(this.VERSION_CHECK_COOKIE, JSON.stringify({ enabled, time: moment().valueOf() }));
  }

  public getAvailableVersions(): Observable<KeptnVersions | undefined> {
    return of(VersionResponseMock);
  }

  public deleteProject(projectName: string): Observable<Record<string, unknown>> {
    return of({
      message: 'ok',
    });
  }

  public createProjectExtended(projectName: string, shipyard: string, data?: IGitDataExtended): Observable<unknown> {
    return of({});
  }

  public createService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public deleteService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return of({
      message: 'ok',
    });
  }

  public getProject(projectName: string): Observable<Project> {
    const projects = [...ProjectsMock];
    return of(projects[0]);
  }

  public getService(projectName: string, stageName: string, serviceName: string): Observable<IService> {
    return of(ProjectsMock[0].stages[0].services[0]);
  }

  public getPlainProject(projectName: string): Observable<Project> {
    return this.getProject(projectName);
  }

  public getProjects(pageSize?: number): Observable<IProjectResult> {
    const projects = [...ProjectsMock];
    const result: IProjectResult = {
      projects: projects,
      totalCount: projects.length,
    };
    return of(result);
  }

  public getUniformRegistrations(uniformDates: { [key: string]: string }): Observable<UniformRegistrationResult[]> {
    return of(UniformRegistrationsMock);
  }

  public getUniformRegistrationInfo(integrationId: string): Observable<UniformRegistrationInfo> {
    return of({
      isControlPlane: true,
      isWebhookService: false,
    });
  }

  public getUniformSubscription(integrationId: string, subscriptionId: string): Observable<UniformSubscription> {
    const subscription = UniformSubscription.fromJSON({
      id: 'df9c0116-28ea-4ee2-8ad7-1fa6f03a8655',
      event: 'sh.keptn.event.deployment.triggered',
      filter: {
        projects: [],
        stages: [],
        services: [],
      },
    });
    return of(subscription);
  }

  public updateUniformSubscription(
    integrationId: string,
    subscription: Partial<UniformSubscription>,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    return of({});
  }
  public createUniformSubscription(
    integrationId: string,
    subscription: Partial<UniformSubscription>,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    return of({});
  }

  public getUniformRegistrationLogs(
    uniformRegistrationId: string,
    pageSize = 100
  ): Observable<UniformRegistrationLogResponse> {
    return of({ logs: UniformRegistrationLogsMock });
  }

  public hasUnreadUniformRegistrationLogs(uniformDates: { [key: string]: string }): Observable<boolean> {
    return of(false);
  }

  public getSecrets(): Observable<{ Secrets: Secret[] }> {
    const secrets = SecretsResponseMock;
    return of(secrets);
  }

  public getSecretsForScope(scope: SecretScope): Observable<Secret[]> {
    const secrets = SecretsResponseMock.Secrets.filter((s) => s.scope === scope);
    return of(secrets);
  }

  public addSecret(secret: Secret): Observable<Record<string, unknown>> {
    return of({});
  }

  public deleteSecret(name: string, scope: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public deleteSubscription(
    integrationId: string,
    subscriptionId: string,
    isWebhookService: boolean
  ): Observable<Record<string, unknown>> {
    return of({});
  }

  public getMetadata(): Observable<IMetadata> {
    return of(MetadataResponseMock);
  }

  public getFileTreeForService(projectName: string, serviceName: string): Observable<FileTree[]> {
    return of(FileTreeMock);
  }

  public getTaskNames(projectName: string): Observable<string[]> {
    return of(['evaluation', 'deployment', 'test', 'release', 'approval', 'rollback', 'get-action', 'action']);
  }

  public getServiceNames(projectName: string): Observable<string[]> {
    return of(['carts', 'carts-db']);
  }

  public getSequences(
    projectName: string,
    pageSize: number,
    sequenceName?: string,
    state?: string,
    fromTime?: string,
    beforeTime?: string,
    keptnContext?: string
  ): Observable<HttpResponse<SequenceResult>> {
    let data = SequencesMock;
    let totalCount = data.length;

    if (pageSize) {
      data = SequencesMock.slice(0, pageSize);
    }

    if (beforeTime) {
      data = SequencesMock.slice(totalCount - 9);
      totalCount = 9;
    }

    const body = {
      totalCount,
      states: data,
    };

    const res = new HttpResponse<SequenceResult>({ body });
    return of(res);
  }

  public getTraces(
    keptnContext: string,
    projectName?: string,
    fromTime?: string
  ): Observable<HttpResponse<EventResult>> {
    const body = {
      pageSize: 100,
      totalCount: 42,
      nextPageKey: 0,
      events: TracesResponseMock,
    };
    const res = new HttpResponse<EventResult>({ body });
    return of(res);
  }

  public getEvaluationResults(
    projectName: string,
    serviceName: string,
    stageName: string,
    fromTime?: string
  ): Observable<EventResult> {
    const result = {
      pageSize: 100,
      totalCount: 5,
      nextPageKey: 0,
      events: EvaluationResultsResponseDataMock,
    };
    return of(result);
  }

  public updateGitUpstreamExtended(projectName: string, data: IGitDataExtended): Observable<unknown> {
    return of({});
  }

  public sendApprovalEvent(
    approval: Trace,
    approve: boolean,
    eventType: EventTypes,
    source: string
  ): Observable<unknown> {
    return of({ keptnContext: '77baf26f-f64d-4a68-9ab5-efde9276ee73' });
  }

  public sendEvaluationInvalidated(evaluation: Trace, reason: string): Observable<unknown> {
    return of({ keptnContext: '77baf26f-f64d-4a68-9ab5-efde9276ee73' });
  }

  public getEvent(type?: string, project?: string, stage?: string, service?: string): Observable<EventResult> {
    const result = {
      pageSize: 1,
      totalCount: 1,
      nextPageKey: 0,
      events: EventResultResponseMock,
    };

    return of(result);
  }

  public sendSequenceControl(project: string, keptnContext: string, state: string): Observable<unknown> {
    return of({});
  }

  public getWebhookConfig(
    subscriptionId: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Observable<WebhookConfig> {
    const config = {
      type: '',
      method: 'POST' as WebhookConfigMethod,
      url: 'https://webhook.site/123456798',
      payload: '{"id":{{.id}}, "shkeptncontext": {{.shkeptncontext}}, "project": {{.data.project}}}',
      header: [],
      sendFinished: false,
      sendStarted: false,
      proxy: '',
      filter: {
        projects: null,
        stages: null,
        services: null,
      },
    };
    return of(config);
  }

  public getServiceStates(projectName: string): Observable<ServiceState[]> {
    return of(ServiceStatesResultResponseMock);
  }

  public getServiceDeployment(projectName: string, keptnContext: string, fromTime?: string): Observable<Deployment> {
    return of(DeploymentResponseMock);
  }

  public getOpenRemediationsOfService(
    projectName: string,
    serviceName: string
  ): Observable<IServiceRemediationInformation> {
    return of({
      stages: [],
    });
  }

  public getIntersectedEvent(
    event: string,
    eventSuffix: string,
    projectName: string,
    stages: string[],
    services: string[]
  ): Observable<Record<string, unknown>> {
    return of({
      data: {},
    });
  }

  public getSequencesFilter(projectName: string): Observable<ISequencesFilter> {
    return of(SequenceFilterMock);
  }

  public triggerSequence(type: string, data: TriggerSequenceData): Observable<TriggerResponse> {
    return of({ keptnContext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a' });
  }
}
/* eslint-enable @typescript-eslint/no-unused-vars */
