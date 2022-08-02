import { Component, HostListener, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { DtToast } from '@dynatrace/barista-components/toast';
import { combineLatest, Observable, of, Subject } from 'rxjs';
import { catchError, filter, finalize, map, mergeMap, startWith, takeUntil } from 'rxjs/operators';
import { IGitDataExtended } from 'shared/interfaces/project';
import { IClientFeatureFlags } from '../../../../../shared/interfaces/feature-flags';
import { PendingChangesComponent } from '../../../_guards/pending-changes.guard';
import { DeleteData, DeleteResult, DeleteType, DeletionProgressEvent } from '../../../_interfaces/delete';
import { KeptnInfo } from '../../../_models/keptn-info';
import { NotificationType } from '../../../_models/notification';
import { Project } from '../../../_models/project';
import { ServerErrors } from '../../../_models/server-error';
import { DataService } from '../../../_services/data.service';
import { EventService } from '../../../_services/event.service';
import { FeatureFlagsService } from '../../../_services/feature-flags.service';
import { NotificationsService } from '../../../_services/notifications.service';
import { FormUtils } from '../../../_utils/form.utils';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message/ktb-project-create-message.component';
import { IGitDataExtendedWithNoUpstream } from './ktb-project-settings-git-extended/ktb-project-settings-git-extended.component';

type DialogState = null | 'unsaved';

enum ProjectSettingsStatus {
  ERROR,
  INIT,
  LOADED,
}

interface ProjectSettingsState {
  projectName?: string;
  resourceServiceEnabled?: boolean;
  gitUpstreamRequired?: boolean;
  gitInputDataExtended?: IGitDataExtended;
  automaticProvisioningMessage?: string;
  state: ProjectSettingsStatus;
}

@Component({
  selector: 'ktb-project-settings',
  templateUrl: './ktb-project-settings.component.html',
  styleUrls: ['./ktb-project-settings.component.scss'],
})
export class KtbProjectSettingsComponent implements OnInit, OnDestroy, PendingChangesComponent {
  private readonly unsubscribe$ = new Subject<void>();
  public ServerErrors = ServerErrors;
  public ProjectSettingsStatus = ProjectSettingsStatus;

  @ViewChild('deleteProjectDialog')
  private deleteProjectDialog?: TemplateRef<MatDialog>;
  public isGitUpstreamInProgress = false;
  public isCreatingProjectInProgress = false;
  private pendingChangesSubject = new Subject<boolean>();
  public isProjectFormTouched = false;
  public shipyardFile?: File;
  public gitDataExtended?: IGitDataExtendedWithNoUpstream;
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl,
  });

  public message = 'You have pending changes. Make sure to save your data before you continue.';
  public unsavedDialogState: DialogState = null;

  public resourceServiceEnabled$ = this.featureFlagService.featureFlags$.pipe(
    map((featureFlags: IClientFeatureFlags) => featureFlags.RESOURCE_SERVICE_ENABLED)
  );

  public projectName$: Observable<string | null> = this.route.paramMap.pipe(map((params) => params.get('projectName')));

  public gitInputDataExtended$: Observable<IGitDataExtended | undefined> = this.projectName$.pipe(
    mergeMap((projectName) => (projectName ? this.dataService.loadPlainProject(projectName) : of(undefined))),
    map((project) => project?.gitCredentials)
  );

  public projectNames$: Observable<string[] | undefined> = this.projectName$.pipe(
    mergeMap((projectName) => (!projectName ? this.dataService.projects : of(undefined))),
    map((projects: Project[] | undefined) => projects?.map((project) => project.projectName))
  );

  public projectCreated$: Observable<boolean> = this.route.queryParams.pipe(map((queryParams) => queryParams.created));

  readonly state$: Observable<ProjectSettingsState> = combineLatest([
    this.dataService.keptnInfo.pipe(filter((info): info is KeptnInfo => !!info)),
    this.dataService.keptnMetadata,
    this.projectName$.pipe(map((projectName) => projectName ?? undefined)),
    this.resourceServiceEnabled$,
    this.gitInputDataExtended$,
    this.projectNames$,
    this.projectCreated$,
  ]).pipe(
    map(
      ([
        keptnInfo,
        metadata,
        projectName,
        resourceServiceEnabled,
        gitInputDataExtended,
        projectNames,
        projectCreated,
      ]) => {
        if (projectName && projectCreated) {
          this.showCreateNotificationAndRedirect(projectName);
        }

        if (projectNames) {
          this.projectNameControl.setValidators([
            Validators.required,
            FormUtils.nameExistsValidator(projectNames),
            Validators.pattern('[a-z]([a-z]|[0-9]|-)*'),
          ]);
        }

        return {
          projectName,
          resourceServiceEnabled,
          gitUpstreamRequired: !metadata?.automaticprovisioning,
          gitInputDataExtended,
          automaticProvisioningMessage: keptnInfo.bridgeInfo.automaticProvisioningMsg,
          state: metadata === null ? ProjectSettingsStatus.ERROR : ProjectSettingsStatus.LOADED,
        };
      }
    ),
    startWith({
      projectName: undefined,
      resourceServiceEnabled: undefined,
      gitUpstreamRequired: undefined,
      gitInputDataExtended: undefined,
      automaticProvisioningMessage: undefined,
      state: ProjectSettingsStatus.INIT,
      projectDeletionData: undefined,
    })
  );

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private toast: DtToast,
    private router: Router,
    private notificationsService: NotificationsService,
    private eventService: EventService,
    private featureFlagService: FeatureFlagsService
  ) {}

  public ngOnInit(): void {
    this.eventService.deletionTriggeredEvent.pipe(takeUntil(this.unsubscribe$)).subscribe((data) => {
      if (data.type === DeleteType.PROJECT) {
        this.deleteProject(data.name);
      }
    });
  }

  private showCreateNotificationAndRedirect(projectName: string): void {
    this.notificationsService.addNotification(
      NotificationType.SUCCESS,
      '',
      {
        component: KtbProjectCreateMessageComponent,
        data: {
          projectName: projectName,
          routerLink: `/project/${projectName}/settings/services/create`,
        },
      },
      10_000
    );
    // Remove query param for not showing notification on reload
    this.router.navigate(['/', 'project', projectName, 'settings', 'project']);
  }

  public updateGitDataExtended(data?: IGitDataExtendedWithNoUpstream): void {
    this.gitDataExtended = data;
    this.projectFormTouched();
  }

  public updateShipyardFile(shipyardFile: File | undefined): void {
    this.shipyardFile = shipyardFile;
    this.projectFormTouched();
  }

  public async createProject(): Promise<void> {
    if (!this.shipyardFile || this.gitDataExtended === undefined) {
      return;
    }
    const fileContent = await FormUtils.readFileContent(this.shipyardFile);
    if (!fileContent) {
      return;
    }
    this.isCreatingProjectInProgress = true;

    const shipyardBase64 = btoa(fileContent);
    const projectName = this.projectNameControl.value;

    this.dataService
      .createProjectExtended(projectName, shipyardBase64, this.gitDataExtended ?? undefined)
      .pipe(
        map(() => true),
        catchError((err) => {
          const errorMessage = err.error || 'please, check the logs of resource-service';
          this.notificationsService.addNotification(
            NotificationType.ERROR,
            `The project could not be created: ${errorMessage}.`
          );
          return of(false);
        }),
        finalize(() => (this.isCreatingProjectInProgress = false))
      )
      .subscribe((success) => {
        if (!success) {
          return;
        }
        this.isProjectFormTouched = false;
        this.router.navigate(['/', 'project', projectName, 'settings', 'project'], {
          queryParams: { created: true },
        });
      });
  }

  private deleteProject(projectName: string): void {
    this.eventService.deletionProgressEvent.next({ isInProgress: true });

    this.dataService
      .deleteProject(projectName)
      .pipe(
        map((): DeletionProgressEvent => {
          return { isInProgress: false, result: DeleteResult.SUCCESS };
        }),
        catchError((err): Observable<DeletionProgressEvent> => {
          const deletionError = 'Project could not be deleted: ' + err.message;
          return of({
            error: deletionError,
            isInProgress: false,
            result: DeleteResult.ERROR,
          });
        })
      )
      .subscribe((progressEvent) => {
        this.eventService.deletionProgressEvent.next(progressEvent);
        if (progressEvent.result === DeleteResult.SUCCESS) {
          this.isProjectFormTouched = false;
          this.router.navigate(['/', 'dashboard']);
        }
      });
  }

  public reject(): void {
    this.pendingChangesSubject.next(false);
    this.hideNotification();
  }

  public reset(): void {
    this.pendingChangesSubject.next(true);
    this.hideNotification();
  }

  public saveAll(projectName?: string): void {
    if (this.isCreateMode(projectName)) {
      this.createProject();
    } else {
      this.updateGitUpstream(projectName);
    }
    this.hideNotification();
  }

  public updateGitUpstream(projectName?: string): void {
    if (!this.gitDataExtended || !projectName) {
      return;
    }
    this.isGitUpstreamInProgress = true;
    this.dataService
      .updateGitUpstream(projectName, this.gitDataExtended)
      .pipe(finalize(() => (this.isGitUpstreamInProgress = false)))
      .subscribe(() => {
        this.notificationsService.addNotification(
          NotificationType.SUCCESS,
          'The Git upstream was changed successfully.'
        );
        this.isProjectFormTouched = false;
      });
  }

  public isProjectFormInvalid(projectName?: string): boolean {
    return this.isCreateMode(projectName) ? this.isProjectCreateFormInvalid() : this.isProjectSettingsFormInvalid();
  }

  public isProjectCreateFormInvalid(): boolean {
    return (
      !this.shipyardFile ||
      this.projectNameForm.invalid ||
      this.gitDataExtended === undefined ||
      this.isCreatingProjectInProgress
    );
  }

  public isProjectSettingsFormInvalid(): boolean {
    return !this.gitDataExtended || this.isGitUpstreamInProgress;
  }

  public projectFormTouched(): void {
    this.isProjectFormTouched = true;
  }

  // @HostListener allows us to also guard against browser refresh, close, etc.
  @HostListener('window:beforeunload', ['$event'])
  public canDeactivate($event?: BeforeUnloadEvent): Observable<boolean> {
    if (this.isProjectFormTouched) {
      if ($event) {
        $event.returnValue = this.message;
      }
      this.showNotification();
      return this.pendingChangesSubject.asObservable();
    } else {
      return of(true);
    }
  }

  public showNotification(): void {
    this.unsavedDialogState = 'unsaved';

    document.querySelector('div[aria-label="Dialog for notifying about unsaved data"]')?.classList.add('shake');
    setTimeout(() => {
      document.querySelector('div[aria-label="Dialog for notifying about unsaved data"]')?.classList.remove('shake');
    }, 500);
  }

  public hideNotification(): void {
    this.unsavedDialogState = null;
  }

  public createProjectDeletionData(projectName?: string): DeleteData {
    return {
      type: DeleteType.PROJECT,
      name: projectName ?? '',
    };
  }

  public isCreateMode(projectName?: string): boolean {
    return !projectName;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
