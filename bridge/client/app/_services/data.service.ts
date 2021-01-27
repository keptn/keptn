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

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private _projects = new BehaviorSubject<Project[]>(null);
  private _roots = new BehaviorSubject<Root[]>(null);
  private _openApprovals = new BehaviorSubject<Trace[]>([]);
  private _keptnInfo = new BehaviorSubject<Object>(null);
  private _rootsLastUpdated: Object = {};
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

  get openApprovals(): Observable<Trace[]> {
    return this._openApprovals.asObservable();
  }

  get keptnInfo(): Observable<any> {
    return this._keptnInfo.asObservable();
  }

  get evaluationResults(): Observable<any> {
    return this._evaluationResults;
  }

  public getRootsLastUpdated(project: Project): Date {
    return this._rootsLastUpdated[project.projectName];
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

  public loadRoots(project: Project) {
    let fromTime: Date = this._rootsLastUpdated[project.projectName];
    this._rootsLastUpdated[project.projectName] = new Date();

    this.apiService.getRoots(project.projectName, null, fromTime ? fromTime.toISOString() : null)
      .pipe(
        debounce(() => timer(10000)),
        map(response => {
          this._rootsLastUpdated[project.projectName] = new Date(response.headers.get("date"));
          return response.body;
        }),
        map(result => result.events||[]),
        mergeMap((roots) =>
          from(roots).pipe(
            mergeMap(
              root => {
                return this.apiService.getTraces(root.shkeptncontext, root.data.project)
                  .pipe(
                    map(response => {
                      this._tracesLastUpdated[root.shkeptncontext] = new Date(response.headers.get("date"));
                      return response.body;
                    }),
                    map(result => result.events||[]),
                    map(this.traceMapper),
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
        project.sequences = [...roots||[], ...project.sequences||[]].sort(DateUtil.compareTraceTimesAsc);
        project.getServices().forEach(service => {
          service.roots = project.sequences.filter(s => s.getService() == service.serviceName);
        });
        this._roots.next(project.sequences);
        roots.forEach(root => {
          this.updateApprovals(root);
        })
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
        root.traces = this.traceMapper([...traces||[], ...root.traces||[]]);
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
        map(traces => traces.map(trace => Trace.fromJSON(trace))),
        map( traces => traces.map((trace: Trace) => {
          if(trace.data.evaluation.indicatorResults){
            trace.data.evaluation.indicatorResults.sort( (resultA, resultB) => resultA.value.metric.localeCompare(resultB.value.metric))
          }
          return trace;
        }))
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

  private updateApprovals(root: Root) {
    if(root.traces.length > 0) {
      this._openApprovals.next(this._openApprovals.getValue().filter(approval => root.traces.indexOf(approval) < 0));
      const approvals = root.getPendingApprovals();
      if (approvals.length !== 0)
        this._openApprovals.next([...this._openApprovals.getValue(), ...approvals].sort(DateUtil.compareTraceTimesAsc));
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

  private traceMapper(traces: Trace[]) {
    return traces
      .map(trace => Trace.fromJSON(trace))
      .sort(DateUtil.compareTraceTimesDesc)
      .reduce((result: Trace[], trace) => {
        if(trace.triggeredid) {
          let trigger = result.find(t => t.id == trace.triggeredid);
          if(trigger)
            trigger.traces.push(trace);
          else
            result.push(trace);
        } else
          result.push(trace);
        return result;
      }, []);
  }
}
