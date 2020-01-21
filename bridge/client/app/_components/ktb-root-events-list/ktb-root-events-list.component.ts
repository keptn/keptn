import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnInit, Output,
  ViewEncapsulation
} from '@angular/core';
import {coerceArray} from "@angular/cdk/coercion";
import {Root} from "../../_models/root";

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

  @Output() readonly selectedEventChange = new EventEmitter<Root>();

  @Input()
  get events(): Root[] {
    return this._events;
  }
  set events(value: Root[]) {
    const newValue = coerceArray(value);
    if (this._events !== newValue) {
      this._events = newValue;
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

  selectEvent(event: Root) {
    this.selectedEvent = event;
    this.selectedEventChange.emit(this.selectedEvent);
  }

}
