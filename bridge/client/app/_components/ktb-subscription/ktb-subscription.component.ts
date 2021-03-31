import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input} from '@angular/core';
import {Subscription} from '../../_models/subscription';
import {KeptnService} from '../../_models/keptn-service';
import {DtTableDataSource} from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-subscription',
  templateUrl: './ktb-subscription.component.html',
  styleUrls: ['./ktb-subscription.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionComponent {
  public _keptnService: KeptnService;
  private defaultTask = 'all';
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

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

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
}
