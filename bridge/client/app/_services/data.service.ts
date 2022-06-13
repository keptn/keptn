import { Injectable } from '@angular/core';
import { BehaviorSubject, forkJoin, Observable, of, Subject } from 'rxjs';
import { catchError, map, mergeMap, take, tap } from 'rxjs/operators';
import { Trace } from '../_models/trace';
import { Project } from '../_models/project';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ApiService } from './api.service';
import moment from 'moment';
import { Sequence } from '../_models/sequence';
import { UniformRegistrationLog } from '../../../shared/interfaces/uniform-registration-log';
import { Secret } from '../_models/secret';
import { HttpResponse } from '@angular/common/http';
import { SequenceResult } from '../_models/sequence-result';
import { EventResult } from '../../../shared/interfaces/event-result';
import { KeptnInfo } from '../_models/keptn-info';
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
import { ISequencesFilter } from '../../../shared/interfaces/sequencesFilter';
import { TriggerResponse, TriggerSequenceData } from '../_models/trigger-sequence';
import { EventData } from '../_components/ktb-evaluation-info/ktb-evaluation-info.component';
import { SecretScope } from '../../../shared/interfaces/secret-scope';
import { IGitDataExtended } from '../_interfaces/git-upstream';
import { getGitData } from '../_utils/git-upstream.utils';
import { ICustomSequences } from '../../../shared/interfaces/custom-sequences';
import { KeptnService } from '../../../shared/models/keptn-service';
import { IMetadata } from '../_interfaces/metadata';

@Injectable({
  providedIn: 'root',
})
export class DataService {
  private _projects = new BehaviorSubject<Project[] | undefined>(undefined);
  private _sequencesUpdated = new Subject<void>();
  private _traces = new BehaviorSubject<Trace[] | undefined>(undefined);
  private _openApprovals = new BehaviorSubject<Trace[]>([]);
  private _keptnInfo = new BehaviorSubject<KeptnInfo | undefined>(undefined);
  private _keptnMetadata = new BehaviorSubject<IMetadata | undefined | null>(undefined); // fetched | not fetched | not existing
  private _sequencesLastUpdated: { [key: string]: Date } = {};
  private _tracesLastUpdated: { [key: string]: Date } = {};
  private _projectName: BehaviorSubject<string> = new BehaviorSubject<string>('');
  private _uniformDates: { [key: string]: string } = this.apiService.uniformLogDates;
  private _hasUnreadUniformRegistrationLogs = new BehaviorSubject<boolean>(false);
  private readonly DEFAULT_SEQUENCE_PAGE_SIZE = 25;
  private readonly DEFAULT_NEXT_SEQUENCE_PAGE_SIZE = 10;
  private _isQualityGatesOnly = new BehaviorSubject<boolean>(false);
  private _evaluationResults = new Subject<EvaluationHistory>();

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

