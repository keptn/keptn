import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { Trace } from '../_models/trace';
import { ApprovalStates } from '../../../shared/models/approval-states';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { IMetadata } from '../_interfaces/metadata';
import moment from 'moment';
import { SequenceResult } from '../_models/sequence-result';
import { Project } from '../_models/project';
import { UniformRegistrationLogResponse } from '../../../shared/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { KeptnInfoResult } from '../../../shared/interfaces/keptn-info-result';
import { KeptnVersions } from '../../../shared/interfaces/keptn-versions';
import { EventResult } from '../_interfaces/event-result';
import { UniformSubscription } from '../_models/uniform-subscription';
import { WebhookConfig } from '../../../shared/interfaces/webhook-config';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { UniformRegistrationResult } from '../../../shared/interfaces/uniform-registration-result';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';
import { KeptnService } from '../../../shared/models/keptn-service';
import { ServiceState } from '../../../shared/models/service-state';
import { Deployment } from '../../../shared/interfaces/deployment';
import { IServiceRemediationInformation } from '../_interfaces/service-remediation-information';
import { EndSessionData } from '../../../shared/interfaces/end-session-data';
import { ISequencesFilter } from '../../../shared/interfaces/sequencesFilter';
import { TriggerResponse, TriggerSequenceData } from '../_models/trigger-sequence';
import { IScopesResult } from '../_interfaces/scopes-result';
import { SecretScope } from '../../../shared/interfaces/secret-scope';
import { ICustomSequences } from '../../../shared/interfaces/custom-sequences';
import { environment } from '../../environments/environment';
import { WindowConfig } from '../../environments/environment.dynamic';
import { IService } from '../../../shared/interfaces/service';
import { IProjectResult } from '../../../shared/interfaces/project-result';
import { IGitDataExtended } from '../../../shared/models/IProject';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  protected _baseUrl: string;
  protected readonly VERSION_CHECK_COOKIE = 'keptn_versioncheck';
  protected readonly ENVIRONMENT_FILTER_COOKIE = 'keptn_environment_filter';
  protected readonly INTEGRATION_DATES = 'keptn_integration_dates';
  protected readonly SEQUENCE_FILTERS_COOKIE = 'keptn_sequence_filters';

  constructor(protected http: HttpClient) {
    this._baseUrl = `./api`;
  }

  public set baseUrl(value: string) {
    this._baseUrl = value;
  }

  public get baseUrl(): string {
    return this._baseUrl;
  }

  public getSequenceFilters(projectName: string): Record<string, string[]> {
    const filters = localStorage.getItem(this.getSequenceFiltersKey(projectName));
    return filters ? JSON.parse(filters) : {};
  }

  public setSequenceFilters(filters: Record<string, string[]>, projectName: string): void {
    localStorage.setItem(this.getSequenceFiltersKey(projectName), JSON.stringify(filters));
  }

  private getSequenceFiltersKey(projectName: string): string {
    return `${this.SEQUENCE_FILTERS_COOKIE}-${projectName}`;
  }

  public get environmentFilter(): { [projectName: string]: { services: string[] } } {
    const item = localStorage.getItem(this.ENVIRONMENT_FILTER_COOKIE);
    const filter = item ? JSON.parse(item) : {};
    return filter instanceof Array ? {} : filter; // old format was an array
  }

  public set environmentFilter(filter: { [projectName: string]: { services: string[] } }) {
    localStorage.setItem(this.ENVIRONMENT_FILTER_COOKIE, JSON.stringify(filter));
  }

  public get uniformLogDates(): { [key: string]: string } {
    const data = localStorage.getItem(this.INTEGRATION_DATES);
    return data ? JSON.parse(data) : {};
  }

  public set uniformLogDates(dates: { [key: string]: string }) {
    localStorage.setItem(this.INTEGRATION_DATES, JSON.stringify(dates));
  }

  public getKeptnInfo(): Observable<KeptnInfoResult> {
    const url = `${this._baseUrl}/bridgeInfo`;
    return this.http.get<KeptnInfoResult>(url);
  }

  public getIntegrationsPage(): Observable<string> {
    const url = `${this._baseUrl}/integrationsPage`;
    return this.http.get<string>(url, { responseType: 'text' as 'json' });
  }

  public isVersionCheckEnabled(): boolean | undefined {
    const item = localStorage.getItem(this.VERSION_CHECK_COOKIE);
    const versionInfo = item ? JSON.parse(item) : undefined;
    let enabled = typeof versionInfo === 'boolean' ? versionInfo : versionInfo?.enabled; // support old format
    if (!enabled && (!versionInfo?.time || moment().subtract(5, 'days').isAfter(versionInfo.time))) {
      enabled = undefined;
    }
    return enabled;
  }

  public setVersionCheck(enabled: boolean): void {
    localStorage.setItem(this.VERSION_CHECK_COOKIE, JSON.stringify({ enabled, time: moment().valueOf() }));
  }

  public getAvailableVersions(): Observable<KeptnVersions | undefined> {
    if (this.isVersionCheckEnabled()) {
      const url = `${this._baseUrl}/version.json`;
      return this.http.get<KeptnVersions>(url);
    } else {
      return of(undefined);
    }
  }

  public deleteProject(projectName: string): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}`;
    return this.http.delete<Record<string, unknown>>(url);
  }

  public createProjectExtended(projectName: string, shipyard: string, data?: IGitDataExtended): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    return this.http.post<unknown>(url, {
      ...data,
      name: projectName,
      shipyard,
    });
  }

  public createService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/service`;
    return this.http.post<Record<string, unknown>>(url, {
      serviceName,
    });
  }

  public deleteService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/service/${serviceName}`;
    return this.http.delete<Record<string, unknown>>(url);
  }

  public getProject(projectName: string): Observable<Project> {
    const url = `${this._baseUrl}/project/${projectName}`;
    const params = {
      approval: 'true',
      remediation: 'true',
    };
    return this.http.get<Project>(url, { params });
  }

  public getService(projectName: string, stageName: string, serviceName: string): Observable<IService> {
    return this.http.get<IService>(
      `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage/${stageName}/service/${serviceName}`
    );
  }

  public getPlainProject(projectName: string): Observable<Project> {
    const url = `${this._baseUrl}/project/${projectName}`;
    return this.http.get<Project>(url);
  }

  public getProjects(pageSize?: number): Observable<IProjectResult> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    const params = {
      disableUpstreamSync: 'true',
      ...(pageSize && { pageSize: pageSize.toString() }),
    };
    return this.http.get<IProjectResult>(url, { params });
  }

  public getUniformRegistrations(uniformDates: { [key: string]: string }): Observable<UniformRegistrationResult[]> {
    const url = `${this._baseUrl}/uniform/registration`;
    return this.http.post<UniformRegistrationResult[]>(url, uniformDates);
  }

  public getUniformRegistrationInfo(integrationId: string): Observable<UniformRegistrationInfo> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/info`;
    return this.http.get<UniformRegistrationInfo>(url);
  }

  public getUniformSubscription(integrationId: string, subscriptionId: string): Observable<UniformSubscription> {
    const url = `${this._baseUrl}/controlPlane/v1/uniform/registration/${integrationId}/subscription/${subscriptionId}`;
    return this.http.get<UniformSubscription>(url);
  }

  public updateUniformSubscription(
    integrationId: string,
    subscription: Partial<UniformSubscription>,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription/${subscription.id}`;
    return this.http.put<Record<string, unknown>>(url, { subscription, webhookConfig });
  }

  public createUniformSubscription(
    integrationId: string,
    subscription: Partial<UniformSubscription>,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription`;
    return this.http.post<Record<string, unknown>>(url, { subscription, webhookConfig });
  }

  public getUniformRegistrationLogs(
    uniformRegistrationId: string,
    pageSize = 100
  ): Observable<UniformRegistrationLogResponse> {
    const url = `${this._baseUrl}/controlPlane/v1/log?integrationId=${uniformRegistrationId}&pageSize=${pageSize}`;
    return this.http.get<UniformRegistrationLogResponse>(url);
  }

  public hasUnreadUniformRegistrationLogs(uniformDates: { [key: string]: string }): Observable<boolean> {
    const url = `${this._baseUrl}/hasUnreadUniformRegistrationLogs`;
    return this.http.post<boolean>(url, uniformDates);
  }

  public getSecrets(): Observable<{ Secrets: Secret[] }> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    return this.http.get<{ Secrets: Secret[] }>(url);
  }

  public getSecretsForScope(scope: SecretScope): Observable<Secret[]> {
    const url = `${this._baseUrl}/secrets/scope/${scope}`;
    return this.http.get<Secret[]>(url);
  }

  public addSecret(secret: Secret): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    return this.http.post<Record<string, unknown>>(url, secret);
  }

  public deleteSecret(name: string, scope: string): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    const params = {
      name,
      scope,
    };
    return this.http.delete<Record<string, unknown>>(url, { params });
  }

  public deleteSubscription(
    integrationId: string,
    subscriptionId: string,
    isWebhookService: boolean
  ): Observable<Record<string, unknown>> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription/${subscriptionId}`;
    return this.http.delete<Record<string, unknown>>(url, {
      params: {
        isWebhookService: String(isWebhookService),
      },
    });
  }

  public getMetadata(): Observable<IMetadata> {
    return this.http.get<IMetadata>(`${this._baseUrl}/v1/metadata`);
  }

  public getFileTreeForService(projectName: string, serviceName: string): Observable<FileTree[]> {
    const url = `${this._baseUrl}/project/${projectName}/service/${serviceName}/files`;
    return this.http.get<FileTree[]>(url);
  }

  public getTaskNames(projectName: string): Observable<string[]> {
    const url = `${this._baseUrl}/project/${projectName}/tasks`;
    return this.http.get<string[]>(url);
  }

  public getServiceNames(projectName: string): Observable<string[]> {
    const url = `${this._baseUrl}/project/${projectName}/services`;
    return this.http.get<string[]>(url);
  }

  public getCustomSequences(projectName: string): Observable<ICustomSequences> {
    const url = `${this._baseUrl}/project/${projectName}/customSequences`;
    return this.http.get<ICustomSequences>(url);
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
    const url = `${this._baseUrl}/controlPlane/v1/sequence/${projectName}`;
    const params = {
      pageSize: pageSize.toString(),
      ...(sequenceName && { name: sequenceName }),
      ...(state && { state }),
      ...(fromTime && { fromTime }),
      ...(beforeTime && { beforeTime }),
      ...(keptnContext && { keptnContext }),
    };

    return this.http.get<SequenceResult>(url, { params, observe: 'response' });
  }

  public getTraces(
    keptnContext: string,
    projectName?: string,
    fromTime?: string,
    type?: EventTypes,
    source?: KeptnService,
    stage?: string,
    pageSize?: number
  ): Observable<HttpResponse<EventResult>> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      keptnContext,
      ...(projectName && { project: projectName }),
      ...(fromTime && { fromTime }),
      ...(type && { type }),
      ...(source && { source }),
      ...(stage && { stage }),
      ...(pageSize && { pageSize }),
    };

    return this.http.get<EventResult>(url, { params, observe: 'response' });
  }

  public getEvaluationResults(
    projectName: string,
    serviceName: string,
    stageName: string,
    fromTime?: string,
    limit?: number
  ): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName} AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
      excludeInvalidated: 'true',
      limit: limit?.toString() || '50',
      ...(fromTime && { fromTime }),
    };
    return this.http.get<EventResult>(url, { params });
  }

  public updateGitUpstreamExtended(projectName: string, data?: IGitDataExtended): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    return this.http.put<unknown>(url, {
      ...data,
      name: projectName,
    });
  }

  public sendApprovalEvent(
    approval: Trace,
    approve: boolean,
    eventType: EventTypes,
    source: string
  ): Observable<unknown> {
    const url = `${this._baseUrl}/v1/event`;

    return this.http.post(url, {
      shkeptncontext: approval.shkeptncontext,
      type: eventType,
      triggeredid: approval.id,
      source: `https://github.com/keptn/keptn/bridge#${source}`,
      data: {
        project: approval.data.project,
        stage: approval.data.stage,
        service: approval.data.service,
        labels: approval.data.labels,
        message: approval.data.message,
        result: approve ? ApprovalStates.APPROVED : ApprovalStates.DECLINED,
        status: 'succeeded',
      },
    });
  }

  public sendEvaluationInvalidated(evaluation: Trace, reason: string): Observable<unknown> {
    const url = `${this._baseUrl}/v1/event`;

    return this.http.post<unknown>(url, {
      shkeptncontext: evaluation.shkeptncontext,
      type: EventTypes.EVALUATION_INVALIDATED,
      triggeredid: evaluation.triggeredid,
      source: 'https://github.com/keptn/keptn/bridge#evaluation.invalidated',
      data: {
        project: evaluation.data.project,
        stage: evaluation.data.stage,
        service: evaluation.data.service,
        evaluation: {
          reason,
        },
      },
    });
  }

  public getEvent(type?: string, project?: string, stage?: string, service?: string): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      pageSize: '1',
      ...(type && { type }),
      ...(project && { project }),
      ...(stage && { stage }),
      ...(service && { service }),
    };

    return this.http.get<EventResult>(url, { params });
  }

  public sendSequenceControl(project: string, keptnContext: string, state: string): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/sequence/${project}/${keptnContext}/control`;

    return this.http.post<unknown>(url, {
      state,
    });
  }

  public getWebhookConfig(
    subscriptionId: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Observable<WebhookConfig> {
    const url = `${this._baseUrl}/uniform/registration/webhook-service/config/${subscriptionId}`;
    const params = {
      projectName,
      ...(stageName && { stageName }),
      ...(serviceName && { serviceName }),
    };
    return this.http.get<WebhookConfig>(url, { params });
  }

  public getServiceStates(projectName: string): Observable<ServiceState[]> {
    return this.http.get<ServiceState[]>(`${this._baseUrl}/project/${projectName}/serviceStates`);
  }

  public getServiceDeployment(projectName: string, keptnContext: string, fromTime?: string): Observable<Deployment> {
    const params = {
      ...(fromTime && { fromTime }),
    };
    return this.http.get<Deployment>(`${this._baseUrl}/project/${projectName}/deployment/${keptnContext}`, { params });
  }

  public getOpenRemediationsOfService(
    projectName: string,
    serviceName: string
  ): Observable<IServiceRemediationInformation> {
    const params = {
      config: 'true',
    };
    return this.http.get<IServiceRemediationInformation>(
      `${this._baseUrl}/project/${projectName}/service/${serviceName}/openRemediations`,
      { params }
    );
  }

  public getIntersectedEvent(
    event: string,
    eventSuffix: string,
    projectName: string,
    stages: string[],
    services: string[]
  ): Observable<Record<string, unknown>> {
    return this.http.post<Record<string, unknown>>(`${this._baseUrl}/intersectEvents`, {
      event,
      eventSuffix,
      projectName,
      stages,
      services,
    });
  }

  public logout(): Observable<EndSessionData | null> {
    return this.http.post<EndSessionData | null>(`./oauth/logout`, {});
  }

  public getSequencesFilter(projectName: string): Observable<ISequencesFilter> {
    return this.http.get<ISequencesFilter>(`${this._baseUrl}/project/${projectName}/sequences/filter`);
  }

  public triggerSequence(type: string, data: TriggerSequenceData): Observable<TriggerResponse> {
    const body = {
      contenttype: 'application/json',
      data,
      type,
      source: 'bridge',
    };
    return this.http.post<TriggerResponse>(`${this._baseUrl}/v1/event`, JSON.stringify(body));
  }

  public getSecretScopes(): Observable<IScopesResult> {
    return this.http.get<IScopesResult>(`${this._baseUrl}/secrets/v1/scope`);
  }

  public getLookAndFeelConfig(): Observable<WindowConfig | undefined> {
    return this.http.get<WindowConfig | undefined>(environment.appConfigUrl);
  }
}
