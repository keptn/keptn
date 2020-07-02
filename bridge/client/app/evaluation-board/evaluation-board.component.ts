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
import {labels, Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {DtCheckboxChange} from "@dynatrace/barista-components/checkbox";

@Component({
  selector: 'app-project-board',
  templateUrl: './evaluation-board.component.html',
  styleUrls: ['./evaluation-board.component.scss']
})
export class EvaluationBoardComponent implements OnInit, OnDestroy {

  public _routeSubs: Subscription = Subscription.EMPTY;
  public error: string = null;

  public contextId: string;
  public root: Root;
  public evaluations: Trace[];

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService) { }

  ngOnInit() {
    this._routeSubs = this.route.params.subscribe(params => {
      if(params["shkeptncontext"]) {
        this.contextId = params["shkeptncontext"];
        this.apiService.getTraces(this.contextId)
          .pipe(
            map(response => response.body),
            map(result => result.events||[]),
            map(traces => traces.map(trace => Trace.fromJSON(trace)).sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime()))
          )
          .subscribe((traces: Trace[]) => {
            if(traces.length > 0) {
              console.log("traces", traces);
              this.root = Root.fromJSON(traces[0]);
              this.root.traces = traces;
              this.evaluations = traces.filter(t => t.type == 'sh.keptn.events.evaluation-done' && (!params["eventselector"] || t.id == params["eventselector"] || t.data.stage == params["eventselector"])) ;
            } else {
              this.error = "contextError";
            }
          }, (err) => {
            this.error = "error";
          });
      }
    });
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats(true);
  }

  getEventLabel(key): string {
    return labels[key] || key;
  }

  ngOnDestroy(): void {
    this._routeSubs.unsubscribe();
  }

}
