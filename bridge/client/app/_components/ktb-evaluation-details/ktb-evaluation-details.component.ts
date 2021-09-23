import Highcharts, { NavigatorXAxisPlotBandsOptions, PointClickEventObject, SeriesColumnOptions, SeriesHeatmapDataOptions, SeriesLineOptions } from 'highcharts';
import { ChangeDetectorRef, Component, Input, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { DtChart, DtChartOptions, DtChartSeries, DtChartSeriesVisibilityChangeEvent } from '@dynatrace/barista-components/chart';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Trace } from '../../_models/trace';
import { EvaluationChartDataItem, EvaluationChartItem } from '../../_models/evaluation-chart-item';
import { HeatmapOptions } from '../../_models/heatmap-options';
import { HeatmapData, HeatmapSeriesOptions } from '../../_models/heatmap-series-options';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { ResultTypes } from '../../../../shared/models/result-types';

// tslint:disable-next-line:no-any
declare var require: any;
const Boost = require('highcharts/modules/boost');
const noData = require('highcharts/modules/no-data-to-display');
const More = require('highcharts/highcharts-more');
const Heatmap = require('highcharts/modules/heatmap');
const Treemap = require('highcharts/modules/treemap');
type SeriesPoint = PointClickEventObject & { series: EvaluationChartItem, point: { evaluationData: Trace } };


