import {
  ChangeDetectionStrategy, ChangeDetectorRef,
  Component,
  Input,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {DtTableDataSource} from '@dynatrace/barista-components/table';
import {Subject} from 'rxjs';
import {Service} from '../../_models/service';
import {DateUtil} from '../../_utils/date.utils';
import {DataService} from '../../_services/data.service';
import {takeUntil} from 'rxjs/operators';
import {Root} from '../../_models/root';

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
  public dataSource: DtTableDataSource<Service> = new DtTableDataSource<Service>();

  @Input()
  get services(): Service[] {
    return this._services;
  }
  set services(value: Service[]) {
    if (this._services !== value) {
      this._services = value;
      this.updateDataSource();
    }
  }

  get pageSize(): number {
    return this._pageSize;
  }
  set pageSize(value: number) {
    if (this._pageSize !== value) {
      this._pageSize = value;
      this.updateDataSource();
    }
  }

  get DEFAULT_PAGE_SIZE(): number {
    return DEFAULT_PAGE_SIZE;
  }

  constructor(public dataService: DataService, public dateUtil: DateUtil, private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
    this.dataService.roots
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.updateDataSource();
      });
  }

  updateDataSource() {
    this.services.sort(this.compare());
    this.dataSource = new DtTableDataSource(this.services.slice(0, this.pageSize));
    this._changeDetectorRef.markForCheck();
  }

  private compare() {
    return (a: Service, b: Service) => {
      if (!a.getRecentRoot()) {
        return 1;
      }
      else if (!b.getRecentRoot()) {
        return -1;
      }
      else {
        return DateUtil.compareTraceTimesAsc(a.getRecentRoot().getLastTrace(), b.getRecentRoot().getLastTrace());
      }
    };
  }

  toggleAllServices() {
    if (this.services.length > this.pageSize) {
      this.pageSize = this.services.length;
    } else if (this.pageSize > DEFAULT_PAGE_SIZE) {
      this.pageSize = DEFAULT_PAGE_SIZE;
    }
  }

  getServiceLink(service: Service) {
    return ['service', service.serviceName, 'context', service.deploymentContext, 'stage', service.stage];
  }

  getSequenceLink(sequence: Root, service: Service) {
    return ['sequence', sequence.shkeptncontext, 'stage', service.stage];
  }

  public toService(row: Service): Service {
    return row;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
