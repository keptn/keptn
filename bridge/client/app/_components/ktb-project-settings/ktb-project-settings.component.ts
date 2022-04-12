import { Component, HostListener, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import {
  GitData,
  KtbProjectSettingsGitComponent,
} from '../ktb-project-settings-git/ktb-project-settings-git.component';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { DtToast } from '@dynatrace/barista-components/toast';
import { NotificationsService } from '../../_services/notifications.service';
import { EventService } from '../../_services/event.service';
import { filter, map, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { FormUtils } from '../../_utils/form.utils';
import { KtbProjectCreateMessageComponent } from '../_status-messages/ktb-project-create-message/ktb-project-create-message.component';
import { NotificationType } from '../../_models/notification';
import { PendingChangesComponent } from '../../_guards/pending-changes.guard';

type DialogState = null | 'unsaved';

@Component({
  selector: 'ktb-project-settings',
  templateUrl: './ktb-project-settings.component.html',
  styleUrls: ['./ktb-project-settings.component.scss'],
})
export class KtbProjectSettingsComponent implements OnInit, OnDestroy, PendingChangesComponent {
  private readonly unsubscribe$ = new Subject<void>();

  @ViewChild('deleteProjectDialog')
  private deleteProjectDialog?: TemplateRef<MatDialog>;

  @ViewChild(KtbProjectSettingsGitComponent)
  private gitSettingsSection?: KtbProjectSettingsGitComponent;

  public projectName?: string;
  public projectDeletionData?: DeleteData;
  public isProjectLoading: boolean | undefined;
  public isCreateMode = false;
  public isGitUpstreamInProgress = false;
  public isCreatingProjectInProgress = false;
  private pendingChangesSubject = new Subject<boolean>();
  public isProjectFormTouched = false;
  public shipyardFile?: File;
  public gitData: GitData = {
    gitFormValid: true,
  };
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl,
  });

  public message = 'You have pending changes. Make sure to save your data before you continue.';
  public unsavedDialogState: DialogState = null;

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private toast: DtToast,
    private router: Router,
    private notificationsService: NotificationsService,
    private eventService: EventService
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      if (!params.projectName) {
        this.isCreateMode = true;
        this.loadProjectsAndSetValidator();
      }

      if (params.projectName) {
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

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
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
      });
  }

  private loadProject(projectName: string): void {
    this.dataService.loadPlainProject(projectName).subscribe((project) => {
      this.gitData = {
        remoteURI: project.gitRemoteURI,
        gitUser: project.gitUser,
      };

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

  public updateGitData(gitData: GitData): void {
    this.gitData.remoteURI = gitData.remoteURI;
    this.gitData.gitUser = gitData.gitUser;
    this.gitData.gitToken = gitData.gitToken;
    this.gitData.gitFormValid = gitData.gitFormValid;
    this.projectFormTouched();
  }

  public updateShipyardFile(shipyardFile: File | undefined): void {
    this.shipyardFile = shipyardFile;
    this.projectFormTouched();
  }

  public setGitUpstream(): void {
    if (this.projectName && this.gitData.remoteURI && this.gitData.gitUser && this.gitData.gitToken) {
      this.isGitUpstreamInProgress = true;
      this.hideNotification();
      this.dataService
        .setGitUpstreamUrl(this.projectName, this.gitData.remoteURI, this.gitData.gitUser, this.gitData.gitToken)
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
              `The Git upstream could not be changed: ${err.error}`
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
        this.dataService
          .createProject(
            projectName,
            shipyardBase64,
            this.gitData.remoteURI,
            this.gitData.gitToken,
            this.gitData.gitUser
          )
          .subscribe(
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
              this.notificationsService.addNotification(
                NotificationType.ERROR,
                `The project could not be created: ${err.error || 'please, check the logs of configuration-service'}.`
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
    this.gitSettingsSection?.reset();
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
      !this.gitData.gitFormValid ||
      this.isCreatingProjectInProgress
    );
  }

  public isProjectSettingsFormInvalid(): boolean {
    return !this.gitData.gitFormValid || this.isGitUpstreamInProgress;
  }

  projectFormTouched(): void {
    this.isProjectFormTouched = true;
  }

  // @HostListener allows us to also guard against browser refresh, close, etc.
  @HostListener('window:beforeunload', ['$event'])
  canDeactivate($event?: BeforeUnloadEvent): Observable<boolean> {
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

  showNotification(): void {
    this.unsavedDialogState = 'unsaved';

    document.querySelector('div[aria-label="Dialog for notifying about unsaved data"]')?.classList.add('shake');
    setTimeout(() => {
      document.querySelector('div[aria-label="Dialog for notifying about unsaved data"]')?.classList.remove('shake');
    }, 500);
  }

  hideNotification(): void {
    this.unsavedDialogState = null;
  }
}
