import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbUniformSubscriptionsComponent } from './ktb-uniform-subscriptions.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { RouterModule } from '@angular/router';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbSubscriptionItemModule } from '../ktb-subscription-item/ktb-subscription-item.module';

@NgModule({
  declarations: [KtbUniformSubscriptionsComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    RouterModule,
    DtOverlayModule,
    KtbSubscriptionItemModule,
  ],
  exports: [KtbUniformSubscriptionsComponent],
})
export class KtbUniformSubscriptionsModule {}
