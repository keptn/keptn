import { Injectable } from '@angular/core';
import { CanDeactivate } from '@angular/router';
import { Observable } from 'rxjs';

export interface ComponentCanDeactivate {
  canDeactivate: () => boolean | Observable<boolean>;
  showNotification: ($event: any) => void;
}

@Injectable({
  providedIn: 'root',
})
export class PendingChangesGuard implements CanDeactivate<ComponentCanDeactivate> {
  canDeactivate(component: ComponentCanDeactivate): boolean | Observable<boolean> {
    const hasPendingChanges = !component.canDeactivate();
    if (hasPendingChanges) {
      // NOTE: this warning message will only be shown when navigating elsewhere within your angular app;
      // when navigating away from your angular app, the browser will show a generic warning message
      // see http://stackoverflow.com/a/42207299/7307355
      component.showNotification(null);
    }
    return component.canDeactivate();
  }
}
