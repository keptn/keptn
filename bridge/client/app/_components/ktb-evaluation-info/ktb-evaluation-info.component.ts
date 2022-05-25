import { Component, Input, NgZone, OnDestroy, TemplateRef, ViewChild } from '@angular/core';
import { DtOverlay, DtOverlayConfig, DtOverlayRef } from '@dynatrace/barista-components/overlay';
import { Trace } from '../../_models/trace';
import { ResultTypes } from '../../../../shared/models/result-types';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { DataService } from '../../_services/data.service';
import { Subject, Subscription } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { DateUtil } from '../../_utils/date.utils';

export interface EventData {
  project: string;
  stage: string;
  service: string;
}

interface EvaluationInfo {
  trace: Trace | undefined;
  showHistory: boolean;
  data: EventData;
}

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss'],
})
export class KtbEvaluationInfoComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private _evaluationResult?: EvaluationResult;
  private eventData?: EventData;
  private _evaluationHistory?: Trace[];
  public readonly evaluationHistoryCount = 5;
  public TraceClass = Trace;
  public isError = false;
  public isWarning = false;
  public isSuccess = false;
  public showHistory = false;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };
  public evaluationsLoaded = false;
  private overlayRef?: DtOverlayRef<unknown>;
  private updateOverlayPositionSubscription = Subscription.EMPTY;

  @ViewChild('overlay', { static: true, read: TemplateRef }) overlayTemplate?: TemplateRef<unknown>;
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
    this.eventData = evaluation?.data;

    if (idBefore !== evaluation?.trace?.id || (!evaluation?.trace && evaluation?.data)) {
      this.fetchEvaluationHistory();
    }
  }

  get evaluationHistory(): Trace[] {
    return (
      this._evaluationHistory ||
      this.evaluation?.data?.evaluationHistory
        ?.filter((evaluation) => evaluation.id !== this.evaluation?.id)
        .slice(0, this.evaluationHistoryCount) ||
      []
    );
  }

  constructor(private dataService: DataService, private ngZone: NgZone, private _dtOverlay: DtOverlay) {}

  private fetchEvaluationHistory(): void {
    const evaluation = this.evaluation;
    let _eventData = this.eventData;
    if (this.evaluation && this.evaluation.data.project && this.evaluation.data.stage && this.evaluation.data.service) {
      _eventData = {
        project: this.evaluation.data.project,
        service: this.evaluation.data.service,
        stage: this.evaluation.data.stage,
      };
    }

    if (this.showHistory && _eventData) {
      // currently the event endpoint does not support skipping entries
      // the other endpoint we have does not support excluding invalidated evaluations
      // we can't use fromTime here if we have a limit. 10 new evaluations and limit to 5 would not pull the new ones
      this.dataService
        .getEvaluationResults(_eventData, this.evaluationHistoryCount + (this.evaluation ? 1 : 0), false)
        .subscribe((traces: Trace[]) => {
          traces.sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
          this.evaluationsLoaded = true;
          // we don't have an evaluation trace if the sequence is currently running or it just doesn't have an evaluation task
          if (evaluation) {
            this._evaluationHistory = undefined;
            evaluation.data.evaluationHistory = traces;
          } else {
            this._evaluationHistory = traces;
          }
        });
    }
  }

  public showEvaluationOverlay(event: MouseEvent, data?: unknown): void {
    if (!this.overlayRef && !this.overlayDisabled && this.overlayTemplate) {
      this.overlayRef = this._dtOverlay.create(event, this.overlayTemplate, { ...this.overlayConfig, data });
      this.updateEvaluationOverlayPosition();
    }
  }

  public updateEvaluationOverlayPosition(): void {
    this.updateOverlayPositionSubscription = this.ngZone.onMicrotaskEmpty
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.overlayRef?.updatePosition();
        // if the content of the overlay changed after initialization the position stayed the same
      });
  }

  public hideEvaluationOverlay(): void {
    if (this.overlayRef) {
      this._dtOverlay.dismiss();
      this.updateOverlayPositionSubscription.unsubscribe();
      this.overlayRef = undefined;
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
