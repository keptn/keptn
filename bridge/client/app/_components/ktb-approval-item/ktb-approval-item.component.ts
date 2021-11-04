import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Trace } from '../../_models/trace';
import { DataService } from '../../_services/data.service';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';

@Component({
  selector: 'ktb-approval-item[event][evaluation]',
  templateUrl: './ktb-approval-item.component.html',
  styleUrls: ['./ktb-approval-item.component.scss'],
})
export class KtbApprovalItemComponent {
  public _event?: Trace;
  public approvalResult?: boolean;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };

  @Input() evaluation?: Trace;
  @Output() approvalSent: EventEmitter<void> = new EventEmitter<void>();

  @Input()
  get event(): Trace | undefined {
    return this._event;
  }

  set event(value: Trace | undefined) {
    if (this._event !== value) {
      this._event = value;
    }
  }

  constructor(private dataService: DataService) {}

  public handleApproval(approval: Trace, result: boolean): void {
    this.dataService.sendApprovalEvent(approval, result).subscribe(() => {
      this.approvalSent.emit();
    });
    this.approvalResult = result;
  }
}
