import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, startWith, switchMap} from "rxjs/operators";
import {Observable, Subscription, timer} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import {Location} from "@angular/common";

import * as moment from 'moment';

import {Root} from "../_models/root";
import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";
import DateUtil from "../_utils/date.utils";
import {Service} from "../_models/service";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  public project$: Observable<Project>;
  public currentRoot: Root;
  public error: string = null;

  private _projectSub: Subscription = Subscription.EMPTY;
  private _routeSubs: Subscription = Subscription.EMPTY;
  private _rootsSubs: Subscription = Subscription.EMPTY;
  private _rootEventsTimer: Subscription = Subscription.EMPTY;
  private _rootEventsTimerInterval = 30;

  private _tracesTimer: Subscription = Subscription.EMPTY;
  private _tracesTimerInterval = 10;

  public projectName: string;
  public serviceName: string;
  public contextId: string;
  public eventId: string;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      if(params["shkeptncontext"]) {
        this.contextId = params["shkeptncontext"];
        this.apiService.getTraces(this.contextId)
          .pipe(
            map(traces => traces.map(trace => Trace.fromJSON(trace)))
          )
          .subscribe((traces: Trace[]) => {
            if(traces.length > 0) {
              if(params["eventselector"]) {
                let trace = traces.find((t: Trace) => t.data.stage == params["eventselector"]);
                if(!trace)
                  trace = traces.reverse().find((t: Trace) => t.type == params["eventselector"]);

                if(trace)
                  this.router.navigate(['/project', trace.data.project, trace.data.service, trace.shkeptncontext, trace.id]);
                else
                  this.error = "trace";
              } else {
                this.router.navigate(['/project', traces[0].data.project, traces[0].data.service, traces[0].shkeptncontext]);
              }
            } else {
              this.error = "trace";
            }
          });
      } else {
        this.projectName = params["projectName"];
        this.serviceName = params["serviceName"];
        this.contextId = params["contextId"];
        this.eventId = params["eventId"];
        this.currentRoot = null;

        this.project$ = this.dataService.projects.pipe(
          map(projects => projects ? projects.find(project => {
            return project.projectName === params['projectName'];
          }) : null)
        );

        this._projectSub = this.project$.subscribe(projects => {
          this.error = false;
        }, error => {
          this.error = true;
        });

        this._rootsSubs.unsubscribe();
        this._rootsSubs = this.dataService.roots.subscribe(roots => {
          if(roots && !this.currentRoot)
            this.currentRoot = roots.find(r => r.shkeptncontext == params["contextId"]);
        });

        this._rootEventsTimer.unsubscribe();
        this._rootEventsTimer = timer(0, this._rootEventsTimerInterval*1000)
          .pipe(
            startWith(0),
            switchMap(() => this.project$),
            filter(project => !!project && !!project.getServices())
          )
          .subscribe(project => {
            project.getServices().forEach(service => {
              this.dataService.loadRoots(project, service);
              if(service.roots && !this.currentRoot)
                this.currentRoot = service.roots.find(r => r.shkeptncontext == params["contextId"]);
            });
          });
      }
    });
  }

  selectRoot(event: any): void {
    this.projectName = event.root.data.project;
    this.serviceName = event.root.data.service;
    this.contextId = event.root.data.shkeptncontext;
    this.eventId = null;
    if(event.stage) {
      let focusEvent = event.root.traces.find(trace => trace.data.stage == event.stage);
      let routeUrl = this.router.createUrlTree(['/project', focusEvent.data.project, focusEvent.data.service, focusEvent.shkeptncontext, focusEvent.id]);
      this.eventId = focusEvent.id;
      this.location.go(routeUrl.toString());
    } else {
      let routeUrl = this.router.createUrlTree(['/project', event.root.data.project, event.root.data.service, event.root.shkeptncontext]);
      this.eventId = event.root.traces[event.root.traces.length-1].id;
      this.location.go(routeUrl.toString());
    }

    this.currentRoot = event.root;
    this.loadTraces(this.currentRoot);
  }

  loadTraces(root: Root): void {
    this._tracesTimer.unsubscribe();
    if(moment().subtract(1, 'day').isBefore(root.time)) {
      this._tracesTimer = timer(0, this._tracesTimerInterval*1000)
        .subscribe(() => {
          this.dataService.loadTraces(root);
        });
    } else {
      this.dataService.loadTraces(root);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats(true);
  }

  getLatestDeployment(project: Project, service: Service, stage?: Stage): Trace {
    let currentService = project.getServices()
      .find(s => s.serviceName == service.serviceName);

    if(currentService.roots)
      return currentService.roots
        .reduce((traces: Trace[], root: Root) => {
          return [...traces, ...root.traces];
        }, [])
        .find(trace => trace.type == 'sh.keptn.events.deployment-finished' && (!stage || (trace.data.stage == stage.stageName && currentService.roots.find(r => r.shkeptncontext == trace.shkeptncontext).isFaulty() != stage.stageName)));
    else
      return null;
  }

  getDeployedServices(project: Project, stage: Stage) {
    return stage.services.filter(service => !!this.getLatestDeployment(project, service, stage));
  }

  getShortImageName(image) {
    let parts = image.split("/");
    return parts[parts.length-1];
  }

  getRootsLastUpdated(project: Project, service: Service): Date {
    return this.dataService.getRootsLastUpdated(project, service);
  }

  getTracesLastUpdated(root: Root): Date {
    return this.dataService.getTracesLastUpdated(root);
  }

  showReloadButton(root: Root) {
    return moment().subtract(1, 'day').isAfter(root.time);
  }

  loadProjects() {
    this.dataService.loadProjects();
  }

  trackStage(index: number, stage: Stage) {
    return stage.stageName;
  }

  ngOnDestroy(): void {
    this._projectSub.unsubscribe();
    this._routeSubs.unsubscribe();
    this._tracesTimer.unsubscribe();
    this._rootEventsTimer.unsubscribe();

  }

}
