import { Component, Input, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Trace } from '../../_models/trace';
import { ApiService } from '../../_services/api.service';

@Component({
  selector: 'ktb-payload-viewer',
  templateUrl: './ktb-payload-viewer.component.html',
  styleUrls: ['./ktb-payload-viewer.component.scss'],
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

  constructor(private apiService: ApiService, private dialog: MatDialog) {}

  showEventPayloadDialog(): void {
    if (this.eventPayloadDialog) {
      this.event = undefined;
      this.loading = true;
      this.apiService.getEvent(this.type, this.project, this.stage, this.service).subscribe(
        (eventResult) => {
          this.event = eventResult.body?.events[0];
          this.loading = false;
        },
        (err) => {
          this.event = undefined;
          this.loading = false;
        }
      );
      this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog);
    }
  }

  closeDialog(): void {
    this.eventPayloadDialogRef?.close();
  }
}
