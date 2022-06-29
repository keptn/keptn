import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { IProjectResult } from '../../../shared/interfaces/project-result';

const baseUrl = `./api`;

@Injectable({
  providedIn: 'root',
})
export class ProjectService {
  constructor(protected http: HttpClient) {}

  public getProjects(pageSize?: number): Observable<IProjectResult> {
    const url = `${baseUrl}/controlPlane/v1/project`;
    const params = {
      disableUpstreamSync: 'true',
      ...(pageSize && { pageSize: pageSize.toString() }),
    };
    return this.http.get<IProjectResult>(url, { params });
  }
}
