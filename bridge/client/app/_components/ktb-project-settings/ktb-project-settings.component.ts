import { Component, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { Subject } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { GitData, KtbProjectSettingsGitComponent } from '../ktb-project-settings-git/ktb-project-settings-git.component';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { DtToast } from '@dynatrace/barista-components/toast';
import { NotificationsService } from '../../_services/notifications.service';
import { EventService } from '../../_services/event.service';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { Project } from '../../_models/project';
import { FormUtils } from '../../_utils/form.utils';
import { NotificationType, TemplateRenderedNotifications } from '../../_models/notification';

@Component({
  selector: 'ktb-project-settings',
  templateUrl: './ktb-project-settings.component.html',
  styleUrls: ['./ktb-project-settings.component.scss'],
})
export class KtbProjectSettingsComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  @ViewChild('deleteProjectDialog')
  private deleteProjectDialog?: TemplateRef<MatDialog>;

  @ViewChild(KtbProjectSettingsGitComponent)
  private gitSettingsSection?: KtbProjectSettingsGitComponent;

  public unsavedDialogState: string | null = null;
  public projectName?: string;
  public projectDeletionData?: DeleteData;
  public isCreateMode = false;
  public isGitUpstreamInProgress = false;
  public isCreatingProjectInProgress = false;
  public shipyardFile?: File;
  public gitData: GitData = {
    gitFormValid: true,
  };
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl,
  });

  constructor(private route: ActivatedRoute,
              private dataService: DataService,
              private toast: DtToast,
              private router: Router,
              private notificationsService: NotificationsService,
              private eventService: EventService) {
  }

  ngOnInit(): void {
    this.route.data.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe((data) => {
      if (data) {
        this.isCreateMode = data.isCreateMode;
      }
    });

    this.dataService.projects
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((projects: Project[] | undefined): projects is Project[] => !!projects),
        map((projects: Project[]) => projects.map(project => project.projectName)),
      ).subscribe((projectNames) => {
      if (this.isCreateMode && this.projectName && projectNames.includes(this.projectName)) {
        this.router.navigate(['/', 'project', this.projectName, 'settings'], {queryParams: {created: true}});
      }
      this.projectNameControl.setValidators([
        Validators.required,
        FormUtils.projectNameExistsValidator(projectNames),
        Validators.pattern('[a-z]([a-z]|[0-9]|-)*'),
      ]);
    });

    this.route.params.pipe(
      map(params => params.projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
      filter((project: Project | undefined): project is Project => !!project),
    ).subscribe(project => {
      if (project.projectName !== this.projectName) {
        this.unsavedDialogState = null;
        this.projectName = project.projectName;

        this.gitData = {
          remoteURI: project.gitRemoteURI,
          gitUser: project.gitUser,
        };

        this.projectDeletionData = {
          type: DeleteType.PROJECT,
          name: this.projectName,
        };
      }
    });

    this.route.queryParams.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe((queryParams) => {
      if (queryParams.created) {
        this.unsavedDialogState = null;
        this.notificationsService.addNotification(NotificationType.Success, TemplateRenderedNotifications.CREATE_PROJECT, undefined, true, {
          projectName: this.projectName,
          routerLink: `/project/${this.projectName}/service`,
        });
        // Remove query param for not showing notification on reload
        this.router.navigate(['/', 'project', this.projectName, 'settings']);
      }
    });

    this.eventService.deletionTriggeredEvent.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(data => {
      if (data.type === DeleteType.PROJECT) {
        this.deleteProject(data.name);
      }
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public updateGitData(gitData: GitData): void {
    this.gitData.remoteURI = gitData.remoteURI;
    this.gitData.gitUser = gitData.gitUser;
    this.gitData.gitToken = gitData.gitToken;
    this.gitData.gitFormValid = gitData.gitFormValid;
    if (gitData.gitFormValid) {
      this.unsavedDialogState = 'unsaved';
    } else {
      this.unsavedDialogState = null;
    }
  }

  public updateShipyardFile(shipyardFile: File | undefined): void {
    this.shipyardFile = shipyardFile;
  }

  public setGitUpstream(): void {
    if (this.projectName && this.gitData.remoteURI && this.gitData.gitUser && this.gitData.gitToken) {
      this.isGitUpstreamInProgress = true;
      this.unsavedDialogState = null;
      this.dataService.setGitUpstreamUrl(this.projectName, this.gitData.remoteURI, this.gitData.gitUser, this.gitData.gitToken)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(() => {
          this.isGitUpstreamInProgress = false;
          this.gitData.gitToken = '';
          this.gitData = {...this.gitData};
          this.notificationsService.addNotification(NotificationType.Success, 'The Git upstream was changed successfully.', 5000);
        }, (err) => {
          console.log(err);
          this.isGitUpstreamInProgress = false;
          this.notificationsService.addNotification(NotificationType.Error, `<div class="long-note align-left p-3">The Git upstream could not be changed:<br/><span class="small">${err.error}</span></div>`);
        });
    }
  }

  public createProject(): void {
    if (this.shipyardFile) {
      this.isCreatingProjectInProgress = true;
      FormUtils.readFileContent(this.shipyardFile).then(fileContent => {
        if (fileContent) {
          const shipyardBase64 = btoa(fileContent);
          const projectName = this.projectNameControl.value;
          this.dataService.createProject(
            projectName, shipyardBase64, this.gitData.remoteURI, this.gitData.gitToken, this.gitData.gitUser,
          ).subscribe(() => {
              this.projectName = projectName;
              this.dataService.loadProjects();
              this.isCreatingProjectInProgress = false;
            },
            () => {
              this.notificationsService.addNotification(NotificationType.Error, 'The project could not be created.', 5000);
              this.isCreatingProjectInProgress = false;
            });
        }
      });
    }
  }

  public deleteProject(projectName: string): void {
    this.eventService.deletionProgressEvent.next({isInProgress: true});

    this.dataService.projectExists(projectName)
      .pipe(
        takeUntil(this.unsubscribe$),
        filter(status => status === false),
      ).subscribe(() => {
      this.router.navigate(['/', 'dashboard']);
    });

    this.dataService.deleteProject(projectName)
      .pipe(
        takeUntil(this.unsubscribe$),
      ).subscribe(() => {
      this.dataService.loadProjects();
      this.eventService.deletionProgressEvent.next({isInProgress: false, result: DeleteResult.SUCCESS});
    }, (err) => {
      const deletionError = 'Project could not be deleted: ' + err.message;
      this.eventService.deletionProgressEvent.next({
        error: deletionError,
        isInProgress: false,
        result: DeleteResult.ERROR,
      });
    });
  }

  public reset(): void {
    this.gitSettingsSection?.reset();
    this.unsavedDialogState = null;
  }

  public saveAll(): void {
    this.setGitUpstream();
    this.unsavedDialogState = null;
  }
}
