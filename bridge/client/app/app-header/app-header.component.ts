import semver from 'semver';

import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Observable, Subject} from 'rxjs';
import {filter, takeUntil} from 'rxjs/operators';

import {Project} from '../_models/project';
import {DataService} from '../_services/data.service';
import {NotificationsService} from '../_services/notifications.service';
import {NotificationType} from '../_models/notification';

@Component({
  selector: 'app-header',
  templateUrl: './app-header.component.html',
  styleUrls: ['./app-header.component.scss']
})
export class AppHeaderComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public projects: Observable<Project[]>;
  public project: Observable<Project>;

  public keptnInfo: any;
  public versionCheckDialogState: string | null;
  public versionCheckReference = '/reference/version_check/';

  constructor(private route: ActivatedRoute, private dataService: DataService, private notificationsService: NotificationsService) { }

  ngOnInit() {
    this.projects = this.dataService.projects;

    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.project = this.dataService.getProject(params.projectName);
      });

    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if (keptnInfo.versionCheckEnabled === undefined) {
          this.showVersionCheckInfoDialog();
        } else if(keptnInfo.versionCheckEnabled) {
          keptnInfo.keptnVersionInvalid = !this.doVersionCheck(keptnInfo.keptnVersion, keptnInfo.availableVersions.cli.stable);
        }
      });
  }

  doVersionCheck(currentVersion, stableVersions:String[]): boolean {
    if(!semver.valid(currentVersion))
      return false;

    const latestVersion = stableVersions[stableVersions.length-1];
    const newerVersions = [];
    let genMessage = (versions, major, minor) => `New Keptn ${versions} available. For details how to upgrade visit <a href="https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/" target="_blank">https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/</a>`;

    stableVersions.forEach(stableVersion => {
      if(semver.lt(currentVersion, stableVersion)) {
        newerVersions.push(stableVersion);
      }
    });

    if(newerVersions.length > 0) {
      let versionsString = newerVersions[0];
      if (newerVersions.length > 1) {
        versionsString = '(' + newerVersions.join(', ') + ')';
      }

      this.notificationsService.addNotification(NotificationType.Info, genMessage(versionsString, semver.major(latestVersion), semver.minor(latestVersion)));
    }

    return true;
  }

  showVersionCheckInfoDialog() {
    if (this.keptnInfo.bridgeInfo.enableVersionCheckFeature)
      this.versionCheckDialogState = 'info';
  }

  acceptVersionCheck(accepted: boolean): void {
    this.dataService.setVersionCheck(accepted);
    if (accepted)
      this.versionCheckDialogState = 'success';

    setTimeout(() => {
      this.versionCheckDialogState = null;
    }, accepted ? 2000 : 0);
  }

  versionCheckClicked(event) {
    this.dataService.setVersionCheck(event.checked);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
