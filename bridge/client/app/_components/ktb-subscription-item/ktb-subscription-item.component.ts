import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, Output } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { DeleteDialogState } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';
import { IUniformSubscription } from '../../../../shared/interfaces/uniform-subscription';
import { formatFilter, getFormattedEvent, isGlobal } from '../../_models/uniform-subscription';

@Component({
  selector: 'ktb-subscription-item[subscription][integrationId][isWebhookService][projectName]',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSubscriptionItemComponent {
  private _subscription?: IUniformSubscription;
  private currentSubscription?: IUniformSubscription;
  public deleteState: DeleteDialogState = null;
  public getFormattedEvent = getFormattedEvent;
  public formatFilter = formatFilter;
  public isGlobal = isGlobal;

  @Output() subscriptionDeleted: EventEmitter<IUniformSubscription> = new EventEmitter<IUniformSubscription>();
  @Input() name?: string;
  @Input() integrationId?: string;
  @Input() isWebhookService = false;
  @Input() projectName = '';

  @Input()
  get subscription(): IUniformSubscription | undefined {
    return this._subscription;
  }

  set subscription(subscription: IUniformSubscription | undefined) {
    if (this._subscription !== subscription) {
      this._subscription = subscription;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(
    private changeDetectorRef: ChangeDetectorRef,
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  public editSubscription(subscription: IUniformSubscription): void {
    this.router.navigate([
      '/',
      'project',
      this.projectName,
      'settings',
      'uniform',
      'integrations',
      this.integrationId,
      'subscriptions',
      subscription.id,
      'edit',
    ]);
  }

  public triggerDeleteSubscription(subscription: IUniformSubscription): void {
    this.currentSubscription = subscription;
    this.deleteState = 'confirm';
  }

  public deleteSubscription(): void {
    if (this.integrationId && this.subscription?.id) {
      this.dataService
        .deleteSubscription(this.integrationId, this.subscription.id, this.isWebhookService)
        .subscribe(() => {
          this.deleteState = 'success';
          this.subscriptionDeleted.emit(this.subscription);
        });
    }
  }
}
