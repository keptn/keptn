import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { Resource } from '../../../shared/interfaces/resource';
import { Stage } from '../_models/stage';
import { ServiceResult } from '../_models/service-result';
import { Trace } from '../_models/trace';
import { ApprovalStates } from '../_models/approval-states';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { Metadata } from '../_models/metadata';
import { Deployment } from '../_models/deployment';
import moment from 'moment';
import { SequenceResult } from '../_models/sequence-result';
import { Project } from '../_models/project';
import { UniformRegistrationLogResponse } from '../../../server/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { KeptnInfoResult } from '../_models/keptn-info-result';
import { KeptnVersions } from '../_models/keptn-versions';
import { EventResult } from '../_interfaces/event-result';
import { ProjectResult } from '../_interfaces/project-result';
import { UniformSubscription } from '../_models/uniform-subscription';
import { WebhookConfig } from '../../../shared/interfaces/webhook-config';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { UniformRegistrationResult } from '../../../shared/interfaces/uniform-registration-result';
import { shareReplay } from 'rxjs/operators';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  private _baseUrl: string;
  private readonly VERSION_CHECK_COOKIE = 'keptn_versioncheck';
  private readonly ENVIRONMENT_FILTER_COOKIE = 'keptn_environment_filter';
  private readonly INTEGRATION_DATES = 'keptn_integration_dates';

  constructor(private http: HttpClient) {
    this._baseUrl = `./api`;
  }

  public set baseUrl(value: string) {
    this._baseUrl = value;
  }

  public get baseUrl(): string {
    return this._baseUrl;
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
    return this.http
      .get<KeptnInfoResult>(url);
  }

  public getIntegrationsPage(): Observable<string> {
    const url = `${this._baseUrl}/integrationsPage`;
    return this.http
      .get<string>(url, {responseType: 'text' as 'json'});
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
    localStorage.setItem(this.VERSION_CHECK_COOKIE, JSON.stringify({enabled, time: moment().valueOf()}));
  }


  public getAvailableVersions(): Observable<KeptnVersions | undefined> {
    if (this.isVersionCheckEnabled()) {
      const url = `${this._baseUrl}/version.json`;
      return this.http
        .get<KeptnVersions>(url);
    } else {
      return of(undefined);
    }
  }

  public deleteProject(projectName: string): Observable<object> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}`;
    return this.http.delete<object>(url);
  }

  /**
   * Creates a new project
   *
   * @param projectName - The unique project name - uniqueness is validated by the backend
   * @param shipyard - The base64 encoded contents of the yaml file
   * @param gitRemoteUrl (optional) - URL of the Git repository for the keptn configurations
   * @param gitToken (optional) - The Git token used for access permissions to the repository
   * @param gitUser (optional) - The username of the Git provider
   * @returns Observable with type unknown of the HttpResponse
   */
  public createProject(projectName: string, shipyard: string, gitRemoteUrl?: string, gitToken?: string, gitUser?: string): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    return this.http.post<unknown>(url, {
      gitRemoteUrl,
      gitToken,
      gitUser,
      name: projectName,
      shipyard,
    });
  }

  public createService(projectName: string, serviceName: string): Observable<object> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/service`;
    return this.http.post<object>(url, {
      serviceName,
    });
  }

  public deleteService(projectName: string, serviceName: string): Observable<object> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/service/${serviceName}`;
    return this.http.delete<object>(url);
  }

  public getProject(projectName: string): Observable<Project> {
    const url = `${this._baseUrl}/project/${projectName}`;
    const params = {
      approval: 'true',
      remediation: 'true',
    };
    return this.http
      .get<Project>(url, {params});
  }

  public getProjects(pageSize?: number): Observable<ProjectResult> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    const params = {
      disableUpstreamSync: 'true',
      ...pageSize && {pageSize: pageSize.toString()},
    };
    return this.http
      .get<ProjectResult>(url, {params});
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

  public updateUniformSubscription(integrationId: string, subscription: Partial<UniformSubscription>, webhookConfig?: WebhookConfig): Observable<object> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription/${subscription.id}`;
    return this.http.put(url, { subscription, webhookConfig });
  }
  public createUniformSubscription(integrationId: string, subscription: Partial<UniformSubscription>, webhookConfig?: WebhookConfig): Observable<object> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription`;
    return this.http.post(url, { subscription, webhookConfig });
  }

  public getUniformRegistrationLogs(uniformRegistrationId: string, pageSize: number = 100): Observable<UniformRegistrationLogResponse> {
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

  public addSecret(secret: Secret): Observable<object> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    return this.http.post(url, secret);
  }

  public deleteSecret(name: string, scope: string): Observable<object> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    const params = {
      name,
      scope,
    };
    return this.http.delete(url, {params});
  }

  public deleteSubscription(integrationId: string, subscriptionId: string, isWebhookService: boolean): Observable<object> {
    const url = `${this._baseUrl}/uniform/registration/${integrationId}/subscription/${subscriptionId}`;
    return this.http.delete(url, {
      params: {
        isWebhookService: String(isWebhookService),
      },
    });
  }

  public getMetadata(): Observable<Metadata> {
    return this.http.get<Metadata>(`${this._baseUrl}/v1/metadata`);
  }


  public getProjectResources(projectName: string): Observable<Resource[]> {
    const url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/resource`;
    return this.http
      .get<Resource[]>(url);
  }

  public getServiceResource(projectName: string, stageName: string, serviceName: string, resourceUri: string): Observable<Resource> {
    const url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service/${serviceName}/resource/${resourceUri}`;
    return this.http
      .get<Resource>(url);
  }

  public getFileTreeForService(projectName: string, serviceName: string): Observable<FileTree[]> {
    const url = `${this._baseUrl}/project/${projectName}/service/${serviceName}/files`;
    return this.http.get<FileTree[]>(url).pipe(shareReplay());
  }

  public getTaskNames(projectName: string): Observable<string[]> {
    const url = `${this._baseUrl}/project/${projectName}/tasks`;
    return this.http.get<string[]>(url);
  }

  public getStages(projectName: string): Observable<Stage[]> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage`;
    return this.http
      .get<Stage[]>(url);
  }

  public getServices(projectName: string, stageName: string, pageSize: number): Observable<ServiceResult> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage/${stageName}/service`;
    const params = {
      pageSize: pageSize.toString(),
    };
    return this.http
      .get<ServiceResult>(url, {params});
  }

  public getOpenRemediations(projectName: string, pageSize: number): Observable<HttpResponse<SequenceResult>> {
    return this.getSequences(projectName, pageSize, 'remediation', 'triggered');
  }

  public getSequences(projectName: string, pageSize: number, sequenceName?: string, state?: string,
                      fromTime?: string, beforeTime?: string, keptnContext?: string): Observable<HttpResponse<SequenceResult>> {
    const url = `${this._baseUrl}/controlPlane/v1/sequence/${projectName}`;
    const params = {
      pageSize: pageSize.toString(),
      ...(sequenceName && {name: sequenceName}),
      ...(state && {state}),
      ...(fromTime && {fromTime}),
      ...(beforeTime && {beforeTime}),
      ...(keptnContext && {keptnContext}),
    };

    return this.http
      .get<SequenceResult>(url, {params, observe: 'response'});
  }

  public getRoots(projectName: string, pageSize: number, serviceName?: string, fromTime?: string,
                  beforeTime?: string, keptnContext?: string): Observable<HttpResponse<EventResult>> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      root: 'true',
      pageSize: pageSize.toString(),
      project: projectName,
      ...serviceName && {serviceName},
      ...fromTime && {fromTime},
      ...beforeTime && {beforeTime},
      ...keptnContext && {keptnContext},
    };

    return this.http
      .get<EventResult>(url, {params, observe: 'response'});
  }

  public getTraces(keptnContext: string, projectName?: string, fromTime?: string): Observable<HttpResponse<EventResult>> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      keptnContext,
      ...projectName && {project: projectName},
      ...fromTime && {fromTime},
    };

    return this.http
      .get<EventResult>(url, {params, observe: 'response'});
  }

  public getDeploymentsOfService(projectName: string, serviceName: string): Observable<Deployment[]> {
    return of([]);
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, fromTime?: string): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `data.project:${projectName} AND data.service:${serviceName} AND data.stage:${stageName}`,
      excludeInvalidated: 'true',
      limit: '50',
      ...fromTime && {fromTime},
    };
    return this.http
      .get<EventResult>(url, {params});
  }

  public getEvaluationResult(shkeptncontext: string): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `shkeptncontext:${shkeptncontext}`,
      limit: '1',
    };
    return this.http
      .get<EventResult>(url, {params});
  }

  public sendGitUpstreamUrl(projectName: string, gitUrl: string, gitUser: string, gitToken: string): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    return this.http.put(url, {
      gitRemoteURL: gitUrl,
      gitToken,
      gitUser,
      name: projectName,
    });
  }

  public sendApprovalEvent(approval: Trace, approve: boolean, eventType: EventTypes, source: string): Observable<unknown> {
    const url = `${this._baseUrl}/v1/event`;

    return this.http
      .post(url, {
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

    return this.http
      .post<unknown>(url, {
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

  public sendSequenceControl(project: string, keptnContext: string, state: string): Observable<unknown> {
    const url = `${this._baseUrl}/controlPlane/v1/sequence/${project}/${keptnContext}/control`;

    return this.http
      .post<unknown>(url, {
        state,
      });
  }

  public getWebhookConfig(subscriptionId: string, projectName: string, stageName?: string, serviceName?: string): Observable<WebhookConfig> {
    const url = `${this._baseUrl}/uniform/registration/webhook-service/config/${subscriptionId}`;
    const params = {
      projectName,
      ...stageName && {stageName},
      ...serviceName && {serviceName},
    };
    return this.http
      .get<WebhookConfig>(url, {params});
  }

  public saveWebhookConfig(config: WebhookConfig): Observable<unknown> {
    const url = `${this._baseUrl}/uniform/registration/webhook-service/config`;

    return this.http
      .post<unknown>(url, {
        config,
      });
  }

}
