import { Component } from '@angular/core';
import {Observable} from 'rxjs';
import {Notification} from '../../_models/notification';
import {NotificationsService} from '../../_services/notifications.service';

@Component({
  selector: 'ktb-notification-bar',
  templateUrl: './ktb-notification-bar.component.html',
  styleUrls: ['./ktb-notification-bar.component.scss']
})
export class KtbNotificationBarComponent {

  public notifications$: Observable<Notification[]>;

  constructor(private notificationsService: NotificationsService) {
    this.notifications$ = this.notificationsService.notifications;
  }

  public hideNotification(notification: Notification) {
    this.notificationsService.removeNotification(notification);
  }
}
