import Highcharts, {
  NavigatorXAxisPlotBandsOptions,
  PointClickEventObject,
  SeriesColumnOptions,
  SeriesHeatmapDataOptions,
  SeriesLineOptions,
} from 'highcharts';
import { ChangeDetectorRef, Component, Input, NgZone, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import {
  DtChart,
  DtChartOptions,
  DtChartSeries,
  DtChartSeriesVisibilityChangeEvent,
} from '@dynatrace/barista-components/chart';
import { Subject } from 'rxjs';
import { take, takeUntil } from 'rxjs/operators';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Trace } from '../../_models/trace';
import { EvaluationChartDataItem, EvaluationChartItem } from '../../_models/evaluation-chart-item';
import { HeatmapOptions } from '../../_models/heatmap-options';
import { HeatmapData, HeatmapSeriesOptions } from '../../_models/heatmap-series-options';
import { IndicatorResult, Target } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';
import { EvaluationHistory } from '../../_interfaces/evaluation-history';
import { AppUtils } from '../../_utils/app.utils';
import Yaml from 'yaml';
import { SloConfig } from '../../../../shared/interfaces/slo-config';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare let require: any;
/* eslint-disable @typescript-eslint/no-var-requires */
const _boostCanvas = require('highcharts/modules/boost-canvas');
const _boost = require('highcharts/modules/boost');
const _noData = require('highcharts/modules/no-data-to-display');
const _more = require('highcharts/highcharts-more');
const _heatmap = require('highcharts/modules/heatmap');
const _treemap = require('highcharts/modules/treemap');
/* eslint-enable @typescript-eslint/no-var-requires */
type SeriesPoint = PointClickEventObject & { series: EvaluationChartItem; point: { evaluationData: Trace } };

_boostCanvas(Highcharts);
_boost(Highcharts);
_noData(Highcharts);
_more(Highcharts);
_noData(Highcharts);
_heatmap(Highcharts);
_treemap(Highcharts);

type SliInfo = {
  score: number;
  warningCount: number;
  failedCount: number;
  passCount: number;
};

type SliInfoDictionary = { [evaluationId: string]: SliInfo | undefined };

interface IHeatmapPoint {
  sliInfo?: {
    passCount: number;
    warningCount: number;
    failedCount: number;
    thresholdPass: number;
    thresholdWarn: number;
    fail: boolean;
    warn: boolean;
  };
  data?: {
    keySli: boolean;
    score: number;
    passTargets: Target[];
    warningTargets: Target[];
  };
  value: number;
  x: number;
  y: number;
  z: number;
  evaluation?: Trace;
  color: string;
}

