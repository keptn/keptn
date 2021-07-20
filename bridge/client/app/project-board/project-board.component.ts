import {Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, switchMap, takeUntil, tap} from 'rxjs/operators';
import {Observable, Subject, timer, combineLatest} from 'rxjs';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import {Project} from '../_models/project';
import {Trace} from '../_models/trace';
import {DataService} from '../_services/data.service';
import {environment} from '../../environments/environment';

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  public logoInvertedUrl = environment?.config?.logoInvertedUrl;
  public project$: Observable<Project | undefined>;
  public contextId?: string;
  private readonly _projectTimerInterval = 30_000;
  public error?: string;

  constructor(private router: Router, private route: ActivatedRoute, private dataService: DataService) {
    const projectName$ = this.route.paramMap.pipe(
      map(params => params.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName)
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
  }

  ngOnInit() {
    this.project$
      .pipe(
        takeUntil(this.unsubscribe$)
      ).subscribe((project: Project | undefined) => {
        if (project === undefined) {
          this.error = 'project';
        } else {
          this.error = undefined;
        }
      }, () => {
        this.error = 'projects';
      });

    if (this.route.snapshot.url[0].path === 'trace') {
      const shkeptncontext$ = this.route.paramMap.pipe(map((params: ParamMap) => params.get('shkeptncontext')));
      const eventselector$ = this.route.paramMap.pipe(map((params: ParamMap) => params.get('eventselector')));
      const traces$ = shkeptncontext$.pipe(
        tap((shkeptncontext: string | null) => {
          this.contextId = shkeptncontext ? shkeptncontext : undefined;
          if (shkeptncontext) {
            this.dataService.loadTracesByContext(shkeptncontext);
          }
        }),
        switchMap(() => this.dataService.traces),
        filter(traces => !!traces),
      );

      combineLatest([traces$, eventselector$])
        .pipe(
          takeUntil(this.unsubscribe$)
        ).subscribe(([traces, eventselector]: [Trace[] | undefined, string | null]) => {
          if (traces?.length) {
            if (eventselector) {
              let trace = traces.find((t: Trace) => t.data.stage === eventselector && !!t.project && !!t.service);
              if (trace) {
                this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext, 'stage', trace.stage]);
              } else {
                trace = traces.reverse().find((t: Trace) => t.type === eventselector && !!t.project && !!t.service);
                if (trace) {
                  this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext, 'event', trace.id]);
                } else {
                  this.error = 'trace';
                }
              }
            } else {
              const trace = traces.find((t: Trace) => !!t.project && !!t.service);
              if (trace) {
                this.router.navigate(['/project', trace.project, 'sequence', trace.shkeptncontext]);
              }
            }
          } else {
            this.error = 'trace';
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
