import { Component, Inject } from '@angular/core';
import { merge, Observable, of, scan, switchMap } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { DataService } from '../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';
import { ProjectSequences } from '../_components/ktb-project-list/ktb-project-list.component';
import { IMetadata } from '../_interfaces/metadata';
import { IProject } from '../../../shared/interfaces/project';
import { Router } from '@angular/router';

const MAX_SEQUENCES = 5;

@Component({
  selector: 'ktb-dashboard-legacy',
  templateUrl: './dashboard-legacy.component.html',
  styleUrls: ['./dashboard-legacy.component.scss'],
})
export class DashboardLegacyComponent {
  private readonly refreshTimer$ = AppUtils.createTimer(0, this.initialDelayMillis);
  public readonly keptnMetadata$ = this.dataService.keptnMetadata.pipe(
    filter((metadata): metadata is IMetadata => metadata != null)
  );
  public readonly projects$: Observable<IProject[] | undefined> = this.dataService.projects;
  public readonly latestSequences$: Observable<ProjectSequences> = this.projects$.pipe(
    switchMap((projects) => (projects ? merge(...this.loadSequences(projects)) : of({}))),
    scan((agg, next) => ({ ...agg, ...next }), {} as ProjectSequences)
  );
  public readonly isQualityGatesOnly$: Observable<boolean> = this.dataService.isQualityGatesOnly;
  public readonly logoInvertedUrl = environment?.config?.logoInvertedUrl;

  constructor(
    private router: Router,
    private dataService: DataService,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    this.checkRefreshProjects();
  }

  checkRefreshProjects(): void {
    const currentNav = this.router.getCurrentNavigation();
    const hadPreviousNavigation = currentNav != null && currentNav.previousNavigation != null;
    if (hadPreviousNavigation) {
      this.refreshProjects();
    }
  }

  private loadSequences(projects: IProject[]): Observable<ProjectSequences>[] {
    return projects.map((project) =>
      this.refreshTimer$.pipe(
        switchMap(() =>
          this.dataService
            .loadLatestSequences(project.projectName, MAX_SEQUENCES)
            .pipe(map((sequences) => ({ [project.projectName]: sequences })))
        )
      )
    );
  }

  refreshProjects(): void {
    this.dataService.loadProjects();
  }
}
