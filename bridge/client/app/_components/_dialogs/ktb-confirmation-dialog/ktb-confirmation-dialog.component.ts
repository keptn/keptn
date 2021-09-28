import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'ktb-deletion-dialog',
  templateUrl: './ktb-confirmation-dialog.component.html',
  styleUrls: [],
})
export class KtbConfirmationDialogComponent {
  constructor(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<KtbConfirmationDialogComponent>
  ) {}

  public confirm(): void {
    this.data.confirmCallback(this.data);
    this.dialogRef.close();
  }
}
