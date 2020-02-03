import {Component, OnDestroy, OnInit} from '@angular/core';
import {filter, map, startWith, switchMap} from "rxjs/operators";
import {Observable, Subscription, timer} from "rxjs";
import {ActivatedRoute} from "@angular/router";

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
  private _rootEventsTimer: Subscription = Subscription.EMPTY;
  private _rootEventsTimerInterval = 30;

  private _tracesTimer: Subscription = Subscription.EMPTY;
  private _tracesTimerInterval = 10;

  constructor(private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      if(params['projectName']) {
        this.currentRoot = null;

        this.project = this.dataService.projects.pipe(
          map(projects => projects.find(project => {
            return project.projectName === params['projectName'];
          }))
        );

        this._rootEventsTimer = timer(0, this._rootEventsTimerInterval*1000)
          .pipe(
            startWith(0),
            switchMap(() => this.project),
            filter(project => !!project && !!project.getServices())
          )
          .subscribe(project => {
            if(project && project.getServices()) {
              project.getServices().forEach(service => {
                this.dataService.loadRoots(project, service);
              });
            }
          });
      }
    });
  }

  loadTraces(root): void {
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

  ngOnDestroy(): void {
    this._routeSubs.unsubscribe();
    this._rootEventsTimer.unsubscribe();
    this._tracesTimer.unsubscribe();
  }

}
