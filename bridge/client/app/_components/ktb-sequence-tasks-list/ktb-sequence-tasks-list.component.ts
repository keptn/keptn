import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";
import {Router} from "@angular/router";
import {Location} from "@angular/common";
import {EventTypes} from "../../_models/event-types";

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
  get focusedEventId(): string {
    return this._focusedEventId;
  }
  set focusedEventId(value: string) {
    if (this._focusedEventId !== value) {
      this._focusedEventId = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private router: Router, private location: Location, private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

  identifyEvent(index, item) {
    return item ? item.time : null;
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats();
  }

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

  isInvalidated(event) {
    return !!this.tasks.find(e => e.isEvaluationInvalidation() && e.triggeredid == event.id);
  }
}
