import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import { FlexModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { ReactiveFormsModule } from '@angular/forms';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbCreateSecretFormComponent],
  imports: [
    CommonModule,
    FlexModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    ReactiveFormsModule,
    DtFormFieldModule,
    KtbPipeModule,
    DtInputModule,
    DtButtonModule,
    DtSelectModule,
    KtbLoadingModule,
    RouterModule,
  ],
  exports: [KtbCreateSecretFormComponent],
})
export class KtbCreateSecretFormModule {}
