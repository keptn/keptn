import { Component, HostListener, Input } from '@angular/core';
import { ComponentCanDeactivate } from '../../_guards/pending-changes.guard';
import { Observable } from 'rxjs';

type DialogState = null | 'unsaved';

@Component({
  selector: 'ktb-pending-changes-notification',
  templateUrl: './ktb-pending-changes-notification.component.html',
  styleUrls: ['./ktb-pending-changes-notification.component.scss'],
})
export class KtbPendingChangesNotificationComponent implements ComponentCanDeactivate {
  @Input() public message = 'You have pending changes. Make sure to save your data before you continue.';
  public dialogState: DialogState = null;

  @Input() onReset: () => void = () => {};
  @Input() onSave: () => void = () => {};
  @Input() isFormInvalid: () => boolean = () => true;
  @Input() canDeactivate: () => Observable<boolean> | boolean = () => true;

  reset(): void {
    this.onReset();
    this.hideNotification();
  }

  save(): void {
    this.onSave();
    this.hideNotification();
  }

  hideNotification(): void {
    this.dialogState = null;
  }

  // @HostListener allows us to also guard against browser refresh, close, etc.
  @HostListener('window:beforeunload', ['$event'])
  showNotification($event: any): void {
    if (!this.canDeactivate()) {
      this.dialogState = 'unsaved';
      $event.returnValue = this.message;
    }
  }
}
