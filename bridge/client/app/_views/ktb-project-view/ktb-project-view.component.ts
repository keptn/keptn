import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { distinctUntilChanged, filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { environment } from '../../../environments/environment';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

@Component({
  selector: 'ktb-project-view',
  templateUrl: './ktb-project-view.component.html',
  styleUrls: ['./ktb-project-view.component.scss'],
})
export class KtbProjectViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public hasProject$: Observable<boolean | undefined>;
  private _errorSubject: BehaviorSubject<string | undefined> = new BehaviorSubject<string | undefined>(undefined);
  public error$: Observable<string | undefined> = this._errorSubject.asObservable();
  public isCreateMode$: Observable<boolean>;
  public hasUnreadLogs$: Observable<boolean>;
  private static readonly uniformLogPollingInterval = 2 * 60_000;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private dataService: DataService,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    this.hasUnreadLogs$ = this.dataService.hasUnreadUniformRegistrationLogs;
    // disable log-polling
    const uniformLogInterval = initialDelayMillis === 0 ? 0 : KtbProjectViewComponent.uniformLogPollingInterval;
    const projectName$ = this.route.paramMap.pipe(
      map((params) => params.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName)
    );

    const uniformLogTimer$ = projectName$.pipe(
      switchMap(() => AppUtils.createTimer(0, uniformLogInterval)),
      takeUntil(this.unsubscribe$)
    );

    uniformLogTimer$.subscribe(() => {
      this.dataService.loadUnreadUniformRegistrationLogs();
    });

    this.hasProject$ = projectName$.pipe(switchMap((projectName) => this.dataService.projectExists(projectName)));

    this.isCreateMode$ = this.route.data.pipe(
      map((data) => !!data.createMode),
      distinctUntilChanged()
    );
  }

  ngOnInit(): void {
    this.hasProject$
      .pipe(
        filter((hasProject) => hasProject !== undefined),
        tap((hasProject) => {
          if (hasProject) {
            this._errorSubject.next(undefined);
          } else {
            this._errorSubject.next('project');
          }
        }),
        takeUntil(this.unsubscribe$)
      )
      .subscribe();
  }

  public loadProjects(): void {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
