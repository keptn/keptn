import {Injectable} from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import {Observable, of} from 'rxjs';
import {map} from 'rxjs/operators';
import {Resource} from '../_models/resource';
import {Stage} from '../_models/stage';
import {ProjectResult} from '../_models/project-result';
import {ServiceResult} from '../_models/service-result';
import {EventResult} from '../_models/event-result';
import {Trace} from '../_models/trace';
import {ApprovalStates} from '../_models/approval-states';
import {EventTypes} from '../_models/event-types';
import {Metadata} from '../_models/metadata';
import {TaskNames} from '../_models/task-names.mock';
import {Deployment} from '../_models/deployment';
import moment from 'moment';
import {SequenceResult} from '../_models/sequence-result';
import {Project} from '../_models/project';
import {UniformRegistration} from '../_models/uniform-registration';
import {UniformRegistrationLogResponse} from '../_models/uniform-registration-log';
import {Secret} from '../_models/secret';
import { KeptnInfoResult } from '../_models/keptn-info-result';
import { KeptnVersions } from '../_models/keptn-versions';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private _baseUrl: string;
  private VERSION_CHECK_COOKIE = 'keptn_versioncheck';
  private ENVIRONMENT_FILTER_COOKIE = 'keptn_environment_filter';

  constructor(private http: HttpClient) {
    this._baseUrl = `./api`;
  }

  public set baseUrl(value: string) {
    this._baseUrl = value;
  }

  public get baseUrl(): string {
    return this._baseUrl;
  }

  public get environmentFilter(): {[projectName: string]: {services: string[]}} {
    const item = localStorage.getItem(this.ENVIRONMENT_FILTER_COOKIE);
    const filter = item ? JSON.parse(item) : {};
    return filter instanceof Array ? {} : filter; // old format was an array
  }

  public set environmentFilter(filter: {[projectName: string]: {services: string[]}}) {
    localStorage.setItem(this.ENVIRONMENT_FILTER_COOKIE, JSON.stringify(filter));
  }

  public getKeptnInfo(): Observable<KeptnInfoResult> {
    const url = `${this._baseUrl}/bridgeInfo`;
    return this.http
      .get<KeptnInfoResult>(url);
  }

  public getKeptnVersion(): Observable<string> {
    const url = `${this._baseUrl}/swagger-ui/swagger.yaml`;
    return this.http
      .get<string>(url, { headers: new HttpHeaders({'Access-Control-Allow-Origin': '*'}) })
      .pipe(
        map(res => res.substring(res.lastIndexOf('version: ') + 9)),
        map(res => res.substring(0, res.indexOf('\n'))),
      );
  }

  public getIntegrationsPage(): Observable<string> {
    const url = `${this._baseUrl}/integrationsPage`;
    return this.http
      .get<string>(url, { responseType: 'text' as 'json' });
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

  /**
   * Creates a new project
   *
   * @param projectName - The unique project name - uniqueness is validated by the backend
   * @param shipyard - The base64 encoded contents of the yaml file
   * @param gitRemoteUrl (optional) - URL of the Git repository for the keptn configurations
   * @param gitToken (optional) - The Git token used for access permissions to the repository
   * @param gitUser (optional) - The username of the Git provider
   * @returns Observable with type any of the HttpResponse
   */
  public createProject(projectName: string, shipyard: string, gitRemoteUrl?: string, gitToken?: string, gitUser?: string): Observable<any> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    return this.http.post<any>(url, {
      gitRemoteUrl: gitRemoteUrl,
      gitToken: gitToken,
      gitUser: gitUser,
      name: projectName,
      shipyard: shipyard
    });
  }

  public getProject(projectName: string): Observable<Project> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}`;
    const params = {
      disableUpstreamSync: 'true'
    };
    return this.http
      .get<Project>(url, {params});
  }

  public getProjects(pageSize?: number): Observable<ProjectResult> {
    const url = `${this._baseUrl}/controlPlane/v1/project`;
    const params = {
      disableUpstreamSync: 'true',
      ... pageSize && {pageSize: pageSize.toString()}
    };

    return this.http
      .get<ProjectResult>(url, {params});
  }

  public getUniformRegistrations(): Observable<UniformRegistration[]> {
    const url = `${this._baseUrl}/controlPlane/v1/uniform/registration`;
    return this.http.get<UniformRegistration[]>(url);
  }

  public getUniformRegistrationLogs(uniformRegistrationId: string, pageSize: number = 100): Observable<UniformRegistrationLogResponse> {
    const url = `${this._baseUrl}/controlPlane/v1/log?integrationId=${uniformRegistrationId}&pageSize=${pageSize}`;
    return this.http.get<UniformRegistrationLogResponse>(url);
  }

  public getSecrets(): Observable<{Secrets: Secret[]}> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    return this.http.get<{Secrets: Secret[]}>(url);
  }

  public addSecret(secret: Secret): Observable<object> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    return this.http.post(url, secret);
  }

  public deleteSecret(name: string, scope: string): Observable<object> {
    const url = `${this._baseUrl}/secrets/v1/secret`;
    const params = {
      name,
      scope
    };
    return this.http.delete(url, {params});
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

  public getTaskNames(projectName: string): Observable<string[]>{
    return of(TaskNames);
  }

  public getStages(projectName: string): Observable<Stage[]> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage`;
    return this.http
      .get<Stage[]>(url);
  }

  public getServices(projectName: string, stageName: string, pageSize: number): Observable<ServiceResult> {
    const url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage/${stageName}/service`;
    const params = {
      pageSize: pageSize.toString()
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
      ...(keptnContext && {keptnContext})
    };

    return this.http
      .get<SequenceResult>(url, { params, observe: 'response' });
  }

  public getRoots(projectName: string, pageSize: number, serviceName?: string, fromTime?: string,
                  beforeTime?: string, keptnContext?: string): Observable<HttpResponse<EventResult>> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      root: 'true',
      pageSize: pageSize.toString(),
      project: projectName,
      ... serviceName && {serviceName},
      ... fromTime && {fromTime},
      ... beforeTime && {beforeTime},
      ... keptnContext && {keptnContext},
    };

    return this.http
      .get<EventResult>(url, { params, observe: 'response'});
  }

  public getTraces(keptnContext: string, projectName?: string, fromTime?: string): Observable<HttpResponse<EventResult>> {
    const url = `${this._baseUrl}/mongodb-datastore/event`;
    const params = {
      pageSize: '100',
      keptnContext,
      ... projectName && {project: projectName},
      ... fromTime && {fromTime}
    };

    return this.http
      .get<EventResult>(url, { params, observe: 'response'});
  }

  public getDeploymentsOfService(projectName: string, serviceName: string): Observable<Deployment[]> {
    return of([]);
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, fromTime?: string): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `data.project:${projectName}%20AND%20data.service:${serviceName}%20AND%20data.stage:${stageName}`,
      excludeInvalidated: 'true',
      limit: '50',
      ... fromTime && {fromTime}
    };
    return this.http
      .get<EventResult>(url, {params});
  }

  public getEvaluationResult(shkeptncontext: string): Observable<EventResult> {
    const url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`;
    const params = {
      filter: `shkeptncontext:${shkeptncontext}`,
      limit: '1'
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
      name: projectName
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
          status: 'succeeded'
        }
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
            reason
          }
        }
      });
  }

}
