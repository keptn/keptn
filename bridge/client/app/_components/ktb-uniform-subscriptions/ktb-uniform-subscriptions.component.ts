import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { IUniformRegistration } from '../../../../shared/interfaces/uniform-registration';
import { IUniformSubscription } from '../../../../shared/interfaces/uniform-subscription';
import { canEditSubscriptions, getSubscriptions, isWebhookService } from '../../_models/uniform-registration';

@Component({
  selector: 'ktb-uniform-subscriptions[uniformRegistration]',
  templateUrl: './ktb-uniform-subscriptions.component.html',
  styleUrls: ['./ktb-uniform-subscriptions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbUniformSubscriptionsComponent {
  private _uniformRegistration?: IUniformRegistration;
  public projectName$: Observable<string | null>;
  public isWebhookService = isWebhookService;
  public canEditSubscriptions = canEditSubscriptions;
  public getSubscriptions = getSubscriptions;

  @Input()
  set uniformRegistration(registration: IUniformRegistration | undefined) {
    if (this._uniformRegistration !== registration) {
      this._uniformRegistration = registration;
      this._changeDetectorRef.markForCheck();
    }
  }
  get uniformRegistration(): IUniformRegistration | undefined {
    return this._uniformRegistration;
  }
  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: ActivatedRoute) {
    this.projectName$ = this.router.paramMap.pipe(map((params) => params.get('projectName')));
  }

  public deleteSubscription(subscription: IUniformSubscription): void {
    if (this.uniformRegistration) {
      const index = this.uniformRegistration.subscriptions.indexOf(subscription);
      if (index >= 0) {
        this.uniformRegistration.subscriptions.splice(index, 1);
        this._changeDetectorRef.markForCheck();
      }
    }
  }
}
