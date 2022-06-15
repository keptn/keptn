import { Component, OnDestroy } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { filter, mergeMap, take } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';

@Component({
  selector: 'ktb-dashboard-legacy',
  templateUrl: './dashboard-legacy.component.html',
  styleUrls: ['./dashboard-legacy.component.scss'],
})
export class DashboardLegacyComponent implements OnDestroy {
  public readonly projects$: Observable<Project[] | undefined> = this.dataService.projects;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly$: Observable<boolean>;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private dataService: DataService) {
    this.isQualityGatesOnly$ = dataService.isQualityGatesOnly;

    this.dataService.keptnInfo
      .pipe(
        filter((keptnInfo) => !!keptnInfo),
        take(1),
        mergeMap(() => this.dataService.loadProjects())
      )
      .subscribe();
  }

  public loadProjects(): void {
    this.dataService.loadProjects().subscribe();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
