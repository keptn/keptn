import { Component, EventEmitter, Input, Output } from '@angular/core';

export type DeleteDialogState = 'confirm' | 'deleting' | 'success' | null;

@Component({
  selector: 'ktb-delete-confirmation[type][name]',
  templateUrl: './ktb-delete-confirmation.component.html',
  styleUrls: ['./ktb-delete-confirmation.component.scss']
})
export class KtbDeleteConfirmationComponent {
  private closeConfirmationDialogTimeout?: ReturnType<typeof setTimeout>;
  private _dialogState: DeleteDialogState = null;

  @Input()
  set dialogState(dialogState: DeleteDialogState) {
    if (this._dialogState !== dialogState) {
      if (this.closeConfirmationDialogTimeout) {
        clearTimeout(this.closeConfirmationDialogTimeout);
      }
      if (dialogState === 'success') {
        this.closeConfirmationDialogTimeout = setTimeout(() => {
          this.closeDialog();
        }, 2000);
      }
      this._dialogState = dialogState;
    }
  }

  get dialogState(): DeleteDialogState {
    return this._dialogState;
  }

  @Input() type?: string;
  @Input() name?: string;
  @Input() deleteMessage?: string;
  @Output() confirmClicked: EventEmitter<void> = new EventEmitter<void>();

  public closeDialog() {
    this.dialogState = null;
  }

  public deleteAction() {
    this.dialogState = 'deleting';
    this.confirmClicked.emit();
  }
}
