/* eslint-disable @typescript-eslint/no-unused-vars */
import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable, of } from 'rxjs';
import { KeptnInfoResult } from '../_models/keptn-info-result';
import moment from 'moment';
import { KeptnVersions } from '../_models/keptn-versions';
import { Project } from '../_models/project';
import { ProjectResult } from '../_interfaces/project-result';
import { UniformRegistrationResult } from '../../../shared/interfaces/uniform-registration-result';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { UniformSubscription } from '../_models/uniform-subscription';
import { WebhookConfig, WebhookConfigMethod } from '../../../shared/interfaces/webhook-config';
import { UniformRegistrationLogResponse } from '../../../shared/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { SecretScope } from '../../../shared/interfaces/secret-scope';
import { Metadata } from '../_models/metadata';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../_interfaces/event-result';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { Trace } from '../_models/trace';
import { ServiceState } from '../../../shared/models/service-state';
import { Deployment } from '../../../shared/interfaces/deployment';
import { IServiceRemediationInformation } from '../_interfaces/service-remediation-information';
import { versionMock } from './_mockData/version.mock';
import { Projects } from './_mockData/projects.mock';
import { UniformRegistrationsMock } from './_mockData/uniform-registrations.mock';
import { UniformRegistrationLogsMock } from './_mockData/uniform-registrations-logs.mock';
import { secretsMock } from './_mockData/secrets.mock';
import { bridgeInfoMock } from './_mockData/bridgeInfo.mock';
import { metadataMock } from './_mockData/metadata.mock';
import { FileTreeMock } from './_mockData/fileTree.mock';
import { SequencesData } from './_mockData/sequences.mock';
import { rootResultMock } from './_mockData/api-responses/root-result.mock';
import { Traces } from './_mockData/traces.mock';
import { evaluationResultsResponseDataMock } from './_mockData/api-responses/evaluation-results.mock';
import { getEventResultMock } from './_mockData/api-responses/event-result.mock';
import { serviceStatesResultMock } from './_mockData/api-responses/service-states-results.mock';
import { deploymentResponseMock } from './_mockData/api-responses/deployment-response.mock';

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
    return of(bridgeInfoMock);
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
    return of(versionMock);
  }

  public deleteProject(projectName: string): Observable<Record<string, unknown>> {
    return of({
      message: 'ok',
    });
  }

  public createProject(
    projectName: string,
    shipyard: string,
    gitRemoteUrl?: string,
    gitToken?: string,
    gitUser?: string
  ): Observable<unknown> {
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
    const projects = [...Projects];
    return of(projects[0]);
  }

  public getPlainProject(projectName: string): Observable<Project> {
    return this.getProject(projectName);
  }

  public getProjects(pageSize?: number): Observable<ProjectResult> {
    const projects = [...Projects];
    const result: ProjectResult = {
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
    const secrets = secretsMock;
    return of(secrets);
  }

  public getSecretsForScope(scope: SecretScope): Observable<Secret[]> {
    const secrets = secretsMock.Secrets.filter((s) => s.scope === scope);
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

  public getMetadata(): Observable<Metadata> {
    return of(metadataMock);
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
    let data = SequencesData;
    let totalCount = 34;

    if (pageSize) {
      data = SequencesData.slice(0, pageSize);
    }

    if (beforeTime) {
      data = SequencesData.slice(25, 34);
      totalCount = 9;
    }

    const body = {
      totalCount,
      states: data,
    };

    const res = new HttpResponse<SequenceResult>({ body });
    return of(res);
  }

  public getRoots(
    projectName: string,
    pageSize: number,
    serviceName?: string,
    fromTime?: string,
    beforeTime?: string,
    keptnContext?: string
  ): Observable<HttpResponse<EventResult>> {
    const body = {
      pageSize: 20,
      totalCount: 1,
      nextPageKey: 0,
      events: rootResultMock,
    };
    const res = new HttpResponse<EventResult>({ body });
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
      events: Traces,
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
      events: evaluationResultsResponseDataMock,
    };
    return of(result);
  }

  public sendGitUpstreamUrl(
    projectName: string,
    gitUrl: string,
    gitUser: string,
    gitToken: string
  ): Observable<unknown> {
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
      events: getEventResultMock,
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
    return of(serviceStatesResultMock);
  }

  public getServiceDeployment(projectName: string, keptnContext: string, fromTime?: string): Observable<Deployment> {
    return of(deploymentResponseMock);
  }

  public getOpenRemediationsOfService(
    projectName: string,
    serviceName: string
  ): Observable<IServiceRemediationInformation> {
    return of({
      stages: [],
    });
  }
}
/* eslint-enable @typescript-eslint/no-unused-vars */
