import {Component, Input} from '@angular/core';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';

import {Trace} from '../../_models/trace';
import {ResultTypes} from '../../_models/result-types';
import {EvaluationResult} from '../../_models/evaluation-result';

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss']
})
export class KtbEvaluationInfoComponent {
  private _evaluationResult?: EvaluationResult;
  public isError = false;
  public isWarning = false;
  public isSuccess = false;
  @Input()
  get evaluationResult(): EvaluationResult | undefined {
    return this._evaluationResult;
  }
  set evaluationResult(result: EvaluationResult | undefined) {
    this._evaluationResult = result;
    if (this.evaluationResult) {
      this.isError = this.evaluationResult.result === ResultTypes.FAILED;
      this.isWarning = this.evaluationResult.result === ResultTypes.WARNING;
      this.isSuccess = this.evaluationResult.result === ResultTypes.PASSED;
    }
  }
  @Input() public evaluation?: Trace;
  @Input() public overlayDisabled?: boolean;
  @Input() public fill?: boolean;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };
}
