import { ChangeDetectorRef, Component, Input, OnInit, ViewChild } from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';
import { SliResult } from '../../_interfaces/sli-result';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { AppUtils } from '../../_utils/app.utils';
import { SloConfig } from '../../../../shared/interfaces/slo-config';
import { DataService } from '../../_services/data.service';
import { Trace } from '../../_models/trace';

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
  private _objectives?: SloConfig['objectives'];
  private _comparedEvents: string[] = [];
  private _projectName = '';
  // either the compared evaluations are fetched on demand if the comparedValue property does not exist,
  //  or it is set through the ktb-evaluation-chart.component  because it already loads the history
  private _comparedIndicatorResults: IndicatorResult[][] = [];
  public maximumAvailableWeight = 1;
  public toSliResult = (row: SliResult): SliResult => row;

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
  set comparedEvents(comparedEvents: string[]) {
    this._comparedEvents = comparedEvents;
    this.updateDataSource();
  }
  get comparedEvents(): string[] {
    return this._comparedEvents;
  }
  @Input()
  set projectName(projectName: string) {
    this._projectName = projectName;
    this.updateDataSource();
  }
  get projectName(): string {
    return this._projectName;
  }

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

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {}

  ngOnInit(): void {
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
    this.maximumAvailableWeight = sliResults.reduce((acc, result) => acc + result.weight, 0);
    this.tableEntries.data = sliResults;
  }

  private assembleTablesEntries(
    indicatorResults: IndicatorResult[],
    fetchedComparedEvaluations = false
  ): SliResult[] | undefined {
    const totalscore = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    const isOld = indicatorResults.some((result) => !!result.targets);
    // splitting of targets into pass and warning was introduced in 0.8
    if (isOld) {
      this.columnNames = ['details', 'name', 'value', 'weight', 'targets', 'result', 'score'];
    } else {
      this.columnNames = ['details', 'name', 'value', 'weight', 'passTargets', 'warningTargets', 'result', 'score'];
    }
    // comparedValue was introduced in 0.12
    const hasComparedValue = indicatorResults.every(
      (indicatorResult) => indicatorResult.value.comparedValue !== undefined
    );
    const loadComparedEvaluations =
      this.comparedEvents.length &&
      this.projectName &&
      (!hasComparedValue || this.comparedIndicatorResults?.length !== this.comparedEvents.length);

    if (loadComparedEvaluations && !fetchedComparedEvaluations) {
      this.dataService.getTracesByIds(this.projectName, this.comparedEvents).subscribe((traces: Trace[]) => {
        this._comparedIndicatorResults = traces.map((trace) => trace.data.evaluation?.indicatorResults ?? []);
        this.updateDataSource(true);
      });
      return undefined;
    }

    return indicatorResults.map((indicatorResult) => {
      const comparedValue = indicatorResult.value.comparedValue ?? this.calculateComparedValue(indicatorResult);
      const compared: Partial<SliResult> = {};
      if (!isNaN(comparedValue)) {
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
        score: totalscore === 0 ? 0 : (indicatorResult.score / totalscore) * this.score,
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

  public calculateComparedValue(indicatorResult: IndicatorResult): number {
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
