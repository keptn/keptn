import { ChangeDetectorRef, Component, Input, OnDestroy, TemplateRef } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Subject } from 'rxjs';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { getSliResultInfo } from './ktb-evaluation-details-utils';
import { IEvaluationSelectionData } from './ktb-evaluation-chart/ktb-evaluation-chart.component';
import { DateUtil } from '../../_utils/date.utils';
import { IEvaluationData } from '../../../../shared/models/trace';

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss'],
})
export class KtbEvaluationDetailsComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public comparedIndicatorResults: IndicatorResult[][] = [];
  public sloDialogRef?: MatDialogRef<string>;
  public invalidateEvaluationDialogRef?: MatDialogRef<Trace | undefined>;
  public _evaluationState: Record<ResultTypes, string> = {
    [ResultTypes.PASSED]: 'recovered',
    [ResultTypes.WARNING]: 'warning',
    [ResultTypes.FAILED]: 'error',
  };
  public selectedEvaluation?: Trace;
  private _evaluationData: IEvaluationSelectionData = { shouldSelect: true };
  public getSliInfo = getSliResultInfo;

  @Input() public showChart = true;
  @Input() public isInvalidated = false;

  // just maps an evaluation to evaluationData
  @Input()
  set evaluation(evaluation: Trace | undefined) {
    this.evaluationData = { evaluation, shouldSelect: true };
  }
  get evaluation(): Trace | undefined {
    return this._evaluationData.evaluation;
  }

  @Input()
  set evaluationData(evaluationData: IEvaluationSelectionData) {
    this._evaluationData = evaluationData;
    this.selectedEvaluation = evaluationData.evaluation;
  }
  get evaluationData(): IEvaluationSelectionData {
    return this._evaluationData;
  }

  constructor(
    private _changeDetectorRef: ChangeDetectorRef,
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil
  ) {}

  public showSloDialog(evaluationData: IEvaluationData, sloDialog: TemplateRef<string>): void {
    this.sloDialogRef = this.dialog.open(sloDialog, {
      data: atob(evaluationData.sloFileContent),
    });
  }

  closeSloDialog(): void {
    this.sloDialogRef?.close();
  }

  copySloPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'slo payload');
  }

  invalidateEvaluationTrigger(invalidateEvaluationDialog: TemplateRef<Trace | undefined>): void {
    this.invalidateEvaluationDialogRef = this.dialog.open(invalidateEvaluationDialog, {
      data: this.evaluationData?.evaluation,
    });
  }

  invalidateEvaluation(evaluation: Trace, reason: string): void {
    this.dataService.invalidateEvaluation(evaluation, reason);
    this.closeInvalidateEvaluationDialog();
  }

  closeInvalidateEvaluationDialog(): void {
    this.invalidateEvaluationDialogRef?.close();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
