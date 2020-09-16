import semver from 'semver';

import {Component, OnDestroy, OnInit} from '@angular/core';
import {Router, RoutesRecognized} from "@angular/router";
import {Observable, Subject, Subscription} from "rxjs";
import {filter, map, takeUntil} from "rxjs/operators";

import {Project} from "../_models/project";
import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";
import {NotificationsService} from "../_services/notifications.service";
import {NotificationType} from "../_models/notification";

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
  public versionCheckReference = "https://keptn.sh/docs/0.7.x/reference/version_check";

  constructor(private router: Router, private dataService: DataService, private apiService: ApiService, private notificationsService: NotificationsService) { }

  ngOnInit() {
    this.projects = this.dataService.projects;

    this.router.events
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(event => {
        if(event instanceof RoutesRecognized) {
          let projectName = event.state.root.children[0].params['projectName'];
          this.project = this.dataService.projects.pipe(
            filter(projects => !!projects),
            map(projects => projects.find(p => {
              return p.projectName === projectName;
            }))
          );
        }
      });

    this.dataService.keptnInfo
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if(keptnInfo.versionCheckEnabled === null) {
          this.showVersionCheckInfoDialog();
        } else if(keptnInfo.versionCheckEnabled) {
          keptnInfo.keptnVersionInvalid = !this.doVersionCheck(keptnInfo.keptnVersion, keptnInfo.availableVersions.cli.stable, keptnInfo.availableVersions.cli.prerelease, "Keptn");
          keptnInfo.bridgeVersionInvalid = !this.doVersionCheck(keptnInfo.bridgeInfo.bridgeVersion, keptnInfo.availableVersions.bridge.stable, keptnInfo.availableVersions.bridge.prerelease, "Keptn Bridge");;
        }
      });
  }

  doVersionCheck(currentVersion, stableVersions, prereleaseVersions, type): boolean {
    if(!semver.valid(currentVersion))
      return false;

    stableVersions.forEach(stableVersion => {
      if(semver.lt(currentVersion, stableVersion)) {
        let genMessage;
        switch(semver.diff(currentVersion, stableVersion)) {
          case 'patch':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a patch version with bug fixes and minor improvements. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
          case 'minor':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a minor update with backwards compatible improvements. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
          case 'major':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a major update and it might contain incompatible changes. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
            break;
        }

        let major = semver.major(stableVersion);
        let minor = semver.minor(stableVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(stableVersion, type, major, minor));
      }
    });
    prereleaseVersions.forEach(prereleaseVersion => {
      if(semver.lt(currentVersion, prereleaseVersion)) {
        let genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a pre-release version with experimental features. For details how to upgrade visit: https://keptn.sh/docs/${major}.${minor}.x/operate/upgrade/`;
        let major = semver.major(prereleaseVersion);
        let minor = semver.minor(prereleaseVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(prereleaseVersion, type, major, minor));
      }
    });

    return true;
  }

  showVersionCheckInfoDialog() {
    if(this.keptnInfo.bridgeInfo.enableVersionCheckFeature)
      this.versionCheckDialogState = 'info';
  }

  acceptVersionCheck(accepted: boolean): void {
    this.dataService.setVersionCheck(accepted);
    if(accepted)
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
