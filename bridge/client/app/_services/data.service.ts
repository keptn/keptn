import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {BehaviorSubject, forkJoin, from, Observable, Subject, of} from "rxjs";
import {filter, map, mergeMap, take, toArray} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {Project} from "../_models/project";
import {Service} from "../_models/service";
import {EventTypes} from "../_models/event-types";

import {ApiService} from "./api.service";
import {DateUtil} from "../_utils/date.utils";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private _projects = new BehaviorSubject<Project[]>(null);
  private _roots = new BehaviorSubject<Root[]>(null);
  private _openApprovals = new BehaviorSubject<Trace[]>([]);
  private _keptnInfo = new BehaviorSubject<any>(null);
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

  public getProject(projectName): Observable<Project> {
    return this.projects.pipe(
      map(projects => projects ? projects.find(project => {
        return project.projectName === projectName;
      }) : null)
    );
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
      versionCheckEnabled: of(this.apiService.isVersionCheckEnabled()),
      metadata: this.apiService.getMetadata()
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
    this.apiService.getProjects(this._keptnInfo.getValue().bridgeInfo.projectsPageSize||50)
      .pipe(
        map(result => result.projects),
        map(projects => projects.map(project => Project.fromJSON(project)))
      ).subscribe((projects: Project[]) => {
      this._projects.next(projects);
    }, (err) => {
      this._projects.next([]);
    });
  }

  public loadProject(projectName) {
    this.apiService.getProject(projectName)
      .pipe(
        map(project => Project.fromJSON(project))
      ).subscribe((project: Project) => {
        let projects = this._projects.getValue();
        let index = projects.findIndex(p => p.projectName == projectName);
        projects.splice((index < 0) ? projects.length : index, (index < 0) ? 0 : 1, project);
        this._projects.next([...projects]);
    });
  }

  public loadRoots(project: Project) {
    let fromTime: Date = this._rootsLastUpdated[project.projectName];
    this._rootsLastUpdated[project.projectName] = new Date();

    from(project.services).pipe(
      mergeMap(
        service => this.apiService.getRoots(project.projectName, service.serviceName, fromTime ? fromTime.toISOString() : null)
          .pipe(
            map(result => result.body.events||[]),
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
      ),
      toArray(),
      map(roots => roots.reduce((result, roots) => result.concat(roots), []))
    ).subscribe((roots: Root[]) => {
      project.sequences = [...roots||[], ...project.sequences||[]].sort(DateUtil.compareTraceTimesAsc);
      project.getServices().forEach(service => {
        service.roots = project.sequences.filter(s => s.getService() == service.serviceName);
      });
      this._roots.next(project.sequences);
      roots.forEach(root => {
        this.updateApprovals(root);
      });
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
        this._roots.next([...this._roots.getValue()]);
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
    this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_STARTED, 'approval.started')
      .pipe(
        mergeMap(()=> this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_FINISHED, 'approval.finished'))
      )
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

  public getOpenApprovals(project: Project, stage: Stage, service?: Service): Trace[]{
    return this._openApprovals.getValue().filter(approval => approval.data.project === project.projectName && approval.data.stage === stage.stageName && (!service || approval.data.service === service.serviceName));
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
      .reduce((result: Trace[], trace: Trace) => {
        let trigger = result.find(t => {
          if(trace.triggeredid)
            return t.id == trace.triggeredid;
        });

        if(trigger)
          trigger.traces.push(trace);
        else
          result.push(trace);

        return result;
      }, []);
  }
}
