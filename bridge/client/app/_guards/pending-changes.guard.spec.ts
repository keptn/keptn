import { TestBed, waitForAsync } from '@angular/core/testing';

import { PendingChangesComponent, PendingChangesGuard } from './pending-changes.guard';
import { Observable, of, Subject } from 'rxjs';

type NotificationState = null | 'unsaved';
class MockComponent implements PendingChangesComponent {
  private pendingChangesSubject = new Subject<boolean>();
  formTouched = false;
  notificationState: NotificationState = null;

  touchForm(): void {
    this.formTouched = true;
  }

  saveForm(): void {
    this.pendingChangesSubject.next(true);
    this.hideNotification();
  }

  rejectChanges(): void {
    this.pendingChangesSubject.next(false);
    this.hideNotification();
  }

  canDeactivate(): Observable<boolean> {
    if (this.formTouched) {
      this.showNotification();
      return this.pendingChangesSubject.asObservable();
    } else {
      return of(true);
    }
  }

  showNotification(): void {
    this.notificationState = 'unsaved';
  }

  hideNotification(): void {
    this.notificationState = null;
  }
}

describe('PendingChangesGuard', () => {
  let mockComponent: MockComponent;
  let guard: PendingChangesGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [MockComponent],
    });
    guard = TestBed.inject(PendingChangesGuard);
    mockComponent = TestBed.get(MockComponent);
  });

  it('expect guard to instantiate', () => {
    expect(guard).toBeTruthy();
  });

  it('can route if unguarded', () => {
    expect(guard.canDeactivate(mockComponent)).toBeTruthy();
  });

  it('will route if guarded and user accepted the dialog', waitForAsync(() => {
    // Mock the behavior of the GuardedComponent:
    mockComponent.touchForm();
    const canDeactivate$ = <Observable<boolean>>guard.canDeactivate(mockComponent);
    canDeactivate$.subscribe((deactivate) => {
      // This is the real test!
      expect(deactivate).toBeTruthy();
    });
    // emulate the accept()
    mockComponent.saveForm();
  }));

  it('will not route if guarded and user rejected the dialog', waitForAsync(() => {
    // Mock the behavior of the GuardedComponent:
    mockComponent.touchForm();
    const canDeactivate$ = <Observable<boolean>>guard.canDeactivate(mockComponent);
    canDeactivate$.subscribe((deactivate) => {
      // This is the real test!
      expect(deactivate).toBeFalsy();
    });
    // emulate the reject()
    mockComponent.rejectChanges();
  }));
});