class HeatmapPoint implements IHeatmapPoint {
  sliInfo?: {
    passCount: number;
    warningCount: number;
    failedCount: number;
    thresholdPass: number;
    thresholdWarn: number;
    fail: boolean;
    warn: boolean;
  };
  data?: {
    keySli: boolean;
    score: number;
    passTargets: Target[];
    warningTargets: Target[];
  };
  value = 0;
  color = '';
  evaluation?: Trace;
  x = 0;
  y = 0;
  z = 0;
}

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss'],
})
export class KtbEvaluationDetailsComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public HeatmapPointClass = HeatmapPoint;
  public comparedIndicatorResults: IndicatorResult[][] = [];
  @Input() public showChart = true;
  @Input() public isInvalidated = false;

  @ViewChild('sloDialog')
  /* eslint-disable @typescript-eslint/no-explicit-any */
  public sloDialog?: TemplateRef<any>;
  public sloDialogRef?: MatDialogRef<any, any>;

  @ViewChild('invalidateEvaluationDialog')
  public invalidateEvaluationDialog?: TemplateRef<any>;
  public invalidateEvaluationDialogRef?: MatDialogRef<any, any>;
  /* eslint-enable @typescript-eslint/no-explicit-any */

  public isHeatmapExtendable = false;
  public isHeatmapExtended = false;
  private heatmapChart?: DtChart;

  @ViewChild('heatmapChart') set heatmap(heatmap: DtChart) {
    this.heatmapChart = heatmap;
  }

  public _evaluationColor: { [key: string]: string } = {
    pass: '#7dc540',
    warning: '#e6be00',
    fail: '#dc172a',
    failed: '#dc172a',
    info: '#f8f8f8',
  };

  public _evaluationState: Map<ResultTypes, string> = new Map<ResultTypes, string>([
    [ResultTypes.PASSED, 'recovered'],
    [ResultTypes.WARNING, 'warning'],
    [ResultTypes.FAILED, 'error'],
  ]);

  public _evaluationData?: Trace;
  public _selectedEvaluationData?: Trace;
  public _comparisonView: string | null = 'heatmap';
  private _metrics: string[] = [];
  public _chartOptions: DtChartOptions = {
    chart: {
      height: 400,
    },
    legend: {
      maxHeight: 70,
    },

    xAxis: {
      type: 'category',
      labels: {
        rotation: -45,
      },
      categories: [],
    },
    yAxis: [
      {
        title: undefined,
        labels: {
          format: '{value}',
        },
        min: 0,
        max: 100,
      },
      {
        title: undefined,
        labels: {
          format: '{value}',
        },
        opposite: true,
      },
    ],
    plotOptions: {
      column: {
        stacking: 'normal',
        pointWidth: 5,
        minPointLength: 2,
        point: {
          events: {
            click: (event: PointClickEventObject): boolean => {
              this._chartSeriesClicked(event as PointClickEventObject & { point: { evaluationData: Trace } });
              return true;
            },
          },
        },
      },
    },
  };
  public _chartSeries: (SeriesColumnOptions | SeriesLineOptions)[] = [];
  public _heatmapOptions: HeatmapOptions = {
    chart: {
      type: 'heatmap',
      height: 400,
    },
    xAxis: [
      {
        categories: [],
        plotBands: [],
        labels: {
          rotation: -45,
        },
        tickPositioner(): number[] {
          const positions = [];
          const labelWidth = 70;
          const ext = this.getExtremes();
          const xMax = Math.round(ext.max);
          const xMin = Math.round(ext.min);
          const maxElements = (document.querySelector('dt-chart')?.clientWidth || labelWidth) / labelWidth;
          const tick = Math.floor(xMax / maxElements) || 1;

          for (let i = xMax; i >= xMin; i -= tick) {
            positions.push(i);
          }
          return positions;
        },
      },
    ],

    yAxis: [
      {
        categories: [],
        title: undefined,
        labels: {
          format: '{value}',
          style: {
            textOverflow: 'ellipsis',
            width: 200,
          },
        },
      },
    ],

    colorAxis: {
      dataClasses: Object.keys(this._evaluationColor)
        .filter((key) => key !== 'failed')
        .map((key) => ({ color: this._evaluationColor[key], name: key })),
    },

    plotOptions: {
      heatmap: {
        point: {
          events: {
            click: (event: PointClickEventObject): boolean => {
              this._heatmapTileClicked(event);
              return true;
            },
          },
        },
      },
    },
  };
  public _heatmapSeries: HeatmapSeriesOptions[] = [];
  private _heatmapSeriesFull: HeatmapSeriesOptions[] = [];
  private _heatmapSeriesReduced: HeatmapSeriesOptions[] = [];
  private _heatmapCategoriesFull: string[] = [];
  private _heatmapCategoriesReduced: string[] = [];
  private _shouldSelectEvaluation = true;
  public updateResults?: EvaluationHistory;

  @Input()
  get evaluationData(): Trace | undefined {
    return this._evaluationData;
  }

  set evaluationData(evaluationData: Trace | undefined) {
    this.setEvaluation({ evaluation: evaluationData, shouldSelect: true });
  }

  @Input()
  set evaluationInfo(evaluationInfo: { evaluation?: Trace; shouldSelect: boolean }) {
    this.setEvaluation(evaluationInfo);
  }

  private setEvaluation(evaluationInfo: { evaluation?: Trace; shouldSelect: boolean }): void {
    if (this._evaluationData?.id !== evaluationInfo.evaluation?.id) {
      this._selectedEvaluationData = undefined;
      this.updateResults = undefined;
      this._evaluationData = evaluationInfo.evaluation;
      this._chartSeries = [];
      this._metrics = ['Score'];
      this._heatmapOptions.yAxis[0].categories = ['Score'];
      this._shouldSelectEvaluation = evaluationInfo.shouldSelect;
      this.evaluationDataChanged();
    } else if (this._evaluationData) {
      this.dataService.loadEvaluationResults(this._evaluationData);
    }
  }

  get heatmapSeries(): DtChartSeries[] {
    return this._heatmapSeries;
  }

  public get heatmapHeight(): number {
    return this.heatmapChart?._chartObject?.plotHeight ?? 0;
  }

  public get heatmapWidth(): number {
    return this.heatmapChart?._chartObject?.chartWidth ?? 0;
  }

  constructor(
    private _changeDetectorRef: ChangeDetectorRef,
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil,
    private zone: NgZone
  ) {}

  public ngOnInit(): void {
    this.dataService.evaluationResults.pipe(takeUntil(this.unsubscribe$)).subscribe((results) => {
      if (this.evaluationData && results.traces?.length) {
        this.parseSloFile(results.traces);
        if (this.evaluationData.data.evaluationHistory?.length) {
          this.updateResults = results;
        } else {
          this.refreshEvaluationBoard(results);
        }
      }
    });
  }

  private evaluationDataChanged(): void {
    if (this._evaluationData) {
      this.dataService.loadEvaluationResults(this._evaluationData);
      if (this.isInvalidated) {
        this.selectEvaluationData(this._evaluationData);
      } else if (!this._selectedEvaluationData && this._evaluationData.data.evaluationHistory) {
        this.setHeatmapDataAfterRender(this._evaluationData.data.evaluationHistory);
      }
    }
  }

  /**
   * If the data for the heatmap is set before the element is rendered, the width of the heatmap exceeds the page width
   * @param data
   * @private
   */
  private setHeatmapDataAfterRender(data: Trace[]): void {
    this.zone.onMicrotaskEmpty
      .asObservable()
      .pipe(take(1), takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.updateChartData(data);
        const trace = data.find((h) => h.shkeptncontext === this._evaluationData?.shkeptncontext);
        this.selectEvaluationData(trace);
      });
  }

  public refreshEvaluationBoard(results: EvaluationHistory): void {
    if (this.evaluationData) {
      if (results.type === 'evaluationHistory' && results.triggerEvent === this.evaluationData) {
        this.evaluationData.data.evaluationHistory = [
          ...(results.traces || []),
          ...(this.evaluationData.data.evaluationHistory || []),
        ].sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
      } else if (
        results.type === 'invalidateEvaluation' &&
        this.evaluationData.data.project === results.triggerEvent.data.project &&
        this.evaluationData.data.service === results.triggerEvent.data.service &&
        this.evaluationData.data.stage === results.triggerEvent.data.stage
      ) {
        this.evaluationData.data.evaluationHistory = this.evaluationData.data.evaluationHistory?.filter(
          (e) => e.id !== results.triggerEvent.id
        );
      }
      this._selectedEvaluationData = this._selectedEvaluationData?.id
        ? this.evaluationData.data.evaluationHistory?.find((h) => h.id === this._selectedEvaluationData?.id)
        : undefined;
      if (this.evaluationData.data.evaluationHistory) {
        this.updateChartData(this.evaluationData.data.evaluationHistory);
      }
    }
    this.updateResults = undefined;
  }

  private parseSloFile(evaluationTraces: Trace[]): void {
    for (const evaluationData of evaluationTraces) {
      if (evaluationData?.data?.evaluation?.sloFileContent && !evaluationData.data.evaluation.sloFileContentParsed) {
        evaluationData.data.evaluation.sloFileContentParsed = Yaml.parse(
          atob(evaluationData.data.evaluation.sloFileContent)
        ) as SloConfig;
        evaluationData.data.evaluation.score_pass =
          evaluationData.data.evaluation.sloFileContentParsed.total_score?.pass?.split('%')[0] ?? '';
        evaluationData.data.evaluation.score_warning =
          evaluationData.data.evaluation.sloFileContentParsed.total_score?.warning?.split('%')[0] ?? '';
        evaluationData.data.evaluation.compare_with =
          evaluationData.data.evaluation.sloFileContentParsed.comparison.compare_with ?? '';
        evaluationData.data.evaluation.include_result_with_score =
          evaluationData.data.evaluation.sloFileContentParsed.comparison.include_result_with_score;
        if (evaluationData.data.evaluation.comparedEvents) {
          evaluationData.data.evaluation.number_of_comparison_results =
            evaluationData.data.evaluation.comparedEvents?.length;
        } else {
          evaluationData.data.evaluation.number_of_comparison_results = 0;
        }
      }
    }
  }

  updateChartData(evaluationHistory: Trace[]): void {
    if (!this._selectedEvaluationData && evaluationHistory) {
      this.selectEvaluationData(evaluationHistory.find((h) => h.id === this._evaluationData?.id));
    }

    if (this.showChart) {
      const chartSeries = this.getChartSeries(evaluationHistory);
      this.sortChartSeries(chartSeries);
      this.updateHeatmapOptions(chartSeries);

      if (this._chartOptions.xAxis && !(this._chartOptions.xAxis instanceof Array)) {
        this._chartOptions.xAxis.categories = this._heatmapOptions.xAxis[0].categories;
      }
      this.setHeatmapData(chartSeries);

      if (this._heatmapSeriesFull[1].data.length > 0) {
        this.setHeatmapReducedSLO();
      }

      this.setSeriesXAxis(chartSeries);
      this._chartSeries = chartSeries;

      if (this.isHeatmapExtendable) {
        this._heatmapSeries = this._heatmapSeriesReduced;
      } else {
        this._heatmapSeries = this._heatmapSeriesFull;
      }
    }

    this.highlightHeatmap();
    this._changeDetectorRef.detectChanges();
  }

  private sortChartSeries(chartSeries: EvaluationChartItem[]): void {
    chartSeries.sort((seriesA, seriesB) => {
      let status;
      if (seriesA.name === 'Score' && seriesB.name === 'Score') {
        status = seriesA.type === 'line' ? 1 : -1;
      } else if (seriesA.name === 'Score') {
        status = -1;
      } else if (seriesB.name === 'Score') {
        status = 1;
      } else {
        status = seriesA.name.localeCompare(seriesB.name);
      }
      return status;
    });
  }

  private setHeatmapReducedSLO(): void {
    const minIdx =
      ((this._heatmapSeriesFull[1].data[this._heatmapSeriesFull[1].data.length - 1] as SeriesHeatmapDataOptions).y ??
        0) - 8;
    const reduced: HeatmapData[] = [];
    for (const series of this._heatmapSeriesFull[1].data) {
      if (series.y >= minIdx) {
        const srs = { ...series };
        srs.y = srs.y - minIdx;
        reduced.push(srs);
      }
    }
    this._heatmapSeriesReduced[1].data = reduced;
  }

  private setSeriesXAxis(chartSeries: EvaluationChartItem[]): void {
    for (const item of chartSeries) {
      for (const data of item.data) {
        data.x = data.evaluationData
          ? this._heatmapOptions.xAxis[0].categories.indexOf(data.evaluationData.getHeatmapLabel())
          : -1;
      }
    }
  }

  private getSliResultInfos(chartSeries: EvaluationChartItem[]): SliInfoDictionary {
    const sliResultInfos: SliInfoDictionary = {};
    for (const chartItem of chartSeries) {
      for (const item of chartItem.data) {
        if (item.evaluationData?.data.evaluation?.indicatorResults && !sliResultInfos[item.evaluationData.id]) {
          const indicatorResults = item.evaluationData.data.evaluation.indicatorResults;
          sliResultInfos[item.evaluationData.id] = this.getSliResultInfo(indicatorResults);
        }
      }
    }
    return sliResultInfos;
  }

  private getSliResultInfo(indicatorResults: IndicatorResult[]): {
    score: number;
    warningCount: number;
    failedCount: number;
    passCount: number;
  } {
    return indicatorResults.reduce(
      (acc, result) => {
        const warning = result.status === ResultTypes.WARNING ? 1 : 0;
        const failed = result.status === ResultTypes.FAILED ? 1 : 0;
        return {
          score: acc.score + result.score,
          warningCount: acc.warningCount + warning,
          failedCount: acc.failedCount + failed,
          passCount: acc.passCount + 1 - warning - failed,
        };
      },
      { score: 0, warningCount: 0, failedCount: 0, passCount: 0 } as SliInfo
    );
  }

  private setHeatmapData(chartSeries: EvaluationChartItem[]): void {
    const sliResultsInfo: SliInfoDictionary = this.getSliResultInfos(chartSeries);
    this._heatmapSeriesReduced = [
      {
        name: 'Score',
        type: 'heatmap',
        rowsize: 0.85,
        turboThreshold: 0,
        data: [],
      },
      {
        name: 'SLOs',
        type: 'heatmap',
        turboThreshold: 0,
        data: [],
      },
    ];

    this._heatmapSeriesFull = [
      {
        name: 'Score',
        type: 'heatmap',
        rowsize: 0.85,
        turboThreshold: 0,
        data:
          chartSeries
            .find((series) => series.name === 'Score')
            ?.data.filter((s): s is EvaluationChartDataItem & { evaluationData: Trace } => !!s.evaluationData)
            .map((s) => {
              const index = this._metrics.indexOf('Score');
              const x = this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData.getHeatmapLabel());
              const dataPoint: IHeatmapPoint = {
                x,
                y: index,
                z: s.y,
                evaluation: s.evaluationData,
                color: this._evaluationColor[s.evaluationData.data.result ?? 'info'],
                value: s.y,
                sliInfo: {
                  warningCount: sliResultsInfo[s.evaluationData.id]?.warningCount ?? 0,
                  failedCount: sliResultsInfo[s.evaluationData.id]?.failedCount ?? 0,
                  passCount: sliResultsInfo[s.evaluationData.id]?.passCount ?? 0,
                  thresholdPass: +(s.evaluationData.data.evaluation?.score_pass ?? 0),
                  thresholdWarn: +(s.evaluationData.data.evaluation?.score_warning ?? 0),
                  fail: s.evaluationData.isFailed(),
                  warn: s.evaluationData.isWarning(),
                },
              };
              const reducedDataPoint = { ...dataPoint };
              reducedDataPoint.y = 9;
              this._heatmapSeriesReduced[0].data.push(reducedDataPoint);
              return dataPoint;
            }) ?? [],
      },
      {
        name: 'SLOs',
        type: 'heatmap',
        turboThreshold: 0,
        data: [...chartSeries].reverse().reduce(
          (r, d) => [
            ...r,
            ...d.data
              .filter((s): s is EvaluationChartDataItem & { indicatorResult: IndicatorResult } => !!s.indicatorResult)
              .map((s) => {
                const index = this._metrics.indexOf(s.indicatorResult.value.metric);
                const x = s.evaluationData
                  ? this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData.getHeatmapLabel())
                  : -1;
                const totalScore = sliResultsInfo[s.evaluationData?.id ?? '']?.score;
                const score = !totalScore
                  ? 0
                  : AppUtils.round(
                      (s.indicatorResult.score / totalScore) * (s.evaluationData?.data.evaluation?.score ?? 1),
                      2
                    );

                return {
                  x,
                  y: index,
                  z: s.indicatorResult.score,
                  color: s.indicatorResult.value.success
                    ? this._evaluationColor[s.indicatorResult.status]
                    : this._evaluationColor.info,
                  data: {
                    keySli: s.indicatorResult.keySli,
                    score: score,
                    passTargets: s.indicatorResult.passTargets,
                    warningTargets: s.indicatorResult.warningTargets,
                  },
                  value: AppUtils.formatNumber(s.indicatorResult.value.value),
                } as IHeatmapPoint;
              }),
          ],
          [] as HeatmapData[]
        ),
      },
    ];
  }

  public formatNumber(num: number): number {
    return AppUtils.formatNumber(num);
  }

  private getChartSeries(evaluationHistory: Trace[]): EvaluationChartItem[] {
    const chartSeries: EvaluationChartItem[] = [];
    for (const evaluation of evaluationHistory) {
      const scoreData = {
        y: evaluation.data.evaluation?.score ?? 0,
        evaluationData: evaluation,
        color: this._evaluationColor[evaluation.data.evaluation?.result ?? 'info'],
        name: evaluation.getChartLabel(),
      };

      let indicatorScoreSeriesColumn = chartSeries.find(
        (series) => series.name === 'Score' && series.type === 'column'
      );
      let indicatorScoreSeriesLine = chartSeries.find((series) => series.name === 'Score' && series.type === 'line');
      if (!indicatorScoreSeriesColumn) {
        indicatorScoreSeriesColumn = {
          metricName: 'Score',
          name: 'Score',
          type: 'column',
          data: [],
          cursor: 'pointer',
          turboThreshold: 0,
        };
        chartSeries.push(indicatorScoreSeriesColumn);
      }
      if (!indicatorScoreSeriesLine) {
        indicatorScoreSeriesLine = {
          name: 'Score',
          metricName: 'Score',
          type: 'line',
          data: [],
          cursor: 'pointer',
          visible: false,
          turboThreshold: 0,
        };
        chartSeries.push(indicatorScoreSeriesLine);
      }

      indicatorScoreSeriesColumn.data.push(scoreData);
      indicatorScoreSeriesLine.data.push(scoreData);

      if (evaluation.data.evaluation?.indicatorResults) {
        evaluation.data.evaluation.indicatorResults.forEach((indicatorResult: IndicatorResult) => {
          const indicatorData = {
            y: indicatorResult.value.value,
            indicatorResult,
            evaluationData: evaluation,
            name: evaluation.getChartLabel(),
          };

          let indicatorChartSeries = chartSeries.find((series) => series.metricName === indicatorResult.value.metric);
          if (!indicatorChartSeries) {
            indicatorChartSeries = {
              metricName: indicatorResult.value.metric,
              name: this.getLastDisplayName(evaluationHistory, indicatorResult.value.metric),
              type: 'line',
              yAxis: 1,
              data: [],
              visible: false,
              turboThreshold: 0,
            };
            chartSeries.push(indicatorChartSeries);
          }
          indicatorChartSeries.data.push(indicatorData);
        });
      }
    }
    return chartSeries;
  }

  private getLastDisplayName(evaluationHistory: Trace[], metric: string): string {
    let displayName = metric;
    if (metric !== 'Score') {
      for (let i = evaluationHistory.length - 1; i >= 0; i--) {
        const result = evaluationHistory[i].data.evaluation?.indicatorResults?.find(
          (indicatorResult) => indicatorResult.value.metric === metric
        );
        if (result) {
          displayName = result.displayName || result.value.metric;
          break;
        }
      }
    }
    return displayName;
  }

  updateHeatmapOptions(chartSeries: EvaluationChartItem[]): void {
    const heatmapCategoriesFull = [...this._heatmapOptions.yAxis[0].categories];
    const heatmapCategoriesReduced = [...this._heatmapOptions.yAxis[0].categories];
    chartSeries.forEach((series, i) => {
      if (!this._metrics.includes(series.metricName)) {
        heatmapCategoriesFull.unshift(series.name);
        if (i <= 10) {
          heatmapCategoriesReduced.unshift(series.name);
        }
        this._metrics.unshift(series.metricName);
      }
      if (series.name === 'Score') {
        this.updateHeatmapScore(series);
      }
    });

    this._heatmapCategoriesFull = heatmapCategoriesFull;
    this._heatmapCategoriesReduced = heatmapCategoriesReduced;

    if (this._heatmapCategoriesFull.length > 10) {
      this.isHeatmapExtendable = true;
      this.isHeatmapExtended = false;
    } else {
      this.isHeatmapExtended = true;
    }
    this._updateHeatmapExtension();
  }

  private updateHeatmapScore(series: EvaluationChartItem): void {
    series.data.sort(this.compareSeriesData);
    this._heatmapOptions.xAxis[0].categories = series.data
      .map((item, index, items) => {
        const duplicateItems = items.filter(
          (c) => c.evaluationData?.getHeatmapLabel() === item.evaluationData?.getHeatmapLabel()
        );
        if (duplicateItems.length > 1) {
          item.label = `${item.evaluationData?.getHeatmapLabel()} (${duplicateItems.indexOf(item) + 1})`;
        } else {
          item.label = item.evaluationData?.getHeatmapLabel();
        }
        return item;
      })
      .map((item) => {
        item.evaluationData?.setHeatmapLabel(item.label ?? '');
        return item.evaluationData?.getHeatmapLabel() ?? '';
      });
  }

  private compareSeriesData(a: EvaluationChartDataItem, b: EvaluationChartDataItem): number {
    return DateUtil.compareTraceTimesDesc(a.evaluationData, b.evaluationData);
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  seriesVisibilityChanged(_: DtChartSeriesVisibilityChangeEvent): void {
    // NOOP
  }

  _chartSeriesClicked(event: PointClickEventObject & { point: { evaluationData: Trace } }): void {
    this.selectEvaluationData(event.point.evaluationData, true);
  }

  _heatmapTileClicked(event: PointClickEventObject): void {
    this.selectEvaluationData(this._heatmapSeries[0].data[event.point.x].evaluation, true);
  }

  selectEvaluationData(evaluation?: Trace, forceSelect = false): void {
    if (this._shouldSelectEvaluation || forceSelect) {
      this._selectedEvaluationData = evaluation;
      this.highlightHeatmap();
    }
  }

  highlightHeatmap(): void {
    if (this._selectedEvaluationData && !this.isInvalidated) {
      this.comparedIndicatorResults = [];
      const secondaryHighlightIndexes = this._selectedEvaluationData?.data.evaluation?.comparedEvents?.map(
        (eventId) => {
          const eventIndex = this._heatmapSeries[0]?.data.findIndex((e) => e.evaluation?.id === eventId);
          this.comparedIndicatorResults.push(
            this._heatmapSeries[0]?.data[eventIndex].evaluation?.data.evaluation?.indicatorResults ?? []
          );
          return eventIndex;
        }
      );
      const plotBands: NavigatorXAxisPlotBandsOptions[] = [];
      const highlightIndex = this._heatmapOptions.xAxis[0].categories.indexOf(
        this._selectedEvaluationData.getHeatmapLabel()
      );
      if (highlightIndex >= 0) {
        plotBands.push({
          className: 'highlight-primary',
          from: highlightIndex - 0.5,
          to: highlightIndex + 0.5,
          zIndex: 100,
        });
      }
      if (secondaryHighlightIndexes) {
        this.setSecondaryHighlight(secondaryHighlightIndexes, plotBands);
      }
      this._heatmapOptions.xAxis[0].plotBands = plotBands;
      if (
        this._selectedEvaluationData.data.evaluation?.number_of_missing_comparison_results &&
        this._selectedEvaluationData?.data.evaluation.comparedEvents?.length !== undefined
      ) {
        this._selectedEvaluationData.data.evaluation.number_of_missing_comparison_results =
          this._selectedEvaluationData?.data.evaluation.comparedEvents.length -
          (this._heatmapOptions.xAxis[0].plotBands?.length - 1);
      }
    } else {
      this._heatmapOptions.xAxis[0].plotBands = [];
      if (this._selectedEvaluationData?.data.evaluation) {
        this._selectedEvaluationData.data.evaluation.number_of_missing_comparison_results = 0;
      }
    }
    this.heatmapChart?._update();
    this._changeDetectorRef.detectChanges();
  }

  private setSecondaryHighlight(
    secondaryHighlightIndices: number[],
    plotBands: NavigatorXAxisPlotBandsOptions[]
  ): void {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const _this = this;
    for (const secondaryHighlightIndex of secondaryHighlightIndices) {
      plotBands.push({
        className: 'highlight-secondary',
        from: secondaryHighlightIndex - 0.5,
        to: secondaryHighlightIndex + 0.5,
        zIndex: 100,
        events: {
          // eslint-disable-next-line @typescript-eslint/no-loop-func
          click(): void {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            const idx = this.options.from + 0.5;
            const evaluation = _this._heatmapSeries[0]?.data[idx]?.evaluation;
            setTimeout(() => {
              _this.selectEvaluationData(evaluation);
            });
          },
        },
      });
    }
  }

  showSloDialog(): void {
    if (this.sloDialog && this._selectedEvaluationData) {
      this.sloDialogRef = this.dialog.open(this.sloDialog, {
        data: atob(this._selectedEvaluationData.data.evaluation?.sloFileContent ?? ''),
      });
    }
  }

  closeSloDialog(): void {
    if (this.sloDialogRef) {
      this.sloDialogRef.close();
    }
  }

  copySloPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'slo payload');
  }

  invalidateEvaluationTrigger(): void {
    if (this.invalidateEvaluationDialog) {
      this.invalidateEvaluationDialogRef = this.dialog.open(this.invalidateEvaluationDialog, {
        data: this._selectedEvaluationData,
      });
    }
  }

  invalidateEvaluation(evaluation: Trace, reason: string): void {
    this.dataService.invalidateEvaluation(evaluation, reason);
    this.closeInvalidateEvaluationDialog();
  }

  closeInvalidateEvaluationDialog(): void {
    if (this.invalidateEvaluationDialogRef) {
      this.invalidateEvaluationDialogRef.close();
    }
  }

  // remove duplicated points like "Score"
  filterPoints(points: SeriesPoint[]): SeriesPoint[] {
    return points.filter(
      (item, index) => index === points.findIndex((subItem) => subItem.series.name === item.series.name)
    );
  }

  public getEvaluationFromPoint(points: SeriesPoint[]): Trace {
    return points[0].point.evaluationData;
  }

  public toggleHeatmap(): void {
    this.isHeatmapExtended = !this.isHeatmapExtended;
    this._updateHeatmapExtension();
  }

  private _updateHeatmapExtension(): void {
    if (this.isHeatmapExtended) {
      this._heatmapSeries = this._heatmapSeriesFull;
      this._heatmapOptions.yAxis[0].categories = this._heatmapCategoriesFull;
      this._heatmapOptions.chart.height = this._heatmapCategoriesFull.length * 28 + 160;
    } else {
      this._heatmapSeries = this._heatmapSeriesReduced;
      this._heatmapOptions.yAxis[0].categories = this._heatmapCategoriesReduced;
      this._heatmapOptions.chart.height = this._heatmapCategoriesReduced.length * 28 + 173;
    }
    if (this.isHeatmapExtendable) {
      this._heatmapOptions.xAxis[0].offset = 40;
    } else {
      this._heatmapOptions.xAxis[0].offset = undefined;
    }

    this.heatmapChart?._update();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
