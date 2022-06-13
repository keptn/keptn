import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { BehaviorSubject, combineLatest, Observable, Subject } from 'rxjs';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { Trace } from '../_models/trace';
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
  public contextId?: string;
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

    if (this.route.snapshot.url[0].path === 'trace') {
      const shkeptncontext$ = this.route.paramMap.pipe(map((params: ParamMap) => params.get('shkeptncontext')));
      const eventselector$ = this.route.paramMap.pipe(map((params: ParamMap) => params.get('eventselector')));
      const traces$ = shkeptncontext$.pipe(
        tap((shkeptncontext: string | null) => {
          this.contextId = shkeptncontext ?? undefined;
          if (shkeptncontext) {
            this.dataService.loadTracesByContext(shkeptncontext);
          }
        }),
        switchMap(() => this.dataService.traces),
        filter((traces) => !!traces)
      );

      combineLatest([traces$, eventselector$])
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(([traces, eventselector]: [Trace[] | undefined, string | null]) => {
          this.navigateToTrace(traces, eventselector);
        });
    }
  }

  public navigateToTrace(traces: Trace[] | undefined, eventselector: string | null): void {
    if (traces?.length) {
      if (eventselector) {
        let trace = this.findTraceForStage(traces, eventselector);
        if (trace) {
          this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext, 'stage', trace.stage]);
          return;
        }
        trace = this.findTraceForEvent(traces, eventselector);
        if (trace) {
          this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext, 'event', trace.id]);
          return;
        }
        this._errorSubject.next('trace');
      } else {
        const trace = this.findTraceForKeptnContext(traces);
        if (trace) {
          this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext]);
        }
      }
    } else {
      this._errorSubject.next('trace');
    }
  }

  private findTraceForKeptnContext(traces: Trace[]): Trace | undefined {
    return traces.find((t: Trace) => !!t.project && !!t.service);
  }

  private findTraceForStage(traces: Trace[], eventselector: string | null): Trace | undefined {
    return traces.find((t: Trace) => t.data.stage === eventselector && !!t.project && !!t.service);
  }

  private findTraceForEvent(traces: Trace[], eventselector: string | null): Trace | undefined {
    return [...traces].reverse().find((t: Trace) => t.type === eventselector && !!t.project && !!t.service);
  }

  public loadProjects(): void {
    this.dataService.loadProjects().subscribe();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
