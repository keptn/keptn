import { Injectable } from '@angular/core';
import { BehaviorSubject, forkJoin, from, Observable, of, Subject } from 'rxjs';
import { catchError, map, mergeMap, switchMap, take, tap, toArray } from 'rxjs/operators';
import { Trace } from '../_models/trace';
import { Stage } from '../_models/stage';
import { Project } from '../_models/project';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ApiService } from './api.service';
import moment from 'moment';
import { Sequence } from '../_models/sequence';
import { UniformRegistrationLog } from '../../../shared/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { Root } from '../_models/root';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../../../shared/interfaces/event-result';
import { KeptnInfo } from '../_models/keptn-info';
import { KeptnInfoResult } from '../../../shared/interfaces/keptn-info-result';
import { UniformRegistration } from '../_models/uniform-registration';
import { UniformSubscription } from '../_models/uniform-subscription';
import { SequenceState } from '../../../shared/models/sequence';
import { WebhookConfig } from '../../../shared/models/webhook-config';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';
import { EvaluationHistory } from '../_interfaces/evaluation-history';
import { Service } from '../_models/service';
import { Deployment } from '../_models/deployment';
import { ServiceState } from '../_models/service-state';
import { ServiceRemediationInformation } from '../_models/service-remediation-information';
import { EndSessionData } from '../../../shared/interfaces/end-session-data';
import { ISequencesMetadata } from '../../../shared/interfaces/sequencesMetadata';
import { TriggerResponse, TriggerSequenceData } from '../_models/trigger-sequence';
import { EventData } from '../_components/ktb-evaluation-info/ktb-evaluation-info.component';
import { SecretScope } from '../../../shared/interfaces/secret-scope';

@Injectable({
  providedIn: 'root',
})
export class DataService {
  protected _projects = new BehaviorSubject<Project[] | undefined>(undefined);
  protected _sequencesUpdated = new Subject<void>();
  protected _traces = new BehaviorSubject<Trace[] | undefined>(undefined);
  protected _openApprovals = new BehaviorSubject<Trace[]>([]);
  protected _keptnInfo = new BehaviorSubject<KeptnInfo | undefined>(undefined);
  protected _rootsLastUpdated: { [key: string]: Date } = {};
  protected _sequencesLastUpdated: { [key: string]: Date } = {};
  protected _tracesLastUpdated: { [key: string]: Date } = {};
  protected _rootTracesLastUpdated: { [key: string]: Date } = {};
  protected _projectName: BehaviorSubject<string> = new BehaviorSubject<string>('');
  protected _uniformDates: { [key: string]: string } = this.apiService.uniformLogDates;
  protected _hasUnreadUniformRegistrationLogs = new BehaviorSubject<boolean>(false);
  protected readonly DEFAULT_SEQUENCE_PAGE_SIZE = 25;
  protected readonly DEFAULT_NEXT_SEQUENCE_PAGE_SIZE = 10;

  protected _isQualityGatesOnly = new BehaviorSubject<boolean>(false);
  protected _evaluationResults = new Subject<EvaluationHistory>();

  public isTriggerSequenceOpen = false;

  constructor(private apiService: ApiService) {}

  get projects(): Observable<Project[] | undefined> {
    return this._projects.asObservable();
  }

