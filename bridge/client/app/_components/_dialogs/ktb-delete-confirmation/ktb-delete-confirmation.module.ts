import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDeleteConfirmationComponent } from './ktb-delete-confirmation.component';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtButtonModule } from '@dynatrace/barista-components/button';

@NgModule({
  declarations: [KtbDeleteConfirmationComponent],
  imports: [CommonModule, DtConfirmationDialogModule, DtButtonModule],
  exports: [KtbDeleteConfirmationComponent],
})
export class KtbDeleteConfirmationModule {}
