import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../_services/data.service';
import { environment } from '../../environments/environment';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';

@Component({
  selector: 'ktb-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss'],
})
export class ProjectBoardComponent implements OnInit, OnDestroy {
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
    const uniformLogInterval = initialDelayMillis === 0 ? 0 : ProjectBoardComponent.uniformLogPollingInterval;
    const projectName$ = this.route.paramMap.pipe(
      map((params) => params.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName)
    );

    const projectTimer$ = projectName$.pipe(
      switchMap((projectName) => AppUtils.createTimer(0, initialDelayMillis).pipe(map(() => projectName))),
      takeUntil(this.unsubscribe$)
    );

    const uniformLogTimer$ = projectName$.pipe(
      switchMap(() => AppUtils.createTimer(0, uniformLogInterval)),
      takeUntil(this.unsubscribe$)
    );

    projectTimer$.subscribe((projectName) => {
      // this is on project-board level because we need the project in environment, service, sequence and settings screen
      // sequence screen because there is a check for the latest deployment context (lastEventTypes)
      this.dataService.loadProject(projectName);
    });

    uniformLogTimer$.subscribe(() => {
      this.dataService.loadUnreadUniformRegistrationLogs();
    });
    this.hasProject$ = projectName$.pipe(switchMap((projectName) => this.dataService.projectExists(projectName)));

    this.isCreateMode$ = this.route.url.pipe(map((urlSegment) => urlSegment[0].path === 'create'));
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
