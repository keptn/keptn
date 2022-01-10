import { Injectable } from '@angular/core';
import { CanDeactivate } from '@angular/router';
import { Observable } from 'rxjs';

export interface PendingChangesComponent {
  canDeactivate: () => boolean | Observable<boolean>;
  showNotification: ($event?: BeforeUnloadEvent) => void;
}

@Injectable({
  providedIn: 'root',
})
export class PendingChangesGuard implements CanDeactivate<PendingChangesComponent> {
  canDeactivate(component: PendingChangesComponent): boolean | Observable<boolean> {
    const hasPendingChanges = component.canDeactivate() !== true;
    if (hasPendingChanges) {
      // NOTE: this warning message will only be shown when navigating elsewhere within your angular app;
      // when navigating away from your angular app, the browser will show a generic warning message
      // see http://stackoverflow.com/a/42207299/7307355
      component.showNotification();
    }
    return component.canDeactivate();
  }
}
