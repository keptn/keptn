import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Input,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {DtTableDataSource} from "@dynatrace/barista-components/table";
import {Subject} from "rxjs";

import {Service} from "../../_models/service";
import {DateUtil} from "../../_utils/date.utils";
import {DataService} from "../../_services/data.service";
import {takeUntil} from "rxjs/operators";

const DEFAULT_PAGE_SIZE = 3;

@Component({
  selector: 'ktb-services-list',
  templateUrl: './ktb-services-list.component.html',
  styleUrls: ['./ktb-services-list.component.scss'],
  host: {
    class: 'ktb-services-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbServicesListComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public _services: Service[] = [];
  public _pageSize: number = DEFAULT_PAGE_SIZE;
  public dataSource: DtTableDataSource<Service>;

  @Input()
  get services(): Service[] {
    return this._services;
  }
  set services(value: Service[]) {
    if (this._services !== value) {
      this._services = value;
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  get pageSize(): number {
    return this._pageSize;
  }
  set pageSize(value: number) {
    if (this._pageSize !== value) {
      this._pageSize = value;
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  get DEFAULT_PAGE_SIZE(): number {
    return DEFAULT_PAGE_SIZE;
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dataService: DataService, public dateUtil: DateUtil) { }

  ngOnInit() {
    this.dataService.roots
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(roots => {
        this.updateDataSource();
        this._changeDetectorRef.markForCheck();
      });
  }

  updateDataSource() {
    let data = this.services.sort((a, b) => {
      if(!a.getRecentSequence())
        return 1;
      else if(!b.getRecentSequence())
        return -1;
      else
        return DateUtil.compareTraceTimesAsc(a.getRecentSequence().getLastTrace(), b.getRecentSequence().getLastTrace());
    });
    this.dataSource = new DtTableDataSource(data.slice(0, this.pageSize));
  }

  toggleAllServices() {
    if(this.services.length > this.pageSize) {
      this.pageSize = this.services.length;
    } else if(this.pageSize > DEFAULT_PAGE_SIZE) {
      this.pageSize = DEFAULT_PAGE_SIZE;
    }
  }

  getServiceLink(service) {
    return ['service', service.serviceName, 'context', service.deploymentContext, 'stage', service.stage]
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
