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

  public addNotification(type: NotificationType, message?: string, componentInfo?: ComponentInfo, time?: number): void {
    const notification = new Notification(type, message, componentInfo, time);
    const notifications = this._notifications.getValue();
    // Check if the notification to add already exists
    const duplicateNotifications = notifications.filter(
      (n) => n.severity === notification.severity && n.message === notification.message
    );

    // Only show notification if it is not shown yet to prevent duplicates (issue #3896 - https://github.com/keptn/keptn/issues/3896)
    if (duplicateNotifications.length === 0) {
      this._notifications.next([...notifications, notification]);
    }
  }

  public removeNotification(notification: Notification): void {
    this._notifications.next(this._notifications.getValue().filter((n) => n !== notification));
  }
}
