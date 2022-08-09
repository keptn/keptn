import { NgModule } from '@angular/core';
import { KtbIntegrationViewRoutingModule } from './ktb-integration-view-routing.module';
import { ReactiveFormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbWebhookSettingsModule } from '../../../_components/ktb-webhook-settings/ktb-webhook-settings.module';
import { KtbPayloadViewerModule } from '../../../_components/ktb-payload-viewer/ktb-payload-viewer.module';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { CommonModule } from '@angular/common';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { KtbExpandableTileModule } from '../../../_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbUniformRegistrationLogsModule } from '../../../_components/ktb-uniform-registration-logs/ktb-uniform-registration-logs.module';
import { KtbUniformSubscriptionsModule } from '../../../_components/ktb-uniform-subscriptions/ktb-uniform-subscriptions.module';
import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription/ktb-modify-uniform-subscription.component';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';

@NgModule({
  declarations: [KtbIntegrationViewComponent, KtbModifyUniformSubscriptionComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    ReactiveFormsModule,
    DtButtonModule,
    DtCheckboxModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtFormFieldModule,
    DtFilterFieldModule,
    DtOverlayModule,
    DtSelectModule,
    DtTableModule,
    KtbExpandableTileModule,
    KtbPipeModule,
    KtbUniformRegistrationLogsModule,
    KtbUniformSubscriptionsModule,
    KtbPayloadViewerModule,
    KtbWebhookSettingsModule,
    KtbLoadingModule,
    KtbIntegrationViewRoutingModule,
    DtConfirmationDialogModule,
  ],
})
export class KtbIntegrationViewModule {}
