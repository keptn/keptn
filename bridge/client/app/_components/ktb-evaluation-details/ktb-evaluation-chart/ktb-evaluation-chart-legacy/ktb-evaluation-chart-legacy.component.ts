import { ChangeDetectorRef, Component, EventEmitter, Input, NgZone, OnDestroy, Output, ViewChild } from '@angular/core';
import { HeatmapData, HeatmapSeriesOptions } from '../../../../_models/heatmap-series-options';
import {
  DtChart,
  DtChartOptions,
  DtChartSeries,
  DtChartSeriesVisibilityChangeEvent,
} from '@dynatrace/barista-components/chart';
import Highcharts, {
  NavigatorXAxisPlotBandsOptions,
  PointClickEventObject,
  SeriesColumnOptions,
  SeriesHeatmapDataOptions,
  SeriesLineOptions,
} from 'highcharts';
import { HeatmapOptions } from '../../../../_models/heatmap-options';
import { Trace } from '../../../../_models/trace';
import { take, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { EvaluationChartDataItem, EvaluationChartItem } from '../../../../_models/evaluation-chart-item';
import { IndicatorResult, Target } from '../../../../../../shared/interfaces/indicator-result';
import { AppUtils } from '../../../../_utils/app.utils';
import { getSliResultInfo, IEvaluationSelectionData, SliInfo, TChartType } from '../../ktb-evaluation-details-utils';
import { DateUtil } from '../../../../_utils/date.utils';

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
  selector: 'ktb-evaluation-chart-legacy[evaluationData][evaluationHistory]',
  templateUrl: './ktb-evaluation-chart-legacy.component.html',
  styleUrls: ['./ktb-evaluation-chart-legacy.component.scss'],
})
export class KtbEvaluationChartLegacyComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private _heatmapSeriesFull: HeatmapSeriesOptions[] = [];
  private _heatmapSeriesReduced: HeatmapSeriesOptions[] = [];
  private _heatmapCategoriesFull: string[] = [];
  private _heatmapCategoriesReduced: string[] = [];
  private heatmapChart?: DtChart;
  private _evaluationData: IEvaluationSelectionData = { shouldSelect: false };
  private _evaluationColor: { [key: string]: string } = {
    pass: '#7dc540',
    warning: '#e6be00',
    fail: '#dc172a',
    failed: '#dc172a',
    info: '#f8f8f8',
  };
  private _metrics: string[] = [];
  private _evaluationHistory: Trace[] = [];
  private selectedEvaluation?: Trace;
  public _heatmapSeries: HeatmapSeriesOptions[] = [];
  public _chartSeries: (SeriesColumnOptions | SeriesLineOptions)[] = [];
  public HeatmapPointClass = HeatmapPoint;
  public isHeatmapExtendable = false;
  public isHeatmapExtended = false;
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

  @Output() selectedEvaluationChange = new EventEmitter<Trace>();

  @ViewChild('heatmapChart') set heatmap(heatmap: DtChart) {
    this.heatmapChart = heatmap;
  }

  @Input()
  set evaluationData(evaluationData: IEvaluationSelectionData) {
    this._evaluationData = evaluationData;
    this.setEvaluation(evaluationData);
  }
  get evaluationData(): IEvaluationSelectionData {
    return this._evaluationData;
  }
  @Input() chartType: TChartType = 'heatmap';
  @Input()
  set evaluationHistory(evaluationHistory: Trace[]) {
    this._evaluationHistory = evaluationHistory;
    this.evaluationDataChanged();
  }
  get evaluationHistory(): Trace[] {
    return this._evaluationHistory;
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

  constructor(private zone: NgZone, public dateUtil: DateUtil, private _changeDetectorRef: ChangeDetectorRef) {}

  private evaluationDataChanged(): void {
    if (!this.selectedEvaluation && this.evaluationHistory) {
      this.setHeatmapDataAfterRender(this.evaluationHistory);
    }
  }

  private selectEvaluationData(evaluation?: Trace, forceSelect = false): void {
    if (this.evaluationData.shouldSelect || forceSelect) {
      this.selectedEvaluation = evaluation;
      this.selectedEvaluationChange.emit(evaluation);
      this.highlightHeatmap();
    }
  }

  private setEvaluation(evaluationData: IEvaluationSelectionData): void {
    if (this._evaluationData.evaluation?.id !== evaluationData.evaluation?.id) {
      this.selectedEvaluation = undefined;
      this.evaluationDataChanged();
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
        const trace = data.find((h) => h.shkeptncontext === this._evaluationData.evaluation?.shkeptncontext);
        this.selectEvaluationData(trace);
      });
  }

  private updateChartData(evaluationHistory: Trace[]): void {
    if (!this.selectedEvaluation && evaluationHistory) {
      this.selectEvaluationData(evaluationHistory.find((h) => h.id === this._evaluationData.evaluation?.id));
    }

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
                  : AppUtils.truncateNumber(
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

  private getSliResultInfos(chartSeries: EvaluationChartItem[]): SliInfoDictionary {
    const sliResultInfos: SliInfoDictionary = {};
    for (const chartItem of chartSeries) {
      for (const item of chartItem.data) {
        if (item.evaluationData?.data.evaluation?.indicatorResults && !sliResultInfos[item.evaluationData.id]) {
          const indicatorResults = item.evaluationData.data.evaluation.indicatorResults;
          sliResultInfos[item.evaluationData.id] = getSliResultInfo(indicatorResults);
        }
      }
    }
    return sliResultInfos;
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

  private updateHeatmapOptions(chartSeries: EvaluationChartItem[]): void {
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
      .map((item, _index, items) => {
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
  public seriesVisibilityChanged(_: DtChartSeriesVisibilityChangeEvent): void {
    // NOOP
  }

  public _chartSeriesClicked(event: PointClickEventObject & { point: { evaluationData: Trace } }): void {
    this.selectEvaluationData(event.point.evaluationData, true);
  }

  public _heatmapTileClicked(event: PointClickEventObject): void {
    this.selectEvaluationData(this._heatmapSeries[0].data[event.point.x].evaluation, true);
  }

  // remove duplicated points like "Score"
  public filterPoints(points: SeriesPoint[]): SeriesPoint[] {
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

  private highlightHeatmap(): void {
    if (this.selectedEvaluation) {
      const plotBands: NavigatorXAxisPlotBandsOptions[] = [];
      const highlightIndex = this._heatmapOptions.xAxis[0].categories.indexOf(
        this.selectedEvaluation.getHeatmapLabel()
      );
      if (highlightIndex >= 0) {
        plotBands.push({
          className: 'highlight-primary',
          from: highlightIndex - 0.5,
          to: highlightIndex + 0.5,
          zIndex: 100,
        });
      }
      this.setSecondaryHighlight(this.selectedEvaluation?.data.evaluation?.comparedEvents, plotBands);
      this._heatmapOptions.xAxis[0].plotBands = plotBands;
    } else {
      this._heatmapOptions.xAxis[0].plotBands = [];
    }
    this.heatmapChart?._update();
    this._changeDetectorRef.detectChanges();
  }

  private setSecondaryHighlight(
    comparedEvents: string[] | undefined,
    plotBands: NavigatorXAxisPlotBandsOptions[]
  ): void {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const _this = this;
    const secondaryHighlightIndices = comparedEvents
      ?.map((eventId: string) => this._heatmapSeries[0]?.data.findIndex((e) => e.evaluation?.id === eventId))
      .filter((eventIndex: number) => eventIndex >= 0);
    if (secondaryHighlightIndices) {
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
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
