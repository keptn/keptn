import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  Inject,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
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
export class KtbDeletionDialogComponent implements OnInit, OnDestroy, AfterViewInit {
  private unsubscribe$ = new Subject<void>();
  public isDeleteInProgress$ = this.eventService.deletionProgressEvent
    .asObservable()
    .pipe(map((evt) => evt.isInProgress));
  public deletionError$ = this.eventService.deletionProgressEvent.asObservable().pipe(map((evt) => evt.error));
  public deletionConfirmationControl = new FormControl('');
  public deletionConfirmationForm = new FormGroup({
    deletionConfirmation: this.deletionConfirmationControl,
  });
  @ViewChild('formInput') formInput: ElementRef | undefined;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: DeleteData,
    public dialogRef: MatDialogRef<KtbDeletionDialogComponent>,
    private eventService: EventService,
    private _changeDetectorRef: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    if (this.data.name) {
      this.deletionConfirmationControl.setValidators([Validators.required, Validators.pattern(this.data.name)]);
    }

    this.eventService.deletionProgressEvent.pipe(takeUntil(this.unsubscribe$)).subscribe((data) => {
      if (data.result === DeleteResult.SUCCESS) {
        this.dialogRef.close();
      }
    });
  }

  ngAfterViewInit(): void {
    this.formInput?.nativeElement.focus();
    this._changeDetectorRef.detectChanges();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public deleteConfirm(): void {
    if (this.deletionConfirmationForm.valid) this.eventService.deletionTriggeredEvent.next(this.data);
  }
}
