import semver from 'semver';
import { DOCUMENT } from '@angular/common';
import { Component, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { NavigationEnd, NavigationStart, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { combineLatest, Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';
import { NotificationsService } from '../_services/notifications.service';
import { NotificationType } from '../_models/notification';
import { environment } from '../../environments/environment';
import { KeptnInfo } from '../_models/keptn-info';
import { DtSwitchChange } from '@dynatrace/barista-components/switch';
import { VersionInfo } from '../../../shared/interfaces/keptn-versions';
import { DtSelect } from '@dynatrace/barista-components/select';
import { IMetadata } from '../_interfaces/metadata';

@Component({
  selector: 'ktb-header',
  templateUrl: './app-header.component.html',
  styleUrls: ['./app-header.component.scss'],
})
export class AppHeaderComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  @ViewChild('projectSelect') projectSelect?: DtSelect<string | undefined>;
  public projects: Observable<Project[] | undefined>;
  public selectedProject: string | undefined;
  public projectBoardView = '';
  public appTitle = environment.config.appTitle;
  public logoUrl = environment.config.logoUrl;
  public logoInvertedUrl = environment.config.logoInvertedUrl;
  public keptnInfo?: KeptnInfo;
  public versionCheckDialogState: string | null = null;
  public versionCheckReference = '/reference/version_check/';
  public metadata?: IMetadata;

  constructor(
    @Inject(DOCUMENT) private _document: Document,
    private router: Router,
    private dataService: DataService,
    private notificationsService: NotificationsService,
    private titleService: Title
  ) {
    this.projects = this.dataService.projects;
    this.router.events.pipe(takeUntil(this.unsubscribe$)).subscribe((event) => {
      if (event instanceof NavigationStart || event instanceof NavigationEnd) {
        this.setProject();
      }
    });
  }

  ngOnInit(): void {
    this.titleService.setTitle(this.appTitle);
    this.setAppFavicon(this.logoInvertedUrl);

    combineLatest([this.dataService.keptnInfo, this.dataService.keptnMetadata])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([keptnInfo, metadata]) => {
        this.keptnInfo = keptnInfo;
        this.metadata = metadata ?? undefined;
        if (!keptnInfo) {
          return;
        }
        if (keptnInfo.versionCheckEnabled === undefined) {
          this.showVersionCheckInfoDialog();
        } else if (
          metadata &&
          keptnInfo.bridgeInfo.enableVersionCheckFeature &&
          keptnInfo.versionCheckEnabled &&
          keptnInfo.availableVersions
        ) {
          keptnInfo.keptnVersionInvalid = !semver.valid(metadata.keptnversion);
          keptnInfo.bridgeVersionInvalid = !semver.valid(keptnInfo.bridgeInfo.bridgeVersion);
          this.doVersionCheck(
            keptnInfo.bridgeInfo.bridgeVersion,
            metadata.keptnversion,
            keptnInfo.availableVersions.bridge,
            keptnInfo.availableVersions.cli
          );
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

  doVersionCheck(
    bridgeVersion: string | undefined,
    cliVersion: string,
    availableBridgeVersions: VersionInfo,
    availableCliVersions: VersionInfo
  ): void {
    if (semver.valid(bridgeVersion) && semver.valid(cliVersion)) {
      const latestVersion = availableCliVersions.stable[availableCliVersions.stable.length - 1];
      const bridgeVersionsString = this.buildVersionString(
        this.getNewerVersions(availableBridgeVersions, bridgeVersion)
      );
      const cliVersionsString = this.buildVersionString(this.getNewerVersions(availableCliVersions, cliVersion));

      if (bridgeVersionsString || cliVersionsString) {
        const versionMessage = `New ${cliVersionsString ? ' Keptn CLI ' + cliVersionsString : ''} ${
          cliVersionsString && bridgeVersionsString ? 'and' : ''
        }
                            ${
                              bridgeVersionsString ? ' Keptn Bridge ' + bridgeVersionsString : ''
                            } available. For details how to upgrade visit
                            <a href="https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(
          latestVersion
        )}.x/operate/upgrade/" target="_blank">
                            https://keptn.sh/docs/${semver.major(latestVersion)}.${semver.minor(
          latestVersion
        )}.x/operate/upgrade/</a>`;
        this.notificationsService.addNotification(NotificationType.INFO, versionMessage, undefined, -1);
      }
    }
  }

  private buildVersionString(versions: VersionInfo): null | string {
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
      prerelease: [],
    };
    if (currentVersion) {
      const stable = availableVersions.stable.filter((stableVersion: string) =>
        semver.lt(currentVersion, stableVersion)
      );
      newerVersions.stable = this.reduceVersions(stable, currentVersion);

      // It is only necessary to check prerelease versions when no stable update is available
      if (newerVersions.stable.length === 0) {
        const prerelease = availableVersions.prerelease.filter((prereleaseVersion: string) =>
          semver.lt(currentVersion, prereleaseVersion)
        );
        newerVersions.prerelease = this.reduceVersions(prerelease, currentVersion);
      }
    }

    return newerVersions;
  }

  private reduceVersions(stable: string[], currentVersion: string): string[] {
    let latestMinor: string | undefined;
    let latestMajor: string | undefined;

    for (const version of stable) {
      if (!latestMajor || semver.gt(version, latestMajor)) {
        latestMajor = version;
      }
      if (semver.minor(version) === semver.minor(currentVersion) && (!latestMinor || semver.gt(version, latestMinor))) {
        latestMinor = version;
      }
    }

    if (latestMajor === latestMinor) {
      latestMajor = undefined;
    }

    return [...(latestMinor ? [latestMinor] : []), ...(latestMajor ? [latestMajor] : [])];
  }

  showVersionCheckInfoDialog(): void {
    if (this.keptnInfo?.bridgeInfo.enableVersionCheckFeature) {
      this.versionCheckDialogState = 'info';
    }
  }

  acceptVersionCheck(accepted: boolean): void {
    this.dataService.setVersionCheck(accepted);
    if (accepted) {
      this.versionCheckDialogState = 'success';
    }

    setTimeout(
      () => {
        this.versionCheckDialogState = null;
      },
      accepted ? 2000 : 0
    );
  }

  versionCheckClicked(event: DtSwitchChange<unknown>): void {
    this.dataService.setVersionCheck(event.checked);
  }

  setAppFavicon(path: string): void {
    this._document.getElementById('appFavicon')?.setAttribute('href', path);
  }

  changeProject(selectedProject: string | undefined): void {
    this.router.navigate(this.getRouterLink(selectedProject as string));
  }

  setProject(): void {
    const urlPieces = this.router.url.split('?')[0].split('/');
    if (urlPieces[1] === 'project') {
      this.selectedProject = urlPieces[2];

      // catch url change and update projectBoardView for the project picker
      if (urlPieces.length > 3) {
        this.projectBoardView = urlPieces[3];
      } else {
        this.projectBoardView = ''; // environment screen
      }
    } else if (urlPieces[1] === 'evaluation') {
      this.dataService.projectName.pipe(takeUntil(this.unsubscribe$)).subscribe((projectName: string) => {
        this.selectedProject = projectName;
      });
    } else {
      this.selectedProject = undefined;
    }

    if (this.projectSelect) {
      this.projectSelect.value = this.selectedProject;
    }
  }

  getKeptnAuthCommand(keptnInfo: KeptnInfo): string {
    return (keptnInfo.bridgeInfo.authMsg || keptnInfo.authCommand) ?? '';
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
