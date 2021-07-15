import {Injectable} from '@angular/core';
import {BehaviorSubject, forkJoin, from, Observable, Subject, of} from 'rxjs';
import {catchError, filter, map, mergeMap, switchMap, take, toArray} from 'rxjs/operators';
import {Trace} from '../_models/trace';
import {Stage} from '../_models/stage';
import {Project} from '../_models/project';
import {EventTypes} from '../_models/event-types';
import {ApiService} from './api.service';
import moment from 'moment';
import {Deployment} from '../_models/deployment';
import {Sequence} from '../_models/sequence';
import {UniformRegistration} from '../_models/uniform-registration';
import {UniformRegistrationLog} from '../_models/uniform-registration-log';
import {Secret} from '../_models/secret';
import { Root } from '../_models/root';
import { DateUtil } from '../_utils/date.utils';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../_models/event-result';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  protected _projects = new BehaviorSubject<Project[]>(null);
  protected _taskNames = new BehaviorSubject<string[]>([]);
  protected _sequences = new BehaviorSubject<Sequence[]>(null);
  protected _roots = new BehaviorSubject<Root[]>(null);
  protected _traces = new BehaviorSubject<Trace[]>(null);
  protected _openApprovals = new BehaviorSubject<Trace[]>([]);
  protected _keptnInfo = new BehaviorSubject<any>(null);
  protected _changedDeployments = new BehaviorSubject<Deployment[]>([]);
  protected _sequences = new BehaviorSubject<Sequence[]>(null);
  protected _rootsLastUpdated: { [key: string]: Date } = {};
  protected _sequencesLastUpdated: { [key: string]: Date } = {};
  protected _tracesLastUpdated: { [key: string]: Date } = {};
  protected _rootTracesLastUpdated: { [key: string]: Date } = {};
  private readonly DEFAULT_SEQUENCE_PAGE_SIZE = 25;
  private readonly DEFAULT_NEXT_SEQUENCE_PAGE_SIZE = 10;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;

  protected _isQualityGatesOnly: BehaviorSubject<boolean> = new BehaviorSubject(false);
  protected _evaluationResults = new Subject();

  constructor(private apiService: ApiService) {
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  get taskNames(): Observable<string[]> {
    return this._taskNames.asObservable();
  }

  get taskNamesTriggered(): Observable<string[]> {
    return this._taskNames.pipe(
      map(tasks => tasks.map(task => task + '.triggered'))
    );
  }

  get sequences(): Observable<Sequence[]> {
    return this._sequences.asObservable();
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

  get changedDeployments(): Observable<Deployment[]> {
    return this._changedDeployments.asObservable();
  }

  get isQualityGatesOnly(): Observable<boolean> {
    return this._isQualityGatesOnly.asObservable();
  }

  public createProject(projectName: string, shipyard: string, gitRemoteUrl?: string, gitToken?: string, gitUser?: string): Observable<any> {
    return this.apiService.createProject(projectName, shipyard, gitRemoteUrl, gitToken, gitUser);
  }

  public getProject(projectName): Observable<Project> {
    return this.projects.pipe(
      map(projects => projects ? projects.find(project => {
        return project.projectName === projectName;
      }) : null)
    );
  }

  public getUniformRegistrations(): Observable<UniformRegistration[]> {
    return this.apiService.getUniformRegistrations();
  }

  public getUniformRegistrationLogs(uniformRegistrationId: string, pageSize?: number): Observable<UniformRegistrationLog[]> {
    return this.apiService.getUniformRegistrationLogs(uniformRegistrationId, pageSize).pipe(
      map((response) => response.logs)
    );
  }

  public getSecrets(): Observable<Secret[]> {
    return this.apiService.getSecrets()
      .pipe(
        map(res => res.Secrets),
        map(secrets => secrets.map(secret => Secret.fromJSON(secret)))
      );
  }

  public addSecret(secret: Secret): Observable<object> {
    return this.apiService.addSecret(Object.assign({}, secret, {
      data: secret.data.reduce((result, item) => Object.assign(result, {[item.key]: item.value}), {})
    }));
  }

  public deleteSecret(name, scope): Observable<object> {
    return this.apiService.deleteSecret(name, scope);
  }

  public getRootsLastUpdated(project: Project): Date {
    return this._rootsLastUpdated[project.projectName];
  }

  public getTracesLastUpdated(sequence: Sequence): Date {
    return this._tracesLastUpdated[sequence.shkeptncontext];
  }

  public setGitUpstreamUrl(projectName: string, gitUrl: string, gitUser: string, gitToken: string): Observable<boolean> {
    return this.apiService.sendGitUpstreamUrl(projectName, gitUrl, gitUser, gitToken).pipe(map(() => {
      this.loadProjects();
      return true;
    }), catchError(() => {
      return of(false);
    }));
  }

  public loadKeptnInfo() {
    // #4165 Get bridge info first to get info if versions.json should be loaded or not
    // Versions should not be loaded if enableVersionCheckFeature is set to false (when ENABLE_VERSION_CHECK is set to false in env)
    this.apiService.getKeptnInfo().subscribe((bridgeInfo) => {
      forkJoin({
        availableVersions: bridgeInfo.enableVersionCheckFeature ? this.apiService.getAvailableVersions() : of(null),
        keptnVersion: this.apiService.getKeptnVersion(),
        versionCheckEnabled: of(this.apiService.isVersionCheckEnabled()),
        metadata: this.apiService.getMetadata()
      }).subscribe((result) => {
        const keptnInfo = {...result, bridgeInfo: {...bridgeInfo}};
        if (keptnInfo.bridgeInfo.showApiToken) {
          if (window.location.href.indexOf('bridge') !== -1) {
            keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf('/bridge'))}/api`;
          }
          else {
            keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf(window.location.pathname))}/api`;
          }

          keptnInfo.bridgeInfo.authCommand = `keptn auth --endpoint=${keptnInfo.bridgeInfo.apiUrl} --api-token=${keptnInfo.bridgeInfo.apiToken}`;

          this._isQualityGatesOnly.next(!keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY'));
        }
        this._keptnInfo.next(keptnInfo);
      }, (err) => {
        this._keptnInfo.error(err);
      });
    });
  }

  public setVersionCheck(enabled: boolean) {
    this.apiService.setVersionCheck(enabled);
    this.loadKeptnInfo();
  }

  public loadProject(projectName: string) {
    this.apiService.getProject(projectName)
      .pipe(
        map(project => Project.fromJSON(project))
      ).subscribe((project: Project) => {
        const projects = this._projects.getValue();
        const existingProject = projects?.find(p => p.projectName === project.projectName);
        if (existingProject){
          const {roots, sequences, ...copyProject} = project;
          Object.assign(existingProject, copyProject);
          this._projects.next(projects);
        }
    });
  }

  public loadProjects() {
    this.apiService.getProjects(this._keptnInfo.getValue().bridgeInfo.projectsPageSize || 50)
      .pipe(
        map(result => result.projects),
        map(projects =>
          projects.map(project => Project.fromJSON(project))
        )
      ).subscribe((projects: Project[]) => {
      const existingProjects = this._projects.getValue();
      projects = projects.map(project => {
        const existingProject = existingProjects?.find(p => p.projectName === project.projectName);
        if (existingProject) {
          return Object.assign(existingProject, project);
        } else {
          return project;
        }
      });
      this._projects.next(projects);
    }, () => {
      this._projects.next([]);
    });
  }

  public loadOpenRemediations(project: Project): void {
    this.apiService.getOpenRemediations(project.projectName, this.MAX_SEQUENCE_PAGE_SIZE).pipe(
      map(response => response.body),
      map(sequenceResult => sequenceResult.states),
      map(sequences => {
        const changedDeployments: Deployment[] = [];
        // remove finished remediations
        for (const service of project.getServices()){
          for (const deployment of service.deployments) {
            for (const stage of deployment.stages) {
              const filteredRemediations = stage.remediations.filter(r => sequences.some(s => s.shkeptncontext === r.shkeptncontext));
              if (filteredRemediations.length !== stage.remediations.length) {
                if (!changedDeployments.some(d => d.shkeptncontext === deployment.shkeptncontext)) {
                  changedDeployments.push(deployment);
                }
                stage.remediations = filteredRemediations;
              }
            }
          }
        }
        return [sequences, changedDeployments];
      }),
      mergeMap(([sequences, changedDeployments]) =>
        from(sequences).pipe(
          mergeMap((sequence: Sequence) => {
            const service = project.getService(sequence.service);
            const sequenceStage = sequence.stages[0].name;
            let result = of(null);
            if (service) {
              const deployment = service.deployments.find(d => d.stages.some(stage => sequence.stages.some(s => s.name === stage.stageName)));
              if (deployment) {
                const stage = deployment.stages.find(s => s.stageName === sequenceStage);
                if (stage) {
                  const existingRemediation = stage.remediations.find(r => r.shkeptncontext === sequence.shkeptncontext);
                  let _resourceContent: Observable<any> = of(null);
                  let _root: Observable<any> = of(null);

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
                  result = forkJoin([_root, _resourceContent]).pipe(switchMap(() => of(deployment)));
                }
              }
            }
            return result;
          }),
          toArray(),
          filter(deployment => !!deployment),
          map((newChangedDeployments: Deployment[]) => {
            const deployments = changedDeployments as Deployment[];
            for (const deployment of newChangedDeployments) {
              if (!deployments.some(d => d.shkeptncontext === deployment.shkeptncontext)) {
                deployments.push(deployment);
              }
            }
            return deployments;
          })
        )
      )
    ).subscribe((deployments: Deployment[]) => {
      this._changedDeployments.next(deployments);
    });
  }

  public loadRoots(project: Project) {
    const fromTime: Date = this._rootsLastUpdated[project.projectName];
    this._rootsLastUpdated[project.projectName] = new Date();

    this.apiService.getRoots(project.projectName, this.DEFAULT_SEQUENCE_PAGE_SIZE, null, fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          const lastUpdated = moment(response.headers.get('date'));
          const lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
          this._rootsLastUpdated[project.projectName] = (lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
          return response.body;
        }),
        map(result => result.events || []),
        mergeMap((roots) => this.rootMapper(roots))
      ).subscribe((roots: Root[]) => {
        project.roots = [...roots || [], ...project.roots || []].sort(DateUtil.compareTraceTimesAsc);
        project.stages.forEach(stage => this.stageRootMapper(stage, project));
        this._roots.next(project.roots);
    });
  }

  public loadSequences(project: Project, fromTime?: Date, beforeTime?: Date, oldSequence?: Sequence): void {
    if (!beforeTime && !fromTime) { // set fromTime if it isn't loadOldSequences
      fromTime = this._sequencesLastUpdated[project.projectName];
    }
    this._sequencesLastUpdated[project.projectName] = new Date();
    this.apiService.getSequences(project.projectName, beforeTime ? this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE : this.DEFAULT_SEQUENCE_PAGE_SIZE, null, null, fromTime?.toISOString(), beforeTime?.toISOString())
      .pipe(
        map(response => {
          this.updateSequencesUpdated(response, project.projectName);
          return response.body;
        }),
        map(body => {
          return [body.states.map(sequence => Sequence.fromJSON(sequence)), body.totalCount ?? body.states.length];
        }),
      ).subscribe(([sequences, totalCount]: [Sequence[], number]) => {
        this.addNewSequences(project, sequences, !!beforeTime, oldSequence);

        if (this.allSequencesLoaded(project.sequences.length, totalCount, fromTime, beforeTime)) {
          project.allSequencesLoaded = true;
        }
        project.stages.forEach(stage => {
          this.stageSequenceMapper(stage, project);
        });
        this._sequences.next(project.sequences);
    });
  }

  private addNewSequences(project: Project, newSequences: Sequence[], areOldSequences: boolean, oldSequence?: Sequence) {
    if (areOldSequences) {
      project.sequences = [...project.sequences || [], ...newSequences || [], ...(oldSequence ? [oldSequence] : [])];
    }
    else {
      project.sequences = [...newSequences || [], ...project.sequences || []];
    }
  }

  private updateSequencesUpdated(response: HttpResponse<SequenceResult>, projectName: string): void {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body.states[0] ? moment(response.body.states[0]?.time) : null;
    this._sequencesLastUpdated[projectName] = (lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
  }

  private updateTracesUpdated(response: HttpResponse<EventResult>, keptnContext: string) {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
    this._tracesLastUpdated[keptnContext] = (lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
  }

  private allSequencesLoaded(sequences: number, totalCount: number, fromTime?: Date, beforeTime?: Date): boolean {
    return fromTime && !beforeTime && sequences >= totalCount || beforeTime && !fromTime && totalCount < this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE;
  }

  public getRoot(projectName: string, shkeptncontext: string): Observable<Root> {
    return this.apiService.getRoots(projectName, 1, null, null, null, shkeptncontext).pipe(
      map(response => response.body.events || []),
      switchMap(roots => this.rootMapper(roots).pipe(
        map(sequences => sequences.pop())
      ))
    );
  }

  public loadOldSequences(project: Project, fromTime?: Date, oldSequence?: Sequence): void {
    this.loadSequences(project, fromTime, new Date(project.sequences[project.sequences.length - 1].time), oldSequence);
  }

  public getSequenceWithTraces(projectName: string, keptnContext: string): Observable<Sequence> {
    return this.apiService.getSequences(projectName, 1, null, null, null, null, keptnContext).pipe(
      map(response => response.body.states || []),
      map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)).shift()),
      switchMap(sequence => this.sequenceMapper([sequence])),
      map(sequences => sequences.shift())
    );
  }

  public updateSequence(projectName: string, keptnContext: string): void {
    this.apiService.getSequences(projectName, 1, null, null, null, null, keptnContext).pipe(
      map(response => response.body.states || []),
      map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)).shift()),
    ).subscribe(sequence => {
      const sequences = this._sequences.getValue();
      const oldSequence = sequences.find(seq => seq.shkeptncontext === keptnContext);
      if (oldSequence) {
        const {traces, ...copySequence} = sequence; // don't overwrite loaded traces
        Object.assign(oldSequence, copySequence);
      }
      this._sequences.next(sequences);
    });
  }

  public loadUntilRoot(project: Project, shkeptncontext: string): void {
    this.getSequenceWithTraces(project.projectName, shkeptncontext).subscribe((sequence: Sequence) => {
      if (sequence) {
        this.loadOldSequences(project, new Date(sequence.time), sequence);
      }
    });
  }

  public loadRootTraces(root: Root) {
    const fromTime: Date = this._rootTracesLastUpdated[root.shkeptncontext];

    this.apiService.getTraces(root.shkeptncontext, root.getProject(), fromTime ? fromTime.toISOString() : null)
      .pipe(
        map(response => {
          const lastUpdated = moment(response.headers.get('date'));
          const lastEvent = response.body.events[0] ? moment(response.body.events[0]?.time) : null;
          this._rootsLastUpdated[root.shkeptncontext] = (lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
          return response.body;
        }),
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        root.traces = Trace.traceMapper([...traces || [], ...root.traces || []]);
        this.getProject(root.getProject()).pipe(take(1))
          .subscribe(project => {
            project.stages.filter(s => root.getStages().includes(s.stageName)).forEach(stage => {
              this.stageRootMapper(stage, project);
            });
          });
        this._roots.next([...this._roots.getValue()]);
      });
  }

  public loadTraces(sequence: Sequence) {
    const fromTime: Date = this._tracesLastUpdated[sequence.shkeptncontext];
    this.apiService.getTraces(sequence.shkeptncontext, sequence.project, fromTime?.toISOString())
      .pipe(
        map(response => {
          this.updateTracesUpdated(response, sequence.shkeptncontext);
          return response.body;
        }),
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        sequence.traces = Trace.traceMapper([...traces || [], ...sequence.traces || []]);
        this.getProject(sequence.project).pipe(take(1))
          .subscribe(project => {
            project.stages.filter(s => sequence.getStages().includes(s.stageName)).forEach(stage => {
              this.stageSequenceMapper(stage, project);
            });
          });
        this._sequences.next([...this._sequences.getValue()]);
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
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        this._traces.next(traces);
      });
  }

  public loadEvaluationResults(event: Trace) {
    let fromTime: Date;
    if (event.data.evaluationHistory) {
      fromTime = new Date(event.data.evaluationHistory[event.data.evaluationHistory.length - 1].time);
    }

    this.apiService.getEvaluationResults(event.data.project, event.data.service, event.data.stage, event.source, fromTime?.toISOString())
      .pipe(
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        this._evaluationResults.next({
          type: 'evaluationHistory',
          triggerEvent: event,
          traces
        });
      });
  }

  public getEvaluationResult(shkeptncontext: string): Observable<Trace> {
    return this.apiService.getEvaluationResult(shkeptncontext)
      .pipe(
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)).find(() => true))
      );
  }

  public sendApprovalEvent(approval: Trace, approve: boolean) {
    this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_STARTED, 'approval.started')
      .pipe(
        mergeMap(() => this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_FINISHED, 'approval.finished'))
      )
      .subscribe(() => {
        const sequence = this._projects.getValue().find(p => p.projectName === approval.data.project)
                        .services.find(s => s.serviceName === approval.data.service)
                        .sequences.find(r => r.shkeptncontext === approval.shkeptncontext);
        this.loadTraces(sequence);
      });
  }

  public invalidateEvaluation(evaluation: Trace, reason: string) {
    this.apiService.sendEvaluationInvalidated(evaluation, reason)
      .pipe(take(1))
      .subscribe(() => {
        this._evaluationResults.next({
          type: 'invalidateEvaluation',
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

  public loadSequences(project: Project, pageSize: number) {
    this.apiService.getSequences(project.projectName, null, null, null, null, pageSize)
      .pipe(
        map(response => response.body),
        map(sequenceResult => sequenceResult.states),
        map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)).sort((sA, sB) => moment(sA.time).isBefore(sB.time) ? -1 : 1))
      ).subscribe((newSequenceStates: Sequence[]) => {
      project.sequenceStates = project.getSequenceStates().map(sequence => newSequenceStates.find(s => s.shkeptncontext == sequence.shkeptncontext) || sequence);
      newSequenceStates.forEach(sequenceState => {
        if (project.getSequenceStates().length == 0 || moment(sequenceState.time).isAfter(project.getSequenceStates()[0].time))
          project.getSequenceStates().unshift(sequenceState);
      });

      this._sequences.next(newSequenceStates);
    });
  }

  private sequenceMapper(sequences: Sequence[]): Observable<Sequence[]> {
    return from(sequences).pipe(
      mergeMap(
        sequence => {
          return this.apiService.getTraces(sequence.shkeptncontext, sequence.project)
            .pipe(
              map(response => {
                this.updateTracesUpdated(response, sequence.shkeptncontext);
                return response.body;
              }),
              map(result => result.events || []),
              map(Trace.traceMapper),
              map(traces => {
                sequence.traces = traces;
                return sequence;
              })
            );
        }
      ),
      toArray(),
    );
  }

  private stageSequenceMapper(stage: Stage, project: Project) {
    stage.services.forEach(service => {
      service.sequences = project.sequences.filter(s => s.service === service.serviceName && s.getStages().includes(stage.stageName));
    });
  }

  private stageRootMapper(stage: Stage, project: Project) {
    stage.services.forEach(service => {
      service.roots = project.roots.filter(s => s.getService() === service.serviceName && s.getStages().includes(stage.stageName));
      service.openApprovals = service.roots.reduce((openApprovals, currentRoot) => {
        const approval = currentRoot.getPendingApproval(stage.stageName);
        if (approval) {
          openApprovals.push(approval);
        }
        return openApprovals;
      }, []);
    });
  }

  private rootMapper(roots: Trace[]): Observable<Root[]> {
    return from(roots).pipe(
      mergeMap(
        root => {
          return this.apiService.getTraces(root.shkeptncontext, root.data.project)
            .pipe(
              map(result => result.body.events || []),
              map(Trace.traceMapper),
              map(traces => ({ ...root, traces}))
            );
        }
      ),
      toArray(),
      map(rs => rs.map(root => Root.fromJSON(root)))
    );
  }
}
