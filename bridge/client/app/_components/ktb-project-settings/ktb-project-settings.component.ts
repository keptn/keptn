import { Component, HostListener, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { DtToast } from '@dynatrace/barista-components/toast';
import { combineLatest, Observable, of, Subject } from 'rxjs';
import { filter, map, startWith, switchMap, takeUntil } from 'rxjs/operators';
import { IGitDataExtended } from 'shared/interfaces/project';
import { IClientFeatureFlags } from '../../../../shared/interfaces/feature-flags';
import { PendingChangesComponent } from '../../_guards/pending-changes.guard';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { IMetadata } from '../../_interfaces/metadata';
import { KeptnInfo } from '../../_models/keptn-info';
import { NotificationType } from '../../_models/notification';
import { Project } from '../../_models/project';
import { ServerErrors } from '../../_models/server-error';
import { DataService } from '../../_services/data.service';
import { EventService } from '../../_services/event.service';
import { FeatureFlagsService } from '../../_services/feature-flags.service';
import { NotificationsService } from '../../_services/notifications.service';
import { AppUtils } from '../../_utils/app.utils';
import { FormUtils } from '../../_utils/form.utils';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message/ktb-project-create-message.component';
import { IGitDataExtendedWithNoUpstream } from './ktb-project-settings-git-extended/ktb-project-settings-git-extended.component';

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
  public gitDataExtended?: IGitDataExtendedWithNoUpstream;
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl,
  });

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
    this.dataService.projects
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
      this.gitInputDataExtendedDefault = project.gitCredentials;
      this.gitInputDataExtended = AppUtils.copyObject(project.gitCredentials); // there should not be a reference. Could
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
      .pipe(switchMap(() => this.dataService.loadProjects()))
      .subscribe(
        () => {
          this.projectName = projectName;
          this.isCreatingProjectInProgress = false;
          this.isProjectFormTouched = false;

          this.router.navigate(['/', 'project', this.projectName, 'settings', 'project'], {
            queryParams: { created: true },
          });
        },
        (err) => {
          const errorMessage = err.error || 'please, check the logs of resource-service';
          this.notificationsService.addNotification(
            NotificationType.ERROR,
            `The project could not be created: ${errorMessage}.`
          );
          this.isCreatingProjectInProgress = false;
        }
      );
  }

  private deleteProject(projectName: string): void {
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
    this.gitInputDataExtended = AppUtils.copyObject(this.gitInputDataExtendedDefault);
    this.pendingChangesSubject.next(true);
    this.hideNotification();
  }

  public saveAll(): void {
    if (this.isCreateMode) {
      this.createProject();
    } else {
      this.updateGitUpstream();
    }
    this.hideNotification();
  }

  public updateGitUpstream(): void {
    if (!this.gitDataExtended || !this.projectName) {
      return;
    }
    this.isGitUpstreamInProgress = true;
    this.dataService.updateGitUpstream(this.projectName, this.gitDataExtended).subscribe(
      () => {
        this.isGitUpstreamInProgress = false;
        this.notificationsService.addNotification(
          NotificationType.SUCCESS,
          'The Git upstream was changed successfully.'
        );
        this.isProjectFormTouched = false;
      },
      () => {
        this.isGitUpstreamInProgress = false;
      }
    );
  }

  public isProjectFormInvalid(): boolean {
    return this.isCreateMode ? this.isProjectCreateFormInvalid() : this.isProjectSettingsFormInvalid();
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

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
