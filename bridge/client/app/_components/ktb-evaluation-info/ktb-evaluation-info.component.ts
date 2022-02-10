import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';

import { Trace } from '../../_models/trace';
import { ResultTypes } from '../../../../shared/models/result-types';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Subject } from 'rxjs/internal/Subject';
import { takeUntil } from 'rxjs/internal/operators/takeUntil';

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss'],
})
export class KtbEvaluationInfoComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
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
  @Input() public showHistory?: boolean;

  public evaluationHistory?: Trace[];

  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    if (this.showHistory && this.evaluation && !this.evaluation.data.evaluationHistory) {
      this.dataService.evaluationResults.pipe(takeUntil(this.unsubscribe$)).subscribe((results) => {
        if (results.type === 'evaluationHistory' && this.evaluation && results.triggerEvent === this.evaluation) {
          this.evaluation.data.evaluationHistory = [
            ...(results.traces || []),
            ...(this.evaluation.data.evaluationHistory || []),
          ].sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
        }
      });
      this.dataService.loadEvaluationResults(this.evaluation, 5);
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
