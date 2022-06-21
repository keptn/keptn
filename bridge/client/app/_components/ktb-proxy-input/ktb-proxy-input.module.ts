import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { KtbIntegerInputModule } from '../../_directives/ktb-integer-input/ktb-integer-input.module';
import { KtbProxyInputComponent } from './ktb-proxy-input.component';

@NgModule({
  declarations: [KtbProxyInputComponent],
  imports: [
    CommonModule,
    DtFormFieldModule,
    DtCheckboxModule,
    DtSelectModule,
    ReactiveFormsModule,
    FlexLayoutModule,
    DtInputModule,
    KtbIntegerInputModule,
  ],
  exports: [KtbProxyInputComponent],
})
export class KtbProxyInputModule {}
