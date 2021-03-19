import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from '../../_models/subscription';
import {map, takeUntil} from 'rxjs/operators';
import {forkJoin, Subject} from 'rxjs';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {DtFilterFieldDefaultDataSource} from '@dynatrace/barista-components/filter-field';
import {Project} from '../../_models/project';
import {KeptnService} from '../../_models/keptn-service';

@Component({
  selector: 'ktb-subscription',
  templateUrl: './ktb-subscription.component.html',
  styleUrls: ['./ktb-subscription.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionComponent implements OnInit, OnDestroy {
  public _keptnService: KeptnService;
  private readonly unsubscribe$ = new Subject<void>();
  public newSubscription: Subscription = new Subscription();
  private defaultTask: string;

  @Input()
  get keptnService(): KeptnService {
    return this._keptnService;
  }
  set keptnService(keptnService: KeptnService) {
    if (this._keptnService !== keptnService) {
      this._keptnService = keptnService;
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
              this.newSubscription.event = this.defaultTask;
              this._changeDetectorRef.markForCheck();
            }
        });
      });
  }

  public addSubscription() {
    this.keptnService.addSubscription(this.newSubscription);
    this.newSubscription = new Subscription();
    this.newSubscription.event = this.defaultTask;
    this._changeDetectorRef.markForCheck();
  }

  public updateSubscriptions() {
    // generate YAML file
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
