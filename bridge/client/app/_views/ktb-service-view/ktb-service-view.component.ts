import {ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, ViewEncapsulation} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";

import {Observable, Subject, Subscription, timer} from "rxjs";
import {filter, startWith, switchMap, take, takeUntil} from "rxjs/operators";

import * as moment from "moment";

import {Root} from "../../_models/root";
import {Project} from "../../_models/project";

import {DataService} from "../../_services/data.service";
import {DateUtil} from "../../_utils/date.utils";
import {DtCheckboxChange} from "@dynatrace/barista-components/checkbox";
import {Location} from "@angular/common";

@Component({
  selector: 'ktb-service-view',
  templateUrl: './ktb-service-view.component.html',
  styleUrls: ['./ktb-service-view.component.scss'],
  host: {
    class: 'ktb-service-view'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbServiceViewComponent implements OnInit {

  private readonly unsubscribe$ = new Subject<void>();
  public project$: Observable<Project>;

  public currentRoot: Root;

  public projectName: string;
  public serviceName: string;
  public contextId: string;
  public eventId: string;

  public selectedStage: string;

  public eventTypes: string[] = [];
  public filterEventTypes: string[] = [];

  private _rootEventsTimerInterval = 30;
  private _tracesTimerInterval = 10;
  private _tracesTimer: Subscription = Subscription.EMPTY;

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private router: Router, private location: Location, private route: ActivatedRoute, public dateUtil: DateUtil) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.projectName = params["projectName"];
        this.serviceName = params["serviceName"];
        this.contextId = params["contextId"];
        this.eventId = params["eventId"];
        this.currentRoot = null;
        this.filterEventTypes = [];

        this.project$ = this.dataService.getProject(params['projectName']);

        this.project$
          .pipe(
            filter(project => !!project && !!project.getServices() && !!project.stages && !!project.sequences),
            take(1)
          )
          .subscribe(project => {
            this.currentRoot = null;
          });

        timer(0, this._rootEventsTimerInterval*1000)
          .pipe(
            startWith(0),
            switchMap(() => this.project$),
            filter(project => !!project && !!project.getServices())
          )
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(project => {
            this.dataService.loadServices(project);
            this.dataService.loadRoots(project);
          });

        this.dataService.roots
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(roots => {
            if(roots) {
              if(!this.currentRoot)
                this.currentRoot = roots.find(r => r.shkeptncontext == params["contextId"]);
              this.eventTypes = this.eventTypes.concat(roots.map(root => root.getLabel())).filter((eventType, i, eventTypes) => eventTypes.indexOf(eventType) === i);
            }
            if(this.currentRoot && !this.eventId)
              this.eventId = this.currentRoot.traces[this.currentRoot.traces.length-1].id;

            this._changeDetectorRef.markForCheck();
          });
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

  getFilteredRoots(roots: Root[]) {
    if(roots)
      return roots.filter(r => this.filterEventTypes.indexOf(r.type) == -1);
  }

  getRootsLastUpdated(project: Project): Date {
    return this.dataService.getRootsLastUpdated(project);
  }

  getTracesLastUpdated(root: Root): Date {
    return this.dataService.getTracesLastUpdated(root);
  }

  showReloadButton(root: Root) {
    return moment().subtract(1, 'day').isAfter(root.time);
  }

  selectStage(stageName: string) {
    this.selectedStage = stageName;
    this._changeDetectorRef.markForCheck();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this._tracesTimer.unsubscribe();
  }
}
