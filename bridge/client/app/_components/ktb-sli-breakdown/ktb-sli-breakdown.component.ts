import { ChangeDetectorRef, Component, Input, OnInit, ViewChild } from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';
import { SliResult } from '../../_models/sli-result';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { AppUtils } from '../../_utils/app.utils';
import { SloConfig } from '../../../../shared/interfaces/slo-config';

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss'],
})
export class KtbSliBreakdownComponent implements OnInit {
  @ViewChild('sortable', { read: DtSort, static: true }) sortable?: DtSort;

  public evaluationState: Map<ResultTypes, string> = new Map<ResultTypes, string>([
    [ResultTypes.PASSED, 'passed'],
    [ResultTypes.WARNING, 'warning'],
    [ResultTypes.FAILED, 'failed'],
  ]);
  public ResultTypes: typeof ResultTypes = ResultTypes;
  private _indicatorResults?: IndicatorResult[];
  private _indicatorResultsFail: IndicatorResult[] = [];
  private _indicatorResultsWarning: IndicatorResult[] = [];
  private _indicatorResultsPass: IndicatorResult[] = [];
  private _score = 0;
  public columnNames: string[] = [];
  public tableEntries: DtTableDataSource<SliResult> = new DtTableDataSource();
  public readonly SliResultClass = SliResult;
  private _comparedIndicatorResults: IndicatorResult[][] = [];

  @Input()
  get indicatorResults(): IndicatorResult[] {
    return [...this._indicatorResultsFail, ...this._indicatorResultsWarning, ...this._indicatorResultsPass];
  }
  set indicatorResults(indicatorResults: IndicatorResult[]) {
    if (this._indicatorResults !== indicatorResults) {
      this._indicatorResults = indicatorResults;
      this._indicatorResultsFail = indicatorResults
        .filter((i) => i.status === ResultTypes.FAILED)
        .sort(this.sortIndicatorResult);
      this._indicatorResultsWarning = indicatorResults
        .filter((i) => i.status === ResultTypes.WARNING)
        .sort(this.sortIndicatorResult);
      this._indicatorResultsPass = indicatorResults
        .filter((i) => i.status !== ResultTypes.FAILED && i.status !== ResultTypes.WARNING)
        .sort(this.sortIndicatorResult);
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input() objectives?: SloConfig['objectives'];

  @Input()
  get comparedIndicatorResults(): IndicatorResult[][] {
    return this._comparedIndicatorResults;
  }
  set comparedIndicatorResults(comparedIndicatorResults: IndicatorResult[][]) {
    this._comparedIndicatorResults = comparedIndicatorResults;
    this.updateDataSource();
  }

  @Input()
  get score(): number {
    return this._score;
  }
  set score(score: number) {
    if (score !== this._score) {
      this._score = score;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {}

  ngOnInit(): void {
    if (this.sortable) {
      this.sortable.sort('score', 'asc');
      this.tableEntries.sort = this.sortable;
    }
  }

  private updateDataSource(): void {
    this.tableEntries.data = this.assembleTablesEntries(this.indicatorResults);
  }

  private assembleTablesEntries(indicatorResults: IndicatorResult[]): SliResult[] {
    const totalscore = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    const isOld = indicatorResults.some((result) => !!result.targets);
    if (isOld) {
      this.columnNames = ['details', 'name', 'value', 'weight', 'targets', 'result', 'score'];
    } else {
      this.columnNames = ['details', 'name', 'value', 'weight', 'passTargets', 'warningTargets', 'result', 'score'];
    }
    return indicatorResults.map((indicatorResult) => {
      const comparedValue = this.calculateComparedValue(indicatorResult);
      const compared: Partial<SliResult> = {};
      if (!!comparedValue) {
        compared.comparedValue = AppUtils.formatNumber(comparedValue);
        compared.calculatedChanges = {
          absolute: AppUtils.formatNumber(indicatorResult.value.value - comparedValue),
          relative: AppUtils.formatNumber((indicatorResult.value.value / (comparedValue || 1)) * 100 - 100),
        };
      }

      return {
        name: indicatorResult.displayName || indicatorResult.value.metric,
        value: indicatorResult.value.message || AppUtils.formatNumber(indicatorResult.value.value),
        result: indicatorResult.status,
        score: totalscore === 0 ? 0 : AppUtils.round((indicatorResult.score / totalscore) * this.score, 2),
        passTargets: indicatorResult.passTargets,
        warningTargets: indicatorResult.warningTargets,
        targets: indicatorResult.targets,
        keySli: indicatorResult.keySli,
        success: indicatorResult.value.success,
        expanded: false,
        weight: this.objectives?.find((obj) => obj.sli === indicatorResult.value.metric)?.weight ?? 1,
        ...compared,
      };
    });
  }

  private calculateComparedValue(indicatorResult: IndicatorResult): number {
    if (indicatorResult.value.comparedValue == undefined) {
      let accSum = 0;
      let accCount = 0;
      for (const comparedIndicatorResult of this.comparedIndicatorResults) {
        const result = comparedIndicatorResult.find((res) => res.value.metric === indicatorResult.value.metric);
        if (result) {
          accSum += result.value.value;
          accCount++;
        }
      }
      return accSum / accCount;
    } else {
      return indicatorResult.value.comparedValue;
    }
  }

  private sortIndicatorResult(resultA: IndicatorResult, resultB: IndicatorResult): number {
    return (resultA.displayName || resultA.value.metric).localeCompare(resultB.displayName || resultB.value.metric);
  }

  public setExpanded(result: SliResult): void {
    if (result.comparedValue !== undefined) {
      result.expanded = !result.expanded;
    }
  }
}
