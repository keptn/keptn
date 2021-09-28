import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input } from '@angular/core';
import { UniformRegistration } from '../../_models/uniform-registration';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { UniformSubscription } from '../../_models/uniform-subscription';

@Component({
  selector: 'ktb-uniform-subscriptions[uniformRegistration]',
  templateUrl: './ktb-uniform-subscriptions.component.html',
  styleUrls: ['./ktb-uniform-subscriptions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbUniformSubscriptionsComponent {
  private _uniformRegistration?: UniformRegistration;
  public projectName$: Observable<string | null>;

  @Input()
  set uniformRegistration(registration: UniformRegistration | undefined) {
    if (this._uniformRegistration !== registration) {
      this._uniformRegistration = registration;
      this._changeDetectorRef.markForCheck();
    }
  }
  get uniformRegistration(): UniformRegistration | undefined {
    return this._uniformRegistration;
  }
  constructor(private _changeDetectorRef: ChangeDetectorRef, private router: ActivatedRoute) {
    this.projectName$ = this.router.paramMap.pipe(map((params) => params.get('projectName')));
  }

  public deleteSubscription(subscription: UniformSubscription): void {
    if (this.uniformRegistration) {
      const index = this.uniformRegistration.subscriptions.indexOf(subscription);
      if (index >= 0) {
        this.uniformRegistration.subscriptions.splice(index, 1);
        this._changeDetectorRef.markForCheck();
      }
    }
  }
}
