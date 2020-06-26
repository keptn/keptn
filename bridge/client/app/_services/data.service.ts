import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {BehaviorSubject, forkJoin, from, Observable, Subject, timer, of} from "rxjs";
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
  private _versionInfo = new BehaviorSubject<Object>({});
  private _rootsLastUpdated: Object = {};
  private _tracesLastUpdated: Object = {};

  private _evaluationResults = new Subject();

  constructor(private http: HttpClient, private apiService: ApiService) {
    this.loadVersionInfo();
    this.loadProjects();
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  get roots(): Observable<Root[]> {
    return this._roots.asObservable();
  }

  get versionInfo(): Observable<any> {
    return this._versionInfo.asObservable();
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

  public loadVersionInfo() {
    forkJoin({
      availableVersions: this.apiService.getAvailableVersions(),
      bridgeVersion: this.apiService.getBridgeVersion(),
      keptnVersion: this.apiService.getKeptnVersion(),
      versionCheckEnabled: of(this.apiService.isVersionCheckEnabled())
    })
    .subscribe((result) => {
      this._versionInfo.next(result);
    }, (err) => {
      this._versionInfo.error(err);
    });
  }

  public setVersionCheck(enabled: boolean) {
    this.apiService.setVersionCheck(enabled);
    this.loadVersionInfo();
  }

  public loadProjects() {
    this.apiService.getProjects()
      .pipe(
        debounce(() => timer(10000)),
        map(result => result.projects),
        mergeMap(projects =>
          from(projects).pipe(
            mergeMap((project) =>
              from(project.stages).pipe(
                mergeMap(
                  stage => this.apiService.getServices(project.projectName, stage.stageName)
                    .pipe(
                      map(result => result.services),
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
        map(result => result.events||[]),
        mergeMap((roots) =>
          from(roots).pipe(
            mergeMap(
              root => {
                let fromTime: Date = this._tracesLastUpdated[root.shkeptncontext];
                this._tracesLastUpdated[root.shkeptncontext] = new Date();

                return this.apiService.getTraces(root.shkeptncontext, root.data.project, fromTime ? fromTime.toISOString() : null)
                  .pipe(
                    map(response => {
                      this._tracesLastUpdated[root.shkeptncontext] = new Date(response.headers.get("date"));
                      return response.body;
                    }),
                    map(result => result.events||[]),
                    map(traces => traces.map(trace => Trace.fromJSON(trace))),
                    map(traces => traces.sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime())),
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

    this.apiService.getTraces(root.shkeptncontext, root.getProject(), fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          this._tracesLastUpdated[root.shkeptncontext] = new Date(response.headers.get("date"));
          return response.body;
        }),
        map(result => result.events||[]),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        root.traces = [...traces||[], ...root.traces||[]].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
      });
  }

  public loadEvaluationResults(event: Trace) {
    let fromTime: Date;
    if(event.data.evaluationHistory)
      fromTime = new Date(event.data.evaluationHistory[event.data.evaluationHistory.length-1].time);

    this.apiService.getEvaluationResults(event.data.project, event.data.service, event.data.stage, event.source, fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(result => result.events||[]),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        event.data.evaluationHistory = [...traces||[], ...event.data.evaluationHistory||[]].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
        this._evaluationResults.next(event);
      });
  }
}
