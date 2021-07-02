import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import {DataService} from "../../_services/data.service";
import {ActivatedRoute} from "@angular/router";
import {filter, map, switchMap, takeUntil} from "rxjs/operators";
import {DtToast} from "@dynatrace/barista-components/toast";
import {GitData} from "../../_components/ktb-project-settings-git/ktb-project-settings-git.component";

@Component({
  selector: 'ktb-settings-view',
  templateUrl: './ktb-settings-view.component.html',
  styleUrls: ['./ktb-settings-view.component.scss']
})
export class KtbSettingsViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private projectName: string;

  public isGitUpstreamInProgress: boolean;
  public gitData: GitData = {};

  constructor(private route: ActivatedRoute, private dataService: DataService, private toast: DtToast) {
  }

  ngOnInit(): void {
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

  setGitUpstream(gitData: GitData) {
    this.isGitUpstreamInProgress = true;
    this.dataService.setGitUpstreamUrl(this.projectName, gitData.remoteURI, gitData.gitUser, gitData.gitToken).subscribe(success => {
      this.isGitUpstreamInProgress = false;
      if (success) {
        this.toast.create('Git Upstream URL set successfully');
      } else {
        this.toast.create('Git Upstream URL could not be set');
      }
    });
  }

}