  get keptnMetadata(): Observable<IMetadata | undefined | null> {
    return this._keptnMetadata.asObservable();
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

  public getService(projectName: string, stageName: string, serviceName: string): Observable<Service> {
    return this.apiService
      .getService(projectName, stageName, serviceName)
      .pipe(map((service) => Service.fromJSON(service)));
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

  public createProjectExtended(projectName: string, shipyard: string, data: IGitDataExtended): Observable<unknown> {
    return this.apiService.createProjectExtended(projectName, shipyard, getGitData(data));
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
    if (this._hasUnreadUniformRegistrationLogs.getValue()) {
      return;
    }
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

  public getTracesLastUpdated(sequence: Sequence): Date | undefined {
    return this._tracesLastUpdated[sequence.shkeptncontext];
  }

  public setGitUpstreamUrl(
    projectName: string,
    gitUrl: string,
    gitToken: string,
    gitUser?: string
  ): Observable<unknown> {
    return this.apiService.sendGitUpstreamUrl(projectName, gitUrl, gitToken, gitUser).pipe(
      tap(() => {
        this.loadProject(projectName);
      })
    );
  }

  public updateGitUpstream(projectName: string, data: IGitDataExtended): Observable<unknown> {
    return this.apiService.updateGitUpstreamExtended(projectName, getGitData(data));
  }

  public loadKeptnInfo(): void {
    // #4165 Get bridge info first to get info if versions.json should be loaded or not
    // Versions should not be loaded if enableVersionCheckFeature is set to false (when ENABLE_VERSION_CHECK is set to false in env)
    this.apiService
      .getKeptnInfo()
      .pipe(
        mergeMap((bridgeInfo) => {
          return forkJoin([
            of(bridgeInfo),
            bridgeInfo.enableVersionCheckFeature
              ? this.apiService.getAvailableVersions().pipe(catchError(() => of(undefined)))
              : of(undefined),
            of(this.apiService.isVersionCheckEnabled()),
            this.apiService.getMetadata().pipe(catchError(() => of(null))),
          ]);
        })
      )
      .subscribe(([bridgeInfo, availableVersions, versionCheckEnabled, metadata]) => {
        const keptnInfo: KeptnInfo = {
          bridgeInfo,
          availableVersions,
          versionCheckEnabled,
        };

        if (keptnInfo.bridgeInfo.showApiToken) {
          keptnInfo.bridgeInfo.apiUrl = `${window.location.href.substring(
            0,
            window.location.href.indexOf(window.location.href.includes('bridge') ? '/bridge' : window.location.pathname)
          )}/api`;

          keptnInfo.authCommand = `keptn auth --endpoint=${keptnInfo.bridgeInfo.apiUrl} --api-token=${keptnInfo.bridgeInfo.apiToken}`;

          this._isQualityGatesOnly.next(!keptnInfo.bridgeInfo.keptnInstallationType?.includes('CONTINUOUS_DELIVERY'));
        }
        this._keptnInfo.next(keptnInfo);
        this._keptnMetadata.next(metadata);
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
          }
        }
      );
  }

  public loadProjects(): Observable<Project[]> {
    return this.apiService.getProjects(this._keptnInfo.getValue()?.bridgeInfo.projectsPageSize || 50).pipe(
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
      ),
      tap((projects) => {
        this._projects.next(projects);
      })
    );
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
        this._sequencesUpdated.next();
      });
  }

  public loadLatestSequences(project: Project, pageSize: number): Observable<Sequence[]> {
    return this.apiService.getSequences(project.projectName, pageSize).pipe(
      map((response) => response.body),
      map((body) => body?.states.map((sequence) => Sequence.fromJSON(sequence)) ?? [])
    );
  }

  public getSequenceFilter(projectName: string): Observable<ISequencesFilter> {
    return this.apiService.getSequencesFilter(projectName);
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
    this.getSequenceByContext(project.projectName, shkeptncontext).subscribe((sequence) => {
      if (sequence) {
        this.loadOldSequences(project, new Date(sequence.time), sequence);
      }
    });
  }

  public getTracesOfSequence(sequence: Sequence): Observable<Trace[]> {
    const fromTime: Date = this._tracesLastUpdated[sequence.shkeptncontext];
    return this.apiService.getTraces(sequence.shkeptncontext, sequence.project, fromTime?.toISOString()).pipe(
      map((response) => {
        this.updateTracesUpdated(response, sequence.shkeptncontext);
        return response.body;
      }),
      map((result) => result?.events || []),
      map((traces) => traces.map((trace) => Trace.fromJSON(trace))),
      map((traces) => Trace.traceMapper([...(traces || []), ...(sequence.traces || [])]))
    );
  }

  public getTracesByContext(
    keptnContext: string,
    type?: EventTypes,
    source?: KeptnService,
    pageSize?: number
  ): Observable<Trace[]> {
    return this.apiService.getTraces(keptnContext, undefined, undefined, type, source, pageSize).pipe(
      map((response) => response.body?.events || []),
      map((traces) => traces.map((trace) => Trace.fromJSON(trace)))
    );
  }

  public loadTracesByContext(keptnContext: string): void {
    this.getTracesByContext(keptnContext).subscribe((traces: Trace[]) => {
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
    return this.apiService.sendApprovalEvent(approval, approve, EventTypes.APPROVAL_FINISHED, 'approval.finished').pipe(
      tap(() => {
        const project = this._projects.getValue()?.find((p) => p.projectName === approval.data.project);
        if (!project?.projectName) {
          return;
        }
        const stage = project.stages.find((st) => st.stageName === approval.data.stage);
        const service = stage?.services.find((sv) => sv.serviceName === approval.data.service);
        const sequence = project.sequences?.find((seq) => seq.shkeptncontext === approval.shkeptncontext);

        if (sequence) {
          // update data of sequence screen
          this.getTracesOfSequence(sequence).subscribe((traces) => {
            sequence.traces = traces;
          });
        }
        if (service) {
          // update data of environment screen
          this.updateServiceApproval(service, approval);
        }
      })
    );
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

  public getCustomSequences(projectName: string): Observable<ICustomSequences> {
    return this.apiService.getCustomSequences(projectName);
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
