import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbNotificationBarComponent } from './ktb-notification-bar.component';
import { By } from '@angular/platform-browser';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';

describe('KtbNotificationBarComponent', () => {
  let service: NotificationsService;
  let fixture: ComponentFixture<KtbNotificationBarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
        DtIconModule,
        DtIconModule.forRoot({
          svgIconLocation: `/assets/icons/{{name}}.svg`,
        }),
      ],
      providers: [NotificationsService],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbNotificationBarComponent);
    service = TestBed.inject(NotificationsService);
    fixture.detectChanges();
  });

  it('should add and remove notifications', () => {
    let notifications = fixture.debugElement.queryAll(By.css('ktb-notification'));
    expect(notifications.length).toBe(0);

    service.addNotification(NotificationType.INFO, 'Information');
    service.addNotification(NotificationType.SUCCESS, 'Success');
    service.addNotification(NotificationType.WARNING, 'Warning');
    service.addNotification(NotificationType.ERROR, 'Error');
    fixture.detectChanges();

    notifications = fixture.debugElement.queryAll(By.css('ktb-notification dt-alert'));
    expect(notifications.length).toBe(4);

    expect(notifications[0].nativeElement.classList).toContain('dt-alert-info');
    expect(notifications[1].nativeElement.classList).toContain('dt-alert-success');
    expect(notifications[2].nativeElement.classList).toContain('dt-alert-warning');
    expect(notifications[3].nativeElement.classList).toContain('dt-alert');
    expect(notifications[3].nativeElement.querySelector('.dt-alert-icon-container dt-icon').classList).toContain(
      'dt-alert-icon'
    );
  });
});
