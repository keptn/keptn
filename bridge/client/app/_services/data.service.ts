import {Injectable} from '@angular/core';
import {BehaviorSubject, from, Observable, Subject, timer} from "rxjs";
import {debounce, map, mergeMap, toArray} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {Project} from "../_models/project";
import {Service} from "../_models/service";

import {ApiService} from "./api.service";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private _projects = new BehaviorSubject<Project[]>(null);
  private _roots = new BehaviorSubject<Root[]>(null);
  private _rootsLastUpdated: Object = {};
  private _tracesLastUpdated: Object = {};

  private _evaluationResults = new Subject();

  constructor(private apiService: ApiService) {
    this.loadProjects();
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  get roots(): Observable<Root[]> {
    return this._roots.asObservable();
  }

  get evaluationResults(): Observable<any> {
    return this._evaluationResults;
  }

  public getRootsLastUpdated(project: Project, service: Service): Date {
    return this._rootsLastUpdated[project.projectName+":"+service.serviceName];
  }

  public getTracesLastUpdated(root: Root): Date {
    return this._tracesLastUpdated[root.shkeptncontext];
  }

  public loadProjects() {
    this.apiService.getProjects()
      .pipe(
        debounce(() => timer(10000)),
        mergeMap(projects =>
          from(projects).pipe(
            mergeMap((project) =>
              from(project.stages).pipe(
                mergeMap(
                  stage => this.apiService.getServices(project.projectName, stage.stageName)
                    .pipe(
                      map(services => services.map(service => Service.fromJSON(service))),
                      map(services => ({ ...stage, services}))
                    )
                ),
                toArray(),
                map(stages => stages.map(stage => Stage.fromJSON(stage))),
                map(stages => {
                  project.stages = project.stages.map(s => stages.find(stage => stage.stageName == s.stageName));
                  return project;
                })
              )
            ),
            toArray(),
            map(val => projects)
          )
        ),
        map(projects => projects.map(project => Project.fromJSON(project)))
      ).subscribe((projects: Project[]) => {
        this._projects.next([...this._projects.getValue() ? this._projects.getValue() : [], ...projects]);
      }, (err) => {
        this._projects.error(err);
      });
  }

  public loadRoots(project: Project, service: Service) {
    let fromTime: Date = this._rootsLastUpdated[project.projectName+":"+service.serviceName];
    this._rootsLastUpdated[project.projectName+":"+service.serviceName] = new Date();

    this.apiService.getRoots(project.projectName, service.serviceName, fromTime ? fromTime.toISOString() : null)
      .pipe(
        debounce(() => timer(10000)),
        map(response => {
          this._rootsLastUpdated[project.projectName+":"+service.serviceName] = new Date(response.headers.get("date"));
          return response.body;
        }),
        mergeMap((roots) =>
          from(roots).pipe(
            mergeMap(
              root => {
                let fromTime: Date = this._tracesLastUpdated[root.shkeptncontext];
                this._tracesLastUpdated[root.shkeptncontext] = new Date();

                return this.apiService.getTraces(root.shkeptncontext, fromTime ? fromTime.toISOString() : null)
                  .pipe(
                    map(response => {
                      this._tracesLastUpdated[root.shkeptncontext] = new Date(response.headers.get("date"));
                      return response.body;
                    }),
                    map(traces => traces.map(trace => Trace.fromJSON(trace))),
                    map(traces => ({ ...root, traces}))
                  )
              }
            ),
            toArray()
          )
        ),
        map(roots => roots.map(root => Root.fromJSON(root)))
      )
      .subscribe((roots: Root[]) => {
        service.roots = [...roots||[], ...service.roots||[]].sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());
        this._roots.next(service.roots);
      });
  }

  public loadTraces(root: Root) {
    let fromTime: Date = this._tracesLastUpdated[root.shkeptncontext];
    this._tracesLastUpdated[root.shkeptncontext] = new Date();

    this.apiService.getTraces(root.shkeptncontext, fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          this._tracesLastUpdated[root.shkeptncontext] = new Date(response.headers.get("date"));
          return response.body;
        }),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        root.traces = [...traces||[], ...root.traces||[]].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
      });
  }

  public loadEvaluationResults(evaluationData, evaluationSource) {
    let fromTime: Date;
    if(evaluationData.evaluationHistory)
      fromTime = new Date(evaluationData.evaluationHistory[evaluationData.evaluationHistory.length-1].time);

    this.apiService.getEvaluationResults(evaluationData.project, evaluationData.service, evaluationData.stage, evaluationSource, fromTime ? fromTime.toISOString() : null)
      .pipe(map(traces => traces.map(trace => Trace.fromJSON(trace))))
      .subscribe((traces: Trace[]) => {
        evaluationData.evaluationHistory = [...traces||[], ...evaluationData.evaluationHistory||[]].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
        this._evaluationResults.next(evaluationData);
      });
  }
}
