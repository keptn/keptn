import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { DeleteData, DeleteResult } from '../../../_interfaces/delete';
import { EventService } from '../../../_services/event.service';
import { map, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-deletion-dialog',
  templateUrl: './ktb-deletion-dialog.component.html',
  styleUrls: [],
})
export class KtbDeletionDialogComponent implements OnInit, OnDestroy {
  private unsubscribe$ = new Subject();
  public isDeleteProjectInProgress$ = this.eventService.deletionProgressEvent.asObservable().pipe(map(evt => evt.isInProgress));
  public deletionError$ = this.eventService.deletionProgressEvent.asObservable().pipe(map(evt => evt.error));
  public deletionConfirmationControl = new FormControl('');
  public deletionConfirmationForm = new FormGroup({
    deletionConfirmation: this.deletionConfirmationControl,
  });

  constructor(@Inject(MAT_DIALOG_DATA) public data: DeleteData, public dialogRef: MatDialogRef<KtbDeletionDialogComponent>, private eventService: EventService) {
  }

  ngOnInit(): void {
    this.deletionConfirmationControl.setValidators([Validators.required, Validators.pattern(this.data.name)]);

    this.eventService.deletionProgressEvent.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(data => {
      if (data.result === DeleteResult.SUCCESS) {
        this.dialogRef.close();
      }
    });
  }

  ngOnDestroy() {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public deleteConfirm(): void {
    this.eventService.deletionTriggeredEvent.next({type: this.data.type, name: this.data.name});
  }
}
