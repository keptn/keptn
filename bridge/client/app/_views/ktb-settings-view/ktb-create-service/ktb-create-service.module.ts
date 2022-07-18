import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbCreateServiceComponent } from './ktb-create-service.component';

@NgModule({
  declarations: [KtbCreateServiceComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FlexLayoutModule,
    DtButtonModule,
    DtInputModule,
    DtFormFieldModule,
    KtbLoadingModule,
  ],
  exports: [KtbCreateServiceComponent],
})
export class KtbCreateServiceModule {}