  get sequencesUpdated(): Observable<void> {
    return this._sequencesUpdated.asObservable();
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

  get evaluationResults(): Observable<EvaluationHistory> {
    return this._evaluationResults;
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

  public setUniformDate(integrationId: string, lastSeen?: string): void {
    this._uniformDates[integrationId] = lastSeen || new Date().toISOString();
    this.apiService.uniformLogDates = this._uniformDates;
  }

  public getUniformDate(id: string): Date | undefined {
    return this._uniformDates[id] ? new Date(this._uniformDates[id]) : undefined;
  }

  public getProject(projectName: string): Observable<Project | undefined> {
    return this.projects.pipe(map((projects) => projects?.find((project) => project.projectName === projectName)));
  }

  public projectExists(projectName: string): Observable<boolean | undefined> {
    return this.projects.pipe(map((projects) => projects?.some((project) => project.projectName === projectName)));
  }

  public createProject(
    projectName: string,
    shipyard: string,
    gitRemoteUrl?: string,
    gitToken?: string,
    gitUser?: string
  ): Observable<unknown> {
    return this.apiService.createProject(projectName, shipyard, gitRemoteUrl, gitToken, gitUser);
  }

  public createService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return this.apiService.createService(projectName, serviceName);
  }

  public deleteService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return this.apiService.deleteService(projectName, serviceName);
  }

  public getUniformRegistrations(): Observable<UniformRegistration[]> {
    return this.apiService
      .getUniformRegistrations(this._uniformDates)
      .pipe(
        map((uniformRegistrations) =>
          uniformRegistrations.map((registration) => UniformRegistration.fromJSON(registration))
        )
      );
  }

  public getUniformRegistrationInfo(integrationId: string): Observable<UniformRegistrationInfo> {
    return this.apiService.getUniformRegistrationInfo(integrationId);
  }

  public getUniformSubscription(integrationId: string, subscriptionId: string): Observable<UniformSubscription> {
    return this.apiService
      .getUniformSubscription(integrationId, subscriptionId)
      .pipe(map((uniformSubscription) => UniformSubscription.fromJSON(uniformSubscription)));
  }

  public updateUniformSubscription(
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    return this.apiService.updateUniformSubscription(integrationId, subscription.reduced, webhookConfig);
  }

  public createUniformSubscription(
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: WebhookConfig
  ): Observable<Record<string, unknown>> {
    return this.apiService.createUniformSubscription(integrationId, subscription.reduced, webhookConfig);
  }

  public getUniformRegistrationLogs(
    uniformRegistrationId: string,
    pageSize?: number
  ): Observable<UniformRegistrationLog[]> {
    return this.apiService
      .getUniformRegistrationLogs(uniformRegistrationId, pageSize)
      .pipe(map((response) => response.logs));
  }

  public loadUnreadUniformRegistrationLogs(): void {
    this.apiService.hasUnreadUniformRegistrationLogs(this._uniformDates).subscribe((status) => {
      this.setHasUnreadUniformRegistrationLogs(status);
    });
  }

  public getSecrets(): Observable<Secret[]> {
    return this.apiService.getSecrets().pipe(
      map((res) => res.Secrets),
      map((secrets) => secrets.map((secret) => Secret.fromJSON(secret)))
    );
  }

  public getSecretsForScope(scope: SecretScope): Observable<Secret[]> {
    return this.apiService.getSecretsForScope(scope);
  }

  public addSecret(secret: Secret): Observable<Record<string, unknown>> {
    return this.apiService.addSecret(
      Object.assign({}, secret, {
        data: secret.data?.reduce((result, item) => Object.assign(result, { [item.key]: item.value }), {}),
      })
    );
  }

  public deleteSecret(name: string, scope: string): Observable<Record<string, unknown>> {
    return this.apiService.deleteSecret(name, scope);
  }

  public deleteSubscription(
    integrationId: string,
    id: string,
    isWebhookService: boolean
  ): Observable<Record<string, unknown>> {
    return this.apiService.deleteSubscription(integrationId, id, isWebhookService);
  }

  public getRootsLastUpdated(project: Project): Date {
    return this._rootsLastUpdated[project.projectName];
  }

  public getTracesLastUpdated(sequence: Sequence): Date {
    return this._tracesLastUpdated[sequence.shkeptncontext];
  }

  public setGitUpstreamUrl(
    projectName: string,
    gitUrl: string,
    gitUser: string,
    gitToken: string
  ): Observable<unknown> {
    return this.apiService.sendGitUpstreamUrl(projectName, gitUrl, gitUser, gitToken).pipe(
      tap(() => {
        this.loadProject(projectName);
      })
    );
  }

