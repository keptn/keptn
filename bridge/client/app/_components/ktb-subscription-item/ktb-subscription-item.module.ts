import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { KtbDeleteConfirmationModule } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.module';

@NgModule({
  declarations: [KtbSubscriptionItemComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtOverlayModule,
    KtbDeleteConfirmationModule,
  ],
  exports: [KtbSubscriptionItemComponent],
})
export class KtbSubscriptionItemModule {}
