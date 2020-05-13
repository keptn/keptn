import * as moment from 'moment';

import {ChangeDetectorRef, Component, Directive, Input, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";
import {MatDialog, MatDialogRef} from "@angular/material/dialog";

@Directive({
  selector: `ktb-event-item-detail, [ktb-event-item-detail], [ktbEventItemDetail]`,
  exportAs: 'ktbEventItemDetail',
})
export class KtbEventItemDetail {}

@Component({
  selector: 'ktb-event-item',
  templateUrl: './ktb-event-item.component.html',
  styleUrls: ['./ktb-event-item.component.scss']
})
export class KtbEventItemComponent implements OnInit {

  public _event: Trace;

  @ViewChild("eventPayloadDialog", {static: false})
  public eventPayloadDialog: TemplateRef<any>;
  public eventPayloadDialogRef: MatDialogRef<any, any>;

  @Input()
  get event(): Trace {
    return this._event;
  }
  set event(value: Trace) {
    if (this._event !== value) {
      this._event = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dialog: MatDialog) { }

  ngOnInit() {
  }

  getCalendarFormat() {
    return DateUtil.getCalendarFormats().sameElse;
  }

  showEventPayloadDialog() {
    this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog, { data: this._event.plainEvent });
  }

  closeEventPayloadDialog() {
    if(this.eventPayloadDialogRef)
      this.eventPayloadDialogRef.close();
  }

  getDuration(start, end) {
    return DateUtil.getDurationFormatted(start, end);
  }

}
