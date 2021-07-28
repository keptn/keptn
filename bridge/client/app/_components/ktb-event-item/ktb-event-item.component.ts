import {ChangeDetectorRef, Component, Directive, Input, TemplateRef, ViewChild} from '@angular/core';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import { Observable, of } from 'rxjs';
import {Project} from '../../_models/project';
import {Trace} from '../../_models/trace';
import {ClipboardService} from '../../_services/clipboard.service';
import {DataService} from '../../_services/data.service';
import {DateUtil} from '../../_utils/date.utils';

@Directive({
  selector: `ktb-event-item-detail, [ktb-event-item-detail], [ktbEventItemDetail]`,
  exportAs: 'ktbEventItemDetail',
})
export class KtbEventItemDetail {
}

@Component({
  selector: 'ktb-event-item',
  templateUrl: './ktb-event-item.component.html',
  styleUrls: ['./ktb-event-item.component.scss'],
})
export class KtbEventItemComponent {

  public project$: Observable<Project | undefined> = of(undefined);
  public _event?: Trace;

  @ViewChild('eventPayloadDialog')
  // tslint:disable-next-line:no-any
  public eventPayloadDialog?: TemplateRef<any>;
  // tslint:disable-next-line:no-any
  public eventPayloadDialogRef?: MatDialogRef<any, any>;

  @Input() public showChartLink = false;
  @Input() public showTime = true;
  @Input() public showLabels = true;

  @Input()
  get event(): Trace | undefined {
    return this._event;
  }

  set event(value: Trace | undefined) {
    if (this._event !== value) {
      this._event = value;
      if (this._event?.project) {
        this.project$ = this.dataService.getProject(this._event?.project);
      }
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef,
              private dataService: DataService,
              private dialog: MatDialog,
              private clipboard: ClipboardService,
              public dateUtil: DateUtil) {
  }

  showEventPayloadDialog() {
    if (this.eventPayloadDialog && this._event) {
      this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog, {data: this._event.plainEvent});
    }
  }

  closeEventPayloadDialog() {
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
