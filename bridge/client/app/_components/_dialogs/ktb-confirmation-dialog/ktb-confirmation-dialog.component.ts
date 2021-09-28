import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'ktb-deletion-dialog',
  templateUrl: './ktb-confirmation-dialog.component.html',
  styleUrls: [],
})
export class KtbConfirmationDialogComponent {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<KtbConfirmationDialogComponent>
  ) {}

  public confirm() {
    this.data.confirmCallback(this.data);
    this.dialogRef.close();
  }
}
