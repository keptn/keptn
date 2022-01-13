import { Injectable } from '@angular/core';
import { CanDeactivate } from '@angular/router';
import { Observable } from 'rxjs';

export interface PendingChangesComponent {
  canDeactivate: ($event?: BeforeUnloadEvent) => Observable<boolean>;
}

@Injectable({
  providedIn: 'root',
})
export class PendingChangesGuard implements CanDeactivate<PendingChangesComponent> {
  canDeactivate(component: PendingChangesComponent): Observable<boolean> {
    return component.canDeactivate();
  }
}
