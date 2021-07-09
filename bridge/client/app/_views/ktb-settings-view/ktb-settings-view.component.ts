import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import {DataService} from "../../_services/data.service";
import {ActivatedRoute} from "@angular/router";
import {filter, map, switchMap, takeUntil} from "rxjs/operators";
import {DtToast} from "@dynatrace/barista-components/toast";
import {GitData} from "../../_components/ktb-project-settings-git/ktb-project-settings-git.component";
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {FormUtils} from "../../_utils/form.utils";

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss']
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private projectName: string;

  public isCreateMode: boolean;
  public isGitUpstreamInProgress: boolean;
  public isCreatingProjectInProgress: boolean;
  public shipyardFile: File;
  public gitData: GitData = {};

  public projectNameControl = new FormControl('');
  public projectNameForm = new FormGroup({
    projectName: this.projectNameControl
  });

  constructor(private route: ActivatedRoute, private dataService: DataService, private toast: DtToast) {
  }

  ngOnInit(): void {
    this.route.data.subscribe((data) => {
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
      this.projectNameControl.setValidators([Validators.required, FormUtils.projectNameExistsValidator(projectNames)]);
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
    this.dataService.createProject(
      this.projectNameControl.value, this.shipyardFile, this.gitData.remoteURI || null, this.gitData.gitToken || null, this.gitData.gitUser || null
    ).subscribe(
      () => {
        this.toast.create('Project successfully created');
        this.isCreatingProjectInProgress = false;
        // TODO Success handling -> navigate to project settings without create mode
      },
      (err) => {
        console.log(err);
        this.toast.create('Project could not be created');
        this.isCreatingProjectInProgress = false;
      });
  }
}
