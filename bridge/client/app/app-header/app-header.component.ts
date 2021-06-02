import semver from 'semver';

import {Component, OnDestroy, OnInit} from '@angular/core';
import {NavigationEnd, Router, RoutesRecognized} from '@angular/router';
import {Observable, Subject} from 'rxjs';
import {filter, map, takeUntil} from 'rxjs/operators';

import {Project} from '../_models/project';
import {DataService} from '../_services/data.service';
import {NotificationsService} from '../_services/notifications.service';
import {NotificationType} from '../_models/notification';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-header',
  templateUrl: './app-header.component.html',
  styleUrls: ['./app-header.component.scss']
})
export class AppHeaderComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public projects: Observable<Project[]>;
  public project: Observable<Project>;
  public projectBoardView = '';
  public appTitle = environment.appTitle;

  public keptnInfo: any;
  public versionCheckDialogState: string | null;
  public versionCheckReference = '/reference/version_check/';

  constructor(private router: Router, private dataService: DataService, private notificationsService: NotificationsService) {}

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
        } else if (event instanceof NavigationEnd) {
          // catch url change and update projectBoardView for the project picker
          const pieces = event.url.split('/');
          if (pieces.length > 3 && pieces[1] === 'project') {
            this.projectBoardView = pieces[3];
          } else {
            this.projectBoardView = ''; // environment screen
          }
        }
      });

    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if (keptnInfo.versionCheckEnabled === undefined) {
          this.showVersionCheckInfoDialog();
        } else if (keptnInfo.bridgeInfo.enableVersionCheckFeature && keptnInfo.versionCheckEnabled) {
          keptnInfo.keptnVersionInvalid = !this.doVersionCheck(
            keptnInfo.bridgeInfo.bridgeVersion,
            keptnInfo.keptnVersion,
            keptnInfo.availableVersions.bridge,
            keptnInfo.availableVersions.cli);
        }
      });
  }

  // Returns a string array that allows routing to the project board view
  getRouterLink(projectName: string): string[] {
    if (this.projectBoardView === '') {
      // unfortunately it is not possible to route directly to the environment screen (default screen)
      return ['/project', projectName];
    }

    return ['/project', projectName, this.projectBoardView];
  }

  doVersionCheck(bridgeVersion, cliVersion, availableBridgeVersions, availableCliVersions): boolean {
    if (!semver.valid(bridgeVersion) || !semver.valid(cliVersion))
      return false;

    const latestVersion = availableCliVersions.stable[availableCliVersions.stable.length - 1];
    const bridgeVersionsString = this.buildVersionString(this.getNewerVersions(availableBridgeVersions, bridgeVersion));
    const cliVersionsString = this.buildVersionString(this.getNewerVersions(availableCliVersions, cliVersion));

    if (bridgeVersionsString || cliVersionsString) {
      let versionMessage = `New ${cliVersionsString ? ' Keptn CLI ' + cliVersionsString : ''} ${cliVersionsString && bridgeVersionsString ? 'and' : ''}
                            ${bridgeVersionsString ? ' Keptn Bridge ' + bridgeVersionsString : ''} available. For details how to upgrade visit
                            <a href="https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(latestVersion)}.x/operate/upgrade/" target="_blank">
                            https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(latestVersion)}.x/operate/upgrade/</a>`;
      this.notificationsService.addNotification(NotificationType.Info, versionMessage);
    }

    return true;
  }

  private buildVersionString(versions) {
    if (versions.stable.length > 0) {
      return versions.stable.join(', ');
    } else if (versions.prerelease.length > 0) {
      return versions.prerelease.join(', ');
    }

    return null;
  }

  private getNewerVersions(availableVersions, currentVersion) {
    const newerVersions = {
      stable: [],
      prerelease: []
    }

    newerVersions.stable = availableVersions.stable.filter((stableVersion) => {
      if (semver.lt(currentVersion, stableVersion)) {
        return stableVersion;
      }
    });

    // It is only necessary to check prerelease versions when no stable update is available
    if (newerVersions.stable.length === 0) {
      newerVersions.prerelease = availableVersions.prerelease.filter((prereleaseVersion) => {
        if (semver.lt(currentVersion, prereleaseVersion)) {
          return prereleaseVersion;
        }
      });
    }

    return newerVersions;
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
