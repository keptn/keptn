import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbConfirmationDialogComponent } from './ktb-confirmation-dialog.component';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbConfirmationDialogComponent],
  entryComponents: [KtbConfirmationDialogComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexLayoutModule,
  ],
  exports: [KtbConfirmationDialogComponent],
})
export class KtbConfirmationDialogModule {}
