import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { Notification, TemplateRenderedNotifications } from '../../_models/notification';
import { NotificationsService } from '../../_services/notifications.service';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-notification-bar',
  templateUrl: './ktb-notification-bar.component.html',
  styleUrls: ['./ktb-notification-bar.component.scss'],
})
export class KtbNotificationBarComponent {
  public notificationTypes = TemplateRenderedNotifications;
  public notifications$: Observable<Notification[]>;

  constructor(private notificationsService: NotificationsService, public location: Location) {
    this.notifications$ = this.notificationsService.notifications;
  }

  public hideNotification(notification: Notification): void {
    this.notificationsService.removeNotification(notification);
  }
}