  public loadKeptnInfo(): void {
    // #4165 Get bridge info first to get info if versions.json should be loaded or not
    // Versions should not be loaded if enableVersionCheckFeature is set to false (when ENABLE_VERSION_CHECK is set to false in env)
    this.apiService.getKeptnInfo().subscribe((bridgeInfo: KeptnInfoResult) => {
      forkJoin({
        availableVersions: bridgeInfo.enableVersionCheckFeature
          ? this.apiService.getAvailableVersions().pipe(catchError(() => of(undefined)))
          : of(undefined),
        versionCheckEnabled: of(this.apiService.isVersionCheckEnabled()),
        metadata: this.apiService.getMetadata(),
      }).subscribe(
        (result) => {
          const keptnInfo: KeptnInfo = { ...result, bridgeInfo: { ...bridgeInfo } };
          if (keptnInfo.bridgeInfo.showApiToken) {
            if (window.location.href.indexOf('bridge') !== -1) {
              keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(
                0,
                window.location.href.indexOf('/bridge')
              )}/api`;
            } else {
              keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(
                0,
                window.location.href.indexOf(window.location.pathname)
              )}/api`;
            }

            keptnInfo.authCommand = `keptn auth --endpoint=${keptnInfo.bridgeInfo.apiUrl} --api-token=${keptnInfo.bridgeInfo.apiToken}`;

            this._isQualityGatesOnly.next(!keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY'));
          }
          this._keptnInfo.next(keptnInfo);
        },
        (err) => {
          this._keptnInfo.error(err);
        }
      );
    });
  }

  public setVersionCheck(enabled: boolean): void {
    this.apiService.setVersionCheck(enabled);
    this.loadKeptnInfo();
  }

  public deleteProject(projectName: string): Observable<Record<string, unknown>> {
    return this.apiService.deleteProject(projectName).pipe(
      tap(() => {
        const projects = this._projects.getValue();
        const projectIdx = projects?.findIndex((p) => p.projectName === projectName) ?? -1;
        if (projectIdx >= 0) {
          projects?.splice(projectIdx, 1);
          this._projects.next(projects);
        }
      })
    );
  }

  public loadPlainProject(projectName: string): Observable<Project> {
    return this.apiService.getPlainProject(projectName).pipe(map((project) => Project.fromJSON(project)));
  }

  public loadProject(projectName: string): void {
    this.apiService
      .getProject(projectName)
      .pipe(map((project) => Project.fromJSON(project)))
      .subscribe(
        (project: Project) => {
          let projects = this._projects.getValue();
          const existingProject = projects?.find((p) => p.projectName === project.projectName);
          if (existingProject) {
            existingProject.update(project);
            existingProject.projectDetailsLoaded = true;
          } else {
            project.projectDetailsLoaded = true;
            projects = [...(projects ?? []), project];
          }
          this._projects.next(projects);
        },
        (err) => {
          if (err.status === 404) {
            const projects = this._projects.getValue();
            const projectIdx = projects?.findIndex((p) => p.projectName === projectName) ?? -1;
            if (projectIdx >= 0) {
              projects?.splice(projectIdx, 1);
              this._projects.next(projects);
            }
          } else {
            this._projects.error(err);
          }
        }
      );
  }

  public loadProjects(): Observable<Project[]> {
    const projects$ = this.apiService.getProjects(this._keptnInfo.getValue()?.bridgeInfo.projectsPageSize || 50).pipe(
      map((result) => (result ? result.projects : [])),
      map((projects) => projects.map((project) => Project.fromJSON(project))),
      map(
        (projects: Project[]) => {
          const existingProjects = this._projects.getValue();
          projects = projects.map((project) => {
            const existingProject = existingProjects?.find((p) => p.projectName === project.projectName);
            if (existingProject) {
              project = existingProject.projectDetailsLoaded
                ? existingProject
                : Object.assign(existingProject, project.reduced);
            }
            return project;
          });
          return projects;
        },
        () => {
          return of([]);
        }
      )
    );

    projects$.subscribe((projects) => {
      this._projects.next(projects);
    });

    return projects$;
  }

