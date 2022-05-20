import { Component, NgZone, OnDestroy, OnInit } from '@angular/core';
import { DtOverlay } from '@dynatrace/barista-components/overlay';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Project } from '../_models/project';
import { DataService } from '../_services/data.service';

@Component({
  selector: 'ktb-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  public readonly projects$: Observable<Project[] | undefined> = this.dataService.projects;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly = false;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private dataService: DataService, private ngZone: NgZone, private _dtOverlay: DtOverlay) {}

  public ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(takeUntil(this.unsubscribe$)).subscribe((isQualityGatesOnly) => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });

    this.loadProjects();
  }

  public loadProjects(): void {
    this.dataService.loadProjects();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
