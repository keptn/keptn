import { Component, Directive, EventEmitter, Input, Output, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Trace } from '../../_models/trace';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';

@Directive({
  selector: `ktb-event-item-detail, [ktb-event-item-detail], [ktbEventItemDetail]`,
  exportAs: 'ktbEventItemDetail',
})
export class KtbEventItemDetailDirective {}

@Component({
  selector: 'ktb-event-item[event]',
  templateUrl: './ktb-event-item.component.html',
  styleUrls: ['./ktb-event-item.component.scss'],
})
export class KtbEventItemComponent {
  public _event?: Trace;

  @ViewChild('eventPayloadDialog')
  /* eslint-disable @typescript-eslint/no-explicit-any */
  public eventPayloadDialog?: TemplateRef<any>;
  public eventPayloadDialogRef?: MatDialogRef<any, any>;
  /* eslint-enable @typescript-eslint/no-explicit-any */

  @Input() public showChartLink = false;
  @Input() public showTime = true;
  @Input() public showLabels = true;
  @Input() public image?: string;
  @Output() public approvalSent: EventEmitter<void> = new EventEmitter<void>();

  @Input()
  get event(): Trace | undefined {
    return this._event;
  }

  set event(value: Trace | undefined) {
    if (this._event !== value) {
      this._event = value;
    }
  }

  constructor(
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil
  ) {}

  showEventPayloadDialog(): void {
    if (this.eventPayloadDialog && this._event) {
      this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog, { data: this._event.plainEvent });
    }
  }

  closeEventPayloadDialog(): void {
    if (this.eventPayloadDialogRef) {
      this.eventPayloadDialogRef.close();
    }
  }

  copyEventPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'event payload');
  }

  isUrl(value: string): boolean {
    try {
      // tslint:disable-next-line:no-unused-expression
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }
}
