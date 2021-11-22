/* eslint-disable @typescript-eslint/no-unused-vars */
import { Injectable } from '@angular/core';
import { DataService } from './data.service';
import { Project } from '../_models/project';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class DataServiceMock extends DataService {
  public getProject(projectName: string): Observable<Project | undefined> {
    if (!this._projects.getValue()?.length) {
      this.loadProjects();
    }
    return this._projects.pipe(map((projects) => projects?.find((project) => project.projectName === projectName)));
  }
}
/* eslint-enable @typescript-eslint/no-unused-vars */
