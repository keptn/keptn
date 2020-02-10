import * as Highcharts from "highcharts";

import * as moment from 'moment';
import {
  ChangeDetectorRef,
  Component,
  Input,
  KeyValueChanges,
  KeyValueDiffer,
  KeyValueDiffers,
  OnInit
} from '@angular/core';
import {DtChartSeriesVisibilityChangeEvent} from "@dynatrace/barista-components/chart";

import {DataService} from "../../_services/data.service";
import DateUtil from "../../_utils/date.utils";

@Component({
  selector: 'ktb-evaluation-details',
  templateUrl: './ktb-evaluation-details.component.html',
  styleUrls: ['./ktb-evaluation-details.component.scss']
})
export class KtbEvaluationDetailsComponent implements OnInit {

  public _evaluationData: any;
  public _evaluationSource: string;

  public _selectedEvaluationData: any;

  private _dataDiffer: KeyValueDiffer<string, any>;
  public _view: string = "singleevaluation";

  public _chartOptions: Highcharts.Options = {
    xAxis: {
      type: 'datetime',
    },
    yAxis: [
      {
        title: null,
        labels: {
          format: '{value}',
        },
        tickInterval: 10,
      }
    ],
    plotOptions: {
      column: {
        stacking: 'normal',
      },
      series: {
        marker: {
          enabled: false,
        },
        point: {
          events: {
            click: (event) => {
              this._chartSeriesClicked(event);
              return true;
            }
          }
        }
      },
    },
  };
  public _chartSeries: Highcharts.IndividualSeriesOptions[] = [
    {
      name: 'Evaluation passed',
      type: 'column',
      data: [],
      color: '#006bb8'
    },
    {
      name: 'Evaluation failed',
      type: 'column',
      data: [],
      color: '#c41425'
    },
    {
      name: 'Evaluation score',
      type: 'line',
      data: [],
      color: '#006bb8'
    },
  ];

  @Input()
  get evaluationData(): any {
    return this._evaluationData;
  }
  set evaluationData(evaluationData: any) {
    if (this._evaluationData !== evaluationData) {
      this._evaluationData = evaluationData;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get evaluationSource(): any {
    return this._evaluationSource;
  }
  set evaluationSource(evaluationSource: any) {
    if (this._evaluationSource !== evaluationSource) {
      this._evaluationSource = evaluationSource;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private differs: KeyValueDiffers) { }

  ngOnInit() {
    this._dataDiffer = this.differs.find(this._evaluationData).create();
  }

  updateChartData(evaluationHistory) {
    let evaluationPassed = [];
    let evaluationFailed = [];
    let evaluationData = [];
    evaluationHistory.forEach((evaluation) => {
      let data = {
        x: moment(evaluation.time).unix()*1000,
        y: evaluation.data.evaluationdetails ? evaluation.data.evaluationdetails.score : 0,
        evaluationData: evaluation
      };
      evaluationData.push(data);
      if(evaluation.data.result == 'pass')
        evaluationPassed.push(data);
      else
        evaluationFailed.push(data);
    });
    this._chartSeries = [
      {
        name: 'Evaluation passed',
        type: 'column',
        data: evaluationPassed,
        color: '#006bb8'
      },
      {
        name: 'Evaluation failed',
        type: 'column',
        data: evaluationFailed,
        color: '#c41425'
      },
      {
        name: 'Evaluation score',
        type: 'line',
        data: evaluationData,
        color: '#006bb8'
      },
    ];
  }

  switchEvaluationView(event) {
    this._view = this._view == "singleevaluation" ? "evaluationcomparison" : "singleevaluation";
    if(this._view == "evaluationcomparison") {
      this.dataService.loadEvaluationResults(this._evaluationData, this._evaluationSource);
      this._changeDetectorRef.markForCheck();
    }
  }

  seriesVisibilityChanged(_: DtChartSeriesVisibilityChangeEvent): void {
    // NOOP
  }

  dataChanged(changes: KeyValueChanges<string, any>) {
    changes.forEachRemovedItem((record) => {
      if(record.key == "evaluationHistory") {
        this.updateChartData(record.currentValue);
      }
    });
    changes.forEachAddedItem((record) => {
      if(record.key == "evaluationHistory") {
        this.updateChartData(record.currentValue);
      }
    });
    changes.forEachChangedItem((record) => {
      if(record.key == "evaluationHistory") {
        this.updateChartData(record.currentValue);
      }
    });
  }

  ngDoCheck(): void {
    const changes = this._dataDiffer.diff(this._evaluationData);
    if (changes) {
      this.dataChanged(changes);
    }
  }

  _chartSeriesClicked(event: Event): boolean {
    this._selectedEvaluationData = event.point.evaluationData.data;
    console.log("data", this._selectedEvaluationData);
    return true;
  }

  getCalendarFormat() {
    return DateUtil.getCalendarFormats().sameElse;
  }

}
