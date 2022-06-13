import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDeletionDialogComponent } from './ktb-deletion-dialog.component';
import { KtbLoadingModule } from '../../ktb-loading/ktb-loading.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { ReactiveFormsModule } from '@angular/forms';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbDeletionDialogComponent],
  imports: [
    CommonModule,
    KtbLoadingModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtFormFieldModule,
    DtInputModule,
    FlexLayoutModule,
    ReactiveFormsModule,
  ],
  exports: [KtbDeletionDialogComponent],
})
export class KtbDeletionDialogModule {}
