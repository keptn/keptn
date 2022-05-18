import { Component, HostListener, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { BehaviorSubject, combineLatest, Observable, of, Subject } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { KtbProjectSettingsGitComponent } from '../ktb-project-settings-git/ktb-project-settings-git.component';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { DtToast } from '@dynatrace/barista-components/toast';
import { NotificationsService } from '../../_services/notifications.service';
import { EventService } from '../../_services/event.service';
import { filter, map, startWith, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { FormUtils } from '../../_utils/form.utils';
import { KtbProjectCreateMessageComponent } from '../_status-messages/ktb-project-create-message/ktb-project-create-message.component';
import { NotificationType } from '../../_models/notification';
import { PendingChangesComponent } from '../../_guards/pending-changes.guard';
import { IClientFeatureFlags } from '../../../../shared/interfaces/feature-flags';
import { IGitData, IGitDataExtended } from '../../_interfaces/git-upstream';
import { AppUtils } from '../../_utils/app.utils';
import { FeatureFlagsService } from '../../_services/feature-flags.service';
import { KeptnInfo } from '../../_models/keptn-info';
import { IMetadata } from '../../_interfaces/metadata';
import { ServerErrors } from '../../_models/server-error';

type DialogState = null | 'unsaved';

enum ProjectSettingsStatus {
  ERROR,
  INIT,
  LOADED,
}

interface ProjectSettingsState {
  gitUpstreamRequired: boolean | undefined;
  automaticProvisioningMessage: string | undefined;
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

  @ViewChild(KtbProjectSettingsGitComponent)
  private gitSettingsSection?: KtbProjectSettingsGitComponent;
  public gitInputDataExtended?: IGitDataExtended;
  public gitInputDataExtendedDefault?: IGitDataExtended;
  public projectName?: string;
  public projectDeletionData?: DeleteData;
  public isProjectLoading: boolean | undefined;
  public isCreateMode = false;
  public isGitUpstreamInProgress = false;
  public isCreatingProjectInProgress = false;
  private pendingChangesSubject = new Subject<boolean>();
  public isProjectFormTouched = false;
  public shipyardFile?: File;
  public gitData: IGitData = {
    gitFormValid: true,
  };
  private gitDataExtended?: IGitDataExtended;
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl,
  });
  public readonly _metadataError$ = new BehaviorSubject<boolean>(false);

  public message = 'You have pending changes. Make sure to save your data before you continue.';
  public unsavedDialogState: DialogState = null;
  public resourceServiceEnabled?: boolean;

  readonly state$: Observable<ProjectSettingsState> = combineLatest([
    this.dataService.keptnInfo,
    this.dataService.keptnMetadata,
  ]).pipe(
    filter((info): info is [KeptnInfo, IMetadata | undefined | null] => !!info[0]),
    map(([keptnInfo, metadata]) => {
      return {
        gitUpstreamRequired: !metadata?.automaticprovisioning,
        automaticProvisioningMessage: keptnInfo.bridgeInfo.automaticProvisioningMsg,
        state: metadata === null ? ProjectSettingsStatus.ERROR : ProjectSettingsStatus.LOADED,
      };
    }),
    startWith({
      gitUpstreamRequired: undefined,
      automaticProvisioningMessage: undefined,
      state: ProjectSettingsStatus.INIT,
    })
  );

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private toast: DtToast,
    private router: Router,
    private notificationsService: NotificationsService,
    private eventService: EventService,
    featureFlagService: FeatureFlagsService
  ) {
    featureFlagService.featureFlags$
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((featureFlags: IClientFeatureFlags) => {
        this.resourceServiceEnabled = featureFlags.RESOURCE_SERVICE_ENABLED;
      });
  }

  public ngOnInit(): void {
    this.route.params.subscribe((params) => {
      if (!params.projectName) {
        this.isCreateMode = true;
        this.isProjectLoading = true;

        this.loadProjectsAndSetValidator();
      } else {
        this.isProjectLoading = true;
        this.isCreateMode = false;
        this.isProjectFormTouched = false;
        this.projectName = params.projectName;

        this.projectDeletionData = {
          type: DeleteType.PROJECT,
          name: this.projectName || '',
        };

        this.loadProject(params.projectName);
      }
    });

    this.route.queryParams.subscribe((queryParams) => {
      if (queryParams.created) {
        this.showCreateNotificationAndRedirect();
      }
    });

    this.eventService.deletionTriggeredEvent.pipe(takeUntil(this.unsubscribe$)).subscribe((data) => {
      if (data.type === DeleteType.PROJECT) {
        this.deleteProject(data.name);
      }
    });
  }

  private loadProjectsAndSetValidator(): void {
    this.dataService
      .loadProjects()
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((projects: Project[] | undefined): projects is Project[] => !!projects),
        map((projects: Project[]) => projects.map((project) => project.projectName))
      )
      .subscribe((projectNames) => {
        this.projectNameControl.setValidators([
          Validators.required,
          FormUtils.nameExistsValidator(projectNames),
          Validators.pattern('[a-z]([a-z]|[0-9]|-)*'),
        ]);
        this.isProjectLoading = false;
      });
  }

  private loadProject(projectName: string): void {
    this.dataService.loadPlainProject(projectName).subscribe((project) => {
      this.gitData = {
        gitRemoteURL: project.gitRemoteURI,
        gitUser: project.gitUser,
      };
      this.gitInputDataExtendedDefault = project.gitUpstream;
      this.gitInputDataExtended = AppUtils.copyObject(project.gitUpstream); // there should not be a reference. Could
      // lead to problems when the form is reset

      this.isProjectLoading = false;
    });
  }

  private showCreateNotificationAndRedirect(): void {
    this.notificationsService.addNotification(
      NotificationType.SUCCESS,
      '',
      {
        component: KtbProjectCreateMessageComponent,
        data: {
          projectName: this.projectName,
          routerLink: `/project/${this.projectName}/settings/services/create`,
        },
      },
      10_000
    );
    // Remove query param for not showing notification on reload
    this.router.navigate(['/', 'project', this.projectName, 'settings', 'project']);
  }

  public updateGitData(gitData: IGitData): void {
    this.gitData.gitRemoteURL = gitData.gitRemoteURL;
    this.gitData.gitUser = gitData.gitUser;
    this.gitData.gitToken = gitData.gitToken;
    this.gitData.gitFormValid = gitData.gitFormValid;
    this.projectFormTouched();
  }

  public updateGitDataExtended(data?: IGitDataExtended): void {
    this.gitDataExtended = data;
    this.projectFormTouched();
  }

  public updateShipyardFile(shipyardFile: File | undefined): void {
    this.shipyardFile = shipyardFile;
    this.projectFormTouched();
  }

  public setGitUpstream(): void {
    if (this.projectName && this.gitData.gitRemoteURL && this.gitData.gitToken) {
      this.isGitUpstreamInProgress = true;
      this.hideNotification();
      this.dataService
        .setGitUpstreamUrl(this.projectName, this.gitData.gitRemoteURL, this.gitData.gitToken, this.gitData.gitUser)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(
          () => {
            this.isGitUpstreamInProgress = false;
            this.gitData.gitToken = '';
            this.gitData = { ...this.gitData };
            this.notificationsService.addNotification(
              NotificationType.SUCCESS,
              'The Git upstream was changed successfully.'
            );

            this.pendingChangesSubject.next(true);
            this.isProjectFormTouched = false;
          },
          (err) => {
            this.isGitUpstreamInProgress = false;
            this.notificationsService.addNotification(
              NotificationType.ERROR,
              `<div class="long-note align-left p-3">The Git upstream could not be changed:<br/><span class="small">${err.error}</span></div>`
            );
          }
        );
    }
  }

  public async createProject(): Promise<void> {
    if (this.shipyardFile) {
      this.isCreatingProjectInProgress = true;
      const fileContent = await FormUtils.readFileContent(this.shipyardFile);
      if (fileContent) {
        const shipyardBase64 = btoa(fileContent);
        const projectName = this.projectNameControl.value;
        const create$ =
          this.resourceServiceEnabled && this.gitDataExtended
            ? this.dataService.createProjectExtended(projectName, shipyardBase64, this.gitDataExtended)
            : this.dataService.createProject(
                projectName,
                shipyardBase64,
                this.gitData.gitRemoteURL,
                this.gitData.gitToken,
                this.gitData.gitUser
              );

        create$.subscribe(
          () => {
            this.projectName = projectName;
            this.dataService.loadProjects().subscribe(() => {
              this.isCreatingProjectInProgress = false;
              this.isProjectFormTouched = false;

              this.router.navigate(['/', 'project', this.projectName, 'settings', 'project'], {
                queryParams: { created: true },
              });
            });
          },
          (err) => {
            const service = this.resourceServiceEnabled ? 'resource-service' : 'configuration-service';
            const errorMessage = err.error || `please, check the logs of ${service}`;
            this.notificationsService.addNotification(
              NotificationType.ERROR,
              `The project could not be created: ${errorMessage}.`
            );
            this.isCreatingProjectInProgress = false;
          }
        );
      }
    }
  }

  public deleteProject(projectName: string): void {
    this.eventService.deletionProgressEvent.next({ isInProgress: true });

    this.dataService.deleteProject(projectName).subscribe(
      () => {
        this.eventService.deletionProgressEvent.next({ isInProgress: false, result: DeleteResult.SUCCESS });
        this.isProjectFormTouched = false;
        this.router.navigate(['/', 'dashboard']);
      },
      (err) => {
        const deletionError = 'Project could not be deleted: ' + err.message;
        this.eventService.deletionProgressEvent.next({
          error: deletionError,
          isInProgress: false,
          result: DeleteResult.ERROR,
        });
      }
    );
  }

  public reject(): void {
    this.pendingChangesSubject.next(false);
    this.hideNotification();
  }

  public reset(): void {
    if (this.resourceServiceEnabled) {
      this.gitInputDataExtended = AppUtils.copyObject(this.gitInputDataExtendedDefault);
    } else {
      this.gitSettingsSection?.reset();
    }
    this.pendingChangesSubject.next(true);
    this.hideNotification();
  }

  public saveAll(): void {
    if (this.isCreateMode) {
      this.createProject();
    } else {
      this.setGitUpstream();
    }
    this.hideNotification();
  }

  public isProjectFormInvalid(): boolean {
    return this.isCreateMode ? this.isProjectCreateFormInvalid() : this.isProjectSettingsFormInvalid();
  }

  public isProjectCreateFormInvalid(): boolean {
    return (
      !this.shipyardFile ||
      this.projectNameForm.invalid ||
      (this.resourceServiceEnabled ? !this.gitDataExtended : !this.gitData.gitFormValid) ||
      this.isCreatingProjectInProgress
    );
  }

  public isProjectSettingsFormInvalid(): boolean {
    return !this.gitData.gitFormValid || this.isGitUpstreamInProgress;
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

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
