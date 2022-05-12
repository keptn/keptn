import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbNotificationBarComponent } from './ktb-notification-bar.component';
import { NotificationsService } from '../../_services/notifications.service';
import { Notification, NotificationType } from '../../_models/notification';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { firstValueFrom } from 'rxjs';

describe('KtbNotificationBarComponent', () => {
  let service: NotificationsService;
  let fixture: ComponentFixture<KtbNotificationBarComponent>;
  let component: KtbNotificationBarComponent;

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
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should add and remove notification', async () => {
    service.addNotification(NotificationType.INFO, 'Information');
    let notifications = await getNotifications();
    expect(notifications.length).toBe(1);
    component.hideNotification(notifications[0]);
    notifications = await getNotifications();
    expect(notifications.length).toBe(0);
  });

  function getNotifications(): Promise<Notification[]> {
    return firstValueFrom(component.notifications$);
  }
});
