import { ChangeDetectorRef, Component, Input, ViewChild } from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';
import { SliResult } from '../../_interfaces/sli-result';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { AppUtils } from '../../_utils/app.utils';
import { SloConfig } from '../../../../shared/interfaces/slo-config';
import { DataService } from '../../_services/data.service';
import { Trace } from '../../_models/trace';
import { finalize } from 'rxjs/operators';

interface IFallbackData {
  comparedEvents: string[];
  projectName: string;
  comparedIndicatorResults?: IndicatorResult[][];
}

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss'],
})
export class KtbSliBreakdownComponent {
  private _sortable?: DtSort;
  @ViewChild('sortable', { read: DtSort, static: false })
  set sortable(sortable: DtSort | undefined) {
    this._sortable = sortable;
    this.updateSort();
  }
  get sortable(): DtSort | undefined {
    return this._sortable;
  }

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
  private _objectives?: SloConfig['objectives'];
  private comparedEvents: string[] = [];
  private projectName = '';
  // either the compared evaluations are fetched on demand if the comparedValue property does not exist,
  //  or it is set through the ktb-evaluation-chart.component  because it already loads the history
  private comparedIndicatorResults?: IndicatorResult[][];
  public maximumAvailableWeight = 1;
  public toSliResult = (row: SliResult): SliResult => row;
  public isLoading = false;

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

