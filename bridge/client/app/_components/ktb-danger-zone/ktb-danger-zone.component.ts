import {Component, Input} from '@angular/core';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {KtbDeletionDialogComponent} from '../_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.component';
import {DeleteData} from '../../_interfaces/delete';

@Component({
  selector: 'ktb-danger-zone',
  templateUrl: './ktb-danger-zone.component.html',
  styleUrls: []
})
export class KtbDangerZoneComponent {
  @Input() data: DeleteData;

  public deletionDialogRef: MatDialogRef<any>;

  constructor(public dialog: MatDialog) { }

  public openDeletionDialog() {
    const data = {
      type: this.data.type,
      name: this.data.name
    };
    this.deletionDialogRef = this.dialog.open(KtbDeletionDialogComponent, {
      data
    });
  }
}
