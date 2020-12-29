import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {map, takeUntil} from "rxjs/operators";
import {Subject, Subscription} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import {Location} from "@angular/common";

import {Root} from "../_models/root";

import {DataService} from "../_services/data.service";
import {ApiService} from "../_services/api.service";
import DateUtil from "../_utils/date.utils";
import {Trace} from "../_models/trace";
import {EventTypes} from "../_models/event-types";
import {EVENT_LABELS} from "../_models/event-labels";

@Component({
  selector: 'app-project-board',
  templateUrl: './evaluation-board.component.html',
  styleUrls: ['./evaluation-board.component.scss']
})
export class EvaluationBoardComponent implements OnInit, OnDestroy {

  private unsubscribe$ = new Subject();

  public error: string = null;
  public contextId: string;
  public root: Root;
  public evaluations: Trace[];

  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: Router, private location: Location, private route: ActivatedRoute, private dataService: DataService, private apiService: ApiService) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        if(params["shkeptncontext"]) {
          this.contextId = params["shkeptncontext"];
          this.apiService.getTraces(this.contextId)
            .pipe(
              map(response => response.body),
              map(result => result.events||[]),
              map(traces => traces.map(trace => Trace.fromJSON(trace)).sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime()))
            )
            .pipe(takeUntil(this.unsubscribe$))
            .subscribe((traces: Trace[]) => {
              if(traces.length > 0) {
                this.root = Root.fromJSON(traces[0]);
                this.root.traces = traces;
                this.evaluations = traces.filter(t => t.type == EventTypes.EVALUATION_FINISHED && (!params["eventselector"] || t.id == params["eventselector"] || t.data.stage == params["eventselector"])) ;
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
    return EVENT_LABELS[key] || key;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
