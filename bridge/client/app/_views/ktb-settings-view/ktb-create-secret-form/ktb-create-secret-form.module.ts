import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';

@NgModule({
  declarations: [KtbCreateSecretFormComponent],
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    FlexLayoutModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtFormFieldModule,
    DtInputModule,
    DtButtonModule,
    DtSelectModule,
    KtbLoadingModule,
    KtbPipeModule,
  ],
  exports: [KtbCreateSecretFormComponent],
})
export class KtbCreateSecretFormModule {}
