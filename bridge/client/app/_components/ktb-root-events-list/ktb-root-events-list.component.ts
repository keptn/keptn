import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnInit, Output,
  ViewEncapsulation
} from '@angular/core';
import {Root} from "../../_models/root";
import DateUtil from "../../_utils/date.utils";
import {Stage} from "../../_models/stage";

@Component({
  selector: 'ktb-root-events-list',
  templateUrl: './ktb-root-events-list.component.html',
  styleUrls: ['./ktb-root-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbRootEventsListComponent implements OnInit {

  public _events: Root[] = [];
  public _selectedEvent: Root = null;

  @Output() readonly selectedEventChange = new EventEmitter<any>();

  @Input()
  get events(): Root[] {
    return this._events;
  }
  set events(value: Root[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): Root {
    return this._selectedEvent;
  }
  set selectedEvent(value: Root) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

  selectEvent(event: Root, stage?: String) {
    this.selectedEvent = event;
    this.selectedEventChange.emit({ root: event, stage: stage });
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats();
  }

  identifyEvent(index, item) {
    return item ? item.time : null;
  }

  getShortImageName(image) {
    let parts = image.split("/");
    return parts[parts.length-1];
  }
}
