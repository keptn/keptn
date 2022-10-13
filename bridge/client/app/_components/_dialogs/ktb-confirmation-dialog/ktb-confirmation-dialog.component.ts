import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { SequenceState } from '../../../_models/sequenceState';

export type SequenceConfirmDialogData = {
  sequence: SequenceState;
  confirmCallback: (params: SequenceConfirmDialogData) => void;
};

@Component({
  selector: 'ktb-deletion-dialog',
  templateUrl: './ktb-confirmation-dialog.component.html',
  styleUrls: [],
})
export class KtbConfirmationDialogComponent {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: SequenceConfirmDialogData,
    public dialogRef: MatDialogRef<KtbConfirmationDialogComponent>
  ) {}

  public confirm(): void {
    this.data.confirmCallback(this.data);
    this.dialogRef.close();
  }
}
