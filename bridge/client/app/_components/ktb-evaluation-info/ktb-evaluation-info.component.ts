import { Component, Inject, Input, OnDestroy } from '@angular/core';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';

import { Trace } from '../../_models/trace';
import { ResultTypes } from '../../../../shared/models/result-types';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { DataService } from '../../_services/data.service';
import { Subject } from 'rxjs/internal/Subject';
import { takeUntil } from 'rxjs/internal/operators/takeUntil';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DateUtil } from '../../_utils/date.utils';

interface EvaluationInfo {
  trace: Trace | undefined;
  showHistory: boolean;
}

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss'],
})
export class KtbEvaluationInfoComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public readonly evaluationHistoryCount = 5;
  private _evaluationResult?: EvaluationResult;
  public isError = false;
  public isWarning = false;
  public isSuccess = false;
  public showHistory = false;
  public evaluationHistory?: Trace[];
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };
  public historyPolling$: Subscription = Subscription.EMPTY;
  public evaluationsLoaded = false;

  @Input() public overlayDisabled?: boolean;
  @Input() public fill?: boolean;
  @Input() public evaluation?: Trace;
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

  @Input()
  public set evaluationInfo(evaluation: EvaluationInfo | undefined) {
    const idBefore = this.evaluation?.id;
    this.evaluation = evaluation?.trace;
    this.evaluationsLoaded = !!evaluation?.trace?.data.evaluationHistory?.length;
    this.showHistory = evaluation?.showHistory ?? false;

    if (idBefore !== evaluation?.trace?.id) {
      this.fetchEvaluationHistory();
    }
  }

  constructor(private dataService: DataService, @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number) {}

  private fetchEvaluationHistory(): void {
    this.historyPolling$.unsubscribe();
    const evaluation = this.evaluation;
    if (this.showHistory && evaluation) {
      this.historyPolling$ = AppUtils.createTimer(0, this.initialDelayMillis)
        .pipe(
          takeUntil(this.unsubscribe$),
          switchMap(() => {
            // currently the event endpoint does not support skipping entries
            // the other endpoint we have does not support excluding invalidated evaluations
            // we can't use fromTime here if we have a limit. 10 new evaluations and limit to 5 would not pull the new ones
            return this.dataService.getEvaluationResults(evaluation, this.evaluationHistoryCount + 1, false);
          })
        )
        .subscribe((traces: Trace[]) => {
          this.evaluationsLoaded = true;
          //TODO: use another place to save the evaluations or change the implementation inside the evaluation-details
          evaluation.data.evaluationHistory = traces.slice(1).sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
        });
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
