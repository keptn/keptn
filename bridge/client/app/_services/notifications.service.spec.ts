import { TestBed } from '@angular/core/testing';
import { NotificationsService } from './notifications.service';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Notification, NotificationType } from '../_models/notification';
import { KtbProjectCreateMessageComponent } from '../_components/_status-messages/ktb-project-create-message/ktb-project-create-message.component';
import { firstValueFrom } from 'rxjs';

describe('NotificationsService', () => {
  let service: NotificationsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });

    service = TestBed.inject(NotificationsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should add notification', async () => {
    service.addNotification(NotificationType.ERROR, 'error1');
    service.addNotification(NotificationType.ERROR, 'error2');
    service.addNotification(NotificationType.ERROR, '', {
      component: KtbProjectCreateMessageComponent,
      data: {
        projectName: 'sockshop',
        routerLink: '/',
      },
    });
    const notifications = await getNotifications();
    expect(notifications.length).toBe(3);
  });

  it('should not add the same notification', async () => {
    service.addNotification(NotificationType.ERROR, 'error1');
    service.addNotification(NotificationType.ERROR, 'error1');
    const notifications = await getNotifications();
    expect(notifications.length).toBe(1);
  });

  it('should not add the same component notification', async () => {
    service.addNotification(NotificationType.ERROR, '', {
      component: KtbProjectCreateMessageComponent,
      data: {
        projectName: 'sockshop',
        routerLink: '/',
      },
    });
    service.addNotification(NotificationType.ERROR, '', {
      component: KtbProjectCreateMessageComponent,
      data: {
        projectName: 'sockshop',
        routerLink: '/',
      },
    });
    const notifications = await getNotifications();
    expect(notifications.length).toBe(1);
  });

  it('should remove notifications', async () => {
    service.addNotification(NotificationType.ERROR, 'error1');
    let notifications = await getNotifications();
    expect(notifications.length).toBe(1);
    service.removeNotification(notifications[0]);
    notifications = await getNotifications();
    expect(notifications.length).toBe(0);
  });

  function getNotifications(): Promise<Notification[]> {
    return firstValueFrom(service.notifications);
  }
});
