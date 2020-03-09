import {Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, startWith, switchMap} from "rxjs/operators";
import {Observable, Subscription, timer} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import {Location} from "@angular/common";

import * as moment from 'moment';

import {Root} from "../_models/root";
import {Project} from "../_models/project";

import {DataService} from "../_services/data.service";
import DateUtil from "../_utils/date.utils";
import {Service} from "../_models/service";

@Component({
  selector: 'app-project-board',
  templateUrl: './project-board.component.html',
  styleUrls: ['./project-board.component.scss']
})
export class ProjectBoardComponent implements OnInit, OnDestroy {

  public project: Observable<Project>;
  public currentRoot: Root;
  public error: boolean = false;

  private _routeSubs: Subscription = Subscription.EMPTY;
  private _rootsSubs: Subscription = Subscription.EMPTY;
  private _rootEventsTimer: Subscription = Subscription.EMPTY;
  private _rootEventsTimerInterval = 30;

  private _tracesTimer: Subscription = Subscription.EMPTY;
  private _tracesTimerInterval = 10;

  public projectName: string;
  public serviceName: string;
  public contextId: string;

  constructor(private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      this.projectName = params["projectName"];
      this.serviceName = params["serviceName"];
      this.contextId = params["contextId"];
      this.currentRoot = null;

      this.project = this.dataService.projects.pipe(
        map(projects => projects ? projects.find(project => {
          return project.projectName === params['projectName'];
        }) : null)
      );

      this._rootsSubs.unsubscribe();
      this._rootsSubs = this.dataService.roots.subscribe(roots => {
        if(roots && !this.currentRoot)
          this.currentRoot = roots.find(r => r.shkeptncontext == params["contextId"]);
      });

      this._rootEventsTimer.unsubscribe();
      this._rootEventsTimer = timer(0, this._rootEventsTimerInterval*1000)
        .pipe(
          startWith(0),
          switchMap(() => this.project),
          filter(project => !!project && !!project.getServices())
        )
        .subscribe(project => {
          project.getServices().forEach(service => {
            this.dataService.loadRoots(project, service);
            if(service.roots && !this.currentRoot)
              this.currentRoot = service.roots.find(r => r.shkeptncontext == params["contextId"]);
          });
        });
    });
  }

  loadTraces(root: Root): void {
    let routeUrl = this.router.createUrlTree(['/project', this.projectName, root.data.service, root.shkeptncontext]);
    this.location.go(routeUrl.toString());

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

  ngOnDestroy(): void {
    this._routeSubs.unsubscribe();
    this._rootEventsTimer.unsubscribe();
    this._tracesTimer.unsubscribe();
  }

}
