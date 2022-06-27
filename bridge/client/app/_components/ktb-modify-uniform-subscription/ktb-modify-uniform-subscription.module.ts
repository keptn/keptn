import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { KtbPayloadViewerModule } from '../ktb-payload-viewer/ktb-payload-viewer.module';
import { KtbWebhookSettingsModule } from '../ktb-webhook-settings/ktb-webhook-settings.module';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbModifyUniformSubscriptionComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtCheckboxModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtFormFieldModule,
    DtFilterFieldModule,
    DtOverlayModule,
    DtSelectModule,
    KtbPayloadViewerModule,
    KtbWebhookSettingsModule,
    KtbLoadingModule,
    FlexLayoutModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  exports: [KtbModifyUniformSubscriptionComponent],
})
export class KtbModifyUniformSubscriptionModule {}