Boost(Highcharts);
noData(Highcharts);
More(Highcharts);
noData(Highcharts);
Heatmap(Highcharts);
Treemap(Highcharts);

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss'],
})
export class KtbEvaluationDetailsComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  public comparedIndicatorResults: IndicatorResult[] = [];
  @Input() public showChart = true;
  @Input() public isInvalidated = false;

  @ViewChild('sloDialog')
  // tslint:disable-next-line:no-any
  public sloDialog?: TemplateRef<any>;
  // tslint:disable-next-line:no-any
  public sloDialogRef?: MatDialogRef<any, any>;

  @ViewChild('invalidateEvaluationDialog')
  // tslint:disable-next-line:no-any
  public invalidateEvaluationDialog?: TemplateRef<any>;
  // tslint:disable-next-line:no-any
  public invalidateEvaluationDialogRef?: MatDialogRef<any, any>;

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
            click: (event: PointClickEventObject) => {
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
    xAxis: [{
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
    }],

    yAxis: [{
      categories: [],
      title: undefined,
      labels: {
        format: '{value}',
        style: {
          textOverflow: 'ellipsis',
          width: 200,
        },
      },
    }],

    colorAxis: {
      dataClasses: Object.keys(this._evaluationColor).filter(key => key !== 'failed').map((key) => {
        return {color: this._evaluationColor[key], name: key};
      }),
    },

    plotOptions: {
      heatmap: {
        point: {
          events: {
            click: (event) => {
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

  @Input()
  get evaluationData(): Trace | undefined {
    return this._evaluationData;
  }

  set evaluationData(evaluationData: Trace | undefined) {
    this.setEvaluation({evaluation: evaluationData, shouldSelect: true});
  }

  @Input()
  set evaluationInfo(evaluationInfo: { evaluation?: Trace, shouldSelect: boolean }) {
    this.setEvaluation(evaluationInfo);
  }

  private setEvaluation(evaluationInfo: { evaluation?: Trace, shouldSelect: boolean }): void {
    if (this._evaluationData !== evaluationInfo.evaluation) {
      this._selectedEvaluationData = evaluationInfo.evaluation?.id === this._evaluationData?.id ? this._selectedEvaluationData : undefined;
      this._evaluationData = evaluationInfo.evaluation;
      this._chartSeries = [];
      this._metrics = ['Score'];
      this._heatmapOptions.yAxis[0].categories = ['Score'];
      this._shouldSelectEvaluation = evaluationInfo.shouldSelect && !this._selectedEvaluationData;
      this.evaluationDataChanged();
      this._changeDetectorRef.markForCheck();
    }
  }

  get heatmapSeries(): DtChartSeries[] {
    // type 'heatmap' does not exist in barista components but in highcharts
    // @ts-ignore
    return this._heatmapSeries as DtChartSeries[];
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService,
              private dialog: MatDialog, private clipboard: ClipboardService, public dateUtil: DateUtil) {
  }

  public ngOnInit(): void {
    this.dataService.evaluationResults
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((results) => {
        if (this.evaluationData) {
          if (results.type === 'evaluationHistory' && results.triggerEvent === this.evaluationData) {
            this.evaluationData.data.evaluationHistory = [...results.traces || [],
              ...this.evaluationData.data.evaluationHistory || []]
              .sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
          } else if (results.type === 'invalidateEvaluation' &&
            this.evaluationData.data.project === results.triggerEvent.data.project &&
            this.evaluationData.data.service === results.triggerEvent.data.service &&
            this.evaluationData.data.stage === results.triggerEvent.data.stage) {
            this.evaluationData.data.evaluationHistory = this.evaluationData.data.evaluationHistory
              ?.filter(e => e.id !== results.triggerEvent.id);
          }
          this._selectedEvaluationData = this._selectedEvaluationData?.id
            ? this.evaluationData.data.evaluationHistory?.find(h => h.id === this._selectedEvaluationData?.id)
            : undefined;
          if (this.evaluationData.data.evaluationHistory) {
            this.updateChartData(this.evaluationData.data.evaluationHistory);
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
        const trace = this._evaluationData.data.evaluationHistory.find(h => h.shkeptncontext === this._evaluationData?.shkeptncontext);
        this.selectEvaluationData(trace);
      }
    }
  }

  private parseSloFile(evaluationData?: Trace): void {
    if (evaluationData?.data?.evaluation?.sloFileContent && !evaluationData.data.evaluation.sloFileContentParsed) {
      evaluationData.data.evaluation.sloFileContentParsed = atob(evaluationData.data.evaluation.sloFileContent);
      evaluationData.data.evaluation.score_pass = evaluationData.data.evaluation.sloFileContentParsed
        .split('total_score:')[1]?.split('pass:')[1]
        ?.split(' ')[1]?.replace(/"/g, '')?.split('%')[0];
      evaluationData.data.evaluation.score_warning = evaluationData.data.evaluation.sloFileContentParsed
        .split('total_score:')[1]?.split('warning:')[1]
        ?.split(' ')[1]?.replace(/"/g, '')?.split('%')[0];
      evaluationData.data.evaluation.compare_with = evaluationData.data.evaluation.sloFileContentParsed
        .split('comparison:')[1]?.split('compare_with:')[1]
        ?.split(' ')[1]?.replace(/"/g, '');
      evaluationData.data.evaluation.include_result_with_score = evaluationData.data.evaluation.sloFileContentParsed
        .split('comparison:')[1]?.split('include_result_with_score:')[1]
        ?.split(' ')[1]?.replace(/"/g, '');
      if (evaluationData.data.evaluation.comparedEvents) {
        evaluationData.data.evaluation.number_of_comparison_results = evaluationData.data.evaluation.comparedEvents?.length;
      } else {
        evaluationData.data.evaluation.number_of_comparison_results = 0;
      }
    }
  }

  updateChartData(evaluationHistory: Trace[]): void {
    if (!this._selectedEvaluationData && evaluationHistory) {
      this.selectEvaluationData(evaluationHistory.find(h => h.id === this._evaluationData?.id));
    }

    if (this.showChart) {
      const chartSeries = this.getChartSeries(evaluationHistory);
      this.sortChartSeries(chartSeries);
      this.updateHeatmapOptions(chartSeries);

      // @ts-ignore
      this._chartOptions.xAxis.categories = this._heatmapOptions.xAxis[0].categories;
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
    const minIdx = ((this._heatmapSeriesFull[1].data[this._heatmapSeriesFull[1].data.length - 1] as SeriesHeatmapDataOptions)
        .y ?? 0
    ) - 8;
    const reduced: HeatmapData[] = [];
    for (const series of this._heatmapSeriesFull[1].data) {
      if (series.y >= minIdx) {
        const srs = {...series};
        srs.y = (srs.y - minIdx);
        reduced.push(srs);
      }
    }
    this._heatmapSeriesReduced[1].data = reduced;
  }

  private setSeriesXAxis(chartSeries: EvaluationChartItem[]): void {
    for (const item of chartSeries) {
      for (const data of item.data) {
        data.x = data.evaluationData ? this._heatmapOptions.xAxis[0].categories.indexOf(data.evaluationData.getHeatmapLabel()) : -1;
      }
    }
  }

  private setHeatmapData(chartSeries: EvaluationChartItem[]): void {
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
        data: chartSeries.find(series => series.name === 'Score')?.data.filter(s => s.evaluationData).map(s => {
          // tslint:disable:no-non-null-assertion
          const index = this._metrics.indexOf('Score');
          const x = this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData!.getHeatmapLabel());
          const dataPoint = {
            x,
            y: index,
            z: s.y,
            evaluation: s.evaluationData,
            color: this._evaluationColor[s.evaluationData!.data.result ?? 'info'],
          };
          const reducedDataPoint = {...dataPoint};
          reducedDataPoint.y = 9;
          this._heatmapSeriesReduced[0].data.push(reducedDataPoint);
          return dataPoint;
          // tslint:enable:no-non-null-assertion
        }) ?? [],
      },
      {
        name: 'SLOs',
        type: 'heatmap',
        turboThreshold: 0,
        data: [...chartSeries].reverse().reduce((r, d) => [...r, ...d.data.filter(s => s.indicatorResult).map(s => {
          // tslint:disable:no-non-null-assertion
          const index = this._metrics.indexOf(s.indicatorResult!.value.metric);
          const x = s.evaluationData ? this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData.getHeatmapLabel()) : -1;
          return {
            x,
            y: index,
            z: s.indicatorResult!.score,
            color: s.indicatorResult!.value.success ? this._evaluationColor[s.indicatorResult!.status] : this._evaluationColor.info,
          };
          // tslint:enable:no-non-null-assertion
        })], [] as HeatmapData[]),
      },
    ];
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

      let indicatorScoreSeriesColumn = chartSeries.find(series => series.name === 'Score' && series.type === 'column');
      let indicatorScoreSeriesLine = chartSeries.find(series => series.name === 'Score' && series.type === 'line');
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

          let indicatorChartSeries = chartSeries.find(series => series.metricName === indicatorResult.value.metric);
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
        const result = evaluationHistory[i].data.evaluation?.indicatorResults
          ?.find(indicatorResult => indicatorResult.value.metric === metric);
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
        const duplicateItems = items.filter(c => c.evaluationData?.getHeatmapLabel() === item.evaluationData?.getHeatmapLabel());
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
      this.parseSloFile(evaluation);
      this._selectedEvaluationData = evaluation;
      this.highlightHeatmap();
    }
  }

  highlightHeatmap(): void {
    if (this._selectedEvaluationData && !this.isInvalidated) {
      const highlightIndex = this._heatmapOptions.xAxis[0].categories.indexOf(this._selectedEvaluationData.getHeatmapLabel());
      const secondaryHighlightIndexes = this._selectedEvaluationData?.data.evaluation?.comparedEvents
        ?.map(eventId => this._heatmapSeries[0]?.data.findIndex(e => e.evaluation?.id === eventId));
      const plotBands: NavigatorXAxisPlotBandsOptions[] = [];
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
      } else {
        this.comparedIndicatorResults = [];
      }
      this._heatmapOptions.xAxis[0].plotBands = plotBands;
      if (this._selectedEvaluationData.data.evaluation?.number_of_missing_comparison_results
        && this._selectedEvaluationData?.data.evaluation.comparedEvents?.length !== undefined) {
        this._selectedEvaluationData.data.evaluation.number_of_missing_comparison_results =
          this._selectedEvaluationData?.data.evaluation.comparedEvents.length - (this._heatmapOptions.xAxis[0].plotBands?.length - 1);
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

  private setSecondaryHighlight(secondaryHighlightIndices: number[], plotBands: NavigatorXAxisPlotBandsOptions[]): void {
    const _this = this;
    const index = secondaryHighlightIndices.find(idx => idx >= 0) ?? -1;
    this.comparedIndicatorResults = index >= 0 ? this._heatmapSeries[0]?.data[index].evaluation?.data.evaluation?.indicatorResults ?? [] : [];
    for (const secondaryHighlightIndex of secondaryHighlightIndices) {
      if (secondaryHighlightIndex >= 0) {
        plotBands.push({
          className: 'highlight-secondary',
          from: secondaryHighlightIndex - 0.5,
          to: secondaryHighlightIndex + 0.5,
          zIndex: 100,
          events: {
            click(): void {
              // @ts-ignore
              const idx = this.options.from + 0.5;
              setTimeout(() => {
                _this.selectEvaluationData(_this._heatmapSeries[0]?.data[idx]?.evaluation);
              });
            },
          },
        });
      }
    }
  }

  showSloDialog(): void {
    if (this.sloDialog && this._selectedEvaluationData) {
      this.sloDialogRef = this.dialog.open(this.sloDialog, {data: this._selectedEvaluationData.data.evaluation?.sloFileContentParsed});
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
      this.invalidateEvaluationDialogRef = this.dialog.open(this.invalidateEvaluationDialog, {data: this._selectedEvaluationData});
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
    return points.filter((item, index) => index === points.findIndex(subItem => subItem.series.name === item.series.name));
  }

  public getEvaluationFromPoint(tooltip: { points: SeriesPoint[] }): Trace {
    return tooltip.points[0].point.evaluationData;
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
