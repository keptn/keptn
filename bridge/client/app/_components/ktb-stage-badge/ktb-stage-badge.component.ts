import { Component, Input } from '@angular/core';
import { Trace } from '../../_models/trace';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { EvaluationBadgeVariant } from '../ktb-evaluation-badge/ktb-evaluation-badge.utils';

@Component({
  selector: 'ktb-stage-badge',
  templateUrl: './ktb-stage-badge.component.html',
  styleUrls: ['./ktb-stage-badge.component.scss'],
})
export class KtbStageBadgeComponent {
  @Input() public evaluation?: Trace;
  @Input() public evaluationResult?: EvaluationResult;
  @Input() public stage?: string;
  @Input() public isSelected?: boolean = undefined;
  @Input() public fillState = EvaluationBadgeVariant.FILL;
  @Input() public error = false;
  @Input() public warning = false;
  @Input() public success = false;
  @Input() public highlight = false;
  @Input() public aborted = false;
}
