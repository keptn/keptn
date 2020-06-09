import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbNotificationBarComponent } from './ktb-notification-bar.component';

describe('KtbNotificationBarComponent', () => {
  let component: KtbNotificationBarComponent;
  let fixture: ComponentFixture<KtbNotificationBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbNotificationBarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbNotificationBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
