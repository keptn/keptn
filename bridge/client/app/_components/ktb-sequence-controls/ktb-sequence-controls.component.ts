import { ChangeDetectorRef, Component, HostBinding, Input } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { Sequence } from '../../_models/sequence';
import { SequenceStateControl } from '../../../../shared/models/sequence';
import {
  KtbConfirmationDialogComponent,
  SequenceConfirmDialogData,
} from '../_dialogs/ktb-confirmation-dialog/ktb-confirmation-dialog.component';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'ktb-sequence-controls',
  templateUrl: './ktb-sequence-controls.component.html',
})
export class KtbSequenceControlsComponent {
  @HostBinding('class') cls = 'ktb-sequence-controls';
  private _sequence?: Sequence;
  private _smallButtons = false;
  public confirmationDialogRef?: MatDialogRef<KtbConfirmationDialogComponent>;

  @Input()
  get sequence(): Sequence | undefined {
    return this._sequence;
  }
  set sequence(sequence: Sequence | undefined) {
    if (this._sequence !== sequence) {
      this._sequence = sequence;
    }
  }

  @Input()
  get smallButtons(): boolean {
    return this._smallButtons;
  }
  set smallButtons(value: boolean) {
    if (this._smallButtons !== value) {
      this._smallButtons = value;
    }
  }

  constructor(
    private _changeDetectorRef: ChangeDetectorRef,
    private dataService: DataService,
    public dialog: MatDialog
  ) {}

  triggerResumeSequence(sequence: Sequence): void {
    this.dataService.sendSequenceControl(sequence, SequenceStateControl.RESUME);
  }

  triggerPauseSequence(sequence: Sequence): void {
    this.dataService.sendSequenceControl(sequence, SequenceStateControl.PAUSE);
  }

  triggerAbortSequence(sequence: Sequence): void {
    const data: SequenceConfirmDialogData = {
      sequence,
      confirmCallback: (params: SequenceConfirmDialogData): void => {
        this.abortSequence(params.sequence);
      },
    };
    this.confirmationDialogRef = this.dialog.open(KtbConfirmationDialogComponent, {
      data,
    });
  }

  abortSequence(sequence: Sequence): void {
    this.dataService.sendSequenceControl(sequence, SequenceStateControl.ABORT);
  }
}
