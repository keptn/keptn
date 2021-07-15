import {Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, switchMap, takeUntil, tap} from "rxjs/operators";
import {Observable, Subject, timer, combineLatest} from "rxjs";
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

  public error: string = null;

  constructor(private router: Router, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    const projectName$ = this.route.params.pipe(
      map(params => params.projectName),
      filter(projectName => projectName)
    );

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
      switchMap(projectName => this.dataService.getProject(projectName))
    );

    this.project$
      .pipe(
        takeUntil(this.unsubscribe$)
      ).subscribe(project => {
        if (project === undefined) {
          this.error = 'project';
        } else {
          this.error = null;
        }
      }, error => {
        this.error = 'projects';
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
                  this.error = "trace";
                }
              }
            } else {
              const trace = traces.find((t: Trace) => !!t.getProject() && !!t.getService());
              this.router.navigate(['/project', trace.getProject(),'sequence', trace.shkeptncontext]);
            }
          } else {
            this.error = "trace";
          }
        });
    }
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
