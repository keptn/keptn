import {ChangeDetectorRef, Component, Input, OnInit, ViewChild} from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';
import {SliResult} from '../../_models/sli-result';

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss']
})
export class KtbSliBreakdownComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable: DtSort;

  public evaluationState = {
    pass: 'passed',
    warning: 'warning',
    fail: 'failed'
  };

  private _indicatorResults: any;
  private _indicatorResultsFail: any = [];
  private _indicatorResultsWarning: any = [];
  private _indicatorResultsPass: any = [];
  private _score: number;
  public columnNames: any = [];
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  private _comparedIndicatorResults: any[];

  @Input()
  get indicatorResults(): any {
    return [...this._indicatorResultsFail, ...this._indicatorResultsWarning, ...this._indicatorResultsPass];
  }
  set indicatorResults(indicatorResults: any) {
    if (this._indicatorResults !== indicatorResults) {
      this._indicatorResults = indicatorResults;
      this._indicatorResultsFail = indicatorResults.filter(i => i.status === 'fail').sort(this.sortIndicatorResult);
      this._indicatorResultsWarning = indicatorResults.filter(i => i.status === 'warning').sort(this.sortIndicatorResult);
      this._indicatorResultsPass = indicatorResults.filter(i => i.status !== 'fail' && i.status !== 'warning').sort(this.sortIndicatorResult);
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get comparedIndicatorResults(): any[] {
    return this._comparedIndicatorResults;
  }
  set comparedIndicatorResults(comparedIndicatorResults: any[]) {
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
    this.sortable.sort('score', 'asc');
    this.tableEntries.sort = this.sortable;
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

  public toSliResult(row: SliResult): SliResult {
    return row;
  }

  private assembleTablesEntries(indicatorResults): SliResult[] {
    const totalscore  = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    const isOld = indicatorResults.some(result => !!result.targets);
    if (isOld) {
      this.columnNames = [
        'details',
        'name',
        'value',
        'targets',
        'result',
        'score'
      ];
    } else {
      this.columnNames = [
        'details',
        'name',
        'value',
        'passTargets',
        'warningTargets',
        'result',
        'score'
      ];
    }
    return indicatorResults.map(indicatorResult =>  {
      const comparedValue = this.comparedIndicatorResults?.find(result => result.value.metric === indicatorResult.value.metric)?.value.value;
      const compared: any = {};
      if (comparedValue) {
        compared.comparedValue = this.formatNumber(comparedValue);
        compared.absoluteChange = this.formatNumber(comparedValue - indicatorResult.value.value);
        compared.relativeChange = this.formatNumber(comparedValue / (indicatorResult.value.value || 1) * 100 - 100);
      }

      return {
        name: indicatorResult.displayName || indicatorResult.value.metric,
        value: indicatorResult.value.message || this.formatNumber(indicatorResult.value.value),
        result: indicatorResult.status,
        score: totalscore === 0 ? 0 : this.round(indicatorResult.score / totalscore * this.score, 2),
        passTargets: indicatorResult.passTargets,
        warningTargets: indicatorResult.warningTargets,
        targets: indicatorResult.targets,
        ...compared,
        keySli: indicatorResult.keySli,
        success: indicatorResult.value.success
      };
    });
  }

  private sortIndicatorResult(resultA: any, resultB: any) {
    return (resultA.displayName || resultA.value.metric).localeCompare(resultB.displayName || resultB.value.metric);
  }

  private round(value: number, places: number): number {
    return +(Math.round(Number(`${value}e+${places}`))  + `e-${places}`);
  }

}
