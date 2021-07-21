import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from 'rxjs';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map, switchMap, takeUntil} from 'rxjs/operators';
import {DtToast} from '@dynatrace/barista-components/toast';
import {GitData} from '../../_components/ktb-project-settings-git/ktb-project-settings-git.component';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {FormUtils} from '../../_utils/form.utils';
import {NotificationType, TemplateRenderedNotifications} from '../../_models/notification';
import {NotificationsService} from '../../_services/notifications.service';
import { Project } from '../../_models/project';

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss'],
  providers: [NotificationsService]
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public projectName?: string;
  public isCreateMode = false;
  public isGitUpstreamInProgress = false;
  public isCreatingProjectInProgress = false;
  public shipyardFile?: File;
  public gitData: GitData = {};
  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl
  });

  constructor(private route: ActivatedRoute, private router: Router, private dataService: DataService, private toast: DtToast, private notificationsService: NotificationsService) {
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
        filter((projects: Project[] | undefined): projects is Project[] => !!projects),
        map((projects: Project[]) => projects.map(project => project.projectName))
      ).subscribe((projectNames) => {
        if (this.isCreateMode && this.projectName && projectNames.includes(this.projectName)) {
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
      filter((project: Project | undefined): project is Project => !!project)
    ).subscribe(project => {
      this.projectName = project.projectName;
      this.gitData.remoteURI = project.gitRemoteURI;
      this.gitData.gitUser = project.gitUser;
    });

    this.route.queryParams.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe((queryParams) => {
      if (queryParams.created) {
        this.notificationsService.addNotification(NotificationType.Success, TemplateRenderedNotifications.CREATE_PROJECT, undefined, true, {projectName: this.projectName, routerLink: `/project/${this.projectName}/service`});
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
    if (this.projectName && this.gitData.remoteURI && this.gitData.gitUser && this.gitData.gitToken) {
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
  }

  public isGitFormValid(): boolean {
    if (!this.gitData.remoteURI && !this.gitData.gitUser && !this.gitData.gitToken) {
      return true;
    }
    return !!(this.gitData.remoteURI?.length && this.gitData.gitUser?.length && this.gitData.gitToken?.length);
  }

  public createProject(): void {
    if (this.shipyardFile) {
      this.isCreatingProjectInProgress = true;
      FormUtils.readFileContent(this.shipyardFile).then(fileContent => {
        if (fileContent) {
          const shipyardBase64 = btoa(fileContent);
          const projectName = this.projectNameControl.value;
          this.dataService.createProject(
            projectName, shipyardBase64, this.gitData.remoteURI, this.gitData.gitToken, this.gitData.gitUser
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
}
