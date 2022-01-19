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

    // Check if the notification to add already exists
    const duplicateNotifications = this._notifications
      .getValue()
      .filter((n) => n.type === notification.type && n.message === notification.message);

    // Only show notification if it is not shown yet to prevent duplicates (issue #3896 - https://github.com/keptn/keptn/issues/3896)
    if (duplicateNotifications.length === 0) {
      this._notifications.next([...this._notifications.getValue(), notification]);
    }
  }

  public removeNotification(notification: Notification): void {
    this._notifications.next(this._notifications.getValue().filter((n) => n !== notification));
  }
}
