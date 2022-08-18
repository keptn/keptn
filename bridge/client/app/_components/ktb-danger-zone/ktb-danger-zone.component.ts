import { Component, Input } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { KtbDeletionDialogComponent } from '../_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.component';
import { DeleteData } from '../../_interfaces/delete';

@Component({
  selector: 'ktb-danger-zone[data]',
  templateUrl: './ktb-danger-zone.component.html',
  styleUrls: ['./ktb-danger-zone.component.scss'],
})
export class KtbDangerZoneComponent {
  @Input() data!: DeleteData;

  constructor(public dialog: MatDialog) {}

  public openDeletionDialog(): void {
    this.dialog.open(KtbDeletionDialogComponent, {
      data: this.data,
      autoFocus: false, // else the close icon will be incorrectly selected
    });
  }
}
