import semver from 'semver';
import {DOCUMENT} from '@angular/common';
import {Component, Inject, OnDestroy, OnInit} from '@angular/core';
import {NavigationEnd, Router, RoutesRecognized} from '@angular/router';
import {Title} from '@angular/platform-browser';
import {Observable, Subject, of} from 'rxjs';
import {filter, takeUntil} from 'rxjs/operators';
import {Project} from '../_models/project';
import {DataService} from '../_services/data.service';
import {NotificationsService} from '../_services/notifications.service';
import {NotificationType} from '../_models/notification';
import {environment} from '../../environments/environment';
import { KeptnInfo } from '../_models/keptn-info';
import { DtSwitchChange } from '@dynatrace/barista-components/switch';
import { VersionInfo } from '../_models/keptn-versions';

@Component({
  selector: 'app-header',
  templateUrl: './app-header.component.html',
  styleUrls: ['./app-header.component.scss']
})
export class AppHeaderComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  public projects: Observable<Project[] | undefined>;
  public project$: Observable<Project | undefined> = of(undefined);
  public projectBoardView = '';
  public appTitle = environment?.config?.appTitle;
  public logoUrl = environment?.config?.logoUrl;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public keptnInfo?: KeptnInfo;
  public versionCheckDialogState: string | null = null;
  public versionCheckReference = '/reference/version_check/';

  constructor(@Inject(DOCUMENT) private _document: HTMLDocument, private router: Router, private dataService: DataService,
              private notificationsService: NotificationsService, private titleService: Title) {
    this.projects = this.dataService.projects;
  }

  ngOnInit() {
    this.titleService.setTitle(this.appTitle);
    this.setAppFavicon(this.logoInvertedUrl);

    this.router.events
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(event => {
        if (event instanceof RoutesRecognized) {
          const projectName = event.state.root.children[0].params.projectName;
          this.project$ = this.dataService.getProject(projectName);
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
      .pipe(
        filter((keptnInfo: KeptnInfo | undefined): keptnInfo is KeptnInfo => !!keptnInfo),
        takeUntil(this.unsubscribe$)
      ).subscribe(keptnInfo => {
        this.keptnInfo = keptnInfo;
        if (keptnInfo.versionCheckEnabled === undefined) {
          this.showVersionCheckInfoDialog();
        } else if (keptnInfo.bridgeInfo.enableVersionCheckFeature && keptnInfo.versionCheckEnabled && keptnInfo.availableVersions) {
          keptnInfo.keptnVersionInvalid = !semver.valid(keptnInfo.metadata.keptnversion);
          keptnInfo.bridgeVersionInvalid = !semver.valid(keptnInfo.bridgeInfo.bridgeVersion);
          this.doVersionCheck(
            keptnInfo.bridgeInfo.bridgeVersion,
            keptnInfo.metadata.keptnversion,
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

  doVersionCheck(bridgeVersion: string | undefined,
                 cliVersion: string,
                 availableBridgeVersions: VersionInfo,
                 availableCliVersions: VersionInfo) {
    if (semver.valid(bridgeVersion) && semver.valid(cliVersion)) {
      const latestVersion = availableCliVersions.stable[availableCliVersions.stable.length - 1];
      const bridgeVersionsString = this.buildVersionString(this.getNewerVersions(availableBridgeVersions, bridgeVersion));
      const cliVersionsString = this.buildVersionString(this.getNewerVersions(availableCliVersions, cliVersion));

      if (bridgeVersionsString || cliVersionsString) {
        const versionMessage = `New ${cliVersionsString ? ' Keptn CLI ' + cliVersionsString : ''} ${cliVersionsString && bridgeVersionsString ? 'and' : ''}
                            ${bridgeVersionsString ? ' Keptn Bridge ' + bridgeVersionsString : ''} available. For details how to upgrade visit
                            <a href="https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(latestVersion)}.x/operate/upgrade/" target="_blank">
                            https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(latestVersion)}.x/operate/upgrade/</a>`;
        this.notificationsService.addNotification(NotificationType.Info, versionMessage);
      }
    }
  }

  private buildVersionString(versions: VersionInfo) {
    if (versions.stable.length > 0) {
      return versions.stable.join(', ');
    } else if (versions.prerelease.length > 0) {
      return versions.prerelease.join(', ');
    }

    return null;
  }

  private getNewerVersions(availableVersions: VersionInfo, currentVersion?: string): VersionInfo {
    const newerVersions: VersionInfo = {
      stable: [],
      prerelease: []
    };
    if (currentVersion) {
      newerVersions.stable = availableVersions.stable.filter((stableVersion: string) => semver.lt(currentVersion, stableVersion));

      // It is only necessary to check prerelease versions when no stable update is available
      if (newerVersions.stable.length === 0) {
        newerVersions.prerelease = availableVersions.prerelease
                                  .filter((prereleaseVersion: string) => semver.lt(currentVersion, prereleaseVersion));
      }
    }

    return newerVersions;
  }

  showVersionCheckInfoDialog() {
    if (this.keptnInfo?.bridgeInfo.enableVersionCheckFeature) {
      this.versionCheckDialogState = 'info';
    }
  }

  acceptVersionCheck(accepted: boolean): void {
    this.dataService.setVersionCheck(accepted);
    if (accepted) {
      this.versionCheckDialogState = 'success';
    }

    setTimeout(() => {
      this.versionCheckDialogState = null;
    }, accepted ? 2000 : 0);
  }

  // tslint:disable-next-line:no-any
  versionCheckClicked(event: DtSwitchChange<any>) {
    this.dataService.setVersionCheck(event.checked);
  }

  setAppFavicon(path: string){
    this._document.getElementById('appFavicon')?.setAttribute('href', path);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