  public loadSequences(project: Project, fromTime?: Date, beforeTime?: Date, oldSequence?: Sequence): void {
    if (!beforeTime && !fromTime) {
      // set fromTime if it isn't loadOldSequences
      fromTime = this._sequencesLastUpdated[project.projectName];
    }
    this._sequencesLastUpdated[project.projectName] = new Date();
    this.apiService
      .getSequences(
        project.projectName,
        beforeTime ? this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE : this.DEFAULT_SEQUENCE_PAGE_SIZE,
        undefined,
        undefined,
        fromTime?.toISOString(),
        beforeTime?.toISOString()
      )
      .pipe(
        map((response) => {
          this.updateSequencesUpdated(response, project.projectName);
          return response.body;
        }),
        map((body) => {
          const count = body?.totalCount ?? body?.states.length ?? 0;
          const sequences = body?.states.map((sequence) => Sequence.fromJSON(sequence)) ?? [];
          return [sequences, count] as [Sequence[], number];
        })
      )
      .subscribe(([sequences, totalCount]: [Sequence[], number]) => {
        this.addNewSequences(project, sequences, !!beforeTime, oldSequence);

        if (this.allSequencesLoaded(project.sequences?.length || 0, totalCount, fromTime, beforeTime)) {
          project.allSequencesLoaded = true;
        }
        project.stages.forEach((stage) => {
          this.stageSequenceMapper(stage, project);
        });
        this._sequencesUpdated.next();
      });
  }

  public loadLatestSequences(project: Project, pageSize: number): Observable<Sequence[]> {
    return this.apiService.getSequences(project.projectName, pageSize).pipe(
      map((response) => response.body),
      map((body) => body?.states.map((sequence) => Sequence.fromJSON(sequence)) ?? [])
    );
  }

  public getSequenceMetadata(projectName: string): Observable<ISequencesMetadata> {
    return this.apiService.getSequencesMetadata(projectName);
  }

  protected addNewSequences(
    project: Project,
    sequences: Sequence[],
    areOldSequences: boolean,
    oldSequence?: Sequence
  ): void {
    let newSequences: Sequence[] = [];
    if (areOldSequences) {
      newSequences = [...(project.sequences || []), ...(sequences || []), ...(oldSequence ? [oldSequence] : [])];
    } else {
      newSequences = [...(sequences || []), ...(project.sequences || [])];
    }
    project.sequences = newSequences.filter(
      (seq, index) => newSequences.findIndex((s) => s.shkeptncontext === seq.shkeptncontext) === index
    );
  }

  private updateSequencesUpdated(response: HttpResponse<SequenceResult>, projectName: string): void {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body?.states[0] ? moment(response.body.states[0]?.time) : null;
    this._sequencesLastUpdated[projectName] = (
      lastEvent && lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated
    ).toDate();
  }

  private updateTracesUpdated(response: HttpResponse<EventResult>, keptnContext: string): void {
    const lastUpdated = moment(response.headers.get('date'));
    const lastEvent = response.body?.events[0] ? moment(response.body.events[0]?.time) : null;
    this._tracesLastUpdated[keptnContext] = (
      lastEvent && lastUpdated.isBefore(lastEvent) ? lastEvent : lastUpdated
    ).toDate();
  }

  protected allSequencesLoaded(sequences: number, totalCount: number, fromTime?: Date, beforeTime?: Date): boolean {
    return (
      (!fromTime && !beforeTime && sequences >= totalCount) ||
      (!!beforeTime && !fromTime && totalCount < this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE)
    );
  }

