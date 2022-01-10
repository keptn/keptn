import { TestBed } from '@angular/core/testing';

import { PendingChangesComponent, PendingChangesGuard } from './pending-changes.guard';
import { Observable, Subject } from 'rxjs';

type NotificationState = null | 'unsaved';
class MockComponent implements PendingChangesComponent {
  private formUnsavedSubject = new Subject<boolean>();
  returnValue: boolean | Observable<boolean> = true;
  notificationState: NotificationState = null;

  formTouched(): void {
    this.returnValue = this.formUnsavedSubject.asObservable();
  }

  saveForm(): void {
    this.formUnsavedSubject.next(true);
  }

  rejectChanges(): void {
    this.formUnsavedSubject.next(false);
  }

  canDeactivate(): boolean | Observable<boolean> {
    return this.returnValue;
  }

  showNotification($event?: BeforeUnloadEvent): void {
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

  it('will route if guarded and user accepted the dialog', () => {
    // Mock the behavior of the GuardedComponent:
    mockComponent.formTouched();
    const canDeactivate$ = <Observable<boolean>>guard.canDeactivate(mockComponent);
    canDeactivate$.subscribe((deactivate) => {
      // This is the real test!
      expect(deactivate).toBeTruthy();
    });
    // emulate the accept()
    mockComponent.saveForm();
  });

  it('will not route if guarded and user rejected the dialog', () => {
    // Mock the behavior of the GuardedComponent:
    mockComponent.formTouched();
    const canDeactivate$ = <Observable<boolean>>guard.canDeactivate(mockComponent);
    canDeactivate$.subscribe((deactivate) => {
      // This is the real test!
      expect(deactivate).toBeFalsy();
    });
    // emulate the reject()
    mockComponent.rejectChanges();
  });
});
