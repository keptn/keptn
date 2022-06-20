import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbUniformRegistrationLogsComponent } from './ktb-uniform-registration-logs.component';
import { UniformRegistrationLogsMock } from '../../_services/_mockData/uniform-registrations-logs.mock';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbUniformRegistrationLogsModule } from './ktb-uniform-registration-logs.module';

describe('KtbUniformRegistrationLogsComponent', () => {
  let component: KtbUniformRegistrationLogsComponent;
  let fixture: ComponentFixture<KtbUniformRegistrationLogsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbUniformRegistrationLogsModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbUniformRegistrationLogsComponent);
    component = fixture.componentInstance;
    component.logs = UniformRegistrationLogsMock;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
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
});
