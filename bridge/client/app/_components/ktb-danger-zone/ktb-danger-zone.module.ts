import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDangerZoneComponent } from './ktb-danger-zone.component';
import { FlexModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { MatDialogModule } from '@angular/material/dialog';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

@NgModule({
  declarations: [KtbDangerZoneComponent],
  imports: [CommonModule, BrowserAnimationsModule, FlexModule, DtButtonModule, MatDialogModule],
  exports: [KtbDangerZoneComponent],
})
export class KtbDangerZoneModule {}
