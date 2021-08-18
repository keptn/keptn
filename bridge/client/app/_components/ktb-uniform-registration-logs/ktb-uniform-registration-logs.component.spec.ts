import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';
import { KtbUniformRegistrationLogsComponent } from './ktb-uniform-registration-logs.component';
import { UniformRegistrationLogsMock } from '../../_models/uniform-registrations-logs.mock';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbUniformRegistrationLogsComponent', () => {
  let component: KtbUniformRegistrationLogsComponent;
  let fixture: ComponentFixture<KtbUniformRegistrationLogsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbUniformRegistrationLogsComponent);
        component = fixture.componentInstance;
        component.logs = UniformRegistrationLogsMock;
        fixture.detectChanges();
      });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should not show logs', () => {
    component.logs = [];
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.uniform-registration-error-log')).toBeFalsy();
    expect(fixture.nativeElement.textContent.trim()).toEqual('No logs for this integration available');
  });

  it('should have 10 logs', () => {
    fixture.detectChanges();
    const rows = fixture.nativeElement.querySelectorAll('.uniform-registration-error-log>div');
    expect(rows.length).toEqual(10);
  });

  it('should show first 2 rows as unread', () => {
    component.lastSeen = new Date('2021-05-10T09:04:05.000Z');
    fixture.detectChanges();
    const firstUnreadRow = fixture.nativeElement.querySelectorAll('.uniform-registration-error-log>div:nth-of-type(1) .notification-indicator');
    const secondUnreadRow = fixture.nativeElement.querySelectorAll('.uniform-registration-error-log>div:nth-of-type(2) .notification-indicator');
    const allIndicators = fixture.nativeElement.querySelectorAll('.uniform-registration-error-log>div .notification-indicator');
    expect(firstUnreadRow).toBeTruthy();
    expect(secondUnreadRow).toBeTruthy();
    expect(allIndicators.length).toEqual(2);
  });

  it('should be unread without initial date', () => {
    component.lastSeen = undefined;
    expect(component.isUnread('2021-05-10T09:04:05.000Z')).toBe(true);
  });

  it('should be unread', () => {
    component.lastSeen = new Date('2021-05-10T09:04:05.000Z');
    expect(component.isUnread('2021-05-10T09:04:05.000Z')).toBe(false);
  });

  it('should be read', () => {
    component.lastSeen = new Date('2021-05-10T09:04:05.000Z');
    expect(component.isUnread('2021-05-10T10:04:05.000Z')).toBe(true);
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
