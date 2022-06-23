import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { ReactiveFormsModule } from '@angular/forms';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbVariableSelectorModule } from '../ktb-variable-selector/ktb-variable-selector.module';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbWebhookSettingsComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtSelectModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtRadioModule,
    DtOverlayModule,
    FlexLayoutModule,
    ReactiveFormsModule,
    KtbVariableSelectorModule,
    DtInputModule,
  ],
  exports: [KtbWebhookSettingsComponent],
})
export class KtbWebhookSettingsModule {}
