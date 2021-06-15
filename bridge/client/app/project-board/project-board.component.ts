import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {filter, startWith, switchMap, takeUntil} from "rxjs/operators";
import {Observable, Subject, timer} from "rxjs";
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

  public error: string = null;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        if (this.route.snapshot.url[0].path === 'trace') {
          this.dataService.loadTracesByContext(params.shkeptncontext);
          this.dataService.traces
            .pipe(filter(traces => !!traces))
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe((traces: Trace[]) => {
              if(traces.length > 0) {
                if(params.eventselector) {
                  let trace = traces.find((t: Trace) => t.data.stage == params.eventselector && !!t.getProject() && !!t.getService());
                  if (trace) {
                    this.router.navigate(['/project', trace.getProject(),'sequence', trace.shkeptncontext, 'stage', trace.getStage()]);
                  } else {
                    trace = traces.reverse().find((t: Trace) => t.type == params.eventselector && !!t.getProject() && !!t.getService());
                    if(trace) {
                      this.router.navigate(['/project', trace.getProject(), 'sequence', trace.shkeptncontext, 'event', trace.id]);
                    } else {
                      this.error = "trace";
                    }
                  }
                } else {
                  let trace = traces.find((t: Trace) => !!t.getProject() && !!t.getService());
                  this.router.navigate(['/project', trace.getProject(),'sequence', trace.shkeptncontext]);
                }
              } else {
                this.error = "trace";
              }
            });
        } else {
          this.contextId = params.contextId;
          this.project$ = this.dataService.getProject(params.projectName);
          this.project$
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(project => {
              if (project === undefined) {
                this.error = 'project';
              } else {
                this.error = null;
              }
              this._changeDetectorRef.markForCheck();
            }, error => {
              this.error = 'projects';
            });

          timer(0, this._rootEventsTimerInterval*1000)
            .pipe(
              startWith(0),
              switchMap(() => this.project$),
              filter(project => !!project && !!project.getServices())
            )
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(project => {
              this.dataService.loadRoots(project);
            });
        }
      });
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
