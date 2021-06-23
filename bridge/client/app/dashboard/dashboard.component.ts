import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {Observable, Subject, timer} from 'rxjs';
import {Project} from '../_models/project';
import {DataService} from '../_services/data.service';
import {environment} from '../../environments/environment';
import {filter, takeUntil} from 'rxjs/operators';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy{
  public projects$: Observable<Project[]>;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly: boolean;

  private readonly _projectTimerInterval = 30 * 1000;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.projects$ = this.dataService.projects;
  }

  public ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(isQualityGatesOnly => {this.isQualityGatesOnly = isQualityGatesOnly});

    timer(this._projectTimerInterval, this._projectTimerInterval)
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
  }
}
