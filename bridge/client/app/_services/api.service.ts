import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders, HttpResponse} from "@angular/common/http";
import {Observable, throwError, of} from "rxjs";
import {catchError, map} from "rxjs/operators";

import {Resource} from "../_models/resource";
import {Stage} from "../_models/stage";
import {ProjectResult} from "../_models/project-result";
import {ServiceResult} from "../_models/service-result";
import {EventResult} from "../_models/event-result";
import {Trace} from "../_models/trace";
import {ApprovalStates} from "../_models/approval-states";
import {EventTypes} from "../_models/event-types";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private _baseUrl: string;
  private defaultHeaders: HttpHeaders = new HttpHeaders({'Content-Type': 'application/json'});

  private VERSION_CHECK_COOKIE = 'keptn_versioncheck';

  set baseUrl(value: string) {
    this._baseUrl = value;
  }

  constructor(private http: HttpClient) {
    this._baseUrl = `/bridge/api`;
  }

  public getBridgeVersion(): Observable<any> {
    let url = `${this._baseUrl}/`;
    return this.http
      .get<any>(url, { headers: this.defaultHeaders })
      .pipe(
        catchError(this.handleError<any>('getBridgeVersion')),
        map(res => res.version),
      );
  }

  public getKeptnVersion(): Observable<any> {
    let url = `${this._baseUrl}/swagger-ui/swagger.yaml`;
    return this.http
      .get<any>(url, { headers: this.defaultHeaders.append('Access-Control-Allow-Origin', '*') })
      .pipe(
        catchError(this.handleError<any>('getKeptnVersion')),
        map(res => res.toString()),
        map(res => res.substring(res.lastIndexOf("version: ")+9)),
        map(res => res.substring(0, res.indexOf("\n"))),
      );
  }

  public isVersionCheckEnabled(): boolean {
    return JSON.parse(localStorage.getItem(this.VERSION_CHECK_COOKIE));
  }

  public setVersionCheck(enabled: boolean): void {
    localStorage.setItem(this.VERSION_CHECK_COOKIE, String(enabled));
  }

  public getAvailableVersions(): Observable<any> {
    if(this.isVersionCheckEnabled()) {
      let url = `${this._baseUrl}/version.json`;
      return this.http
        .get<any>(url, { headers: this.defaultHeaders })
        .pipe(catchError(this.handleError<any>('getAvailableVersions')));
    } else {
      return of(null);
    }
  }

  public getProjects(): Observable<ProjectResult> {
    let url = `${this._baseUrl}/configuration-service/v1/project?disableUpstreamSync=true`;
    return this.http
      .get<ProjectResult>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<ProjectResult>('getProjects')));
  }

  public getProjectResources(projectName): Observable<Resource[]> {
    let url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/resource`;
    return this.http
      .get<Resource[]>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<Resource[]>('getProjectResources')));
  }

  public getStages(projectName): Observable<Stage[]> {
    let url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/stage`;
    return this.http
      .get<Stage[]>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<Stage[]>('getStages')));
  }

  public getServices(projectName, stageName): Observable<ServiceResult> {
    let url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service`;
    return this.http
      .get<ServiceResult>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<ServiceResult>('getServices')));
  }

  public getRoots(projectName: string, serviceName: string, fromTime?: string): Observable<HttpResponse<EventResult>> {
    let url = `${this._baseUrl}/mongodb-datastore/event?root=true&pageSize=20&project=${projectName}&service=${serviceName}`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<EventResult>(url, { headers: this.defaultHeaders, observe: 'response' })
      .pipe(catchError(this.handleError<HttpResponse<EventResult>>('getRoots')));
  }

  public getTraces(contextId: string, projectName?: string, fromTime?: string): Observable<HttpResponse<EventResult>> {
    let url = `${this._baseUrl}/mongodb-datastore/event?pageSize=100&keptnContext=${contextId}`;
    if(projectName)
      url += `&project=${projectName}`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<EventResult>(url, { headers: this.defaultHeaders, observe: 'response' })
      .pipe(catchError(this.handleError<HttpResponse<EventResult>>('getTraces')));
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, source: string, fromTime?: string) {
    let url = `${this._baseUrl}/mongodb-datastore/event?type=sh.keptn.events.evaluation-done&project=${projectName}&service=${serviceName}&stage=${stageName}&source=${source}&pageSize=50`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<EventResult>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<EventResult>('getEvaluationResults')));
  }

  public sendApprovalEvent(approval: Trace, approve: boolean) {
    let url = `${this._baseUrl}/v1/event`;
    return this.http
      .post<any>(url, {
        "shkeptncontext": approval.shkeptncontext,
        "type": EventTypes.APPROVAL_FINISHED,
        "triggeredid": approval.id,
        "source": "https://github.com/keptn/keptn/bridge#approval.finished",
        "data": Object.assign(approval.data, {
          "approval": {
            "result": approve ? ApprovalStates.APPROVED : ApprovalStates.DECLINED,
            "status": "succeeded"
          }
        })
      }, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<any>('sendApprovalEvent')));
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      // TODO: handel error and show to the user?!
      this.log(`${operation} failed: ${error.message}`);
      return throwError(error);
    };
  }

  private log(message: string) {
    console.log(message);
  }

}
