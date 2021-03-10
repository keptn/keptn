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
  public dataSource: DtTableDataSource<Service>;

  @Input()
  get services(): Service[] {
    return this._services;
  }
  set services(value: Service[]) {
    if (this._services !== value) {
      this._services = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dataService: DataService, public dateUtil: DateUtil) { }

  ngOnInit() {
    this.dataSource = new DtTableDataSource(this.services);
    this.dataService.roots
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(roots => {
        this._changeDetectorRef.markForCheck();
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
