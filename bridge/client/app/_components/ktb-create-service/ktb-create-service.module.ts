import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbCreateServiceComponent } from './ktb-create-service.component';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { ReactiveFormsModule } from '@angular/forms';
import { FlexModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtInputModule } from '@dynatrace/barista-components/input';

@NgModule({
  declarations: [KtbCreateServiceComponent],
  imports: [
    CommonModule,
    KtbLoadingModule,
    DtFormFieldModule,
    ReactiveFormsModule,
    FlexModule,
    DtButtonModule,
    DtInputModule,
  ],
  exports: [KtbCreateServiceComponent],
})
export class KtbCreateServiceModule {}
