import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, Output } from '@angular/core';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { DeleteDialogState } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';

@Component({
  selector: 'ktb-subscription-item[subscription][integrationId][isWebhookService][projectName]',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSubscriptionItemComponent {
  private _subscription?: UniformSubscription;
  private currentSubscription?: UniformSubscription;
  public deleteState: DeleteDialogState = null;

  @Output() subscriptionDeleted: EventEmitter<UniformSubscription> = new EventEmitter<UniformSubscription>();
  @Input() name?: string;
  @Input() integrationId?: string;
  @Input() isWebhookService = false;
  @Input() projectName = '';

  @Input()
  get subscription(): UniformSubscription | undefined {
    return this._subscription;
  }

  set subscription(subscription: UniformSubscription | undefined) {
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

  public editSubscription(subscription: UniformSubscription): void {
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

  public triggerDeleteSubscription(subscription: UniformSubscription): void {
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
