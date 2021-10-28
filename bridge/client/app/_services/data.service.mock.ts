/* eslint-disable @typescript-eslint/no-unused-vars */
import { Injectable } from '@angular/core';
import { DataService } from './data.service';
import { Project } from '../_models/project';
import { KeptnInfo } from './_mockData/keptnInfo.mock';
import { Projects } from './_mockData/projects.mock';
import { Traces } from './_mockData/traces.mock';
import { Evaluations } from './_mockData/evaluations.mock';
import { Trace } from '../_models/trace';
import { map } from 'rxjs/operators';
import { Observable, of } from 'rxjs';
import { Sequence } from '../_models/sequence';
import { UniformRegistrationsMock } from '../_models/uniform-registrations.mock';
import { UniformRegistrationLog } from '../../../server/interfaces/uniform-registration-log';
import { UniformRegistrationLogsMock } from '../_models/uniform-registrations-logs.mock';
import { SequencesData } from './_mockData/sequences.mock';
import { UniformRegistration } from '../_models/uniform-registration';
import { UniformSubscription } from '../_models/uniform-subscription';
import { WebhookConfig } from '../../../shared/models/webhook-config';
import { AppUtils } from '../_utils/app.utils';
import { WebhookConfigMock } from './_mockData/webhook-config.mock';
import { FileTreeMock } from '../_models/fileTree.mock';
import { FileTree } from '../../../shared/interfaces/resourceFileTree';
import { UniformRegistrationInfo } from '../../../shared/interfaces/uniform-registration-info';
import { Secret } from '../_models/secret';
import { SecretScope } from '../../../shared/interfaces/secret-scope';

@Injectable({
  providedIn: 'root',
})
export class DataServiceMock extends DataService {
  public loadKeptnInfo(): void {
    this._keptnInfo.next(KeptnInfo);
  }

  public loadProjects(): Observable<Project[]> {
    const projects = Projects.map((project) => Project.fromJSON(project));
    this._projects.next(projects);
    return of(projects);
  }

  public loadProject(projectName: string): void {
    this._projects.next([...Projects]);
  }

  public loadPlainProject(projectName: string): Observable<Project> {
    return of(Projects[0]);
  }

  public loadSequences(project: Project, fromTime?: Date, beforeTime?: Date, oldSequence?: Sequence): void {
    let totalCount;
    let sequences: Sequence[];
    if (beforeTime) {
      const currentIndex = project.sequences?.length || 0;
      sequences = SequencesData.slice(currentIndex, currentIndex + this.DEFAULT_NEXT_SEQUENCE_PAGE_SIZE);
      totalCount = sequences.length;
    } else {
      totalCount = SequencesData.length;
      sequences = SequencesData.slice(0, this.DEFAULT_SEQUENCE_PAGE_SIZE);
    }
    this.addNewSequences(project, sequences, !!beforeTime, oldSequence);

    if (this.allSequencesLoaded(project.sequences?.length || 0, totalCount, fromTime, beforeTime)) {
      project.allSequencesLoaded = true;
    }
    project.stages.forEach((stage) => {
      this.stageSequenceMapper(stage, project);
    });
    this._sequencesUpdated.next();
  }

  public getProject(projectName: string): Observable<Project | undefined> {
    if (!this._projects.getValue()?.length) {
      this.loadProjects();
    }
    return this._projects.pipe(map((projects) => projects?.find((project) => project.projectName === projectName)));
  }

  public createProject(
    projectName: string,
    shipyard: string,
    gitRemoteUrl?: string,
    gitToken?: string,
    gitUser?: string
  ): Observable<unknown> {
    return of(undefined);
  }

  public deleteProject(projectName: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public loadTraces(sequence: Sequence): void {
    sequence.traces = [...(Traces || []), ...(sequence.traces || [])];
    this._sequencesUpdated.next();
  }

  public loadTracesByContext(shkeptncontext: string): void {
    this._traces.next(Traces.filter((t) => t.shkeptncontext === shkeptncontext));
  }

  public loadEvaluationResults(event: Trace): void {
    this._evaluationResults.next({
      type: 'evaluationHistory',
      triggerEvent: event,
      traces: [Evaluations],
    });
  }

  public setGitUpstreamUrl(
    projectName: string,
    gitUrl: string,
    gitUser: string,
    gitToken: string
  ): Observable<boolean> {
    this.loadProjects();
    return of(true);
  }

  public getUniformRegistrations(): Observable<UniformRegistration[]> {
    const copyUniform = AppUtils.copyObject(UniformRegistrationsMock);
    return of(copyUniform.map((registration) => UniformRegistration.fromJSON(registration)));
  }

  public getUniformRegistrationLogs(): Observable<UniformRegistrationLog[]> {
    return of(UniformRegistrationLogsMock);
  }

  public deleteSubscription(integrationId: string, subscriptionId: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public getTaskNames(projectName: string): Observable<string[]> {
    return of(['approval', 'deployment', 'test']);
  }

  public updateUniformSubscription(
    integrationId: string,
    subscription: UniformSubscription
  ): Observable<Record<string, unknown>> {
    return of({});
  }

  public createUniformSubscription(
    integrationId: string,
    subscription: UniformSubscription
  ): Observable<Record<string, unknown>> {
    return of({});
  }

  public createService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public deleteService(projectName: string, serviceName: string): Observable<Record<string, unknown>> {
    return of({});
  }

  public getWebhookConfig(projectName: string, stageName?: string, serviceName?: string): Observable<WebhookConfig> {
    return of(WebhookConfigMock);
  }

  public getFileTreeForService(projectName: string, serviceName: string): Observable<FileTree[]> {
    return of(FileTreeMock);
  }

  public getUniformRegistrationInfo(integrationId: string): Observable<UniformRegistrationInfo> {
    const registration = UniformRegistrationsMock.find((r) => r.id === integrationId);
    return of({
      isWebhookService: registration?.isWebhookService ?? false,
      isControlPlane: registration?.metadata.location === 'control-plane' ?? false,
    });
  }

  public getSecretsForScope(scope: SecretScope): Observable<Secret[]> {
    const secrets = [new Secret(), new Secret()];
    secrets[0].name = 'SecretA';
    secrets[0].keys = ['key1', 'key2', 'key3'];
    secrets[1].name = 'SecretB';
    secrets[1].keys = ['key1', 'key2', 'key3'];

    return of(secrets);
  }
}
/* eslint-enable @typescript-eslint/no-unused-vars */
