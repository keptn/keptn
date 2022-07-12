import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDangerZoneComponent } from './ktb-danger-zone.component';
import { FlexModule } from '@angular/flex-layout';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { MatDialogModule } from '@angular/material/dialog';
import { KtbDeletionDialogModule } from '../_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.module';

@NgModule({
  declarations: [KtbDangerZoneComponent],
  imports: [CommonModule, FlexModule, DtButtonModule, MatDialogModule, KtbDeletionDialogModule],
  exports: [KtbDangerZoneComponent],
})
export class KtbDangerZoneModule {}
