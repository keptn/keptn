import * as Highcharts from "highcharts";

declare var require: any;
const Boost = require('highcharts/modules/boost');
const noData = require('highcharts/modules/no-data-to-display');
const More = require('highcharts/highcharts-more');
const Heatmap = require("highcharts/modules/heatmap");
const Treemap = require("highcharts/modules/treemap");


Boost(Highcharts);
noData(Highcharts);
More(Highcharts);
noData(Highcharts);
Heatmap(Highcharts);
Treemap(Highcharts);

import * as moment from 'moment';
import {ChangeDetectorRef, Component, Input, OnDestroy, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {MatDialog, MatDialogRef} from "@angular/material/dialog";
import {DtChart, DtChartSeriesVisibilityChangeEvent} from "@dynatrace/barista-components/chart";

import {Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";

import {ClipboardService} from '../../_services/clipboard.service';
import {DataService} from "../../_services/data.service";
import {DateUtil} from "../../_utils/date.utils";

import {Trace} from "../../_models/trace";

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss']
})
export class KtbEvaluationDetailsComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  public comparedIndicatorResults: any[] = [];
  @Input() public showChart = true;
  @Input() public isInvalidated = false;

  @ViewChild('sloDialog')
  public sloDialog: TemplateRef<any>;
  public sloDialogRef: MatDialogRef<any, any>;

  @ViewChild('invalidateEvaluationDialog')
  public invalidateEvaluationDialog: TemplateRef<any>;
  public invalidateEvaluationDialogRef: MatDialogRef<any, any>;

  public isHeatmapExtendable = false;
  public isHeatmapExtended = false;

  private heatmapChart: DtChart;

  @ViewChild('heatmapChart') set heatmap(heatmap: DtChart) {
    this.heatmapChart = heatmap;
  }

  public _evaluationColor = {
    pass: '#7dc540',
    warning: '#e6be00',
    fail: '#dc172a',
    failed: '#dc172a',
    info: '#f8f8f8'
  };

  public _evaluationState = {
    pass: 'recovered',
    warning: 'warning',
    fail: 'error',
    failed: 'error'
  };

  public _evaluationData: Trace;
  public _selectedEvaluationData: Trace;
  public _comparisonView: string = 'heatmap';
  private _metrics: string[];

  public _chartOptions: Highcharts.Options = {
    chart: {
      height: 400
    },
    legend: {
      maxHeight: 70
    },

    xAxis: {
      type: 'category',
      labels: {
        rotation: -45
      },
      categories: [],
    },
    yAxis: [
      {
        title: null,
        labels: {
          format: '{value}',
        },
        min: 0,
        max: 100,
      },
      {
        title: null,
        labels: {
          format: '{value}',
        },
        opposite: true
      },
    ],
    plotOptions: {
      column: {
        stacking: 'normal',
        pointWidth: 5,
        minPointLength: 2,
        point: {
          events: {
            click: (event) => {
              this._chartSeriesClicked(event);
              return true;
            }
          }
        },
      },
    },
  };
  public _chartSeries: Highcharts.SeriesOptions[] = [];

  public _heatmapOptions: Highcharts.Options = {
    chart: {
      type: 'heatmap',
      height: 400
    },
    xAxis: [{
      categories: [],
      plotBands: [],
      labels: {
        rotation: -45
      },
      tickPositioner: function() {
        const positions = [],
          labelWidth = 70,
          ext = this.getExtremes(),
          xMax = Math.round(ext.max),
          xMin = Math.round(ext.min),
          maxElements = (document.querySelector('dt-chart')?.clientWidth || labelWidth) / labelWidth,
          tick = Math.floor(xMax / maxElements) || 1;

        for (let i = xMax; i >= xMin; i -= tick) {
          positions.push(i);
        }
        return positions;
      }
    }],

    yAxis: [{
      categories: [],
      title: null,
      labels: {
        format: '{value}',
        style: {
          textOverflow: 'ellipsis',
          width: 200,
        }
      }
    }],

    colorAxis: {
      dataClasses: Object.keys(this._evaluationColor).filter(key => key != 'failed').map((key) => {
        return {color: this._evaluationColor[key], name: key}
      })
    },

    plotOptions: {
      heatmap: {
        point: {
          events: {
            click: (event) => {
              this._heatmapTileClicked(event);
              return true;
            }
          }
        },
      },
    },
  };
  public _heatmapSeries: Highcharts.SeriesHeatmapOptions[] = [];
  private _heatmapSeriesFull: Highcharts.SeriesHeatmapOptions[] = [];
  private _heatmapSeriesReduced: Highcharts.SeriesHeatmapOptions[] = [];
  private _heatmapCategoriesFull: string[];
  private _heatmapCategoriesReduced: string[];

  @Input()
  get evaluationData(): Trace {
    return this._evaluationData;
  }

  set evaluationData(evaluationData: Trace) {
    if (this._evaluationData !== evaluationData) {
      this._evaluationData = evaluationData;
      this._chartSeries = [];
      this._heatmapSeries = [];
      this._metrics = ['Score'];
      this._heatmapOptions.yAxis[0].categories = ['Score'];
      this._selectedEvaluationData = null;
      this.evaluationDataChanged();
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private dialog: MatDialog, private clipboard: ClipboardService, public dateUtil: DateUtil) {
  }

  ngOnInit() {
    this.dataService.evaluationResults
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((results) => {
        if (results.type == "evaluationHistory" && results.triggerEvent == this.evaluationData) {
          this.evaluationData.data.evaluationHistory = [...results.traces || [], ...this.evaluationData.data.evaluationHistory || []].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
          this.updateChartData(this.evaluationData.data.evaluationHistory);
        } else if (results.type == "invalidateEvaluation" &&
          this.evaluationData.data.project == results.triggerEvent.data.project &&
          this.evaluationData.data.service == results.triggerEvent.data.service &&
          this.evaluationData.data.stage == results.triggerEvent.data.stage) {
          this.evaluationData.data.evaluationHistory = this.evaluationData.data.evaluationHistory.filter(e => e.id != results.triggerEvent.id);
          this._selectedEvaluationData = null;
          this.updateChartData(this.evaluationData.data.evaluationHistory);
        }
      });
  }

  private evaluationDataChanged() {
    if (this._evaluationData) {
      this.dataService.loadEvaluationResults(this._evaluationData);
      if (this.isInvalidated)
        this.selectEvaluationData(this._evaluationData);
      else if (!this._selectedEvaluationData && this._evaluationData.data.evaluationHistory)
        this.selectEvaluationData(this._evaluationData.data.evaluationHistory.find(h => h.shkeptncontext === this._evaluationData.shkeptncontext));
    }
  }

  private parseSloFile(evaluationData) {
    if (evaluationData && evaluationData.data && evaluationData.data.evaluation.sloFileContent && !evaluationData.data.evaluation.sloFileContentParsed) {
      evaluationData.data.evaluation.sloFileContentParsed = atob(evaluationData.data.evaluation.sloFileContent);
      evaluationData.data.evaluation.score_pass = evaluationData.data.evaluation.sloFileContentParsed.split("total_score:")[1]?.split("pass:")[1]?.split(" ")[1]?.replace(/\"/g, "")?.split("%")[0];
      evaluationData.data.evaluation.score_warning = evaluationData.data.evaluation.sloFileContentParsed.split("total_score:")[1]?.split("warning:")[1]?.split(" ")[1]?.replace(/\"/g, "")?.split("%")[0];
      evaluationData.data.evaluation.compare_with = evaluationData.data.evaluation.sloFileContentParsed.split("comparison:")[1]?.split("compare_with:")[1]?.split(" ")[1]?.replace(/\"/g, "");
      evaluationData.data.evaluation.include_result_with_score = evaluationData.data.evaluation.sloFileContentParsed.split("comparison:")[1]?.split("include_result_with_score:")[1]?.split(" ")[1]?.replace(/\"/g, "");
      if (evaluationData.data.evaluation.comparedEvents !== null) {
        evaluationData.data.evaluation.number_of_comparison_results = evaluationData.data.evaluation.comparedEvents?.length;
      } else {
        evaluationData.data.evaluation.number_of_comparison_results = 0;
      }
    }
  }

  updateChartData(evaluationHistory) {
    const chartSeries = [];

    if (!this._selectedEvaluationData && evaluationHistory)
      this.selectEvaluationData(evaluationHistory.find(h => h.id === this._evaluationData.id));

    if (this.showChart) {
      evaluationHistory.forEach((evaluation) => {
        let scoreData = {
          y: evaluation.data.evaluation ? evaluation.data.evaluation.score : 0,
          evaluationData: evaluation,
          color: this._evaluationColor[evaluation.data.evaluation.result],
          name: evaluation.getChartLabel(),
        };

        let indicatorScoreSeriesColumn = chartSeries.find(series => series.name == 'Score' && series.type == 'column');
        let indicatorScoreSeriesLine = chartSeries.find(series => series.name == 'Score' && series.type == 'line');
        if (!indicatorScoreSeriesColumn) {
          indicatorScoreSeriesColumn = {
            metricName: 'Score',
            name: 'Score',
            type: 'column',
            data: [],
            cursor: 'pointer',
            turboThreshold: 0
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
            turboThreshold: 0
          };
          chartSeries.push(indicatorScoreSeriesLine);
        }

        indicatorScoreSeriesColumn.data.push(scoreData);
        indicatorScoreSeriesLine.data.push(scoreData);

        if (evaluation.data.evaluation.indicatorResults) {
          evaluation.data.evaluation.indicatorResults.forEach((indicatorResult) => {
            const indicatorData = {
              y: indicatorResult.value.value,
              indicatorResult: indicatorResult,
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
                turboThreshold: 0
              };
              chartSeries.push(indicatorChartSeries);
            }
            indicatorChartSeries.data.push(indicatorData);
          });
        }
      });
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

      this.updateHeatmapOptions(chartSeries);

      // @ts-ignore
      this._chartOptions.xAxis.categories = this._heatmapOptions.xAxis[0].categories;
      this._heatmapSeriesReduced = [
        {
          name: 'Score',
          type: 'heatmap',
          rowsize: 0.85,
          turboThreshold: 0,
          data: []
        },
        {
          name: 'SLOs',
          type: 'heatmap',
          turboThreshold: 0,
          data: []
        }
      ];
      this._heatmapSeriesFull = [
        {
          name: 'Score',
          type: 'heatmap',
          rowsize: 0.85,
          turboThreshold: 0,
          data: chartSeries.find(series => series.name === 'Score').data.map((s) => {
            const index = this._metrics.indexOf('Score');
            const x = this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData.getHeatmapLabel());
            const dataPoint = {
              x: x,
              y: index,
              z: s.y,
              evaluation: s.evaluationData,
              color: this._evaluationColor[s.evaluationData.data.result],
            };
            const reducedDataPoint = {...dataPoint};
            reducedDataPoint.y = 9;
            this._heatmapSeriesReduced[0].data.push(reducedDataPoint);
            return dataPoint;
          })
        },
        {
          name: 'SLOs',
          type: 'heatmap',
          turboThreshold: 0,
          data: [...chartSeries].reverse().reduce((r, d, i) => [...r, ...d.data.filter(s => s.indicatorResult).map((s) => {
            const index = this._metrics.indexOf(s.indicatorResult.value.metric);
            const x = this._heatmapOptions.xAxis[0].categories.indexOf(s.evaluationData.getHeatmapLabel());
            const dataPoint = {
              x: x,
              y: index,
              z: s.indicatorResult.score,
              color: s.indicatorResult.value.success ? this._evaluationColor[s.indicatorResult.status] : this._evaluationColor['info']
            };
            if(i < 9) {
              this._heatmapSeriesReduced[1].data.push(dataPoint);
            }
            return dataPoint;
          })], [])
        },
      ];


      chartSeries.forEach(item => {
        item.data.forEach(data => {
          data.x = this._heatmapOptions.xAxis[0].categories.indexOf(data.evaluationData.getHeatmapLabel());
        });
      });
      this._chartSeries = chartSeries;

      if(this.isHeatmapExtendable) {
        this._heatmapSeries = this._heatmapSeriesReduced;
      } else {
        this._heatmapSeries = this._heatmapSeriesFull;
      }
    }

    this.highlightHeatmap();
    this._changeDetectorRef.detectChanges();
  }

  private getLastDisplayName(evaluationHistory, metric): string {
    let displayName = metric;
    if (metric !== 'Score') {
      for (let i = evaluationHistory.length - 1; i >= 0; i--) {
        const result = evaluationHistory[i].data.evaluation.indicatorResults?.find(indicatorResult => indicatorResult.value.metric === metric);
        if (result) {
          displayName = result.displayName || result.value.metric;
          break;
        }
      }
    }
    return displayName;
  }

  updateHeatmapOptions(chartSeries) {
    const heatmapCategoriesFull = [...this._heatmapOptions.yAxis[0].categories];
    const heatmapCategoriesReduced = [...this._heatmapOptions.yAxis[0].categories];
    chartSeries.forEach((series, i) => {
      if (!this._metrics.includes(series.metricName)) {
        heatmapCategoriesFull.unshift(series.name);
        if(i <= 10) {
          heatmapCategoriesReduced.unshift(series.name);
        }
        this._metrics.unshift(series.metricName);
      }
      if (series.name == "Score") {
        let categories = series.data
          .sort((a, b) => moment(a.evaluationData.time).unix() - moment(b.evaluationData.time).unix())
          .map((item, index, items) => {
            let duplicateItems = items.filter(c => c.evaluationData.getHeatmapLabel() == item.evaluationData.getHeatmapLabel());
            if (duplicateItems.length > 1)
              item.label = `${item.evaluationData.getHeatmapLabel()} (${duplicateItems.indexOf(item) + 1})`;
            else
              item.label = item.evaluationData.getHeatmapLabel();
            return item;
          })
          .map((item) => {
            item.evaluationData.setHeatmapLabel(item.label);
            return item.evaluationData.getHeatmapLabel();
          });

        this._heatmapOptions.xAxis[0].categories = categories;
      }
    });

    this._heatmapCategoriesFull = heatmapCategoriesFull;
    this._heatmapCategoriesReduced = heatmapCategoriesReduced;


    if(this._heatmapCategoriesFull.length > 10) {
      this.isHeatmapExtendable = true;
      this.isHeatmapExtended = false;
    } else {
      this.isHeatmapExtended = true;
    }
    this._updateHeatmapExtension();
  }



  seriesVisibilityChanged(_: DtChartSeriesVisibilityChangeEvent): void {
    // NOOP
  }

  _chartSeriesClicked(event) {
    this.selectEvaluationData(event.point.evaluationData);
  }

  _heatmapTileClicked(event) {
    this.selectEvaluationData(this._heatmapSeries[0].data[event.point.x]['evaluation']);
  }

  selectEvaluationData(evaluation) {
    this.parseSloFile(evaluation);
    this._selectedEvaluationData = evaluation;
    this.highlightHeatmap();
  }

  highlightHeatmap() {
    if (this._selectedEvaluationData && !this.isInvalidated) {
      const _this = this;
      const highlightIndex = this._heatmapOptions.xAxis[0].categories.indexOf(this._selectedEvaluationData.getHeatmapLabel());
      const secondaryHighlightIndexes = this._selectedEvaluationData?.data.evaluation.comparedEvents?.map(eventId => this._heatmapSeries[0]?.data.findIndex(e => e['evaluation'].id == eventId));
      const plotBands = [];
      if (highlightIndex >= 0)
        plotBands.push({
          className: 'highlight-primary',
          from: highlightIndex - 0.5,
          to: highlightIndex + 0.5,
          zIndex: 100
        });
      if(secondaryHighlightIndexes) {
        const index = secondaryHighlightIndexes.find(index => index >= 0);
        this.comparedIndicatorResults = this._heatmapSeries[0]?.data[index]['evaluation'].data.evaluation.indicatorResults ?? [];

        secondaryHighlightIndexes.forEach(highlightIndex => {
          if (highlightIndex >= 0)
            plotBands.push({
              className: 'highlight-secondary',
              from: highlightIndex - 0.5,
              to: highlightIndex + 0.5,
              zIndex: 100,
              events: {
                click: function () {
                  let index = this.options.from + 0.5;
                  setTimeout(() => {
                    _this.selectEvaluationData(_this._heatmapSeries[0].data[index]['evaluation']);
                  });
                }
              }
            });
        });
      }
      else {
        this.comparedIndicatorResults = [];
      }
      this._heatmapOptions.xAxis[0].plotBands = plotBands;
      this._selectedEvaluationData.data.evaluation.number_of_missing_comparison_results = this._selectedEvaluationData?.data.evaluation.comparedEvents?.length - (this._heatmapOptions.xAxis[0].plotBands?.length - 1);
    } else {
      this._heatmapOptions.xAxis[0].plotBands = [];
      if (this._selectedEvaluationData) {
        this._selectedEvaluationData.data.evaluation.number_of_missing_comparison_results = 0;
      }
    }
    this.heatmapChart?._update();
    this._changeDetectorRef.detectChanges();
  }

  showSloDialog() {
    this.sloDialogRef = this.dialog.open(this.sloDialog, {data: this._selectedEvaluationData.data.evaluation.sloFileContentParsed});
  }

  closeSloDialog() {
    if (this.sloDialogRef) {
      this.sloDialogRef.close();
    }
  }

  copySloPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'slo payload');
  }

  invalidateEvaluationTrigger() {
    this.invalidateEvaluationDialogRef = this.dialog.open(this.invalidateEvaluationDialog, {data: this._selectedEvaluationData});
  }

  invalidateEvaluation(evaluation, reason) {
    this.dataService.invalidateEvaluation(evaluation, reason);
    this.closeInvalidateEvaluationDialog();
  }

  closeInvalidateEvaluationDialog() {
    if (this.invalidateEvaluationDialogRef)
      this.invalidateEvaluationDialogRef.close();
  }

  // remove duplicated points like "Score"
  filterPoints(points: any[]): any[] {
    return points.filter((item, index) => index === points.findIndex(subItem => subItem.series.name === item.series.name));
  }

  public toggleHeatmap() {
    this.isHeatmapExtended = !this.isHeatmapExtended;
    this._updateHeatmapExtension();
  }

  private _updateHeatmapExtension() {
    if(this.isHeatmapExtended) {
      this._heatmapSeries = this._heatmapSeriesFull;
      this._heatmapOptions.yAxis[0].categories = this._heatmapCategoriesFull;
      this._heatmapOptions.chart.height = this._heatmapCategoriesFull.length * 28 + 160;
    } else {
      this._heatmapSeries = this._heatmapSeriesReduced;
      this._heatmapOptions.yAxis[0].categories = this._heatmapCategoriesReduced;
      this._heatmapOptions.chart.height = this._heatmapCategoriesReduced.length * 28 + 173;
    }
    if(this.isHeatmapExtendable) {
      this._heatmapOptions.xAxis[0].offset = 40;
    } else {
      this._heatmapOptions.xAxis[0].offset = undefined;
    }

    this.heatmapChart._update();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
