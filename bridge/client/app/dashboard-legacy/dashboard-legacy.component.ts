import { Component, Inject } from '@angular/core';
import { merge, Observable, of, scan, switchMap } from 'rxjs';
import { filter, map, take } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';
import { ProjectSequences } from '../_components/ktb-project-list/ktb-project-list.component';
import { IMetadata } from '../_interfaces/metadata';

const MAX_SEQUENCES = 5;

@Component({
  selector: 'ktb-dashboard-legacy',
  templateUrl: './dashboard-legacy.component.html',
  styleUrls: ['./dashboard-legacy.component.scss'],
})
export class DashboardLegacyComponent {
  private readonly refreshTimer$ = AppUtils.createTimer(0, this.initialDelayMillis);
  private readonly keptnInfo$ = this.dataService.keptnInfo.pipe(
    filter((keptnInfo) => !!keptnInfo),
    take(1)
  );
  public readonly keptnMetadata$ = this.dataService.keptnMetadata.pipe(
    filter((metadata): metadata is IMetadata | null => metadata !== undefined)
  );
  public readonly projects$: Observable<Project[] | undefined> = this.dataService.projects;
  public readonly latestSequences$: Observable<ProjectSequences> = this.projects$.pipe(
    switchMap((projects) => (projects ? merge(...this.loadSequences(projects)) : of({}))),
    scan((agg, next) => ({ ...agg, ...next }), {} as ProjectSequences)
  );
  public readonly isQualityGatesOnly$: Observable<boolean> = this.dataService.isQualityGatesOnly;
  public readonly logoInvertedUrl = environment?.config?.logoInvertedUrl;

  constructor(private dataService: DataService, @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number) {
    this.keptnInfo$.subscribe(() => this.loadProjects());
  }

  private loadSequences(projects: Project[]): Observable<ProjectSequences>[] {
    return projects.map((project) =>
      this.refreshTimer$.pipe(
        switchMap(() =>
          this.dataService
            .loadLatestSequences(project, MAX_SEQUENCES)
            .pipe(map((sequences) => ({ [project.projectName]: sequences })))
        )
      )
    );
  }

  public loadProjects(): void {
    this.dataService.loadProjects().subscribe();
  }
}
