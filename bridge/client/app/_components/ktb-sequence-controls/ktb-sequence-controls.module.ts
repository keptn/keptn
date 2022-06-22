import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbSequenceControlsComponent } from './ktb-sequence-controls.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

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
  ],
  exports: [KtbSequenceControlsComponent],
})
export class KtbSequenceControlsModule {}
