import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbPendingChangesNotificationComponent } from './ktb-pending-changes-notification.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';

describe('KtbPendingChangesNotificationComponent', () => {
  let component: KtbPendingChangesNotificationComponent;
  let fixture: ComponentFixture<KtbPendingChangesNotificationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbPendingChangesNotificationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
