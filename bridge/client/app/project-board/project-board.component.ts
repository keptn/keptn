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
import {DtCheckboxChange} from "@dynatrace/barista-components/checkbox";
import {EVENT_LABELS} from "../_models/event-labels";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  public project$: Observable<Project>;
  public openApprovals$: Observable<Trace[]>;

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

  public view: string = 'services';
  public selectedStage: Stage = null;

  public eventTypes: string[] = [];
  public filterEventTypes: string[] = [];

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      if(params["shkeptncontext"]) {
        this.contextId = params["shkeptncontext"];
        this.apiService.getTraces(this.contextId)
          .pipe(
            map(response => response.body),
            map(result => result.events||[]),
            map(traces => traces.map(trace => Trace.fromJSON(trace)))
          )
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
        this.openApprovals$ = this.dataService.openApprovals;

        this._projectSub = this.project$.subscribe(project => {
          if(project === undefined)
            this.error = 'project';
          this._changeDetectorRef.markForCheck();
        }, error => {
          this.error = 'projects';
        });

        this._rootsSubs.unsubscribe();
        this._rootsSubs = this.dataService.roots.subscribe(roots => {
          if(roots) {
            if(!this.currentRoot)
              this.currentRoot = roots.find(r => r.shkeptncontext == params["contextId"]);
            this.eventTypes = this.eventTypes.concat(roots.map(r => r.type)).filter((r, i, a) => a.indexOf(r) === i);
          }
          if(this.currentRoot && !this.eventId)
            this.eventId = this.currentRoot.traces[this.currentRoot.traces.length-1].id;
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
            });
          });
      }
    });
  }

  selectRoot(event: any): void {
    this.projectName = event.root.getProject();
    this.serviceName = event.root.getService();
    this.contextId = event.root.data.shkeptncontext;
    this.eventId = null;
    if(event.stage) {
      let focusEvent = event.root.traces.find(trace => trace.data.stage == event.stage);
      let routeUrl = this.router.createUrlTree(['/project', focusEvent.getProject(), focusEvent.getService(), focusEvent.shkeptncontext, focusEvent.id]);
      this.eventId = focusEvent.id;
      this.location.go(routeUrl.toString());
    } else {
      let routeUrl = this.router.createUrlTree(['/project', event.root.getProject(), event.root.getService(), event.root.shkeptncontext]);
      this.eventId = event.root.traces[event.root.traces.length-1].id;
      this.location.go(routeUrl.toString());
    }

    this.currentRoot = event.root;
    this.loadTraces(this.currentRoot);
  }

  selectDeployment(deployment: Trace, project: Project) {
    this.selectRoot({
      root: project.getServices().find(service => service.serviceName === deployment.data.service).roots.find(root => root.shkeptncontext === deployment.shkeptncontext),
      stage: deployment.data.stage
    });
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

  selectView(view) {
    this.view = view;
  }

  filterEvents(event: DtCheckboxChange<string>, eventType: string): void {
    let index = this.filterEventTypes.indexOf(eventType);
    if(index == -1) {
      this.filterEventTypes.push(eventType);
    } else {
      this.filterEventTypes.splice(index, 1);
    }
  }

  isFilteredEvent(eventType: string) {
    return this.filterEventTypes.indexOf(eventType) == -1;
  }

  getEventLabel(key): string {
    return EVENT_LABELS[key] || key;
  }

  getFilteredRoots(roots: Root[]) {
    if(roots)
      return roots.filter(r => this.filterEventTypes.indexOf(r.type) == -1);
  }

  selectStage(stage) {
    this.selectedStage = stage;
  }

  countOpenApprovals(openApprovals: Trace[], stage: Stage, service?: Service) {
    return openApprovals.filter(approval => approval == openApprovals.find(a => a.data.stage == stage.stageName && a.data.service == approval.data.service && (!service || a.data.service == service.serviceName))).length;
  }

  getOpenApprovals(openApprovals: Trace[], stage: Stage, service: Service) {
    return openApprovals.filter(approval => approval.data.stage == stage.stageName && approval.data.service == service.serviceName);
  }

  approveDeployment(approval) {
    this.dataService.sendApprovalEvent(approval, true);
  }

  declineDeployment(approval) {
    this.dataService.sendApprovalEvent(approval, false);
  }

  ngOnDestroy(): void {
    this._projectSub.unsubscribe();
    this._routeSubs.unsubscribe();
    this._tracesTimer.unsubscribe();
    this._rootEventsTimer.unsubscribe();
  }

}
