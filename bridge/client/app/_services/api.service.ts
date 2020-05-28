import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders, HttpResponse} from "@angular/common/http";
import {Observable, throwError, of} from "rxjs";
import {catchError} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Project} from "../_models/project";
import {Resource} from "../_models/resource";
import {Stage} from "../_models/stage";
import {Service} from "../_models/service";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private baseUrl: string;
  private defaultHeaders: HttpHeaders = new HttpHeaders({'Content-Type': 'application/json'});

  private VERSION_CHECK_COOKIE = 'keptn_versioncheck';

  set baseUrl(value: string) {
    this._baseUrl = value;
  }

  constructor(private http: HttpClient) {
    this._baseUrl = `/api`;
  }

  public getBridgeVersion(): Observable<any> {
    let url = `${this.baseUrl}/api/`;
    return this.http
      .get<any>(url, { headers: this.defaultHeaders })
      .pipe(
        catchError(this.handleError<any>('getBridgeVersion')),
        map(res => res.version),
      );
  }

  public getKeptnVersion(): Observable<any> {
    let url = `${this.baseUrl}/api/swagger-ui/swagger.yaml`;
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
      let url = `${this.baseUrl}/api/version.json`;
      return this.http
        .get<any>(url, { headers: this.defaultHeaders })
        .pipe(catchError(this.handleError<any>('getAvailableVersions')));
    } else {
      return of(null);
    }
  }

  public getProjects(): Observable<Project[]> {
    let url = `${this._baseUrl}/configuration-service/v1/project?DisableUpstreamSync=true`;
    return this.http
      .get<Project[]>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<Project[]>('getProjects')));
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

  public getServices(projectName, stageName): Observable<Service[]> {
    let url = `${this._baseUrl}/configuration-service/v1/project/${projectName}/stage/${stageName}/service`;
    return this.http
      .get<Service[]>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<Service[]>('getServices')));
  }

  public getRoots(projectName: string, serviceName: string, fromTime?: String): Observable<HttpResponse<Root[]>> {
    let url = `${this._baseUrl}/mongodb-datastore/event?pageSize=20&project=${projectName}&service=${serviceName}`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<Root[]>(url, { headers: this.defaultHeaders, observe: 'response' })
      .pipe(catchError(this.handleError<HttpResponse<Root[]>>('getRoots')));
  }

  public getTraces(contextId: string, projectName?: string, fromTime?: String): Observable<HttpResponse<Trace[]>> {
    let url = `${this._baseUrl}/mongodb-datastore/v1/event?pageSize=100&keptnContext=${contextId}`;
    if(projectName)
      url += `&project=${projectName}`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<Trace[]>(url, { headers: this.defaultHeaders, observe: 'response' })
      .pipe(catchError(this.handleError<HttpResponse<Trace[]>>('getTraces')));
  }

  public getEvaluationResults(projectName: string, serviceName: string, stageName: string, source: string, fromTime?: String) {
    let url = `${this._baseUrl}/mongodb-datastore/v1/event?type=sh.keptn.events.evaluation-done&project=${projectName}&service=${serviceName}&stage=${stageName}&source=${source}&pageSize=50`;
    if(fromTime)
      url += `&fromTime=${fromTime}`;
    return this.http
      .get<Trace[]>(url, { headers: this.defaultHeaders })
      .pipe(catchError(this.handleError<Trace[]>('getEvaluationResults')));
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
