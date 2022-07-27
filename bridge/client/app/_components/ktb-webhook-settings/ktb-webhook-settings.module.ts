import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { KtbVariableSelectorModule } from '../ktb-variable-selector/ktb-variable-selector.module';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';

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
