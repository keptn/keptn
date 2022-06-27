import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbSequenceControlsComponent } from './ktb-sequence-controls.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { KtbConfirmationDialogModule } from '../_dialogs/ktb-confirmation-dialog/ktb-confirmation-dialog.module';

@NgModule({
  declarations: [KtbSequenceControlsComponent],
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    MatDialogModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtButtonModule,
    KtbConfirmationDialogModule,
  ],
  exports: [KtbSequenceControlsComponent],
})
export class KtbSequenceControlsModule {}
