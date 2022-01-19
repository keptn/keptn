import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ComponentInfo, Notification, NotificationType } from '../_models/notification';

@Injectable({
  providedIn: 'root',
})
export class NotificationsService {
  private _notifications = new BehaviorSubject<Notification[]>([]);

  get notifications(): Observable<Notification[]> {
    return this._notifications.asObservable();
  }

  public addNotification(
    type: NotificationType,
    message?: string,
    componentInfo?: ComponentInfo,
    time?: number,
    showOnTop = false
  ): void {
    const notification = new Notification(type, message, componentInfo, time, showOnTop);
    const notifications = this._notifications.getValue();
    // Check if the notification to add already exists
    const duplicateNotifications = notifications.filter(
      (n) => n.type === notification.type && n.message === notification.message
    );

    // Only show notification if it is not shown yet to prevent duplicates (issue #3896 - https://github.com/keptn/keptn/issues/3896)
    if (duplicateNotifications.length === 0) {
      const notificationsOnTop: Notification[] = [];
      const notificationsOnBottom: Notification[] = [];
      for (const n of notifications) {
        if (n.showOnTop) {
          notificationsOnTop.push(n);
        } else {
          notificationsOnBottom.push(n);
        }
      }
      // new notifications should be shown on top, whereas notifications like "update keptn" should always stay on top
      this._notifications.next([...notificationsOnTop, notification, ...notificationsOnBottom]);
    }
  }

  public removeNotification(notification: Notification): void {
    this._notifications.next(this._notifications.getValue().filter((n) => n !== notification));
  }
}
