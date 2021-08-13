import { ChangeDetectorRef, Component, Input, OnInit, ViewChild } from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';
import { SliResult } from '../../_models/sli-result';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss']
})
export class KtbSliBreakdownComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable?: DtSort;

  public evaluationState: Map<ResultTypes, string> = new Map<ResultTypes, string>([
    [ResultTypes.PASSED, 'passed'],
    [ResultTypes.WARNING, 'warning'],
    [ResultTypes.FAILED, 'failed']
  ]);
  public ResultTypes: typeof ResultTypes = ResultTypes;
  private _indicatorResults?: IndicatorResult[];
  private _indicatorResultsFail: IndicatorResult[] = [];
  private _indicatorResultsWarning: IndicatorResult[] = [];
  private _indicatorResultsPass: IndicatorResult[] = [];
  private _score = 0;
  public columnNames: string[] = [];
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  public readonly SliResultClass = SliResult;
  private _comparedIndicatorResults: IndicatorResult[] = [];

  @Input()
  get indicatorResults(): IndicatorResult[] {
    return [...this._indicatorResultsFail, ...this._indicatorResultsWarning, ...this._indicatorResultsPass];
  }
  set indicatorResults(indicatorResults: IndicatorResult[]) {
    if (this._indicatorResults !== indicatorResults) {
      this._indicatorResults = indicatorResults;
      this._indicatorResultsFail = indicatorResults.filter(i => i.status === ResultTypes.FAILED).sort(this.sortIndicatorResult);
      this._indicatorResultsWarning = indicatorResults.filter(i => i.status === ResultTypes.WARNING).sort(this.sortIndicatorResult);
      this._indicatorResultsPass = indicatorResults.filter(i => i.status !== ResultTypes.FAILED && i.status !== ResultTypes.WARNING)
                                  .sort(this.sortIndicatorResult);
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get comparedIndicatorResults(): IndicatorResult[] {
    return this._comparedIndicatorResults;
  }
  set comparedIndicatorResults(comparedIndicatorResults: IndicatorResult[]) {
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

  constructor(private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    if (this.sortable) {
      this.sortable.sort('score', 'asc');
      this.tableEntries.sort = this.sortable;
    }
  }

  private updateDataSource() {
    this.tableEntries.data = this.assembleTablesEntries(this.indicatorResults);
  }

  private formatNumber(value: number) {
    let n = value;
    if (n < 1) {
      n = Math.floor(n * 1000) / 1000;
    } else if (n < 100) {
      n = Math.floor(n * 100) / 100;
    } else if (n < 1000) {
      n = Math.floor(n * 10) / 10;
    } else {
      n = Math.floor(n);
    }

    return n;
  }

  private assembleTablesEntries(indicatorResults: IndicatorResult[]): SliResult[] {
    const totalscore  = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    const isOld = indicatorResults.some(result => !!result.targets);
    if (isOld) {
      this.columnNames = [
        'details',
        'name',
        'value',
        'weight',
        'targets',
        'result',
        'score'
      ];
    } else {
      this.columnNames = [
        'details',
        'name',
        'value',
        'weight',
        'passTargets',
        'warningTargets',
        'result',
        'score'
      ];
    }
    return indicatorResults.map(indicatorResult =>  {
      const comparedValue = this.comparedIndicatorResults
                            ?.find(result => result.value.metric === indicatorResult.value.metric)
                            ?.value.value;
      const compared: Partial<SliResult> = {};
      if (comparedValue) {
        compared.comparedValue = this.formatNumber(comparedValue);
        compared.calculatedChanges = {
          absolute: this.formatNumber(comparedValue - indicatorResult.value.value),
          relative: this.formatNumber(comparedValue / (indicatorResult.value.value || 1) * 100 - 100)
        };
      }

      return {
        name: indicatorResult.displayName || indicatorResult.value.metric,
        value: indicatorResult.value.message || this.formatNumber(indicatorResult.value.value),
        result: indicatorResult.status,
        score: totalscore === 0 ? 0 : this.round(indicatorResult.score / totalscore * this.score, 2),
        passTargets: indicatorResult.passTargets,
        warningTargets: indicatorResult.warningTargets,
        targets: indicatorResult.targets,
        keySli: indicatorResult.keySli,
        success: indicatorResult.value.success,
        expanded: false,
        weight: indicatorResult.score,
        ...compared,
      };
    });
  }

  private sortIndicatorResult(resultA: IndicatorResult, resultB: IndicatorResult) {
    return (resultA.displayName || resultA.value.metric).localeCompare(resultB.displayName || resultB.value.metric);
  }

  private round(value: number, places: number): number {
    return +(Math.round(Number(`${value}e+${places}`))  + `e-${places}`);
  }

  public setExpanded(result: SliResult): void {
    if (result.comparedValue !== undefined) {
      result.expanded = !result.expanded;
    }
  }

}
