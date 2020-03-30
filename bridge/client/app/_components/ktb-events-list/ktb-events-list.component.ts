import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit, ViewEncapsulation} from '@angular/core';
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";
import {Router} from "@angular/router";
import {Location} from "@angular/common";

@Component({
  selector: 'ktb-events-list',
  templateUrl: './ktb-events-list.component.html',
  styleUrls: ['./ktb-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbEventsListComponent implements OnInit {

  public _events: Trace[] = [];
  public _focusedEventId: string;

  @Input()
  get events(): Trace[] {
    return this._events;
  }
  set events(value: Trace[]) {
    if (this._events !== value) {
      this._events = value;
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
    let routeUrl = this.router.createUrlTree(['/project', event.data.project, event.data.service, event.shkeptncontext, event.id]);
    this.location.go(routeUrl.toString());
  }

}
