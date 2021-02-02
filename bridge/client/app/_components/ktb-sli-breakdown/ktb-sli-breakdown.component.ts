import {ChangeDetectorRef, Component, Input, OnInit, ViewChild} from '@angular/core';
import {DtOverlayConfig} from "@dynatrace/barista-components/overlay";
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss']
})
export class KtbSliBreakdownComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable: DtSort;

  public _evaluationColor = {
    'pass': '#7dc540',
    'warning': '#fd8232',
    'fail': '#dc172a',
    'info': '#f8f8f8'
  };

  public _evaluationState = {
    'pass': 'recovered',
    'warning': 'warning',
    'fail': 'error'
  };

  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
    originY: 'center',
  };

  public _indicatorResults: any;
  public _indicatorResultsFail: any = [];
  public _indicatorResultsWarning: any = [];
  public _indicatorResultsPass: any = [];
  public tableEntries: DtTableDataSource<object>;

  @Input()
  get indicatorResults(): any {
    return [...this._indicatorResultsFail, ...this._indicatorResultsWarning, ...this._indicatorResultsPass];
  }
  set indicatorResults(indicatorResults: any) {
    if (this._indicatorResults !== indicatorResults) {
      this._indicatorResults = indicatorResults;
      this._indicatorResultsFail = indicatorResults.filter(i => i.status === 'fail');
      this._indicatorResultsWarning = indicatorResults.filter(i => i.status === 'warning');
      this._indicatorResultsPass = indicatorResults.filter(i => i.status !== 'fail' && i.status !== 'warning');
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
    this.tableEntries = new DtTableDataSource(this.assembleTablesEntries(this.indicatorResults));

    // Set the dtSort reference on the dataSource, so it can react to sorting.
    this.tableEntries.sort = this.sortable;
  }

  formatNumber(number) {
    let n = number;
    if(n < 1)
      n = Math.floor(n*1000)/1000;
    else if(n < 100)
      n = Math.floor(n*100)/100;
    else if(n < 1000)
      n = Math.floor(n*10)/10;
    else
      n = Math.floor(n);

    return n;
  }

  assembleTablesEntries(indicatorRestults): any {
    var tableEntries = [];

    let totalscore  = 0;
    indicatorRestults.forEach(result => totalscore += result.score);

    
    for (let indicatorRestult of indicatorRestults) {
      let name = indicatorRestult.value.metric;
      let value = this.formatNumber(indicatorRestult.value.value);
      let result = indicatorRestult.status;
      let score = (indicatorRestult.score / totalscore).toFixed(2);
      let criteria = "";

      for (let target of indicatorRestult.targets) {
        criteria += target.criteria + ' '
      }

      let entry = {
        name: name,
        value: value,
        result: result,
        score: score,
        criteria: criteria
      }

      tableEntries.push(entry);
    } 

    return tableEntries;
  }

}
