import semver from 'semver';

import {Component, OnDestroy, OnInit} from '@angular/core';
import {Router, RoutesRecognized} from "@angular/router";
import {Observable, Subscription} from "rxjs";
import {filter, map} from "rxjs/operators";

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

  private routeSub: Subscription = Subscription.EMPTY;
  private versionInfoSub: Subscription = Subscription.EMPTY;

  public projects: Observable<Project[]>;
  public project: Observable<Project>;

  public versionInfo: any;
  public versionCheckDialogState: string | null;
  public versionCheckReference = "https://keptn.sh/docs/0.7.x/reference/bridge/version_check";

  constructor(private router: Router, private dataService: DataService, private apiService: ApiService, private notificationsService: NotificationsService) { }

  ngOnInit() {
    this.projects = this.dataService.projects;

    this.routeSub = this.router.events.subscribe(event => {
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

    this.versionInfoSub = this.dataService.versionInfo.subscribe(versionInfo => {
      this.versionInfo = versionInfo;
      if(versionInfo.versionCheckEnabled === null) {
        this.showVersionCheckInfoDialog();
      } else if(versionInfo.versionCheckEnabled) {
        if(semver.valid(versionInfo.keptnVersion)) {
          if(versionInfo.availableVersions.cli)
            this.doVersionCheck(versionInfo.keptnVersion, versionInfo.availableVersions.cli.stable, versionInfo.availableVersions.cli.prerelease, "Keptn");
        } else {
          versionInfo.keptnVersionInvalid = true;
        }
        if(semver.valid(versionInfo.bridgeVersion)) {
          if(versionInfo.availableVersions.bridge)
            this.doVersionCheck(versionInfo.bridgeVersion, versionInfo.availableVersions.bridge.stable, versionInfo.availableVersions.bridge.prerelease, "Keptn Bridge");
        } else {
          versionInfo.bridgeVersionInvalid = true;
        }
      }
    });
  }

  doVersionCheck(currentVersion, stableVersions, prereleaseVersions, type) {
    stableVersions.forEach(stableVersion => {
      if(semver.lt(currentVersion, stableVersion)) {
        let genMessage;
        switch(semver.diff(currentVersion, stableVersion)) {
          case 'patch':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a patch with bug fixes for your current version. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.0/operate/upgrade/`;
            break;
          case 'minor':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a minor update with backwards compatible improvements. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.0/operate/upgrade/`;
            break;
          case 'major':
            genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a major update and it might contain incompatible changes. For details how to upgrade visit https://keptn.sh/docs/${major}.${minor}.0/operate/upgrade/`;
            break;
        }

        let major = semver.major(stableVersion);
        let minor = semver.minor(stableVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(stableVersion, type, major, minor));
      }
    });
    prereleaseVersions.forEach(prereleaseVersion => {
      if(semver.lt(currentVersion, prereleaseVersion)) {
        let genMessage = (version, type, major, minor) => `New ${type} ${version} available. This is a pre-release version that might contain experimental features. For details how to upgrade visit: https://keptn.sh/docs/${major}.${minor}.0/operate/upgrade/`;
        let major = semver.major(prereleaseVersion);
        let minor = semver.minor(prereleaseVersion);
        this.notificationsService.addNotification(NotificationType.Info, genMessage(prereleaseVersion, type, major, minor));
      }
    });
  }

  showVersionCheckInfoDialog() {
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
    this.routeSub.unsubscribe();
    this.versionInfoSub.unsubscribe();
  }

}
