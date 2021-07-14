import {Component, OnDestroy, OnInit} from '@angular/core';
import {catchError, filter, map, startWith, switchMap, takeUntil, tap} from "rxjs/operators";
import {Observable, Subject, timer, combineLatest, BehaviorSubject, of} from "rxjs";

import {ActivatedRoute, Router} from "@angular/router";

import {Project} from "../_models/project";
import {Trace} from "../_models/trace";

import {DataService} from "../_services/data.service";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  public logoInvertedUrl = environment?.config?.logoInvertedUrl;

  public project$: Observable<Project>;
  public contextId: string;
  private _rootEventsTimerInterval = 30;
  private readonly _projectTimerInterval = 30_000;

  private _errorSubject: BehaviorSubject<string> = new BehaviorSubject<string>(null);
  public error$: Observable<string> = this._errorSubject.asObservable();
  public isCreateMode$: Observable<boolean>;

  constructor(private router: Router, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    const projectName$ = this.route.params.pipe(
      map(params => params.projectName),
      filter(projectName => projectName)
    );

    this.isCreateMode$ = this.route.url.pipe(map(urlSegment => {
      return urlSegment[0].path === 'create';
    }));

    const timer$ = projectName$.pipe(
      switchMap((projectName) => timer(0, this._projectTimerInterval).pipe(map(() => projectName))),
      takeUntil(this.unsubscribe$)
    );
    timer$.subscribe(projectName => {
      this.dataService.loadProject(projectName);
      // this is on project-board level because we need the project in environment, service, sequence and settings screen
      // sequence screen because there is a check for the latest deployment context (lastEventTypes)
    });

    this.project$ = projectName$.pipe(
      switchMap(projectName => {return this.dataService.getProject(projectName)}),
      tap(project => {
        if (project === undefined) {
          this._errorSubject.next('project');
        } else {
          this._errorSubject.next(null);
        }
      }),
      catchError(() => {
        this._errorSubject.next('projects');
        return of(null);
      }));

    timer(0, this._rootEventsTimerInterval*1000)
      .pipe(
        startWith(0),
        switchMap(() => this.project$),
        filter(project => !!project && !!project.getServices()),
        takeUntil(this.unsubscribe$)
      ).subscribe(project => {
        this.dataService.loadRoots(project);
      });

    if (this.route.snapshot.url[0].path === 'trace') {
      const shkeptncontext$ = this.route.params.pipe(map(params => params.shkeptncontext));
      const eventselector$ = this.route.params.pipe(map(params => params.eventselector));
      const traces$ = shkeptncontext$.pipe(
        tap(shkeptncontext => {
          this.contextId = shkeptncontext;
          this.dataService.loadTracesByContext(shkeptncontext);
        }),
        switchMap(() => this.dataService.traces),
        filter(traces => !!traces),
      );

      combineLatest([traces$, eventselector$])
        .pipe(
          takeUntil(this.unsubscribe$)
        ).subscribe(([traces, eventselector]) => {
          if(traces.length > 0) {
            if(eventselector) {
              let trace = traces.find((t: Trace) => t.data.stage == eventselector && !!t.getProject() && !!t.getService());
              if (trace) {
                this.router.navigate(['/project', trace.getProject(),'sequence', trace.shkeptncontext, 'stage', trace.getStage()]);
              } else {
                trace = traces.reverse().find((t: Trace) => t.type == eventselector && !!t.getProject() && !!t.getService());
                if(trace) {
                  this.router.navigate(['/project', trace.getProject(), 'sequence', trace.shkeptncontext, 'event', trace.id]);
                } else {
                  this._errorSubject.next('trace');
                }
              }
            } else {
              const trace = traces.find((t: Trace) => !!t.getProject() && !!t.getService());
              this.router.navigate(['/project', trace.getProject(),'sequence', trace.shkeptncontext]);
            }
          } else {
            this._errorSubject.next('trace');
          }
        });
    }
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
