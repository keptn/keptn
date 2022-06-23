import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbTriggerSequenceComponent } from './ktb-trigger-sequence.component';
import { FlexModule } from '@angular/flex-layout';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbDateInputModule } from '../ktb-date-input/ktb-date-input.module';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { FormsModule } from '@angular/forms';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbTriggerSequenceComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtFormFieldModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInputModule,
    DtOverlayModule,
    DtRadioModule,
    DtSelectModule,
    FlexModule,
    FormsModule,
    KtbDateInputModule,
    KtbLoadingModule,
    RouterModule,
  ],
  exports: [KtbTriggerSequenceComponent],
})
export class KtbTriggerSequenceModule {}
