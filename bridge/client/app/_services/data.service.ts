import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {BehaviorSubject, forkJoin, from, Observable, Subject, timer, of} from "rxjs";
import {debounce, filter, map, mergeMap, take, toArray} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {Project} from "../_models/project";
import {Service} from "../_models/service";

import {ApiService} from "./api.service";
import {EventTypes} from "../_models/event-types";
import DateUtil from "../_utils/date.utils";
import {Sequence} from "../_models/sequence";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private _projects = new BehaviorSubject<Project[]>(null);
  private _roots = new BehaviorSubject<Root[]>(null);
  private _sequences = new BehaviorSubject<Sequence[]>(null);
  private _openApprovals = new BehaviorSubject<Trace[]>([]);
  private _keptnInfo = new BehaviorSubject<Object>(null);
  private _rootsLastUpdated: Object = {};
  private _sequencesLastUpdated: Object = {};
  private _tracesLastUpdated: Object = {};

  private _evaluationResults = new Subject();

  constructor(private http: HttpClient, private apiService: ApiService) {
    this.loadKeptnInfo();
    this.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(take(1))
      .subscribe(keptnInfo => {
        this.loadProjects();
      });
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  get roots(): Observable<Root[]> {
    return this._roots.asObservable();
  }

  get sequences(): Observable<Sequence[]> {
    return this._sequences.asObservable();
  }

  get openApprovals(): Observable<Trace[]> {
    return this._openApprovals.asObservable();
  }

  get keptnInfo(): Observable<any> {
    return this._keptnInfo.asObservable();
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

  public loadKeptnInfo() {
    forkJoin({
      availableVersions: this.apiService.getAvailableVersions(),
      bridgeInfo: this.apiService.getKeptnInfo(),
      keptnVersion: this.apiService.getKeptnVersion(),
      versionCheckEnabled: of(this.apiService.isVersionCheckEnabled())
    }).subscribe((result) => {
        if(result.bridgeInfo.showApiToken) {
          if(window.location.href.indexOf('bridge') != -1)
            result.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf('/bridge'))}/api`;
          else
            result.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf(window.location.pathname))}/api`;

          result.bridgeInfo.authCommand = `keptn auth --endpoint=${result.bridgeInfo.apiUrl} --api-token=${result.bridgeInfo.apiToken}`;
        }
        this._keptnInfo.next(result);
      }, (err) => {
        this._keptnInfo.error(err);
      });
  }

  public setVersionCheck(enabled: boolean) {
    this.apiService.setVersionCheck(enabled);
    this.loadKeptnInfo();
  }

  public loadProjects() {
    // @ts-ignore
    this.apiService.getProjects(this._keptnInfo.getValue().bridgeInfo.projectsPageSize||50)
      .pipe(
        debounce(() => timer(10000)),
        map(result => result.projects),
        mergeMap(projects =>
          from(projects).pipe(
            mergeMap((project) =>
              from(project.stages).pipe(
                mergeMap(
                  // @ts-ignore
                  stage => this.apiService.getServices(project.projectName, stage.stageName, this._keptnInfo.getValue().bridgeInfo.servicesPageSize||50)
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
      this._projects.next([]);
    });
  }

  public loadServices(project: Project) {
    from(project.stages).pipe(
      mergeMap(
        // @ts-ignore
        stage => this.apiService.getServices(project.projectName, stage.stageName, this._keptnInfo.getValue().bridgeInfo.servicesPageSize||50)
          .pipe(
            map(result => result.services),
            map(services => services.map(service => Service.fromJSON(service))),
            map(services => ({ ...stage, services}))
          )
      ),
      toArray(),
      map(stages => stages.map(stage => Stage.fromJSON(stage)))
    ).subscribe((stages: Stage[]) => {
      project.stages.forEach((stage: Stage) => {
        stage.services.forEach((service: Service) => {
          service.deployedImage = stages.find(s => s.stageName == stage.stageName).services.find(s => s.serviceName == service.serviceName).deployedImage;
        });
      });
    }, (err) => {
      this._projects.next([]);
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
        service.roots = [...roots||[], ...service.roots||[]].sort(DateUtil.compareTraceTimes);
        this._roots.next(service.roots);
        roots.forEach(root => {
          this.updateApprovals(root);
        })
      });
  }

  public loadSequences(project: Project) {
    let fromTime: Date = this._sequencesLastUpdated[project.projectName];
    this._sequencesLastUpdated[project.projectName] = new Date();

    this.apiService.getRoots(project.projectName, null, fromTime ? fromTime.toISOString() : null)
      .pipe(
        debounce(() => timer(10000)),
        map(response => {
          this._sequencesLastUpdated[project.projectName] = new Date(response.headers.get("date"));
          return response.body;
        }),
        map(result => result.events||[]),
        mergeMap((sequences) =>
          from(sequences).pipe(
            mergeMap(
              sequence => {
                let fromTime: Date = this._tracesLastUpdated[sequence.shkeptncontext];
                this._tracesLastUpdated[sequence.shkeptncontext] = new Date();

                return this.apiService.getTraces(sequence.shkeptncontext, sequence.data.project, fromTime ? fromTime.toISOString() : null)
                  .pipe(
                    map(response => {
                      this._tracesLastUpdated[sequence.shkeptncontext] = new Date(response.headers.get("date"));
                      return response.body;
                    }),
                    map(result => result.events||[]),
                    map(traces => traces.map(trace => Trace.fromJSON(trace))),
                    map(traces => traces.sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime())),
                    map(traces => ({ ...sequence, traces}))
                  )
              }
            ),
            toArray()
          )
        ),
        map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)))
      )
      .subscribe((sequences: Sequence[]) => {
        project.sequences = [...sequences||[], ...project.sequences||[]].sort(DateUtil.compareTraceTimes);
        this._sequences.next(sequences);
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
        this.updateApprovals(root);
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
        this._evaluationResults.next({
          type: "evaluationHistory",
          triggerEvent: event,
          traces: traces
        });
      });
  }

  public sendApprovalEvent(approval: Trace, approve: boolean) {
    this.apiService.sendApprovalEvent(approval, approve)
      .pipe(take(1))
      .subscribe(() => {
        let root = this._projects.getValue().find(p => p.projectName == approval.data.project).services.find(s => s.serviceName == approval.data.service).roots.find(r => r.shkeptncontext == approval.shkeptncontext);
        this.loadTraces(root);
      });
  }

  private updateApprovals(root) {
    if(root.traces.length > 0) {
      this._openApprovals.next(this._openApprovals.getValue().filter(approval => root.traces.indexOf(approval) < 0));
      if(root.traces[root.traces.length-1].type == EventTypes.APPROVAL_TRIGGERED)
        this._openApprovals.next([...this._openApprovals.getValue(), root.traces[root.traces.length-1]].sort(DateUtil.compareTraceTimes));
    }
  }

  public invalidateEvaluation(evaluation: Trace, reason: string) {
    this.apiService.sendEvaluationInvalidated(evaluation, reason)
      .pipe(take(1))
      .subscribe(() => {
        this._evaluationResults.next({
          type: "invalidateEvaluation",
          triggerEvent: evaluation
        });
      });
  }
}
