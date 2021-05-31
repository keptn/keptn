import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';

import {Trace} from '../../_models/trace';
import {ResultTypes} from '../../_models/result-types';

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss']
})
export class KtbEvaluationInfoComponent implements OnInit, OnDestroy {
  private _evaluationResult: {result: ResultTypes, score: number};
  public isError: boolean;
  public isWarning: boolean;
  public isSuccess: boolean;
  @Input()
  get evaluationResult() {
    return this._evaluationResult;
  }
  set evaluationResult(result: {result: ResultTypes, score: number}) {
    this._evaluationResult = result;
    if (this.evaluationResult) {
      this.isError = this.evaluationResult.result === ResultTypes.FAILED;
      this.isWarning = this.evaluationResult.result === ResultTypes.WARNING;
      this.isSuccess = this.evaluationResult.result === ResultTypes.PASSED;
    }
  }
  @Input() public evaluation: Trace;
  @Input() public overlayDisabled: boolean;
  @Input() public fill: boolean | undefined;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  constructor() { }

  ngOnInit() {
  }

  ngOnDestroy(): void {
  }

}
