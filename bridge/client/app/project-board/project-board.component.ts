import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, startWith, switchMap, take, takeUntil} from "rxjs/operators";
import {Observable, Subject, timer} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";

import {Project} from "../_models/project";
import {Trace} from "../_models/trace";

import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";

import {Location} from "@angular/common";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  public project$: Observable<Project>;
  public contextId: string;
  private _rootEventsTimerInterval = 30;

  public error: string = null;
  public view: string = 'environment';

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService, private location: Location) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        if (this.route.snapshot.url[0].path === 'trace') {
          this.apiService.getTraces(params.shkeptncontext)
            .pipe(
              map(response => response.body),
              map(result => result.events||[]),
              map(traces => traces.map(trace => Trace.fromJSON(trace)))
            )
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
              if (project === undefined)
                this.error = 'project';
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

          if(this.route.snapshot.url.length > 2) {
            this.view = this.route.snapshot.url[2].path;
          } else {
            this.view = 'environment';
          }
        }
      });
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  selectView(view: string, projectName: string) {
    this.setView(view, '/project', projectName);
  }

  setView(view: string, ...urlCommands: string[]) {
    if (this.view !== view) {
      this.router.navigate(urlCommands);
    }
  }

  redirectView(view: string, projectName: string) {
    this.setView(view, '/project', projectName, view);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