  public getRoot(projectName: string, shkeptncontext: string): Observable<Root | undefined> {
    return this.apiService.getRoots(projectName, 1, undefined, undefined, undefined, shkeptncontext).pipe(
      map((response) => response.body?.events || []),
      switchMap((roots) => this.rootMapper(roots).pipe(map((sequences) => sequences.pop())))
    );
  }

  public loadOldSequences(project: Project, fromTime?: Date, oldSequence?: Sequence): void {
    if (project.sequences) {
      this.loadSequences(
        project,
        fromTime,
        new Date(project.sequences[project.sequences.length - 1].time),
        oldSequence
      );
    }
  }

  public getSequenceWithTraces(projectName: string, keptnContext: string): Observable<Sequence | undefined> {
    return this.getSequenceByContext(projectName, keptnContext).pipe(
      switchMap((sequence) => (sequence ? this.sequenceMapper([sequence]) : [])),
      map((sequences) => sequences.shift())
    );
  }

  public updateSequence(projectName: string, keptnContext: string): void {
    this.getSequenceByContext(projectName, keptnContext).subscribe((sequence) => {
      const project = this._projects.getValue()?.find((p) => p.projectName === projectName);
      const sequences = project?.sequences;
      const oldSequence = sequences?.find((seq) => seq.shkeptncontext === keptnContext);
      if (oldSequence && sequence) {
        const { traces, ...copySequence } = sequence; // don't overwrite loaded traces
        Object.assign(oldSequence, copySequence);
      }
      this._sequencesUpdated.next();
    });
  }

  public loadUntilRoot(project: Project, shkeptncontext: string): void {
    this.getSequenceWithTraces(project.projectName, shkeptncontext).subscribe((sequence) => {
      if (sequence) {
        this.loadOldSequences(project, new Date(sequence.time), sequence);
      }
    });
  }

