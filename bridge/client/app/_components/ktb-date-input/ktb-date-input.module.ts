import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtDatepickerModule } from '@dynatrace/barista-components/experimental/datepicker';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbDatetimePickerComponent, KtbDatetimePickerDirective } from './ktb-datetime-picker.component';
import { KtbTimeInputComponent } from './ktb-time-input.component';

@NgModule({
  declarations: [KtbDatetimePickerComponent, KtbTimeInputComponent, KtbDatetimePickerDirective],
  imports: [
    CommonModule,
    RouterModule,
    FlexModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtDatepickerModule,
    DtButtonModule,
    DtFormFieldModule,
    DtInputModule,
    ReactiveFormsModule,
  ],
  exports: [KtbDatetimePickerComponent, KtbTimeInputComponent, KtbDatetimePickerDirective],
})
export class KtbDateInputModule {}
