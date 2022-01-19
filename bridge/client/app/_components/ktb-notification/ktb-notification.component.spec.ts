import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbNotificationComponent } from './ktb-notification.component';

describe('KtbNotificationComponent', () => {
  let component: KtbNotificationComponent;
  let fixture: ComponentFixture<KtbNotificationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbNotificationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