  public loadTraces(sequence: Sequence): void {
    const fromTime: Date = this._tracesLastUpdated[sequence.shkeptncontext];
    this.apiService
      .getTraces(sequence.shkeptncontext, sequence.project, fromTime?.toISOString())
      .pipe(
        map((response) => {
          this.updateTracesUpdated(response, sequence.shkeptncontext);
          return response.body;
        }),
        map((result) => result?.events || []),
        map((traces) => traces.map((trace) => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        sequence.traces = Trace.traceMapper([...(traces || []), ...(sequence.traces || [])]);
        this.getProject(sequence.project)
          .pipe(take(1))
          .subscribe((project) => {
            if (project) {
              project.stages
                .filter((s) => sequence.getStages().includes(s.stageName))
                .forEach((stage) => {
                  this.stageSequenceMapper(stage, project);
                });
            }
          });
        this._sequencesUpdated.next();
      });
  }

  public loadTracesByContext(shkeptncontext: string): void {
    this.apiService
      .getTraces(shkeptncontext)
      .pipe(
        map((response) => response.body),
        map((result) => result?.events || []),
        map((traces) => traces.map((trace) => Trace.fromJSON(trace)))
      )
      .subscribe((traces: Trace[]) => {
        this._traces.next(traces);
      });
  }

  public getEvent(type?: string, project?: string, stage?: string, service?: string): Observable<Trace | undefined> {
    return this.apiService.getEvent(type, project, stage, service).pipe(map((result) => result.events[0]));
  }

  public getSequenceByContext(projectName: string, shkeptncontext: string): Observable<Sequence | undefined> {
    return this.apiService
      .getSequences(projectName, 1, undefined, undefined, undefined, undefined, shkeptncontext)
      .pipe(
        map((response) => response.body?.states || []),
        map((sequences) => sequences.map((sequence) => Sequence.fromJSON(sequence)).shift())
      );
  }

  public getEvaluationResults(event: Trace | EventData, limit?: number, useFromTime = true): Observable<Trace[]> {
    let fromTime: Date | undefined;
    let eventData: EventData | undefined;
    if (event instanceof Trace) {
      const time = event.data.evaluationHistory?.[event.data.evaluationHistory.length - 1]?.time;
      if (time && useFromTime) {
        fromTime = new Date(time);
      }
      if (event.data.project && event.data.service && event.data.stage) {
        eventData = {
          project: event.data.project,
          service: event.data.service,
          stage: event.data.stage,
        };
      }
    } else {
      eventData = event;
    }
    if (eventData) {
      return this.apiService
        .getEvaluationResults(eventData.project, eventData.service, eventData.stage, fromTime?.toISOString(), limit)
        .pipe(
          map((result) => result.events || []),
          map((traces) => traces.map((trace) => Trace.fromJSON(trace)))
        );
    } else {
      return of([]);
    }
  }

  public loadEvaluationResults(event: Trace): void {
    this.getEvaluationResults(event).subscribe((traces: Trace[]) => {
      if (traces.length) {
        this._evaluationResults.next({
          type: 'evaluationHistory',
          triggerEvent: event,
          traces,
        });
      }
    });
  }

  public sendApprovalEvent(approval: Trace, approve: boolean): Observable<unknown> {
    const approval$ = this.apiService.sendApprovalEvent(
      approval,
      approve,
      EventTypes.APPROVAL_FINISHED,
      'approval.finished'
    );

    approval$.subscribe(() => {
      const project = this._projects.getValue()?.find((p) => p.projectName === approval.data.project);
      if (project?.projectName) {
        const stage = project.stages.find((st) => st.stageName === approval.data.stage);
        const service = stage?.services.find((sv) => sv.serviceName === approval.data.service);
        const sequence = service?.sequences.find((seq) => seq.shkeptncontext === approval.shkeptncontext);

        if (sequence) {
          // update data of sequence screen
          this.loadTraces(sequence);
        }
        if (service) {
          // update data of environment screen
          this.updateServiceApproval(service, approval);
        }
      }
    });
    return approval$;
  }

  private updateServiceApproval(service: Service, approval: Trace): void {
    const approvalIndex = service.openApprovals.findIndex((a) => a.trace.id === approval.id);
    if (approvalIndex >= 0 && approval.data.project) {
      service.openApprovals.splice(approvalIndex, 1);
      this.apiService
        .getSequences(
          approval.data.project,
          1,
          undefined,
          undefined,
          undefined,
          undefined,
          service.getLatestEvent()?.keptnContext
        )
        .subscribe((response) => {
          const seq = response.body?.states[0];
          if (seq) {
            service.latestSequence = Sequence.fromJSON(seq);
          }
        });
    }
  }

  public sendSequenceControl(sequence: Sequence, state: string): void {
    sequence.setState(SequenceState.UNKNOWN);
    this.apiService.sendSequenceControl(sequence.project, sequence.shkeptncontext, state).subscribe(() => {
      this.updateSequence(sequence.project, sequence.shkeptncontext);
    });
  }

  public invalidateEvaluation(evaluation: Trace, reason: string): void {
    this.apiService
      .sendEvaluationInvalidated(evaluation, reason)
      .pipe(take(1))
      .subscribe(() => {
        this._evaluationResults.next({
          type: 'invalidateEvaluation',
          triggerEvent: evaluation,
        });
      });
  }

  public getTaskNames(projectName: string): Observable<string[]> {
    return this.apiService
      .getTaskNames(projectName)
      .pipe(map((taskNames) => taskNames.sort((taskA, taskB) => taskA.localeCompare(taskB))));
  }

  public getServiceNames(projectName: string): Observable<string[]> {
    return this.apiService
      .getServiceNames(projectName)
      .pipe(map((serviceNames) => serviceNames.sort((serviceA, serviceB) => serviceA.localeCompare(serviceB))));
  }

  public getCustomSequenceNames(projectName: string): Observable<string[]> {
    return this.apiService.getCustomSequenceNames(projectName);
  }

  public getWebhookConfig(
    subscriptionId: string,
    projectName: string,
    stageName?: string,
    serviceName?: string
  ): Observable<WebhookConfig> {
    return this.apiService.getWebhookConfig(subscriptionId, projectName, stageName, serviceName);
  }

  public getFileTreeForService(projectName: string, serviceName: string): Observable<FileTree[]> {
    return this.apiService.getFileTreeForService(projectName, serviceName);
  }

  private sequenceMapper(sequences: Sequence[]): Observable<Sequence[]> {
    return from(sequences).pipe(
      mergeMap((sequence) =>
        this.apiService.getTraces(sequence.shkeptncontext, sequence.project).pipe(
          map((response) => {
            this.updateTracesUpdated(response, sequence.shkeptncontext);
            return response.body;
          }),
          map((result) => result?.events || []),
          map(Trace.traceMapper),
          map((traces) => {
            sequence.traces = traces;
            return sequence;
          })
        )
      ),
      toArray()
    );
  }

  protected stageSequenceMapper(stage: Stage, project: Project): void {
    stage.services.forEach((service) => {
      if (project.sequences) {
        service.sequences = project.sequences.filter(
          (s) => s.service === service.serviceName && s.getStages().includes(stage.stageName)
        );
      }
    });
  }

  private rootMapper(roots: Trace[]): Observable<Root[]> {
    return from(roots).pipe(
      mergeMap((root) =>
        this.apiService.getTraces(root.shkeptncontext, root.data.project).pipe(
          map((result) => result.body?.events || []),
          map(Trace.traceMapper),
          map((traces) => ({ ...root, traces }))
        )
      ),
      toArray(),
      map((rs) => rs.map((root) => Root.fromJSON(root)))
    );
  }

  public getServiceStates(projectName: string): Observable<ServiceState[]> {
    return this.apiService
      .getServiceStates(projectName)
      .pipe(map((serviceStates) => serviceStates.map((state) => ServiceState.fromJSON(state))));
  }

  public getServiceDeployment(projectName: string, keptnContext: string, fromTime?: string): Observable<Deployment> {
    return this.apiService.getServiceDeployment(projectName, keptnContext, fromTime).pipe(
      map((deployment) => {
        return Deployment.fromJSON(deployment);
      })
    );
  }

  public getOpenRemediationsOfService(
    projectName: string,
    serviceName: string
  ): Observable<ServiceRemediationInformation> {
    return this.apiService.getOpenRemediationsOfService(projectName, serviceName).pipe(
      map((serviceRemediationInformation) => {
        return ServiceRemediationInformation.fromJSON(serviceRemediationInformation);
      })
    );
  }

  public getIntersectedEvent(
    event: string,
    eventSuffix: string,
    projectName: string,
    stages: string[],
    services: string[]
  ): Observable<Record<string, unknown>> {
    return this.apiService.getIntersectedEvent(event, eventSuffix, projectName, stages, services);
  }

  public logout(): Observable<EndSessionData | null> {
    return this.apiService.logout();
  }

  public triggerDelivery(data: TriggerSequenceData): Observable<TriggerResponse> {
    const type = EventTypes.PREFIX + data.stage + EventTypes.DELIVERY_TRIGGERED_SUFFIX;

    return this.apiService.triggerSequence(type, data);
  }

  public triggerEvaluation(data: TriggerSequenceData): Observable<TriggerResponse> {
    const type = EventTypes.PREFIX + data.stage + EventTypes.EVALUATION_TRIGGERED_SUFFIX;
    return this.apiService.triggerSequence(type, data);
  }

  public triggerCustomSequence(data: TriggerSequenceData, sequence: string): Observable<TriggerResponse> {
    const type = EventTypes.PREFIX + data.stage + '.' + sequence + '.triggered';
    return this.apiService.triggerSequence(type, data);
  }

  public getSecretScopes(): Observable<string[]> {
    return this.apiService.getSecretScopes().pipe(map((result) => result.scopes));
  }
}
