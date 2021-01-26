import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
<<<<<<< HEAD
import {Router} from "@angular/router";
import {Location} from "@angular/common";

import {Trace} from "../../_models/trace";
import {DateUtil} from "../../_utils/date.utils";
=======
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";
import {Router} from "@angular/router";
import {Location} from "@angular/common";
import {EventTypes} from "../../_models/event-types";
>>>>>>> create component for sequence overview, sequence list and task details

@Component({
  selector: 'ktb-sequence-tasks-list',
  templateUrl: './ktb-sequence-tasks-list.component.html',
  styleUrls: ['./ktb-sequence-tasks-list.component.scss'],
  host: {
    class: 'ktb-sequence-tasks-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceTasksListComponent implements OnInit {

  public _tasks: Trace[] = [];
<<<<<<< HEAD
  public _stage: String;
=======
>>>>>>> create component for sequence overview, sequence list and task details
  public _focusedEventId: string;

  @Input()
  get tasks(): Trace[] {
    return this._tasks;
  }
  set tasks(value: Trace[]) {
    if (this._tasks !== value) {
      this._tasks = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
<<<<<<< HEAD
  get stage(): String {
    return this._stage;
  }
  set stage(value: String) {
    if (this._stage !== value) {
      this._stage = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
=======
>>>>>>> create component for sequence overview, sequence list and task details
  get focusedEventId(): string {
    return this._focusedEventId;
  }
  set focusedEventId(value: string) {
    if (this._focusedEventId !== value) {
      this._focusedEventId = value;
      this._changeDetectorRef.markForCheck();
    }
  }

<<<<<<< HEAD
  constructor(private router: Router, private location: Location, private _changeDetectorRef: ChangeDetectorRef, public dateUtil: DateUtil) { }
=======
  constructor(private router: Router, private location: Location, private _changeDetectorRef: ChangeDetectorRef) { }
>>>>>>> create component for sequence overview, sequence list and task details

  ngOnInit() {
  }

  identifyEvent(index, item) {
    return item ? item.time : null;
  }

<<<<<<< HEAD
=======
  getCalendarFormats() {
    return DateUtil.getCalendarFormats();
  }

>>>>>>> create component for sequence overview, sequence list and task details
  private currentScrollElement;
  scrollIntoView(element) {
    if(element != this.currentScrollElement) {
      this.currentScrollElement = element;
      setTimeout(() => {
        element.scrollIntoView({ behavior: 'smooth' });
      }, 0);
    }
    return true;
  }

  focusEvent(event) {
    if(event.getProject() && event.getService()) {
      let routeUrl = this.router.createUrlTree(['/project', event.getProject(), event.getService(), event.shkeptncontext, event.id]);
      this.location.go(routeUrl.toString());
    }
  }

<<<<<<< HEAD
  getTasksByStage(tasks: Trace[], stage: String) {
    return tasks.filter(t => t.data?.stage == stage);
  }

=======
>>>>>>> create component for sequence overview, sequence list and task details
  isInvalidated(event) {
    return !!this.tasks.find(e => e.isEvaluationInvalidation() && e.triggeredid == event.id);
  }
}
