import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders, HttpResponse} from "@angular/common/http";
import {Observable, of} from "rxjs";
import {map} from "rxjs/operators";

import {Resource} from "../_models/resource";
import {Stage} from "../_models/stage";
import {ProjectResult} from "../_models/project-result";
import {ServiceResult} from "../_models/service-result";
import {EventResult} from "../_models/event-result";
import {Trace} from "../_models/trace";
import {ApprovalStates} from "../_models/approval-states";
import {EventTypes} from "../_models/event-types";
import {Metadata} from '../_models/metadata';
import {Project} from "../_models/project";
import {KeptnService} from '../_models/keptn-service';
import {KeptnServicesMock} from '../_models/keptn-services.mock';
import {TaskNames} from '../_models/task-names.mock';
import {Deployment} from '../_models/deployment';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private _baseUrl: string;
  private VERSION_CHECK_COOKIE = 'keptn_versioncheck';

  constructor(private http: HttpClient) {
    this._baseUrl = `./api`;
  }

  public set baseUrl(value: string) {
    this._baseUrl = value;
  }

  public get baseUrl() {
    return this._baseUrl;
  }

  public getKeptnInfo(): Observable<any> {
    let url = `${this._baseUrl}/bridgeInfo`;
    return this.http
      .get<any>(url)
  }

  public getKeptnVersion(): Observable<any> {
    let url = `${this._baseUrl}/swagger-ui/swagger.yaml`;
    return this.http
      .get<any>(url, { headers: new HttpHeaders({'Access-Control-Allow-Origin': '*'}) })
      .pipe(
        map(res => res.toString()),
        map(res => res.substring(res.lastIndexOf("version: ")+9)),
        map(res => res.substring(0, res.indexOf("\n"))),
      );
  }

  public getIntegrationsPage(): Observable<any> {
    let url = `${this._baseUrl}/integrationsPage`;
    return this.http
      .get<any>(url, { responseType: 'text' as 'json' });
  }

  public isVersionCheckEnabled(): boolean | undefined {
    const versionInfo = JSON.parse(localStorage.getItem(this.VERSION_CHECK_COOKIE));
    let enabled = typeof versionInfo === 'boolean' ? versionInfo : versionInfo?.enabled; // support old format
    if (!enabled && (!versionInfo?.time || moment().subtract(5, 'days').isAfter(versionInfo.time))) {
      enabled = undefined;
    }
    return enabled;
  }

  public setVersionCheck(enabled: boolean): void {
    localStorage.setItem(this.VERSION_CHECK_COOKIE, JSON.stringify({enabled, time: moment().valueOf()}));
  }

  public getAvailableVersions(): Observable<any> {
    if(this.isVersionCheckEnabled()) {
      let url = `${this._baseUrl}/version.json`;
      return this.http
        .get<any>(url);
    } else {
      return of(null);
    }
  }

  public getProjects(pageSize?: number): Observable<ProjectResult> {
    let url = `${this._baseUrl}/controlPlane/v1/project?disableUpstreamSync=true`;
    if(pageSize)
      url += `&pageSize=${pageSize}`;
    return this.http
      .get<ProjectResult>(url);
  }

  public getProject(projectName: string): Observable<Project> {
    let url = `${this._baseUrl}/controlPlane/v1/project/${projectName}`;
    return this.http.get<Project>(url);
  }

  public getKeptnServices(projectName: string): Observable<KeptnService[]> {
    return of(KeptnServicesMock);
  }

  public getMetadata(): Observable<Metadata> {
    return this.http.get<Metadata>(`${this._baseUrl}/v1/metadata`);
  }

  public getProjectResources(projectName): Observable<Resource[]> {
    let url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/resource`;
    return this.http
      .get<Resource[]>(url);
  }

  public getTaskNames(projectName: string): Observable<string[]>{
    return of(TaskNames);
  }

  public getStages(projectName): Observable<Stage[]> {
    let url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage`;
    return this.http
      .get<Stage[]>(url);
  }

  public getServices(projectName: string, stageName: string, pageSize: number): Observable<ServiceResult> {
    let url = `${this._baseUrl}/controlPlane/v1/project/${projectName}/stage/${stageName}/service?pageSize=${pageSize}`;
    return this.http
      .get<ServiceResult>(url);
  }

  public getRoots(projectName: string, pageSize: number, serviceName?: string, fromTime?: string, beforeTime?: string, shkeptncontext?: string): Observable<HttpResponse<EventResult>> {
    let url = `${this._baseUrl}/mongodb-datastore/event?root=true&pageSize=${pageSize}&project=${projectName}`;
    if (serviceName) {
      url += `&service=${serviceName}`;
    }
    if (fromTime) {
      url += `&fromTime=${fromTime}`;
    }
    if (beforeTime) {
      url += `&beforeTime=${beforeTime}`;
    }
    if (shkeptncontext) {
      url += `&keptnContext=${shkeptncontext}`;
    }

    return this.http
      .get<EventResult>(url, { observe: 'response' });
  }

  public getTraces(contextId: string, projectName?: string, fromTime?: string): Observable<HttpResponse<EventResult>> {
    let url = `${this._baseUrl}/mongodb-datastore/event?pageSize=100&keptnContext=${contextId}`;
    if(projectName)
      url += `&project=${projectName}`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<EventResult>(url, { observe: 'response' });
  }

  public getDeploymentsOfService(projectName: string, serviceName: string): Observable<Deployment[]> {
    return of([]);
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, source: string, fromTime?: string) {
    let url = `${this._baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}?filter=data.project:${projectName}%20AND%20data.service:${serviceName}%20AND%20data.stage:${stageName}&excludeInvalidated=true&limit=50`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<EventResult>(url);
  }

  public sendApprovalEvent(approval: Trace, approve: boolean, eventType: EventTypes, source: string ) {
    const url = `${this._baseUrl}/v1/event`;

    return this.http
      .post<any>(url, {
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

  public sendEvaluationInvalidated(evaluation: Trace, reason: string) {
    let url = `${this._baseUrl}/v1/event`;

    return this.http
      .post<any>(url, {
        "shkeptncontext": evaluation.shkeptncontext,
        "type": EventTypes.EVALUATION_INVALIDATED,
        "triggeredid": evaluation.triggeredid,
        "source": "https://github.com/keptn/keptn/bridge#evaluation.invalidated",
        "data": {
          "project": evaluation.data.project,
          "stage": evaluation.data.stage,
          "service": evaluation.data.service,
          "evaluation": {
            "reason": reason
          }
        }
      });
  }

}
