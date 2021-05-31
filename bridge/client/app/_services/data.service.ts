import {Injectable} from '@angular/core';
import {BehaviorSubject, forkJoin, from, Observable, Subject, of} from "rxjs";
import {catchError, map, mergeMap, switchMap, take, toArray} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {Project} from "../_models/project";
import {Service} from "../_models/service";
import {EventTypes} from "../_models/event-types";

import {ApiService} from "./api.service";
import {DateUtil} from "../_utils/date.utils";

import * as moment from 'moment';
import {KeptnService} from '../_models/keptn-service';
import {Deployment} from '../_models/deployment';
import {Sequence} from '../_models/sequence';
import {Resource} from '../_models/resource';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  protected _projects = new BehaviorSubject<Project[]>(null);
  protected _taskNames = new BehaviorSubject<string[]>([]);
  protected _roots = new BehaviorSubject<Root[]>(null);
  protected _traces = new BehaviorSubject<Trace[]>(null);
  protected _openApprovals = new BehaviorSubject<Trace[]>([]);
  protected _keptnInfo = new BehaviorSubject<any>(null);
  public _remediationsUpdated: Subject<void> = new Subject<void>();
  protected _rootsLastUpdated: Object = {};
  protected _tracesLastUpdated: Object = {};
  private readonly DEFAULT_SEQUENCE_PAGE_SIZE = 25;
  private readonly DEFAULT_NEXT_SEQUENCE_PAGE_SIZE = 10;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;

  protected _evaluationResults = new Subject();

  constructor(private apiService: ApiService) {
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  get taskNames(): Observable<string[]> {
    return  this._taskNames.asObservable();
  }

  get taskNamesTriggered(): Observable<string[]> {
    return this._taskNames.pipe(
      map(tasks => tasks.map(task => task + '.triggered'))
    );
  }

  get roots(): Observable<Root[]> {
    return this._roots.asObservable();
  }

  get traces(): Observable<Trace[]> {
    return this._traces.asObservable();
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

  public getKeptnServices(projectName: string): Observable<KeptnService[]> {
    return this.apiService.getKeptnServices(projectName).pipe(
      map(services => services.map(service => KeptnService.fromJSON(service))),
      map(services => services.sort((serviceA, serviceB) => serviceA.name.localeCompare(serviceB.name)))
    );
  }

  public getRootsLastUpdated(project: Project): Date {
    return this._rootsLastUpdated[project.projectName];
  }

  public getTracesLastUpdated(root: Root): Date {
    return this._tracesLastUpdated[root.shkeptncontext];
  }

  public setGitUpstreamUrl(projectName: string, gitUrl: string, gitUser: string, gitToken: string): Observable<boolean> {
    return this.apiService.sendGitUpstreamUrl(projectName, gitUrl, gitUser, gitToken).pipe(map(res => {
      this.loadProjects();
      return true;
    }), catchError((err) => {
      return of(false);
    }));
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
        map(projects =>
          projects.map(project => {
            project.stages = project.stages.map(stage => {
              stage.services = stage.services.map(service => {
                service.stage = stage.stageName;
                return Service.fromJSON(service);
              });
              return Stage.fromJSON(stage);
            });
            return Project.fromJSON(project);
          })
        )
      ).subscribe((projects: Project[]) => {
      this._projects.next(projects);
    }, (err) => {
      this._projects.next([]);
    });
  }

  public loadOpenRemediations(project: Project): void {
    this.apiService.getOpenRemediations(project.projectName).pipe(
      map(response => response.body),
      map(sequenceResult => sequenceResult.states),
      map(sequences => {
        // remove finished remediations
        for (const service of project.getServices()){
          for (const deployment of service.deployments) {
            for (const stage of deployment.stages) {
              stage.remediations = stage.remediations.filter(r => sequences.some(s => s.shkeptncontext === r.shkeptncontext));
            }
          }
        }
        return sequences;
      }),
      mergeMap((sequences) =>
        from(sequences).pipe(
          mergeMap((sequence) => {
            const service = project.getService(sequence.service);
            const sequenceStage = sequence.stages[0].name;
            let result = of(null);
            if (service) {
              const deployment = service.deployments.find(d => d.stages.some(stage => sequence.stages.some(s => s.name === stage.stageName)));
              if (deployment) {
                const stage = deployment.stages.find(s => s.stageName === sequenceStage);
                if (stage) {
                  const existingRemediation = stage.remediations.find(r => r.shkeptncontext === sequence.shkeptncontext);
                  let _root: Observable<any> = of(null);
                  let _resourceContent: Observable<any> = of(null);

                  // update existing remediation
                  if (existingRemediation) {
                    Object.assign(existingRemediation, Sequence.fromJSON(sequence));
                  }
                  else {
                    const remediation = Sequence.fromJSON(sequence);
                    stage.remediations.push(remediation);
                    if (!remediation.problemTitle) {
                      _root = this.getRoot(project.projectName, remediation.shkeptncontext).pipe(
                        map(root => {
                          remediation.problemTitle = root.getProblemTitle();
                        }));
                    }
                  }

                  if (!stage?.config) {
                    _resourceContent = this.apiService.getServiceResource(project.projectName, sequenceStage, deployment.service, 'remediation.yaml').pipe(
                      map(resource => {
                        stage.config = atob(resource.resourceContent);
                        return stage;
                      })
                    );
                  }
                  result = forkJoin([_root, _resourceContent]);
                }
              }
            }
            return result;
          }),
          toArray()
        )
      )
    ).subscribe(() => {
      this._remediationsUpdated.next();
    });
  }

  public loadRoots(project: Project) {
    let fromTime: Date = this._rootsLastUpdated[project.projectName];
    this._rootsLastUpdated[project.projectName] = new Date();

    this.apiService.getRoots(project.projectName, this.DEFAULT_SEQUENCE_PAGE_SIZE, null, fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          let lastUpdated = moment(response.headers.get("date"));
          let lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
          this._rootsLastUpdated[project.projectName] = lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated;
          return response.body;
        }),
        map(result => result.events||[]),
        mergeMap((roots) => this.rootMapper(roots))
      ).subscribe((roots: Root[]) => {
        if(!project.sequences?.length && roots.length < this.DEFAULT_SEQUENCE_PAGE_SIZE) {
          project.allSequencesLoaded = true;
        }
        project.sequences = [...roots||[], ...project.sequences||[]].sort(DateUtil.compareTraceTimesAsc);
        project.stages.forEach(stage => this.stageMapper(stage, project));
        this._roots.next(project.sequences);
    });
  }

  public loadOldRoots(project: Project, fromRoot?: Root) {
    this.apiService.getRoots(project.projectName, fromRoot ? this.MAX_SEQUENCE_PAGE_SIZE : this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE, null, fromRoot?.time ? new Date(fromRoot.time).toISOString() : null, new Date(project.sequences[project.sequences.length - 1].time).toISOString()).pipe(
      map(response => response.body.events || []),
      mergeMap((roots) => this.rootMapper(roots))
    ).subscribe((roots: Root[]) => {
      if (!fromRoot && roots.length < this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE) {
        project.allSequencesLoaded = true;
      }
      if (roots.length !== 0 || fromRoot) {
        project.sequences = [...(project.sequences || []), ...(roots || []), ...(fromRoot ? [fromRoot] : [])].sort(DateUtil.compareTraceTimesAsc);
        project.stages.forEach(stage => this.stageMapper(stage, project));
      }
      this._roots.next(project.sequences);
    });
  }

  public getRoot(projectName: string, shkeptncontext: string): Observable<Root> {
    return this.apiService.getRoots(projectName, 1, null, null, null, shkeptncontext).pipe(
      map(response => response.body.events || []),
      switchMap(roots => this.rootMapper(roots).pipe(
        map(sequences => sequences.pop())
      ))
    )
  }

  public loadUntilRoot(project: Project, shkeptncontext: string): void {
    this.getRoot(project.projectName, shkeptncontext).subscribe((root: Root) => {
      if (root) {
        this.loadOldRoots(project, root);
      }
    })
  }

  public loadTraces(root: Root) {
    let fromTime: Date = this._tracesLastUpdated[root.shkeptncontext];

    this.apiService.getTraces(root.shkeptncontext, root.getProject(), fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          let lastUpdated = moment(response.headers.get("date"));
          let lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
          this._tracesLastUpdated[root.shkeptncontext] = lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated;
          return response.body;
        }),
        map(result => result.events||[]),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        root.traces = Trace.traceMapper([...traces||[], ...root.traces||[]]);
        this.getProject(root.getProject()).pipe(take(1))
          .subscribe(project => {
            project.stages.filter(s => root.getStages().includes(s.stageName)).forEach(stage => {
              stage.services.filter(s => root.getService() == s.serviceName).forEach(service => {
                service.openApprovals = service.roots.reduce((openApprovals, root) => {
                  const approval = root.getPendingApproval(stage.stageName);
                  if(approval) {
                    openApprovals.push(approval);
                  }
                  return openApprovals;
                }, []);
              });
            });
          });
        this._roots.next([...this._roots.getValue()]);
      });
  }

  public getDeploymentsOfService(projectName: string, serviceName: string): Observable<Deployment[]> {
    return this.apiService.getDeploymentsOfService(projectName, serviceName).pipe(
      map(deployments => deployments.map(deployment => Deployment.fromJSON(deployment)))
    );
  }

  public loadTracesByContext(shkeptncontext: string) {
    this.apiService.getTraces(shkeptncontext)
      .pipe(
        map(response => response.body),
        map(result => result.events||[]),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        this._traces.next(traces);
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
    this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_STARTED, 'approval.started')
      .pipe(
        mergeMap(()=> this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_FINISHED, 'approval.finished'))
      )
      .subscribe(() => {
        let root = this._projects.getValue().find(p => p.projectName == approval.data.project).services.find(s => s.serviceName == approval.data.service).roots.find(r => r.shkeptncontext == approval.shkeptncontext);
        this.loadTraces(root);
      });
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

  public loadTaskNames(projectName: string) {
    this.apiService.getTaskNames(projectName)
      .pipe(
        map(taskNames => taskNames.sort((taskA, taskB) => taskA.localeCompare(taskB)))
      )
      .subscribe(taskNames => {
      this._taskNames.next(taskNames);
    });
  }

  private rootMapper(roots: Trace[]): Observable<Root[]> {
    return from(roots).pipe(
      mergeMap(
        root => {
          return this.apiService.getTraces(root.shkeptncontext, root.data.project)
            .pipe(
              map(response => {
                const lastUpdated = moment(response.headers.get('date'));
                const lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
                this._tracesLastUpdated[root.shkeptncontext] = lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated;
                return response.body;
              }),
              map(result => result.events || []),
              map(Trace.traceMapper),
              map(traces => ({ ...root, traces}))
            );
        }
      ),
      toArray(),
      map(rs => rs.map(root => Root.fromJSON(root)))
    );
  }

  private stageMapper(stage: Stage, project: Project) {
    stage.services.forEach(service => {
      service.roots = project.sequences.filter(s => s.getService() === service.serviceName && s.getStages().includes(stage.stageName));
      service.openApprovals = service.roots.reduce((openApprovals, root) => {
        const approval = root.getPendingApproval(stage.stageName);
        if(approval) {
          openApprovals.push(approval);
        }
        return openApprovals;
      }, []);
    });
  }
}
