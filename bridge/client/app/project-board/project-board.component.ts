import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {map, takeUntil} from "rxjs/operators";
import {Observable, Subject} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";

import {Project} from "../_models/project";
import {Trace} from "../_models/trace";

import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public project$: Observable<Project>;
  public contextId: string;

  public error: string = null;
  public view: string = 'services';

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        if(params["shkeptncontext"]) {
          this.apiService.getTraces(params["shkeptncontext"])
            .pipe(
              map(response => response.body),
              map(result => result.events||[]),
              map(traces => traces.map(trace => Trace.fromJSON(trace)))
            )
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe((traces: Trace[]) => {
              if(traces.length > 0) {
                if(params["eventselector"]) {
                  let trace = traces.find((t: Trace) => t.data.stage == params["eventselector"] && !!t.getProject() && !!t.getService());
                  if(!trace)
                    trace = traces.reverse().find((t: Trace) => t.type == params["eventselector"] && !!t.getProject() && !!t.getService());

                  if(trace)
                    this.router.navigate(['/project', trace.getProject(), trace.getService(), trace.shkeptncontext, trace.id]);
                  else
                    this.error = "trace";
                } else {
                  let trace = traces.find((t: Trace) => !!t.getProject() && !!t.getService());
                  this.router.navigate(['/project', trace.getProject(), trace.getService(), trace.shkeptncontext]);
                }
              } else {
                this.error = "trace";
              }
            });
        } else {
          this.contextId = params["contextId"];
          this.project$ = this.dataService.getProject(params['projectName']);
          this.project$
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe(project => {
              if(project === undefined)
                this.error = 'project';
              this._changeDetectorRef.markForCheck();
            }, error => {
              this.error = 'projects';
            });
        }
      });
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  selectView(view) {
    this.view = view;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
