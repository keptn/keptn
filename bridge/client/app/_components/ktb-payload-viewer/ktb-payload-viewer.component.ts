import { Component, Input, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Trace } from '../../_models/trace';
import { DataService } from '../../_services/data.service';

@Component({
  selector: 'ktb-payload-viewer',
  templateUrl: './ktb-payload-viewer.component.html',
  styleUrls: [],
})
export class KtbPayloadViewerComponent {
  @Input()
  public buttonTitle = 'Show event';

  @Input()
  public type: string | undefined;

  @Input()
  public stage: string | undefined;

  @Input()
  public service: string | undefined;

  @Input()
  public project: string | undefined;

  @ViewChild('eventPayloadDialog')
  public eventPayloadDialog?: TemplateRef<ViewChild>;
  public eventPayloadDialogRef?: MatDialogRef<ViewChild>;

  public event: Trace | undefined;

  public loading = false;
  public error: Error | undefined;

  constructor(private dataService: DataService, private dialog: MatDialog) {}

  showEventPayloadDialog(): void {
    if (this.eventPayloadDialog) {
      this.event = undefined;
      this.loading = true;
      this.dataService.getEvent(this.type, this.project, this.stage, this.service).subscribe(
        (event) => {
          this.event = event;
          this.loading = false;
          this.error = undefined;
        },
        (err) => {
          this.event = undefined;
          this.loading = false;
          this.error = err;
        }
      );
      this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog);
    }
  }

  closeDialog(): void {
    this.eventPayloadDialogRef?.close();
  }
}
