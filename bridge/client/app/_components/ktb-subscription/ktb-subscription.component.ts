import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from '../../_models/subscription';
import {map, takeUntil} from 'rxjs/operators';
import {forkJoin, Subject} from 'rxjs';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {DtFilterFieldDefaultDataSource} from '@dynatrace/barista-components/filter-field';
import {Project} from '../../_models/project';
import {KeptnService} from '../../_models/keptn-service';
import {DtTableDataSource} from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-subscription',
  templateUrl: './ktb-subscription.component.html',
  styleUrls: ['./ktb-subscription.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionComponent implements OnInit, OnDestroy {
  public _keptnService: KeptnService;
  private readonly unsubscribe$ = new Subject<void>();
  private defaultTask: string;
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();

  @Input()
  get keptnService(): KeptnService {
    return this._keptnService;
  }
  set keptnService(keptnService: KeptnService) {
    if (this._keptnService !== keptnService) {
      this._keptnService = keptnService;
      this.tableEntries.data = keptnService.subscriptions;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.dataService.taskNames
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(tasks => {
            if (tasks.length !== 0) {
              this.defaultTask = tasks[0] + '.triggered';
              this._changeDetectorRef.markForCheck();
            }
        });
      });
  }

  public addSubscription() {
    const newSubscription = new Subscription();
    newSubscription.event = this.defaultTask;
    newSubscription.expanded = true;
    this.keptnService.addSubscription(newSubscription);
    this.updateDataSource();
  }

  public deleteSubscription(rowIndex: number) {
    this.keptnService.deleteSubscription(rowIndex);
    this.updateDataSource();
  }

  private updateDataSource() {
    this.tableEntries.data = this.keptnService.subscriptions;
    this._changeDetectorRef.markForCheck();
  }

  public updateSubscriptions() {
    // generate YAML file
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
