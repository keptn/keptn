import { Component, Input, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { DataService } from '../../_services/data.service';

@Component({
  selector: 'ktb-payload-viewer',
  templateUrl: './ktb-payload-viewer.component.html',
  styleUrls: ['./ktb-payload-viewer.component.scss'],
})
export class KtbPayloadViewerComponent {

  @Input()
  public buttonTitle = 'Show event';

  @Input()
  public eventType;

  @Input

  @ViewChild('eventPayloadDialog')
  public eventPayloadDialog?: TemplateRef<ViewChild>;
  public eventPayloadDialogRef?: MatDialogRef<ViewChild>;

  constructor(private dataService: DataService, private dialog: MatDialog) {
  }

  showEventPayloadDialog(): void {
    if(this.eventPayloadDialog) {
      this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog, {
        data: JSON.stringify({}, null, 2)
      });
    }
  }

  closeDialog(): void {
    this.eventPayloadDialogRef?.close();
  }
}
