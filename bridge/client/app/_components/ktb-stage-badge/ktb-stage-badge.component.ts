import {Component, Input} from '@angular/core';

import {Trace} from '../../_models/trace';
import {Stage} from '../../_models/stage';
import {EvaluationResult} from '../../../../shared/interfaces/evaluation-result';

@Component({
  selector: 'ktb-stage-badge',
  templateUrl: './ktb-stage-badge.component.html',
  styleUrls: ['./ktb-stage-badge.component.scss'],
})
export class KtbStageBadgeComponent {

  @Input() public evaluation?: Trace;
  @Input() public evaluationResult?: EvaluationResult;
  @Input() public stage?: Stage;
  @Input() public isSelected?: boolean = undefined;
  @Input() public fill?: boolean = undefined;
  @Input() public error = false;
  @Input() public warning = false;
  @Input() public success = false;
  @Input() public highlight = false;
}
