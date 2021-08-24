import { Component, Inject, NgZone, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';
import { environment } from '../../environments/environment';
import { takeUntil } from 'rxjs/operators';
import { DtOverlay } from '@dynatrace/barista-components/overlay';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';

@Component({
  selector: 'ktb-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  public projects$: Observable<Project[] | undefined>;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly = false;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private dataService: DataService, private ngZone: NgZone, private _dtOverlay: DtOverlay, @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number) {
    this.projects$ = this.dataService.projects;
  }

  public ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(isQualityGatesOnly => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });

    AppUtils.createTimer(this.initialDelayMillis, this.initialDelayMillis)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.loadProjects();
      });
  }

  public loadProjects() {
    this.dataService.loadProjects();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