  @Input()
  set objectives(objectives: SloConfig['objectives'] | undefined) {
    this._objectives = objectives;
    this.updateDataSource();
  }
  get objectives(): SloConfig['objectives'] | undefined {
    return this._objectives;
  }
  @Input()
  set fallBackData(data: IFallbackData) {
    this.comparedEvents = data.comparedEvents;
    this.projectName = data.projectName;
    this.comparedIndicatorResults = data.comparedIndicatorResults;

    this.updateDataSource();
  }
  get fallBackData(): IFallbackData {
    return {
      comparedEvents: this.comparedEvents,
      projectName: this.projectName,
      comparedIndicatorResults: this.comparedIndicatorResults,
    };
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

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {}

  private updateSort(): void {
    if (this.sortable) {
      this.sortable.sort('score', 'asc');
      this.tableEntries.sort = this.sortable;
    }
  }

  private updateDataSource(fetchedComparedResults = false): void {
    const sliResults = this.assembleTablesEntries(this.indicatorResults, fetchedComparedResults);
    if (!sliResults) {
      return;
    }
    // max reachable weight is actually the max reachable score. max weight = 100% score
    this.maximumAvailableWeight = sliResults
      .filter((sli) => sli.result !== ResultTypes.INFO)
      .reduce((acc, result) => acc + result.weight, 0);
    this.tableEntries.data = sliResults;
    this.updateSort();
  }

  private assembleTablesEntries(
    indicatorResults: IndicatorResult[],
    fetchedComparedEvaluations = false
  ): SliResult[] | undefined {
    const totalscore = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    this.setColumnNames(indicatorResults);

    // comparedValue was introduced in 0.12
    const comparedValueMissing = indicatorResults.some(
      (indicatorResult) =>
        indicatorResult.value.comparedValue === undefined || indicatorResult.value.comparedValue === null
    );
    const loadComparedEvaluations =
      this.comparedEvents.length &&
      this.projectName &&
      (comparedValueMissing ||
        (this.comparedIndicatorResults && this.comparedIndicatorResults.length !== this.comparedEvents.length));

    if (loadComparedEvaluations && !fetchedComparedEvaluations) {
      this.isLoading = true;
      this.dataService
        .getTracesByIds(this.projectName, this.comparedEvents)
        .pipe(
          finalize(() => {
            this.isLoading = false;
          })
        )
        .subscribe((traces: Trace[]) => {
          this.comparedIndicatorResults = traces.map((trace) => trace.data.evaluation?.indicatorResults ?? []);
          this.updateDataSource(true);
        });
      return undefined;
    }

    const sliResults = indicatorResults.map<SliResult>((indicatorResult) => ({
      name: indicatorResult.displayName || indicatorResult.value.metric,
      value: indicatorResult.value.message || AppUtils.formatNumber(indicatorResult.value.value),
      result: indicatorResult.status,
      score: totalscore === 0 ? 0 : (indicatorResult.score / totalscore) * this.score,
      passTargets: indicatorResult.passTargets,
      warningTargets: indicatorResult.warningTargets,
      targets: indicatorResult.targets,
      keySli: indicatorResult.keySli,
      success: indicatorResult.value.success,
      expanded: false,
      weight: this.objectives?.find((obj) => obj.sli === indicatorResult.value.metric)?.weight ?? 1,
      ...this.getComparedValues(indicatorResult),
    }));

    return this.getUniqueSliResult(sliResults);
  }

  public getUniqueSliResult(sliResults: SliResult[]): SliResult[] {
    return sliResults.map((sliResult, index) => ({
      ...sliResult,
      name: this.getUniqueSliName(sliResult.name, index, sliResults),
    }));
  }

  private getUniqueSliName(sliName: string, sliResultIndex: number, sliResults: SliResult[]): string {
    const duplicates = sliResults.filter((result) => result.name === sliName);
    if (duplicates.length > 1) {
      const previousDuplicates = sliResults.filter(
        (result, index) => index < sliResultIndex && result.name === sliName
      ).length;
      return `${sliName} (${previousDuplicates + 1})`;
    }
    return sliName;
  }

  private setColumnNames(indicatorResults: IndicatorResult[]): void {
    const isOld = indicatorResults.some((result) => !!result.targets);
    // splitting of targets into pass and warning was introduced in 0.8
    if (isOld) {
      this.columnNames = ['details', 'name', 'value', 'weight', 'targets', 'result', 'score'];
    } else {
      this.columnNames = ['details', 'name', 'value', 'weight', 'passTargets', 'warningTargets', 'result', 'score'];
    }
  }

  private getComparedValues(indicatorResult: IndicatorResult): Partial<SliResult> {
    const comparedValue = indicatorResult.value.comparedValue ?? this.calculateComparedValue(indicatorResult);
    const compared: Partial<SliResult> = {};
    if (!isNaN(comparedValue)) {
      compared.comparedValue = AppUtils.formatNumber(comparedValue);
      compared.calculatedChanges = {
        absolute: AppUtils.formatNumber(indicatorResult.value.value - comparedValue),
        relative: AppUtils.formatNumber(this.getRelativeChange(indicatorResult, comparedValue)),
      };
    }
    return compared;
  }

  private getRelativeChange(indicatorResult: IndicatorResult, comparedValue: number): number {
    if (indicatorResult.value.value === 0 && comparedValue === 0) {
      return 0;
    }
    return (indicatorResult.value.value / (comparedValue || 1)) * 100 - 100;
  }

  public calculateComparedValue(indicatorResult: IndicatorResult): number {
    let accSum = 0;
    let accCount = 0;
    for (const comparedIndicatorResult of this.comparedIndicatorResults ?? []) {
      const result = comparedIndicatorResult.find((res) => res.value.metric === indicatorResult.value.metric);
      if (result) {
        accSum += result.value.value;
        accCount++;
      }
    }
    return accSum / accCount;
  }

  private sortIndicatorResult(resultA: IndicatorResult, resultB: IndicatorResult): number {
    return (resultA.displayName || resultA.value.metric).localeCompare(resultB.displayName || resultB.value.metric);
  }

  public setExpanded(result: SliResult): void {
    if (result.comparedValue !== undefined) {
      result.expanded = !result.expanded;
    }
  }

  public getRelativeSliWeight(result: SliResult): number {
    return this.maximumAvailableWeight !== 0 ? (result.weight / this.maximumAvailableWeight) * 100 : 0;
  }
}
