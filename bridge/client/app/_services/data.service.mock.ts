import { Injectable } from '@angular/core';
import { DataService } from './data.service';
import { ApiService } from './api.service';
import { Root } from '../_models/root';
import { Project } from '../_models/project';
import { DateUtil } from '../_utils/date.utils';
import { KeptnInfo } from './_mockData/keptnInfo.mock';
import { Projects } from './_mockData/projects.mock';
import { RootEvents } from './_mockData/roots.mock';
import { Traces } from './_mockData/traces.mock';
import { Evaluations } from './_mockData/evaluations.mock';
import { Trace } from '../_models/trace';
import { map } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DataServiceMock extends DataService {
  constructor(apiService: ApiService) {
    super(apiService);
  }

  public loadKeptnInfo() {
    this._keptnInfo.next(KeptnInfo);
  }

  public loadProjects() {
    this._projects.next(Projects);
  }

  public loadProject(projectName: string) {
    this._projects.next([...Projects]);
  }

  public getProject(projectName: string): Observable<Project | undefined> {
    this.loadProjects();
    return this._projects.pipe(map(projects => {
      return projects?.find(project => project.projectName === projectName);
    }));
  }

  public deleteProject(projectName: string): Observable<object> {
    return of({});
  }

  public loadTraces(sequence: Sequence) {
    sequence.traces = [...Traces || [], ...sequence.traces || []];
    this._sequences.next([...this._sequences.getValue() || []]);
  }

  public loadTracesByContext(shkeptncontext: string) {
    this._traces.next(Traces.filter(t => t.shkeptncontext === shkeptncontext));
  }

  public loadEvaluationResults(event: Trace) {
    this._evaluationResults.next({
      type: 'evaluationHistory',
      triggerEvent: event,
      traces: [Evaluations]
    });
  }

  public setGitUpstreamUrl(projectName: string, gitUrl: string, gitUser: string, gitToken: string): Observable<boolean> {
    this.loadProjects();
    return of(true);
  }
}

