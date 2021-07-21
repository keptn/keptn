import {Component, OnDestroy, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {Subject} from 'rxjs';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map, switchMap, take, takeUntil} from 'rxjs/operators';
import {DtToast} from '@dynatrace/barista-components/toast';
import {MatDialog} from '@angular/material/dialog';
import {GitData} from '../../_components/ktb-project-settings-git/ktb-project-settings-git.component';
import {FormUtils} from '../../_utils/form.utils';
import {NotificationType, TemplateRenderedNotifications} from '../../_models/notification';
import {NotificationsService} from '../../_services/notifications.service';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import {EventService} from '../../_services/event.service';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss'],
  providers: [NotificationsService]
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public projectName: string;

  @ViewChild('deleteProjectDialog')
  private deleteProjectDialog: TemplateRef<MatDialog>;

  public projectDeletionData: DeleteData;

  public isCreateMode: boolean;
  public isGitUpstreamInProgress: boolean;
  public isCreatingProjectInProgress: boolean;
  public shipyardFile: File;
  public gitData: GitData = {};

  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl
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
      takeUntil(this.unsubscribe$)
    ).subscribe((data) => {
      if (data) {
        this.isCreateMode = data.isCreateMode;
      }
    });

    this.dataService.projects
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((projects) => !!projects),
        map((projects) => projects ? projects.map(project => project.projectName) : null)
      ).subscribe((projectNames) => {
      if (this.isCreateMode && projectNames.includes(this.projectName)) {
        this.router.navigate(['/', 'project', this.projectName, 'settings'], {queryParams: {created: true}});
      }
      this.projectNameControl.setValidators([
        Validators.required,
        FormUtils.projectNameExistsValidator(projectNames),
        Validators.pattern('[a-z]([a-z]|[0-9]|-)*')
      ]);
    });

    this.route.params.pipe(
      map(params => params.projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
      filter(project => !!project)
    ).subscribe(project => {
      this.projectName = project.projectName;
      this.gitData.remoteURI = project.gitRemoteURI;
      this.gitData.gitUser = project.gitUser;

      this.projectDeletionData = {
        type: DeleteType.PROJECT,
        name: this.projectName
      };
    });

    this.route.queryParams.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe((queryParams) => {
      if (queryParams.created) {
        this.notificationsService.addNotification(NotificationType.Success, TemplateRenderedNotifications.CREATE_PROJECT, null, true, {
          projectName: this.projectName,
          routerLink: `/project/${this.projectName}/service`
        });
      }
    });

    this.eventService.deletionTriggeredEvent.pipe(
      takeUntil(this.unsubscribe$)
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
  }

  public setGitUpstream(): void {
    this.isGitUpstreamInProgress = true;
    this.dataService.setGitUpstreamUrl(this.projectName, this.gitData.remoteURI, this.gitData.gitUser, this.gitData.gitToken)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(success => {
        this.isGitUpstreamInProgress = false;
        if (success) {
          this.toast.create('Git Upstream URL set successfully');
        } else {
          this.toast.create('Git Upstream URL could not be set');
        }
      });
  }

  public isGitFormValid(): boolean {
    if (!this.gitData.remoteURI && !this.gitData.gitUser && !this.gitData.gitToken) {
      return true;
    }
    return this.gitData.remoteURI.length > 0 && this.gitData.gitUser.length > 0 && this.gitData.gitToken.length > 0;
  }

  public createProject(): void {
    this.isCreatingProjectInProgress = true;

    FormUtils.readFileContent(this.shipyardFile).then(fileContent => {
      const shipyardBase64 = btoa(fileContent);
      const projectName = this.projectNameControl.value;
      this.dataService.createProject(
        projectName, shipyardBase64, this.gitData.remoteURI || null, this.gitData.gitToken || null, this.gitData.gitUser || null
      ).subscribe(() => {
          this.projectName = projectName;
          this.dataService.loadProjects();
          this.isCreatingProjectInProgress = false;
        },
        () => {
          this.notificationsService.addNotification(NotificationType.Error, 'The project could not be created.', 5000);
          this.isCreatingProjectInProgress = false;
        });
    });
  }

  public deleteProject(projectName) {
    this.eventService.deletionProgressEvent.next({error: null, isInProgress: true, result: null});

    this.dataService.projects
      .pipe(take(1))
      .subscribe(() => {
        this.router.navigate(['/', 'dashboard']);
      });

    this.dataService.deleteProject(projectName)
      .pipe(take(1))
      .subscribe(() => {
        this.dataService.loadProjects();
        this.eventService.deletionProgressEvent.next({error: null, isInProgress: false, result: DeleteResult.SUCCESS});
      }, (err) => {
        const deletionError = 'Project could not be deleted: ' + err.message;
        this.eventService.deletionProgressEvent.next({error: deletionError, isInProgress: false, result: DeleteResult.ERROR});
      });
  }
}
