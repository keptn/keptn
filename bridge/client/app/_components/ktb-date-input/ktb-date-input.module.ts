import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDatetimePickerComponent, KtbDatetimePickerDirective } from './ktb-datetime-picker.component';
import { KtbTimeInputComponent } from './ktb-time-input.component';
import { FlexModule } from '@angular/flex-layout';
import { DtDatepickerModule } from '@dynatrace/barista-components/experimental/datepicker';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { ReactiveFormsModule } from '@angular/forms';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [KtbDatetimePickerComponent, KtbTimeInputComponent, KtbDatetimePickerDirective],
  imports: [
    CommonModule,
    RouterModule,
    HttpClientModule,
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
