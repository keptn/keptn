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
        } else if (keptnInfo.versionCheckEnabled) {
          keptnInfo.keptnVersionInvalid = !this.doVersionCheck(keptnInfo.keptnVersion, keptnInfo.availableVersions.cli.stable, keptnInfo.availableVersions.cli.prerelease, 'Keptn');
          keptnInfo.bridgeVersionInvalid = !this.doVersionCheck(keptnInfo.bridgeInfo.bridgeVersion, keptnInfo.availableVersions.bridge.stable, keptnInfo.availableVersions.bridge.prerelease, 'Keptn Bridge');;
        }
      });
  }

  doVersionCheck(currentVersion, stableVersions, prereleaseVersions, type): boolean {
    if(!semver.valid(currentVersion))
      return false;

    stableVersions.forEach(stableVersion => {
      if (semver.lt(currentVersion, stableVersion)) {
        let genMessage;
        switch (semver.diff(currentVersion, stableVersion)) {
          case 'patch':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a patch version with bug fixes and minor improvements. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
          case 'minor':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a minor update with backwards compatible improvements. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
          case 'major':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a major update and it might contain incompatible changes. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
          default:
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. It might contain incompatible changes. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
        }

        const major = semver.major(stableVersion);
        const minor = semver.minor(stableVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(stableVersion, type, major, minor));
      }
    });
    prereleaseVersions.forEach(prereleaseVersion => {
      if (semver.lt(currentVersion, prereleaseVersion)) {
        const genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a pre-release version with experimental features. For details how to upgrade visit: https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
        const major = semver.major(prereleaseVersion);
        const minor = semver.minor(prereleaseVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(prereleaseVersion, type, major, minor));
      }
    });

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
