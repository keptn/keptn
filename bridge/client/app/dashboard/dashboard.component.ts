import { Component, OnDestroy } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';

@Component({
  selector: 'ktb-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnDestroy {
  public readonly projects$: Observable<Project[] | undefined> = this.dataService.projects;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly$: Observable<boolean>;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private dataService: DataService) {
    this.isQualityGatesOnly$ = dataService.isQualityGatesOnly;

    this.dataService.keptnInfo
      .pipe(
        filter((keptnInfo) => !!keptnInfo),
        take(1)
      )
      .subscribe(() => {
        this.loadProjects();
      });
  }

  public loadProjects(): void {
    this.dataService.loadProjects();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
