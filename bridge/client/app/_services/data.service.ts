import { Injectable } from '@angular/core';
import { BehaviorSubject, forkJoin, from, Observable, Subject, of } from 'rxjs';
import { map, mergeMap, switchMap, take, tap, toArray } from 'rxjs/operators';
import { Trace } from '../_models/trace';
import { Stage } from '../_models/stage';
import { Project } from '../_models/project';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ApiService } from './api.service';
import moment from 'moment';
import { Deployment } from '../_models/deployment';
import { Sequence } from '../_models/sequence';
import { UniformRegistrationLog } from '../../../server/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { Root } from '../_models/root';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../../../shared/interfaces/event-result';
import { KeptnInfo } from '../_models/keptn-info';
import { KeptnInfoResult } from '../_models/keptn-info-result';
import { DeploymentStage } from '../_models/deployment-stage';
import { UniformRegistration } from '../_models/uniform-registration';
import { SequenceState } from '../../../shared/models/sequence';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  protected _projects = new BehaviorSubject<Project[] | undefined>(undefined);
  protected _taskNames = new BehaviorSubject<string[]>([]);
  protected _sequences = new BehaviorSubject<Sequence[] | undefined>(undefined);
  protected _traces = new BehaviorSubject<Trace[] | undefined>(undefined);
  protected _openApprovals = new BehaviorSubject<Trace[]>([]);
  protected _keptnInfo = new BehaviorSubject<KeptnInfo | undefined>(undefined);
  protected _changedDeployments = new BehaviorSubject<Deployment[]>([]);
  protected _rootsLastUpdated: { [key: string]: Date } = {};
  protected _sequencesLastUpdated: { [key: string]: Date } = {};
  protected _tracesLastUpdated: { [key: string]: Date } = {};
  protected _rootTracesLastUpdated: { [key: string]: Date } = {};
  protected _projectName: BehaviorSubject<string> = new BehaviorSubject<string>('');
  protected _uniformDates: {[key: string]: string} = this.apiService.uniformLogDates;
  protected _hasUnreadUniformRegistrationLogs = new BehaviorSubject<boolean>(false);
  protected readonly DEFAULT_SEQUENCE_PAGE_SIZE = 25;
  protected readonly DEFAULT_NEXT_SEQUENCE_PAGE_SIZE = 10;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;

  protected _isQualityGatesOnly = new BehaviorSubject<boolean>(false);
  protected _evaluationResults = new Subject<{ type: string, triggerEvent: Trace, traces?: Trace[] }>();

  constructor(private apiService: ApiService) {
  }

  get projects(): Observable<Project[] | undefined> {
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

  get sequences(): Observable<Sequence[] | undefined> {
    return this._sequences.asObservable();
  }

  get traces(): Observable<Trace[] | undefined> {
    return this._traces.asObservable();
  }

  get openApprovals(): Observable<Trace[]> {
    return this._openApprovals.asObservable();
  }

  get keptnInfo(): Observable<KeptnInfo | undefined> {
    return this._keptnInfo.asObservable();
  }

  get evaluationResults(): Observable<{ type: string, triggerEvent: Trace, traces?: Trace[] }> {
    return this._evaluationResults;
  }

  get changedDeployments(): Observable<Deployment[]> {
    return this._changedDeployments.asObservable();
  }

  get isQualityGatesOnly(): Observable<boolean> {
    return this._isQualityGatesOnly.asObservable();
  }

  get projectName(): Observable<string> {
    return this._projectName.asObservable();
  }

  get hasUnreadUniformRegistrationLogs(): Observable<boolean> {
    return this._hasUnreadUniformRegistrationLogs.asObservable();
  }

  public setHasUnreadUniformRegistrationLogs(status: boolean): void {
    this._hasUnreadUniformRegistrationLogs.next(status);
  }

  public setProjectName(projectName: string): void {
    this._projectName.next(projectName);
  }

  public setUniformDate(integrationId: string, lastSeen?: string) {
    this._uniformDates[integrationId] = lastSeen || new Date().toISOString();
    this.apiService.uniformLogDates = this._uniformDates;
  }

  public getUniformDate(id: string): Date | undefined {
    return this._uniformDates[id] ? new Date(this._uniformDates[id]) : undefined;
  }

  public getProject(projectName: string): Observable<Project | undefined> {
    return this.projects.pipe(
      map(projects => projects?.find(project => project.projectName === projectName))
    );
  }

  public projectExists(projectName: string): Observable<boolean | undefined> {
    return this.projects.pipe(map((projects) => projects?.some(project => project.projectName === projectName)));
  }

  public createProject(projectName: string, shipyard: string, gitRemoteUrl?: string, gitToken?: string, gitUser?: string): Observable<unknown> {
    return this.apiService.createProject(projectName, shipyard, gitRemoteUrl, gitToken, gitUser);
  }

  public getUniformRegistrations(): Observable<UniformRegistration[]> {
    return this.apiService.getUniformRegistrations(this._uniformDates).pipe(
      map(uniformRegistrations => uniformRegistrations.map(registration => UniformRegistration.fromJSON(registration)))
    );
  }

  public getUniformRegistrationLogs(uniformRegistrationId: string, pageSize?: number): Observable<UniformRegistrationLog[]> {
    return this.apiService.getUniformRegistrationLogs(uniformRegistrationId, pageSize).pipe(
      map((response) => response.logs)
    );
  }

  public loadUnreadUniformRegistrationLogs(): void {
    this.apiService.hasUnreadUniformRegistrationLogs(this._uniformDates).subscribe(status => {
      this.setHasUnreadUniformRegistrationLogs(status);
    });
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

  public deleteSecret(name: string, scope: string): Observable<object> {
    return this.apiService.deleteSecret(name, scope);
  }

  public getRootsLastUpdated(project: Project): Date {
    return this._rootsLastUpdated[project.projectName];
  }

  public getTracesLastUpdated(sequence: Sequence): Date {
    return this._tracesLastUpdated[sequence.shkeptncontext];
  }

  public setGitUpstreamUrl(projectName: string, gitUrl: string, gitUser: string, gitToken: string): Observable<unknown> {
    return this.apiService.sendGitUpstreamUrl(projectName, gitUrl, gitUser, gitToken).pipe(tap(() => {
      this.loadProject(projectName);
    }));
  }

  public loadKeptnInfo(): void {
    // #4165 Get bridge info first to get info if versions.json should be loaded or not
    // Versions should not be loaded if enableVersionCheckFeature is set to false (when ENABLE_VERSION_CHECK is set to false in env)
    this.apiService.getKeptnInfo().subscribe((bridgeInfo: KeptnInfoResult) => {
      forkJoin({
        availableVersions: bridgeInfo.enableVersionCheckFeature ? this.apiService.getAvailableVersions() : of(undefined),
        versionCheckEnabled: of(this.apiService.isVersionCheckEnabled()),
        metadata: this.apiService.getMetadata()
      }).subscribe((result) => {
        const keptnInfo: KeptnInfo = {...result, bridgeInfo: {...bridgeInfo}};
        if (keptnInfo.bridgeInfo.showApiToken) {
          if (window.location.href.indexOf('bridge') !== -1) {
            keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf('/bridge'))}/api`;
          } else {
            keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(0, window.location.href.indexOf(window.location.pathname))}/api`;
          }

          keptnInfo.authCommand = `keptn auth --endpoint=${keptnInfo.bridgeInfo.apiUrl} --api-token=${keptnInfo.bridgeInfo.apiToken}`;

          this._isQualityGatesOnly.next(!keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY'));
        }
        this._keptnInfo.next(keptnInfo);
      }, (err) => {
        this._keptnInfo.error(err);
      });
    });
  }

  public setVersionCheck(enabled: boolean): void {
    this.apiService.setVersionCheck(enabled);
    this.loadKeptnInfo();
  }

  public deleteProject(projectName: string): Observable<object> {
    return this.apiService.deleteProject(projectName);
  }

  public loadProject(projectName: string): void {
    this.apiService.getProject(projectName)
      .pipe(
        map(project => Project.fromJSON(project))
      ).subscribe((project: Project) => {
      const projects = this._projects.getValue();
      const existingProject = projects?.find(p => p.projectName === project.projectName);
      if (existingProject) {
        const {sequences, ...copyProject} = project;
        Object.assign(existingProject, copyProject);
        this._projects.next(projects);
      }
    }, err => {
      if (err.status === 404) {
        const projects = this._projects.getValue();
        const projectIdx = projects?.findIndex(p => p.projectName === projectName) ?? -1;
        if (projectIdx >= 0) {
          projects?.splice(projectIdx, 1);
          this._projects.next(projects);
        }
      } else {
        this._projects.error(err);
      }
    });
  }

  public loadProjects(): void {
    this.apiService.getProjects(this._keptnInfo.getValue()?.bridgeInfo.projectsPageSize || 50)
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
          const {sequences, ...copyProject} = project;
          return Object.assign(existingProject, copyProject);
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
      map(sequenceResult => sequenceResult?.states ?? []),
      map((sequences: Sequence[]): [Sequence[], Deployment[]] => {
        const changedDeployments: Deployment[] = [];
        // remove finished remediations
        for (const service of project.getServices()) {
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
            let result: Observable<null | Deployment> = of(null);
            if (service) {
              const deployment = service.deployments.find(d => d.stages.some(stage => sequence.stages.some(s => s.name === stage.stageName)));
              if (deployment) {
                const stage = deployment.stages.find(s => s.stageName === sequenceStage);
                if (stage) {
                  const existingRemediation = stage.remediations.find(r => r.shkeptncontext === sequence.shkeptncontext);
                  let _resourceContent: Observable<DeploymentStage | undefined> = of(undefined);
                  let _root: Observable<Root | undefined> = of(undefined);

                  // update existing remediation
                  if (existingRemediation) {
                    Object.assign(existingRemediation, Sequence.fromJSON(sequence));
                  } else {
                    const remediation = Sequence.fromJSON(sequence);
                    stage.remediations.push(remediation);
                    if (!remediation.problemTitle) {
                      _root = this.getRoot(project.projectName, remediation.shkeptncontext).pipe(
                        tap(root => {
                          remediation.problemTitle = root?.getProblemTitle();
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
          map(deployments => deployments.filter((deployment: Deployment | null): deployment is Deployment => !!deployment)),
          map((newChangedDeployments: Deployment[]) => {
            const deployments = changedDeployments;
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

  public loadSequences(project: Project, fromTime?: Date, beforeTime?: Date, oldSequence?: Sequence): void {
    if (!beforeTime && !fromTime) { // set fromTime if it isn't loadOldSequences
      fromTime = this._sequencesLastUpdated[project.projectName];
    }
    this._sequencesLastUpdated[project.projectName] = new Date();
    this.apiService.getSequences(project.projectName, beforeTime ? this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE : this.DEFAULT_SEQUENCE_PAGE_SIZE, undefined, undefined, fromTime?.toISOString(), beforeTime?.toISOString())
      .pipe(
        map(response => {
          this.updateSequencesUpdated(response, project.projectName);
          return response.body;
        }),
        map(body => {
          const count = body?.totalCount ?? body?.states.length ?? 0;
          const sequences = body?.states.map(sequence => Sequence.fromJSON(sequence)) ?? [];
          return [sequences, count] as [Sequence[], number];
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

  public loadLatestSequences(project: Project, pageSize: number): Observable<Sequence[]> {
    return this.apiService.getSequences(project.projectName, pageSize)
      .pipe(
        map(response => response.body),
        map(body => body?.states.map(sequence => Sequence.fromJSON(sequence)) ?? []),
      );
  }

  protected addNewSequences(project: Project, newSequences: Sequence[], areOldSequences: boolean, oldSequence?: Sequence): void {
    if (areOldSequences) {
      project.sequences = [...project.sequences || [], ...newSequences || [], ...(oldSequence ? [oldSequence] : [])];
    } else {
      project.sequences = [...newSequences || [], ...project.sequences || []];
    }
  }

  private updateSequencesUpdated(response: HttpResponse<SequenceResult>, projectName: string): void {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body?.states[0] ? moment(response.body.states[0]?.time) : null;
    this._sequencesLastUpdated[projectName] = (lastEvent && lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
  }

  private updateTracesUpdated(response: HttpResponse<EventResult>, keptnContext: string): void {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body?.events[0] ? moment(response.body.events[0]?.time) : null;
    this._tracesLastUpdated[keptnContext] = (lastEvent && lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated).toDate();
  }

  protected allSequencesLoaded(sequences: number, totalCount: number, fromTime?: Date, beforeTime?: Date): boolean {
    return !!fromTime && !beforeTime && sequences >= totalCount
      || !!beforeTime && !fromTime && totalCount < this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE;
  }

  public getRoot(projectName: string, shkeptncontext: string): Observable<Root | undefined> {
    return this.apiService.getRoots(projectName, 1, undefined, undefined, undefined, shkeptncontext).pipe(
      map(response => response.body?.events || []),
      switchMap(roots => this.rootMapper(roots).pipe(
        map(sequences => sequences.pop())
      ))
    );
  }

  public loadOldSequences(project: Project, fromTime?: Date, oldSequence?: Sequence): void {
    this.loadSequences(project, fromTime, new Date(project.sequences[project.sequences.length - 1].time), oldSequence);
  }

  public getSequenceWithTraces(projectName: string, keptnContext: string): Observable<Sequence | undefined> {
    return this.apiService.getSequences(projectName, 1, undefined, undefined, undefined, undefined, keptnContext).pipe(
      map(response => response.body?.states || []),
      map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)).shift()),
      switchMap(sequence => sequence ? this.sequenceMapper([sequence]) : []),
      map(sequences => sequences.shift())
    );
  }

  public updateSequence(projectName: string, keptnContext: string): void {
    this.apiService.getSequences(projectName, 1, undefined, undefined, undefined, undefined, keptnContext).pipe(
      map(response => response.body?.states || []),
      map(sequences => sequences.map(sequence => Sequence.fromJSON(sequence)).shift()),
    ).subscribe(sequence => {
      const sequences = this._sequences.getValue();
      const oldSequence = sequences?.find(seq => seq.shkeptncontext === keptnContext);
      if (oldSequence && sequence) {
        const {traces, ...copySequence} = sequence; // don't overwrite loaded traces
        Object.assign(oldSequence, copySequence);
      }
      this._sequences.next(sequences);
    });
  }

  public loadUntilRoot(project: Project, shkeptncontext: string): void {
    this.getSequenceWithTraces(project.projectName, shkeptncontext).subscribe(sequence => {
      if (sequence) {
        this.loadOldSequences(project, new Date(sequence.time), sequence);
      }
    });
  }

  public loadTraces(sequence: Sequence): void {
    const fromTime: Date = this._tracesLastUpdated[sequence.shkeptncontext];
    this.apiService.getTraces(sequence.shkeptncontext, sequence.project, fromTime?.toISOString())
      .pipe(
        map(response => {
          this.updateTracesUpdated(response, sequence.shkeptncontext);
          return response.body;
        }),
        map(result => result?.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        sequence.traces = Trace.traceMapper([...traces || [], ...sequence.traces || []]);
        this.getProject(sequence.project).pipe(take(1))
          .subscribe(project => {
            if (project) {
              project.stages.filter(s => sequence.getStages().includes(s.stageName)).forEach(stage => {
                this.stageSequenceMapper(stage, project);
              });
            }
          });
        this._sequences.next([...this._sequences.getValue() ?? []]);
      });
  }

  public loadTracesByContext(shkeptncontext: string): void {
    this.apiService.getTraces(shkeptncontext)
      .pipe(
        map(response => response.body),
        map(result => result?.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        this._traces.next(traces);
      });
  }

  public loadEvaluationResults(event: Trace): void {
    let fromTime: Date | undefined;
    const time = event.data.evaluationHistory?.[event.data.evaluationHistory.length - 1]?.time;
    if (time) {
      fromTime = new Date(time);
    }
    if (event.data.project && event.data.service && event.data.stage) {
      this.apiService.getEvaluationResults(event.data.project, event.data.service, event.data.stage, fromTime?.toISOString())
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
  }

  public getEvaluationResult(shkeptncontext: string): Observable<Trace | undefined> {
    return this.apiService.getEvaluationResult(shkeptncontext)
      .pipe(
        map(result => result.events || []),
        map(traces => traces.map(trace => Trace.fromJSON(trace)).find(() => true))
      );
  }

  public sendApprovalEvent(approval: Trace, approve: boolean): void {
    this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_FINISHED, 'approval.finished')
      .subscribe(() => {
        const sequence = this._projects.getValue()?.find(p => p.projectName === approval.data.project)
          ?.getServices().find(s => s.serviceName === approval.data.service)
          ?.sequences.find(r => r.shkeptncontext === approval.shkeptncontext);
        if (sequence) {
          this.loadTraces(sequence);
        }
      });
  }

  public sendSequenceControl(sequence: Sequence, state: string): void {
    sequence.state = SequenceState.UNKNOWN;
    this.apiService.sendSequenceControl(sequence.project, sequence.shkeptncontext, state)
      .subscribe(() => {
        this.updateSequence(sequence.project, sequence.shkeptncontext);
      });
  }

  public invalidateEvaluation(evaluation: Trace, reason: string): void {
    this.apiService.sendEvaluationInvalidated(evaluation, reason)
      .pipe(take(1))
      .subscribe(() => {
        this._evaluationResults.next({
          type: 'invalidateEvaluation',
          triggerEvent: evaluation
        });
      });
  }

  public loadTaskNames(projectName: string): void {
    this.apiService.getTaskNames(projectName)
      .pipe(
        map(taskNames => taskNames.sort((taskA, taskB) => taskA.localeCompare(taskB)))
      )
      .subscribe(taskNames => {
        this._taskNames.next(taskNames);
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
              map(result => result?.events || []),
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

  protected stageSequenceMapper(stage: Stage, project: Project): void {
    stage.services.forEach(service => {
      service.sequences = project.sequences.filter(s => s.service === service.serviceName && s.getStages().includes(stage.stageName));
    });
  }

  private rootMapper(roots: Trace[]): Observable<Root[]> {
    return from(roots).pipe(
      mergeMap(
        root => {
          return this.apiService.getTraces(root.shkeptncontext, root.data.project)
            .pipe(
              map(result => result.body?.events || []),
              map(Trace.traceMapper),
              map(traces => ({...root, traces}))
            );
        }
      ),
      toArray(),
      map(rs => rs.map(root => Root.fromJSON(root)))
    );
  }
}
