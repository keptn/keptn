import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, Output } from '@angular/core';
import { Trace } from '../../_models/trace';
import { DataService } from '../../_services/data.service';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { EventTypes } from '../../../../shared/interfaces/event-types';
import { KeptnService } from '../../../../shared/models/keptn-service';

@Component({
  selector: 'ktb-approval-item[event]',
  templateUrl: './ktb-approval-item.component.html',
  styleUrls: ['./ktb-approval-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbApprovalItemComponent {
  public _event?: Trace;
  public approvalResult?: boolean;
  public evaluation?: Trace;
  public evaluationExists = true;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };

  @Output() approvalSent: EventEmitter<void> = new EventEmitter<void>();

  @Input()
  get event(): Trace | undefined {
    return this._event;
  }

  set event(value: Trace | undefined) {
    if (this._event === value) {
      return;
    }
    this.evaluation = undefined;
    this._event = value;
    if (!value) {
      return;
    }
    this.loadEvaluation(value);
  }

  constructor(private dataService: DataService, private changeDetectorRef_: ChangeDetectorRef) {}

  private loadEvaluation(trace: Trace): void {
    this.dataService
      .getTracesByContext(
        trace.shkeptncontext,
        EventTypes.EVALUATION_FINISHED,
        KeptnService.LIGHTHOUSE_SERVICE,
        trace.stage,
        1
      )
      .subscribe((evaluation) => {
        this.evaluation = evaluation[0];
        this.evaluationExists = !!this.evaluation;
        this.changeDetectorRef_.markForCheck();
      });
  }

  public handleApproval(approval: Trace, result: boolean): void {
    this.dataService.sendApprovalEvent(approval, result).subscribe(() => {
      this.approvalSent.emit();
    });
    this.approvalResult = result;
    this.changeDetectorRef_.markForCheck();
  }
}
