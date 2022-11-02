import { Component, Input, TemplateRef } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { getSliResultInfo, IEvaluationSelectionData } from './ktb-evaluation-details-utils';
import { DateUtil } from '../../_utils/date.utils';

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss'],
})
export class KtbEvaluationDetailsComponent {
  private sloDialogRef?: MatDialogRef<string>;
  public comparedIndicatorResults?: IndicatorResult[][];
  public invalidateEvaluationDialogRef?: MatDialogRef<Trace | undefined>;
  public _evaluationState: Record<ResultTypes, string> = {
    [ResultTypes.PASSED]: 'recovered',
    [ResultTypes.WARNING]: 'warning',
    [ResultTypes.FAILED]: 'error',
    [ResultTypes.INFO]: 'info',
  };
  public selectedEvaluation?: Trace;
  private _evaluationData: IEvaluationSelectionData = { shouldSelect: true };
  public getSliInfo = getSliResultInfo;
  public numberOfMissingComparisonResults = 0;

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
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil
  ) {}

  public showSloDialog(sloFileContent: string, sloDialog: TemplateRef<string>): void {
    this.sloDialogRef = this.dialog.open(sloDialog, {
      data: atob(sloFileContent),
    });
  }

  public closeSloDialog(): void {
    this.sloDialogRef?.close();
  }

  public copySloPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'slo payload');
  }

  public invalidateEvaluationTrigger(invalidateEvaluationDialog: TemplateRef<Trace | undefined>): void {
    this.invalidateEvaluationDialogRef = this.dialog.open(invalidateEvaluationDialog, {
      data: this.evaluationData?.evaluation,
    });
  }

  public invalidateEvaluation(evaluation: Trace, reason: string): void {
    this.dataService.invalidateEvaluation(evaluation, reason);
    this.closeInvalidateEvaluationDialog();
  }

  public closeInvalidateEvaluationDialog(): void {
    this.invalidateEvaluationDialogRef?.close();
  }
}
