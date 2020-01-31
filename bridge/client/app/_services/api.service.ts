import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {BehaviorSubject, Observable, of, throwError} from "rxjs";
import {catchError, map, retry} from "rxjs/operators";

import {environment} from "../../environments/environment";

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

  private baseUrl: string = environment.apiUrl;
  private headers: HttpHeaders = new HttpHeaders({'Content-Type': 'application/json'});

  constructor(private http: HttpClient) {
  }

  public getProjects(): Observable<Project[]> {
    let url = `${this.baseUrl}/api/project`;
    return this.http
      .get<Project[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Project[]>('getProjects')));
  }

  public getProjectResources(projectName): Observable<Resource[]> {
    let url = `${this.baseUrl}/api/project/${projectName}/resource`;
    return this.http
      .get<Resource[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Resource[]>('getProjectResources')));
  }

  public getStages(projectName): Observable<Stage[]> {
    let url = `${this.baseUrl}/api/project/${projectName}/stage`;
    return this.http
      .get<Stage[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Stage[]>('getStages')));
  }

  public getServices(projectName, stageName): Observable<Service[]> {
    let url = `${this.baseUrl}/api/project/${projectName}/stage/${stageName}/service`;
    return this.http
      .get<Service[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Service[]>('getServices')));
  }

  public getRoots(projectName: string, serviceName: string, fromTime?: String): Observable<Root[]> {
    let url = `${this.baseUrl}/api/roots/${projectName}/${serviceName}`;
    if(fromTime)
      url += `?fromTime=${fromTime}`;
    return this.http
      .get<Root[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Root[]>('getRoots')));
  }

  public getTraces(contextId: string, fromTime?: String): Observable<Trace[]> {
    let url = `${this.baseUrl}/api/traces/${contextId}`;
    if(fromTime)
      url += `?fromTime=${fromTime}`;
    return this.http
      .get<Trace[]>(url, { headers: this.headers })
      .pipe(catchError(this.handleError<Trace[]>('getTraces')));
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
