import { Injectable } from '@angular/core';
import { CanDeactivate } from '@angular/router';
import { Observable, of } from 'rxjs';

export interface PendingChangesComponent {
  canDeactivate: ($event?: BeforeUnloadEvent) => Observable<boolean>;
}

@Injectable({
  providedIn: 'root',
})
export class PendingChangesGuard implements CanDeactivate<PendingChangesComponent> {
  canDeactivate(component: PendingChangesComponent | null): Observable<boolean> {
    // null because of lazy loading
    return component?.canDeactivate() ?? of(true);
  }
}
