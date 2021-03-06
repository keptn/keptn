import {ChangeDetectorRef, Component, NgZone, OnDestroy, OnInit} from '@angular/core';
import {Observable, Subject, timer} from 'rxjs';
import {Project} from '../_models/project';
import {DataService} from '../_services/data.service';
import {environment} from '../../environments/environment';
import {takeUntil} from 'rxjs/operators';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: []
})
export class DashboardComponent implements OnInit, OnDestroy{
  public projects$: Observable<Project[]>;
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public isQualityGatesOnly: boolean;

  private readonly _projectTimerInterval = 30 * 1000;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private _changeDetectorRef: ChangeDetectorRef,
              private dataService: DataService,
              private ngZone: NgZone) {
    this.projects$ = this.dataService.projects;
  }

  public ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(isQualityGatesOnly => {this.isQualityGatesOnly = isQualityGatesOnly});

    // If we don't run this outside angular e2e tests will fail
    // because Protractor waits for async tasks to complete - in case of timer they do not finish so the tests time out
    // https://github.com/angular/protractor/blob/master/docs/timeouts.md#waiting-for-angular
    this.ngZone.runOutsideAngular(() => {
      timer(this._projectTimerInterval, this._projectTimerInterval)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(() => {
          this.ngZone.run(() => {
            this.loadProjects();
          });
        });
    });
  }

  public loadProjects() {
    this.dataService.loadProjects();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
