import { Component, Input, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Trace } from '../../_models/trace';
import { DataService } from '../../_services/data.service';
import { catchError, finalize } from 'rxjs/operators';
import { of } from 'rxjs';

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
    if (!this.eventPayloadDialog) {
      return;
    }

    this.event = undefined;
    this.error = undefined;
    this.loading = true;
    this.dataService
      .getEvent(this.type, this.project, this.stage, this.service)
      .pipe(
        catchError((err) => {
          this.error = err;
          return of(undefined);
        }),
        finalize(() => (this.loading = false))
      )
      .subscribe((event) => (this.event = event));
    this.eventPayloadDialogRef = this.dialog.open(this.eventPayloadDialog);
  }

  closeDialog(): void {
    this.eventPayloadDialogRef?.close();
  }
}
